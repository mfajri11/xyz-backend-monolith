package http

import (
	"bytes"
	"encoding/json"
	"net/http"
)

var (
	DefaultHttpClient *HTTPClient
)

func init() {
	DefaultHttpClient = &HTTPClient{cl: http.DefaultClient}
}

type Option func(*http.Request)

func WithHeader(key, value string) Option {
	return func(r *http.Request) {
		r.Header.Set(key, value)
	}
}

func WithBasicAuth(username, password string) Option {
	return func(r *http.Request) {
		r.SetBasicAuth(username, password)
	}
}

func WithBearerToken(token string) Option {
	return func(r *http.Request) {
		r.Header.Set("Authorization", "Bearer "+token)
	}
}

func WithHeaders(headers map[string]string) Option {
	return func(r *http.Request) {
		for k, v := range headers {
			r.Header.Set(k, v)
		}
	}
}

func Get(url string, options ...Option) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	for _, option := range options {
		option(req)
	}
	return DefaultHttpClient.cl.Do(req)
}

func Post(url string, body interface{}, options ...Option) (*http.Response, error) {
	var (
		err       error
		bodyBytes []byte
	)
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	for _, option := range options {
		option(req)
	}

	return DefaultHttpClient.cl.Do(req)
}
