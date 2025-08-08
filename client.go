package xyo

import "net/http"

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
