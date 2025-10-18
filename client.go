package xyo

import (
	"net/http"
)

const apiBasePath string = "https://api.xyo.financial"

// ClientConfig should be used to for a Client creation via NewClient method
type ClientConfig struct {
	// APIKey can be obtained from client dashboard located at: https://xyo.financial/dashboard
	APIKey     string
	httpClient *httpClient
	apiBaseURL string
}

type Client interface {
	Enrichment
}

type httpClient struct {
	request func(req *http.Request) (*http.Response, error)
}

type client struct {
	config *ClientConfig
}

// NewClient will accept ClientConfig struct where APIKey is defined
// Client is required to access Enrichment Services through SDK
func NewClient(opt *ClientConfig) Client {
	return &client{
		config: &ClientConfig{
			APIKey: opt.APIKey,
			httpClient: &httpClient{
				request: http.DefaultClient.Do,
			},
			apiBaseURL: apiBasePath,
		},
	}
}
