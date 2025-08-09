package xyo

import (
	"fmt"
	"net/http"
)

type ClientConfig struct {
	APIKey string
}

type Client interface {
	Health() error

	Enrichment
}

type httpClient struct {
	request func(req *http.Request) (*http.Response, error)
}

type internalClient struct {
	httpClient *httpClient
	config     *ClientConfig
}

// NewClient will accept ClientConfig struct where APIKey is defined
// Client is required to access Enrichment Services through SDK
func NewClient(opt *ClientConfig) Client {
	return &internalClient{
		httpClient: &httpClient{
			request: http.DefaultClient.Do,
		},
		config: opt,
	}
}

// Health is an indicator of overall API Status
// This could be done before sending a request or unusual response from API
func (c *internalClient) Health() error {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.xyo.financial/healthz",
		nil,
	)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.request(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return fmt.Errorf("health check failed with status code %d", resp.StatusCode)
}
