package transport

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestResponse = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_client_request_total",
		Help: "HTTP client request",
	}, []string{"host", "code", "method"})
)

type metrics struct {
	next http.RoundTripper
}

func NewMetrics(next http.RoundTripper) *metrics {
	return &metrics{
		next: next,
	}
}

func (p *metrics) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := p.next.RoundTrip(req)

	statusCode := "0"
	if err == nil {
		statusCode = fmt.Sprintf("%d", res.StatusCode)
	}

	m := requestResponse.With(prometheus.Labels{
		"code":   statusCode,
		"method": req.Method,
		"host":   req.URL.Host,
	})
	m.Inc()

	return res, err
}

func (c *metrics) Describe(ch chan<- *prometheus.Desc) {
	requestResponse.Describe(ch)
}

func (c *metrics) Collect(ch chan<- prometheus.Metric) {
	requestResponse.Collect(ch)
}
