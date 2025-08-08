package xyo

import (
	"fmt"
	"net/http"
)

type ClientConfig struct {
	APIKey string
}

type Client struct {
	httpClient *http.Client
	config     *ClientConfig
}

func NewClient(opt *ClientConfig) *Client {
	return &Client{
		httpClient: http.DefaultClient,
		config:     opt,
	}
}

func (c *Client) Health() (err error) {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.xyo.financial/healthz",
		nil,
	)
	if err != nil {
		return
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return fmt.Errorf("health check failed with status code %d", resp.StatusCode)
}
