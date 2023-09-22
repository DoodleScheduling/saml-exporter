//go:build unit

package transport

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/tj/assert"
)

func TestMetrics(t *testing.T) {
	requestResponse.Reset()
	okMock := NewMock(&http.Response{
		StatusCode: 200,
	}, nil)

	metrics := NewMetrics(okMock)
	res, err := metrics.RoundTrip(&http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{},
	})
	assert.Equal(t, 200, res.StatusCode)
	assert.NoError(t, err)

	reg := prometheus.NewRegistry()
	reg.MustRegister(metrics)
	assert.NoError(t, testutil.GatherAndCompare(reg, strings.NewReader(`
	# HELP http_client_request_total HTTP client request
	# TYPE http_client_request_total counter
	http_client_request_total{code="200",host="",method="GET"} 1
	`)))

	requestResponse.Reset()
	errMock := NewMock(&http.Response{
		StatusCode: 500,
	}, errors.New("random error"))

	metrics = NewMetrics(errMock)
	res, err = metrics.RoundTrip(&http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{},
	})
	assert.Equal(t, 500, res.StatusCode)
	assert.Error(t, err)

	reg = prometheus.NewRegistry()
	reg.MustRegister(metrics)
	assert.NoError(t, testutil.GatherAndCompare(reg, strings.NewReader(`
	# HELP http_client_request_total HTTP client request
	# TYPE http_client_request_total counter
	http_client_request_total{code="0",host="",method="POST"} 1
`)))
}
