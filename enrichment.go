package xyo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type EnrichmentRequest struct {
	Content     string `json:"content"`
	CountryCode string `json:"countryCode"`
}

type EnrichmentResponse struct {
	Merchant    string   `json:"merchant"`
	Description string   `json:"description"`
	Categories  []string `json:"categories"`
	Logo        string   `json:"logo"`
}

type EnrichTransactionCollectionResponse struct {
	ID   string `json:"id"`
	Link string `json:"link"`
}

type EnrichmentCollectionStatus string

const (
	EnrichmentCollectionStatusReady   EnrichmentCollectionStatus = "READY"
	EnrichmentCollectionStatusFailure EnrichmentCollectionStatus = "FAILED"
	EnrichmentCollectionStatusPending EnrichmentCollectionStatus = "PENDING"
)

type EnrichmentCollectionStatusResponse struct {
	Status EnrichmentCollectionStatus `json:"status"`
}

type Enrichment interface {
	EnrichTransaction(enrichmentReq *EnrichmentRequest) (*EnrichmentResponse, error)
	EnrichTransactionCollection(enrichmentReq []*EnrichmentRequest) (*EnrichTransactionCollectionResponse, error)
	EnrichTransactionCollectionStatus(ID string) (EnrichmentCollectionStatus, error)
}

// EnrichTransaction will generate enriched transaction
func (c *internalClient) EnrichTransaction(enrichmentReq *EnrichmentRequest) (*EnrichmentResponse, error) {
	requestBody, err := json.Marshal(enrichmentReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.xyo.financial/v1/ai/finance/enrichment/transaction",
		bytes.NewReader(requestBody),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))

	resp, err := c.httpClient.request(req)
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

// EnrichTransactionCollection will produce an ID and a download link for bulk transaction enrichment
func (c *internalClient) EnrichTransactionCollection(enrichmentReq []*EnrichmentRequest) (*EnrichTransactionCollectionResponse, error) {
	requestBody, err := json.Marshal(enrichmentReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.xyo.financial/v1/ai/finance/enrichment/transactions",
		bytes.NewReader(requestBody),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := c.httpClient.request(req)
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

// EnrichTransactionCollectionStatus will return the status of bulk transactions enrichment request with a given ID
// Check EnrichmentCollectionStatus for a possible Status Value
func (c *internalClient) EnrichTransactionCollectionStatus(ID string) (EnrichmentCollectionStatus, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("https://api.xyo.financial/v1/ai/finance/enrichment/transactions/status/%s", ID),
		nil,
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := c.httpClient.request(req)
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
