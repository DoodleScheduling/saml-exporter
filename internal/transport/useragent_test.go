//go:build unit

package transport

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/tj/assert"
)

func TestUserAgent(t *testing.T) {
	okMock := NewMock(&http.Response{
		StatusCode: 200,
	}, nil)

	agent := NewUserAgent("test-agent", okMock)
	req := &http.Request{
		URL: &url.URL{},
	}

	res, err := agent.RoundTrip(req)
	assert.Equal(t, 200, res.StatusCode)
	assert.NoError(t, err)
	assert.Equal(t, "test-agent", req.Header.Get("User-Agent"))

	errMock := NewMock(&http.Response{
		StatusCode: 500,
	}, errors.New("random error"))

	agent = NewUserAgent("test-agent", errMock)
	req = &http.Request{
		URL: &url.URL{},
	}

	res, err = agent.RoundTrip(req)
	assert.Equal(t, 500, res.StatusCode)
	assert.Error(t, err)
	assert.Equal(t, "test-agent", req.Header.Get("User-Agent"))
}
