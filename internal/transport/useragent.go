package transport

import (
	"net/http"
)

type userAgent struct {
	name string
	next http.RoundTripper
}

func NewUserAgent(name string, next http.RoundTripper) *userAgent {
	return &userAgent{
		name: name,
		next: next,
	}
}

func (p *userAgent) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Header == nil {
		req.Header = make(http.Header)
	}

	req.Header.Set("User-Agent", p.name)
	res, err := p.next.RoundTrip(req)
	return res, err
}
