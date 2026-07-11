package xyo

import (
	"net/http"
	"time"
)

const defaultAPIBaseURL = "https://api.xyo.financial"
const defaultTimeout = 30 * time.Second

// ClientConfig configures the XYO API client.
type ClientConfig struct {
	// APIKey is required. Obtain from https://xyo.financial/dashboard
	APIKey string

	// BaseURL overrides the default API base URL. Leave empty to use the default.
	// Useful for pointing at staging or a local mock server during testing.
	BaseURL string

	// HTTPClient overrides the default HTTP client. Leave nil to use a client
	// with a 30-second timeout — the production-safe default.
	HTTPClient *http.Client
}

// Client is the entry point for interacting with the XYO.Financial enrichment API.
type Client interface {
	Enrichment
}

// httpClient wraps the transport function to allow injection in tests.
type httpClient struct {
	request func(req *http.Request) (*http.Response, error)
}

type client struct {
	apiKey     string
	apiBaseURL string
	http       *httpClient
}

// NewClient creates a new XYO API client from the provided configuration.
// It panics if config is nil or APIKey is empty — both are programming errors,
// not runtime conditions.
func NewClient(config *ClientConfig) Client {
	if config == nil {
		panic("xyo: NewClient called with nil config")
	}
	if config.APIKey == "" {
		panic("xyo: NewClient called with empty APIKey")
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = defaultAPIBaseURL
	}

	httpCl := config.HTTPClient
	if httpCl == nil {
		httpCl = &http.Client{
			Timeout: defaultTimeout,
		}
	}

	return &client{
		apiKey:     config.APIKey,
		apiBaseURL: baseURL,
		http:       &httpClient{request: httpCl.Do},
	}
}
