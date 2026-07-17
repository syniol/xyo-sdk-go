package xyo

import (
	"errors"
	"net/http"
	"strings"
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

type client struct {
	apiKey     string
	apiBaseURL string
	http       *http.Client
}

// NewClient creates a new XYO API client from the provided configuration.
// It returns an error if config is nil or APIKey is empty.
func NewClient(config *ClientConfig) (Client, error) {
	if config == nil {
		return nil, errors.New("xyo: NewClient called with nil config")
	}
	if config.APIKey == "" {
		return nil, errors.New("xyo: NewClient called with empty APIKey")
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = defaultAPIBaseURL
	}
	baseURL = strings.TrimRight(baseURL, "/")

	httpCl := config.HTTPClient
	if httpCl == nil {
		httpCl = &http.Client{
			Timeout: defaultTimeout,
		}
	}

	return &client{
		apiKey:     config.APIKey,
		apiBaseURL: baseURL,
		http:       httpCl,
	}, nil
}
