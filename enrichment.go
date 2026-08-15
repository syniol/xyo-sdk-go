package xyo

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/xyo-financial/sdk-go/v2/openapi"
)

// EnrichmentRequest is the request payload for single and bulk transaction enrichment.
type EnrichmentRequest struct {
	// Content is the payment description, maximum 128 characters.
	Content string `json:"content"`
	// CountryCode is the ISO 3166-1 alpha-2 two-character country code.
	CountryCode string `json:"countryCode"`
}

// EnrichmentResponse is the result of a single payment transaction enrichment.
type EnrichmentResponse struct {
	// Merchant is the name of the merchant.
	Merchant string `json:"merchant"`
	// Description is a brief description of the merchant.
	Description string `json:"description"`
	// Categories lists categories fitting the description of the merchant.
	Categories []string `json:"categories"`
	// Logo is a base64-encoded PNG or JPEG representing the merchant logo.
	Logo string `json:"logo"`
	// Location describes the country and city. May be empty if the API returns null.
	Location string `json:"location"`
	// Address describes the exact address of purchase. May be empty if the API returns null.
	Address string `json:"address"`
}

// EnrichTransactionCollectionResponse is the result of a bulk enrichment submission.
type EnrichTransactionCollectionResponse struct {
	// ID is the work ID for the enrichment request.
	ID string `json:"id"`
	// Link is the URL to the downloadable tar.gz results archive.
	Link string `json:"link"`
}

// EnrichmentCollectionStatus represents the processing state of a bulk enrichment job.
type EnrichmentCollectionStatus string

const (
	EnrichmentCollectionStatusReady   EnrichmentCollectionStatus = "READY"
	EnrichmentCollectionStatusFailed  EnrichmentCollectionStatus = "FAILED"
	EnrichmentCollectionStatusPending EnrichmentCollectionStatus = "PENDING"
)

// Enrichment defines the contract for all XYO transaction enrichment operations.
type Enrichment interface {
	// EnrichTransaction enriches a single payment transaction synchronously.
	EnrichTransaction(ctx context.Context, req *EnrichmentRequest) (*EnrichmentResponse, error)

	// EnrichTransactions submits a bulk enrichment request asynchronously
	// and returns a job ID and download link.
	EnrichTransactions(ctx context.Context, reqs []*EnrichmentRequest) (*EnrichTransactionCollectionResponse, error)

	// GetEnrichmentStatus returns the processing status of a bulk enrichment job.
	GetEnrichmentStatus(ctx context.Context, id string) (EnrichmentCollectionStatus, error)

	// --- Backward-compatible aliases (retained for existing integrations) ---

	// EnrichTransactionCollection is an alias for EnrichTransactions.
	EnrichTransactionCollection(ctx context.Context, reqs []*EnrichmentRequest) (*EnrichTransactionCollectionResponse, error)

	// EnrichTransactionCollectionStatus is an alias for GetEnrichmentStatus.
	EnrichTransactionCollectionStatus(ctx context.Context, id string) (EnrichmentCollectionStatus, error)

	// DownloadEnrichmentCollection downloads and decodes a bulk enrichment result tarball.
	DownloadEnrichmentCollection(ctx context.Context, downloadURL string) ([]*EnrichmentResponse, error)
}

// EnrichTransaction enriches a single payment transaction.
func (c *client) EnrichTransaction(ctx context.Context, req *EnrichmentRequest) (*EnrichmentResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("xyo: EnrichTransaction: request is nil")
	}

	genReq := openapi.NewEnrichmentRequest(req.Content, req.CountryCode)
	resp, httpResp, err := c.apiClient.EnrichmentAPI.EnrichTransaction(ctx).EnrichmentRequest(*genReq).Execute()
	if err != nil {
		return nil, parseOpenAPIError(err, "EnrichTransaction", httpResp)
	}

	return &EnrichmentResponse{
		Merchant:    resp.GetMerchant(),
		Description: resp.GetDescription(),
		Categories:  resp.GetCategories(),
		Logo:        resp.GetLogo(),
		Location:    resp.GetLocation(),
		Address:     resp.GetAddress(),
	}, nil
}

// EnrichTransactions submits a bulk enrichment request asynchronously.
func (c *client) EnrichTransactions(ctx context.Context, reqs []*EnrichmentRequest) (*EnrichTransactionCollectionResponse, error) {
	items := make([]openapi.EnrichTransactionsRequestInner, 0, len(reqs))
	for i, req := range reqs {
		if req == nil {
			return nil, fmt.Errorf("xyo: EnrichTransactions: request at index %d is nil", i)
		}
		items = append(items, openapi.EnrichTransactionsRequestInner{
			Content:     &req.Content,
			CountryCode: &req.CountryCode,
		})
	}

	resp, httpResp, err := c.apiClient.EnrichmentAPI.EnrichTransactions(ctx).
		EnrichTransactionsRequestInner(items).Execute()
	if err != nil {
		return nil, parseOpenAPIError(err, "EnrichTransactions", httpResp)
	}

	return &EnrichTransactionCollectionResponse{
		ID:   resp.GetId(),
		Link: resp.GetLink(),
	}, nil
}

