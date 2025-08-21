package xyo

import (
	"net/http"
)

type ClientConfig struct {
	APIKey string
}

type Client interface {
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
