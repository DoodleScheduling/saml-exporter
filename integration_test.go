//go:build integration

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/tj/assert"
)

type testContainer struct {
	testcontainers.Container
	URI string
}

type integrationTest struct {
	name            string
	image           string
	expectedMetrics map[string]int
}

func TestMain(t *testing.T) {
	tests := []integrationTest{
		{
			name:  "integration test using kristophjunge/test-saml-idp",
			image: "kristophjunge/test-saml-idp:latest",
			expectedMetrics: map[string]int{
				"http_client_request":      1,
				"saml_x509_cert_not_after": 2,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			executeIntegrationTest(t, test)
		})
	}
}

func executeIntegrationTest(t *testing.T, test integrationTest) {
	container, err := setupContainer(context.TODO(), test.image)
	assert.NoError(t, err)

	os.Setenv("URL", container.URI)
	go func() {
		main()
	}()

	//binding is blocking, do this async but wait 200ms for tcp port to be open
	time.Sleep(200 * time.Millisecond)
	resp, err := http.Get("http://localhost:9412/metrics")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	d := expfmt.NewDecoder(resp.Body, expfmt.ResponseFormat(resp.Header))
	metricsFound := 0
	famMetrinsFound := 0

	for {
		fam := io_prometheus_client.MetricFamily{}
		if err = d.Decode(&fam); err != nil {
			break
		}

		if i, ok := test.expectedMetrics[fam.GetName()]; ok {
			metricsFound += i
			famMetrinsFound++
		}
	}

	assert.Equal(t, 3, metricsFound)
	assert.Len(t, test.expectedMetrics, famMetrinsFound)

	//tear down http server and unregister collector
	assert.NoError(t, srv.Shutdown(context.TODO()))
	prometheus.Unregister(promCollector)
}

func setupContainer(ctx context.Context, image string) (*testContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        image,
		ExposedPorts: []string{"8080/tcp"},
		WaitingFor:   wait.ForListeningPort("8080"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return nil, err
	}

	ip, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "8080")
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("http://%s:%s/simplesaml/saml2/idp/metadata.php", ip, mappedPort.Port())

	return &testContainer{Container: container, URI: uri}, nil
}