const (
	// DefaultMaxTarEntries is the maximum number of entries processed from a bulk enrichment tarball.
	DefaultMaxTarEntries = 50000
	// DefaultMaxEntryBytes is the maximum allowed size (in bytes) for a single JSON file within the tarball (10 MiB).
	DefaultMaxEntryBytes = 10 * 1024 * 1024
	// DefaultMaxArchiveBytes is the maximum total uncompressed bytes allowed across the tarball stream (100 MiB).
	DefaultMaxArchiveBytes = 100 * 1024 * 1024
)

// GetEnrichmentStatus returns the processing status of a bulk enrichment job.
func (c *client) GetEnrichmentStatus(ctx context.Context, id string) (EnrichmentCollectionStatus, error) {
	if id == "" {
		return "", fmt.Errorf("xyo: GetEnrichmentStatus: id cannot be empty")
	}

	resp, httpResp, err := c.apiClient.EnrichmentAPI.GetEnrichmentStatus(ctx, id).Execute()
	if err != nil {
		return "", parseOpenAPIError(err, "GetEnrichmentStatus", httpResp)
	}

	return EnrichmentCollectionStatus(resp.GetStatus()), nil
}

// --- Backward-compatible alias methods ---

// EnrichTransactionCollection is an alias for EnrichTransactions.
func (c *client) EnrichTransactionCollection(ctx context.Context, reqs []*EnrichmentRequest) (*EnrichTransactionCollectionResponse, error) {
	return c.EnrichTransactions(ctx, reqs)
}

// EnrichTransactionCollectionStatus is an alias for GetEnrichmentStatus.
func (c *client) EnrichTransactionCollectionStatus(ctx context.Context, id string) (EnrichmentCollectionStatus, error) {
	return c.GetEnrichmentStatus(ctx, id)
}

// DownloadEnrichmentCollection downloads and decodes a bulk enrichment result tarball
// from the URL returned by EnrichTransactions.
func (c *client) DownloadEnrichmentCollection(ctx context.Context, downloadURL string) ([]*EnrichmentResponse, error) {
	if downloadURL == "" {
		return nil, fmt.Errorf("xyo: DownloadEnrichmentCollection: downloadURL cannot be empty")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return nil, fmt.Errorf("xyo: DownloadEnrichmentCollection: build request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if userAgent := c.apiClient.GetConfig().UserAgent; userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}

	// Only attach Authorization header if download URL matches the configured API host
	parsedDownloadURL, parseDownloadErr := url.Parse(downloadURL)
	if parseDownloadErr == nil && parsedDownloadURL.Host != "" {
		parsedBaseURL, parseBaseErr := url.Parse(c.apiBaseURL)
		if parseBaseErr == nil && parsedBaseURL.Host != "" {
			if strings.EqualFold(parsedDownloadURL.Host, parsedBaseURL.Host) {
				if authHeader, ok := c.apiClient.GetConfig().DefaultHeader["Authorization"]; ok {
					req.Header.Set("Authorization", authHeader)
				}
			}
		}
	}

	httpCl := c.apiClient.GetConfig().HTTPClient
	if httpCl == nil {
		httpCl = &http.Client{Timeout: defaultTimeout}
	}

	resp, err := httpCl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("xyo: DownloadEnrichmentCollection: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		var errResp ErrorResponse
		if json.Unmarshal(body, &errResp) == nil && len(errResp.Errors) > 0 {
			errResp.HTTPStatusCode = resp.StatusCode
			return nil, fmt.Errorf("xyo: DownloadEnrichmentCollection: %w", &errResp)
		}
		return nil, fmt.Errorf("xyo: DownloadEnrichmentCollection: status %d", resp.StatusCode)
	}

	// Limit total compressed stream and decompressed stream to prevent decompression bombs
	limitedBody := io.LimitReader(resp.Body, DefaultMaxArchiveBytes)
	gzReader, err := gzip.NewReader(limitedBody)
	if err != nil {
		return nil, fmt.Errorf("xyo: DownloadEnrichmentCollection: gzip stream: %w", err)
	}
	defer func() { _ = gzReader.Close() }()

	limitedTarStream := io.LimitReader(gzReader, DefaultMaxArchiveBytes)
	tarReader := tar.NewReader(limitedTarStream)
	var results []*EnrichmentResponse

	for entryCount := 0; ; entryCount++ {
		if entryCount >= DefaultMaxTarEntries {
			return nil, fmt.Errorf("xyo: DownloadEnrichmentCollection: tarball contains too many entries (exceeds limit of %d)", DefaultMaxTarEntries)
		}

		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("xyo: DownloadEnrichmentCollection: tar next: %w", err)
		}
		if header.Typeflag != tar.TypeReg {
			continue
		}
		if header.Size > DefaultMaxEntryBytes {
			return nil, fmt.Errorf("xyo: DownloadEnrichmentCollection: entry %q exceeds maximum allowed size (%d bytes > %d bytes)", header.Name, header.Size, DefaultMaxEntryBytes)
		}

		var result EnrichmentResponse
		entryReader := io.LimitReader(tarReader, DefaultMaxEntryBytes)
		if err := json.NewDecoder(entryReader).Decode(&result); err != nil {
			return nil, fmt.Errorf("xyo: DownloadEnrichmentCollection: decode json from %s: %w", header.Name, err)
		}
		results = append(results, &result)
	}

	return results, nil
}
