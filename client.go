package xyo

import (
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/xyo-financial/sdk-go/v2/openapi"
)

const defaultAPIBaseURL = "https://api.xyo.financial"
const defaultTimeout = 30 * time.Second

// DefaultEnterpriseTransport provides a high-throughput, connection-pooled HTTP transport
// tuned specifically for high-concurrency enterprise microservices.
//
// Unlike Go's http.DefaultTransport which restricts MaxIdleConnsPerHost to 2,
// DefaultEnterpriseTransport allocates up to 100 idle persistent connections per host,
// preventing TCP TIME_WAIT socket thrashing and ephemeral port exhaustion under heavy load.
var DefaultEnterpriseTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	MaxIdleConnsPerHost:   100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

// NewDefaultHTTPClient returns a production-ready *http.Client configured with a 30-second
// timeout and the connection-pooled DefaultEnterpriseTransport.
func NewDefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout:   defaultTimeout,
		Transport: DefaultEnterpriseTransport.Clone(),
	}
}

// ClientConfig configures the XYO API client.
type ClientConfig struct {
	// APIKey is required. Obtain from https://xyo.financial/dashboard
	APIKey string

	// BaseURL overrides the default API base URL. Leave empty to use the default.
	// Useful for pointing at staging or a local mock server during testing.
	BaseURL string

	// HTTPClient overrides the default HTTP client. Leave nil to use NewDefaultHTTPClient()
	// with a 30-second timeout and enterprise connection pooling.
	HTTPClient *http.Client
}

// Config is an alias for ClientConfig.
type Config = ClientConfig

// Client is the entry point for interacting with the XYO.Financial enrichment API.
type Client interface {
	Enrichment
}

type client struct {
	apiBaseURL string
	apiClient  *openapi.APIClient
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
		httpCl = NewDefaultHTTPClient()
	}

	cfg := openapi.NewConfiguration()
	cfg.Debug = false
	cfg.UserAgent = DefaultUserAgent
	cfg.Servers = openapi.ServerConfigurations{
		{
			URL: baseURL,
		},
	}
	cfg.HTTPClient = httpCl
	cfg.AddDefaultHeader("Authorization", "Bearer "+config.APIKey)

	apiClient := openapi.NewAPIClient(cfg)

	return &client{
		apiBaseURL: baseURL,
		apiClient:  apiClient,
	}, nil
}
