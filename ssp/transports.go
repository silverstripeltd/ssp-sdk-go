package ssp

import (
	"net/http"
)

// BasicAuthTransport is a RoundTripper that injects BasicAuth credentials. It's based on a similar RoundTripper
// found in the github.com/google/go-github package.
type BasicAuthTransport struct {
	Username  string
	Password  string
	Transport http.RoundTripper
}

func (t *BasicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = cloneRequest(req)
	req.SetBasicAuth(t.Username, t.Password)
	return t.transport().RoundTrip(req)
}

func (t *BasicAuthTransport) Client() *http.Client {
	return &http.Client{Transport: t}
}

func (t *BasicAuthTransport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}

func cloneRequest(r *http.Request) *http.Request {
	r2 := new(http.Request)
	*r2 = *r
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return r2
}
