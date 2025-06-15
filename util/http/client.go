package http

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type HTTPClient struct {
	cl *http.Client
}

func NewClient() *HTTPClient {
	return &HTTPClient{cl: http.DefaultClient}
}

func NewClientWithHTTPClient(cl *http.Client) *HTTPClient {
	return &HTTPClient{cl: cl}
}

func (c *HTTPClient) Get(url string, options ...Option) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	for _, option := range options {
		option(req)
	}
	return c.cl.Do(req)
}

func (c *HTTPClient) Post(url string, body interface{}, options ...Option) (*http.Response, error) {
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

	return c.cl.Do(req)
}
