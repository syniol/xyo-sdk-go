package xyo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/syniol/xyo-sdk-go/internal"
)

// EnrichmentRequest is a request data structure used for single and collection enrichment
type EnrichmentRequest struct {
	// Content is a maximum of 128 characters long payment description
	Content string `field:"required" json:"content"`
	// CountryCode ISO 3166-1 alpha-2 (Two characters format)
	CountryCode string `field:"required" json:"countryCode"`
}

// EnrichmentResponse is a result of payment transaction enrichment
type EnrichmentResponse struct {
	// Merchant is a name of merchant
	Merchant string `field:"required" json:"merchant"`
	// Description A brief description about the merchant
	Description string `field:"required" json:"description"`
	// Categories any type of categories fitting the description of the Merchant
	Categories []string `field:"required" json:"categories"`
	// Logo is base64 encoded png or jpeg representing the logo of Merchant
	Logo string `field:"required" json:"logo"`
	// Location describes the country, city. This is an optional field that could be null
	Location string `field:"optional" json:"location"`
	// Address describes exact address of purchase. This is an optional field that could be null
	Address string `field:"optional" json:"address"`
}

// EnrichTransactionCollectionResponse is a result of bulk enrichment
type EnrichTransactionCollectionResponse struct {
	// ID is a work ID for an enrichment request
	ID string `field:"required" json:"id"`
	// Link is a downloadable tar.gz Compressed file
	Link string `field:"required" json:"link"`
}

// EnrichmentCollectionStatus represents the status of EnrichTransactionCollectionResponse
// Currently there are three possible associated enum values for the status
type EnrichmentCollectionStatus string

const (
	EnrichmentCollectionStatusReady   EnrichmentCollectionStatus = "READY"
	EnrichmentCollectionStatusFailed  EnrichmentCollectionStatus = "FAILED"
	EnrichmentCollectionStatusPending EnrichmentCollectionStatus = "PENDING"
)

// EnrichmentCollectionStatusResponse provides a status of bulk enrichment
type EnrichmentCollectionStatusResponse struct {
	// Status could be READY, PENDING, FAILED
	Status EnrichmentCollectionStatus `field:"required" json:"status"`
}

type Enrichment interface {
	EnrichTransaction(enrichmentReq *EnrichmentRequest) (*EnrichmentResponse, error)
	EnrichTransactionCollection(enrichmentReq []*EnrichmentRequest) (*EnrichTransactionCollectionResponse, error)
	EnrichTransactionCollectionStatus(ID string) (EnrichmentCollectionStatus, error)
}

// EnrichTransaction should be used for a single transaction enrichment
func (c *client) EnrichTransaction(enrichmentReq *EnrichmentRequest) (*EnrichmentResponse, error) {
	requestBody, err := json.Marshal(enrichmentReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/v1/ai/finance/enrichment/transaction", c.config.apiBaseURL),
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

// EnrichTransactionCollection should be used for bulk enrichment request
// Unlike EnrichTransaction this method won;t produce the enrichment result instantly
// However, it allows you to download the enrichment results when it becomes available using download link in response
// Status of Downloadable enrichment result can be queried using EnrichTransactionCollectionStatus
func (c *client) EnrichTransactionCollection(enrichmentReq []*EnrichmentRequest) (*EnrichTransactionCollectionResponse, error) {
	requestBody, err := json.Marshal(enrichmentReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/v1/ai/finance/enrichment/transactions", c.config.apiBaseURL),
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

// EnrichTransactionCollectionStatus returns the status of request from EnrichTransactionCollection
// ID is the value of ID taken from EnrichTransactionCollection response
func (c *client) EnrichTransactionCollectionStatus(ID string) (EnrichmentCollectionStatus, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/v1/ai/finance/enrichment/transactions/status/%s", c.config.apiBaseURL, ID),
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
