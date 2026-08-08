package xyo

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/xyo-financial/sdk-go/v2/openapi"
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

// Config is an alias for ClientConfig.
type Config = ClientConfig


// Client is the entry point for interacting with the XYO.Financial enrichment API.
type Client interface {
	Enrichment
}

type client struct {
	apiKey    string
	apiClient *openapi.APIClient
	http      *http.Client
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

	cfg := openapi.NewConfiguration()
	cfg.Servers = openapi.ServerConfigurations{
		{
			URL: baseURL,
		},
	}
	cfg.HTTPClient = httpCl
	cfg.AddDefaultHeader("Authorization", "Bearer "+config.APIKey)

	apiClient := openapi.NewAPIClient(cfg)

	return &client{
		apiKey:    config.APIKey,
		apiClient: apiClient,
		http:      httpCl,
	}, nil
}
