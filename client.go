package xyo

import "net/http"

type ClientConfig struct {
	APIKey string
}

type clientConnection struct {
	httpClient *http.Client
	config     ClientConfig
}

type Client struct {
	clientConnection

	enrichmentRequester
}

func NewClient(opt *ClientConfig) *Client {
	return &Client{}
}
