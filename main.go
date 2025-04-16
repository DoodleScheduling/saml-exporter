package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	"github.com/doodlescheduling/saml-exporter/internal/collector"
	"github.com/doodlescheduling/saml-exporter/internal/transport"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sethvargo/go-envconfig"
	flag "github.com/spf13/pflag"
	"go.uber.org/zap"
)

type Config struct {
	URL []string `env:"URL"`
	Log struct {
		Level    string `env:"LOG_LEVEL, default=info"`
		Encoding string `env:"LOG_ENCODING, default=json"`
	}
	Bind        string `env:"BIND, default=:9412"`
	MetricsPath string `env:"METRICS_PATH, default=/metrics"`
	HealthPath  string `env:"HEALTH_PATH, default=/healthz"`
	UserAgent   string `env:"USER_AGENT, default=saml-exporter (go-http-client)"`
}

var (
	config        = &Config{}
	srv           *http.Server
	promCollector *collector.Collector
)

func init() {
	flag.StringVarP(&config.Log.Level, "log-level", "l", "", "Define the log level (default is warning) [debug,info,warn,error]")
	flag.StringVarP(&config.Log.Encoding, "log-encoding", "e", "", "Define the log format (default is json) [json,console]")
	flag.StringVarP(&config.Bind, "bind", "b", "", "Address to bind http server (default is :9412)")
	flag.StringVarP(&config.MetricsPath, "metrics-path", "", "", "Metric path (default is /metrics)")
	flag.StringVarP(&config.HealthPath, "health-path", "", "", "Health probe path (default is /healthz)")
	flag.StringVarP(&config.UserAgent, "user-agent", "", "", "HTTP client user-agent header")
}

func main() {
	ctx := context.Background()
	if err := envconfig.Process(ctx, config); err != nil {
		log.Fatal(err)
	}

	flag.Parse()

	logger, err := buildLogger()
	must(err)
	if len(flag.Args()) > 0 {
		config.URL = flag.Args()
	}

	urls, err := buildURL()
	must(err)

	prometheus.MustRegister(version.NewCollector("saml_exporter"))

	client := http.DefaultClient
	metricsRoundTripper := transport.NewMetrics(http.DefaultTransport)
	client.Transport = transport.NewLogger(logger, transport.NewUserAgent(config.UserAgent, metricsRoundTripper))
	promCollector = collector.New(logger, http.DefaultClient, urls, metricsRoundTripper)
	prometheus.MustRegister(promCollector)

	logger.Info("starting http server...", "bind", config.Bind)
	srv = buildHTTPServer(prometheus.DefaultGatherer)
	err = srv.ListenAndServe()

	// Only panic if we have a net error
	if _, ok := err.(*net.OpError); ok {
		panic(err)
	} else {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
	}
}

func buildURL() ([]*url.URL, error) {
	if len(config.URL) == 0 {
		return nil, errors.New("at least one url to a saml metadata endpoint is required")
	}

	var result []*url.URL

	for _, u := range config.URL {
		res, err := url.Parse(u)
		if err != nil {
			return nil, err
		}

		result = append(result, res)
	}

	return result, nil
}

func buildLogger() (logr.Logger, error) {
	logOpts := zap.NewDevelopmentConfig()
	logOpts.Encoding = config.Log.Encoding

	err := logOpts.Level.UnmarshalText([]byte(config.Log.Level))
	if err != nil {
		return logr.Discard(), err
	}

	zapLog, err := logOpts.Build()
	if err != nil {
		return logr.Discard(), err
	}

	return zapr.NewLogger(zapLog), nil
}

// Run executes a blocking http server. Starts the http listener with the metrics and healthz endpoints.
func buildHTTPServer(reg prometheus.Gatherer) *http.Server {
	mux := http.NewServeMux()

	if config.MetricsPath != "/" {
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, fmt.Sprintf("Use the %s endpoint", config.MetricsPath), http.StatusOK)
		})
	}

	mux.HandleFunc(config.HealthPath, func(w http.ResponseWriter, r *http.Request) { http.Error(w, "OK", http.StatusOK) })
	mux.HandleFunc(config.MetricsPath, func(w http.ResponseWriter, r *http.Request) {
		promhttp.HandlerFor(reg, promhttp.HandlerOpts{}).ServeHTTP(w, r)
	})

	srv := http.Server{Addr: config.Bind, Handler: mux}
	return &srv
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
