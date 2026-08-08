package xyo

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
	for _, req := range reqs {
		if req == nil {
			continue
		}
		content := req.Content
		cc := req.CountryCode
		items = append(items, openapi.EnrichTransactionsRequestInner{
			Content:     &content,
			CountryCode: &cc,
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

// GetEnrichmentStatus returns the processing status of a bulk enrichment job.
func (c *client) GetEnrichmentStatus(ctx context.Context, id string) (EnrichmentCollectionStatus, error) {
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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return nil, fmt.Errorf("xyo: DownloadEnrichmentCollection: build request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("xyo: DownloadEnrichmentCollection: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("xyo: DownloadEnrichmentCollection: status %d", resp.StatusCode)
	}

	gzReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("xyo: DownloadEnrichmentCollection: gzip stream: %w", err)
	}
	defer func() { _ = gzReader.Close() }()

	tarReader := tar.NewReader(gzReader)
	var results []*EnrichmentResponse

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("xyo: DownloadEnrichmentCollection: tar next: %w", err)
		}
		if header.Typeflag != tar.TypeReg {
			continue
		}
		var result EnrichmentResponse
		if err := json.NewDecoder(tarReader).Decode(&result); err != nil {
			return nil, fmt.Errorf("xyo: DownloadEnrichmentCollection: decode json from %s: %w", header.Name, err)
		}
		results = append(results, &result)
	}

	return results, nil
}
