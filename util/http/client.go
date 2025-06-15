package http

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type HTTPClient struct {
	cl      Doer
	BaseURL string
	apiKey  string
	appID   string
}

type HTTPClientMock struct {
	doer Doer
}

func NewClient(baseURL, apiKey, appID string) *HTTPClient {
	return &HTTPClient{
		cl:      http.DefaultClient,
		BaseURL: baseURL,
		apiKey:  apiKey,
		appID:   appID,
	}
}

func NewClientWithHTTPClient(cl *http.Client, baseURL, apiKey, appID string) *HTTPClient {
	return &HTTPClient{
		cl:      http.DefaultClient,
		BaseURL: baseURL,
		apiKey:  apiKey,
		appID:   appID,
	}
}

func (c *HTTPClient) GetAPIKey() string {
	return c.apiKey
}

func (c *HTTPClient) GetAPIID() string {
	return c.appID
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
