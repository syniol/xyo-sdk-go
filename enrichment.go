package xyo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/syniol/xyo-sdk-go/internal"
	"net/http"
)

// EnrichmentRequest is a request data structure used for single and collection enrichment
// Content is a maximum of 128 characters long payment description
// CountryCode ISO 3166-1 alpha-2 (Two characters format)
type EnrichmentRequest struct {
	Content     string `json:"content"`
	CountryCode string `json:"countryCode"`
}

// EnrichmentResponse is a result of payment transaction enrichment
// Merchant is a name of merchant
// Description A brief description about the merchant
// Categories any type of categories fitting the description of the Merchant
// Logo is base64 encoded png or jpeg representing the logo of Merchant
type EnrichmentResponse struct {
	Merchant    string   `json:"merchant"`
	Description string   `json:"description"`
	Categories  []string `json:"categories"`
	Logo        string   `json:"logo"`
}

// EnrichTransactionCollectionResponse is a result of bulk enrichment
// ID is a work ID for an enrichment request
// Link is a downloadable tar.gz Compressed file
type EnrichTransactionCollectionResponse struct {
	ID   string `json:"id"`
	Link string `json:"link"`
}

// EnrichmentCollectionStatus represents the status of EnrichTransactionCollectionResponse
// Currently there are three possible associated enum values for the status
// READY, PENDING, FAILED
type EnrichmentCollectionStatus string

const (
	EnrichmentCollectionStatusReady   EnrichmentCollectionStatus = "READY"
	EnrichmentCollectionStatusFailed  EnrichmentCollectionStatus = "FAILED"
	EnrichmentCollectionStatusPending EnrichmentCollectionStatus = "PENDING"
)

type EnrichmentCollectionStatusResponse struct {
	Status EnrichmentCollectionStatus `json:"status"`
}

type Enrichment interface {
	// EnrichTransaction should be used for a single transaction enrichment
	EnrichTransaction(enrichmentReq *EnrichmentRequest) (*EnrichmentResponse, error)

	// EnrichTransactionCollection should be used for bulk enrichment request
	// Unlike EnrichTransaction this method won;t produce the enrichment result instantly
	// However, it allows you to download the enrichment results when it becomes available using download link in response
	// Status of Downloadable enrichment result can be queried using EnrichTransactionCollectionStatus
	EnrichTransactionCollection(enrichmentReq []*EnrichmentRequest) (*EnrichTransactionCollectionResponse, error)

	// EnrichTransactionCollectionStatus returns the status of request from EnrichTransactionCollection
	// ID is the value of ID taken from EnrichTransactionCollection response
	EnrichTransactionCollectionStatus(ID string) (EnrichmentCollectionStatus, error)
}

func (c *client) EnrichTransaction(enrichmentReq *EnrichmentRequest) (*EnrichmentResponse, error) {
	requestBody, err := json.Marshal(enrichmentReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/v1/ai/finance/enrichment/transaction", ApiBasePath),
		bytes.NewReader(requestBody),
	)
	if err != nil {
		return nil, err
	}

	internal.MandatoryAPIHeaders(req, c.config.APIKey)

	resp, err := c.config.httpClient.request(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("EnrichTransaction returned status code %d", resp.StatusCode)
	}

	var enrichmentResponse EnrichmentResponse
	err = json.NewDecoder(resp.Body).Decode(&enrichmentResponse)

	return &enrichmentResponse, err
}

func (c *client) EnrichTransactionCollection(enrichmentReq []*EnrichmentRequest) (*EnrichTransactionCollectionResponse, error) {
	requestBody, err := json.Marshal(enrichmentReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/v1/ai/finance/enrichment/transactions", ApiBasePath),
		bytes.NewReader(requestBody),
	)
	if err != nil {
		return nil, err
	}

	internal.MandatoryAPIHeaders(req, c.config.APIKey)

	resp, err := c.config.httpClient.request(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("enrich transaction collection returned status code: %d", resp.StatusCode)
	}

	var enrichTransactionCollectionResponse EnrichTransactionCollectionResponse
	err = json.NewDecoder(resp.Body).Decode(&enrichTransactionCollectionResponse)

	return &enrichTransactionCollectionResponse, err
}

func (c *client) EnrichTransactionCollectionStatus(ID string) (EnrichmentCollectionStatus, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/v1/ai/finance/enrichment/transactions/status/%s", ApiBasePath, ID),
		nil,
	)
	if err != nil {
		return "", err
	}

	internal.MandatoryAPIHeaders(req, c.config.APIKey)

	resp, err := c.config.httpClient.request(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("enrich transaction collection status returned status code: %d", resp.StatusCode)
	}

	var response EnrichmentCollectionStatusResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	return response.Status, nil
}
