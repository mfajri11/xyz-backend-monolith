package http

import "net/http"

type DoFunc func(*http.Request) (*http.Response, error)

func (f DoFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func NewMock(do DoFunc) *HTTPClient {
	return &HTTPClient{cl: do}
}
