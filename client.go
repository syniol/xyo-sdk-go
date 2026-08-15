package xyo

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/xyo-financial/sdk-go/v2/openapi"
)

const defaultAPIBaseURL = "https://api.xyo.financial"
const defaultTimeout = 30 * time.Second

// DefaultEnterpriseTransport provides a high-throughput, connection-pooled HTTP transport
// tuned specifically for high-concurrency enterprise microservices.
//
// Explicitly enforces minimum TLS 1.2 protocol version (PCI-DSS 4.0 §4.2.1 compliance).
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
	TLSClientConfig: &tls.Config{
		MinVersion: tls.VersionTLS12,
	},
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

	// APIKeySupplier provides a dynamic API key supplier function for runtime secret rotation (NIST SP 800-57).
	// If specified, it takes precedence over static APIKey.
	APIKeySupplier func() string

	// BaseURL overrides the default API base URL. Leave empty to use XYO_API_BASE_URL env or default https://api.xyo.financial.
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
	// Close releases any idle persistent HTTP connections held by the client.
	Close() error
}

type client struct {
	apiBaseURL  string
	keySupplier func() string
	apiClient   *openapi.APIClient
	httpCl      *http.Client
}

type authRoundTripper struct {
	base        http.RoundTripper
	apiBaseURL  string
	keySupplier func() string
}

func (a *authRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	clone := req.Clone(req.Context())
	if a.keySupplier != nil {
		if key := a.keySupplier(); key != "" {
			// Only attach Authorization header if target host matches configured API base URL host (prevents token leakage)
			attach := true
			if a.apiBaseURL != "" && req.URL != nil && req.URL.Host != "" {
				if baseURL, err := url.Parse(a.apiBaseURL); err == nil && baseURL.Host != "" {
					if !strings.EqualFold(req.URL.Host, baseURL.Host) {
						attach = false
					}
				}
			}
			if attach {
				clone.Header.Set("Authorization", "Bearer "+key)
			}
		}
	}
	base := a.base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(clone)
}

// Close gracefully closes any idle persistent HTTP connections.
func (c *client) Close() error {
	if c.httpCl != nil {
		if t, ok := c.httpCl.Transport.(*http.Transport); ok {
			t.CloseIdleConnections()
		} else if at, ok := c.httpCl.Transport.(*authRoundTripper); ok {
			if t, ok := at.base.(*http.Transport); ok {
				t.CloseIdleConnections()
			}
		}
	}
	return nil
}

// NewClient creates a new XYO API client from the provided configuration.
// It returns an error if config is nil or APIKey is empty.
func NewClient(config *ClientConfig) (Client, error) {
	if config == nil {
		return nil, errors.New("xyo: NewClient called with nil config")
	}

	keySupplier := config.APIKeySupplier
	if keySupplier == nil {
		if config.APIKey == "" {
			return nil, errors.New("xyo: NewClient called with empty APIKey")
		}
		key := config.APIKey
		keySupplier = func() string { return key }
	} else if key := keySupplier(); key == "" && config.APIKey == "" {
		return nil, errors.New("xyo: NewClient called with empty APIKey")
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		if envURL := os.Getenv("XYO_API_BASE_URL"); envURL != "" {
			baseURL = envURL
		} else {
			baseURL = defaultAPIBaseURL
		}
	}
	baseURL = strings.TrimRight(baseURL, "/")

	rawHTTPCl := config.HTTPClient
	if rawHTTPCl == nil {
		rawHTTPCl = NewDefaultHTTPClient()
	}

	// Wrap HTTP transport with dynamic auth header interceptor (NIST SP 800-57)
	transport := rawHTTPCl.Transport
	if transport == nil {
		transport = DefaultEnterpriseTransport.Clone()
	}
	wrappedTransport := &authRoundTripper{
		base:        transport,
		apiBaseURL:  baseURL,
		keySupplier: keySupplier,
	}
	httpCl := &http.Client{
		Timeout:   rawHTTPCl.Timeout,
		Transport: wrappedTransport,
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

	apiClient := openapi.NewAPIClient(cfg)

	return &client{
		apiBaseURL:  baseURL,
		keySupplier: keySupplier,
		apiClient:   apiClient,
		httpCl:      rawHTTPCl,
	}, nil
}
