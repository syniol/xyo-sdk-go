package xyo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/syniol/xyo-sdk-go/internal"
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

// enrichmentCollectionStatusResponse is an internal deserialisation wrapper.
type enrichmentCollectionStatusResponse struct {
	Status EnrichmentCollectionStatus `json:"status"`
}

// Enrichment defines the contract for all XYO transaction enrichment operations.
type Enrichment interface {
	EnrichTransaction(req *EnrichmentRequest) (*EnrichmentResponse, error)
	EnrichTransactionCollection(reqs []*EnrichmentRequest) (*EnrichTransactionCollectionResponse, error)
	EnrichTransactionCollectionStatus(id string) (EnrichmentCollectionStatus, error)
}

// apiError reads the response body and returns a formatted error string.
// It always closes the body.
func apiError(resp *http.Response, op string) error {
	var body string
	if resp.Body != nil {
		defer resp.Body.Close()
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		body = string(bytes.TrimSpace(b))
	}
	if body != "" {
		return fmt.Errorf("xyo: %s: status %d: %s", op, resp.StatusCode, body)
	}
	return fmt.Errorf("xyo: %s: status %d", op, resp.StatusCode)
}

// EnrichTransaction enriches a single payment transaction.
func (c *client) EnrichTransaction(enrichmentReq *EnrichmentRequest) (*EnrichmentResponse, error) {
	requestBody, err := json.Marshal(enrichmentReq)
	if err != nil {
		return nil, fmt.Errorf("xyo: enrich transaction: marshal request: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/v1/ai/finance/enrichment/transaction", c.apiBaseURL),
		bytes.NewReader(requestBody),
	)
	if err != nil {
		return nil, fmt.Errorf("xyo: enrich transaction: build request: %w", err)
	}

	internal.MandatoryAPIHeaders(req, c.apiKey)

	resp, err := c.http.request(req)
	if err != nil {
		return nil, fmt.Errorf("xyo: enrich transaction: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp, "enrich transaction")
	}
	defer resp.Body.Close()

	var result EnrichmentResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("xyo: enrich transaction: decode response: %w", err)
	}

	return &result, nil
}

// EnrichTransactionCollection submits a bulk enrichment request.
// Unlike EnrichTransaction, this method won't produce enrichment results instantly.
// It returns a job ID and download link; use EnrichTransactionCollectionStatus to
// poll for completion.
func (c *client) EnrichTransactionCollection(enrichmentReq []*EnrichmentRequest) (*EnrichTransactionCollectionResponse, error) {
	requestBody, err := json.Marshal(enrichmentReq)
	if err != nil {
		return nil, fmt.Errorf("xyo: enrich transaction collection: marshal request: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/v1/ai/finance/enrichment/transactions", c.apiBaseURL),
		bytes.NewReader(requestBody),
	)
	if err != nil {
		return nil, fmt.Errorf("xyo: enrich transaction collection: build request: %w", err)
	}

	internal.MandatoryAPIHeaders(req, c.apiKey)

	resp, err := c.http.request(req)
	if err != nil {
		return nil, fmt.Errorf("xyo: enrich transaction collection: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp, "enrich transaction collection")
	}
	defer resp.Body.Close()

	var result EnrichTransactionCollectionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("xyo: enrich transaction collection: decode response: %w", err)
	}

	return &result, nil
}

// EnrichTransactionCollectionStatus returns the processing status of a bulk enrichment job.
// id is the ID returned by EnrichTransactionCollection.
func (c *client) EnrichTransactionCollectionStatus(id string) (EnrichmentCollectionStatus, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/v1/ai/finance/enrichment/transactions/status/%s", c.apiBaseURL, id),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("xyo: enrich transaction collection status: build request: %w", err)
	}

	internal.MandatoryAPIHeaders(req, c.apiKey)

	resp, err := c.http.request(req)
	if err != nil {
		return "", fmt.Errorf("xyo: enrich transaction collection status: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", apiError(resp, "enrich transaction collection status")
	}
	defer resp.Body.Close()

	var result enrichmentCollectionStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("xyo: enrich transaction collection status: decode response: %w", err)
	}

	return result.Status, nil
}
