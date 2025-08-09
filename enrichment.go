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
	Merchant    string
	Description string
	Categories  []string
	Logo        string
}

type EnrichTransactionCollectionResponse struct {
	ID   string `json:"id"`
	Link string `json:"link"`
}

type EnrichmentCollectionStatus string

const (
	EnrichmentCollectionSuccess EnrichmentCollectionStatus = "READY"
	EnrichmentCollectionFailure EnrichmentCollectionStatus = "FAILED"
	EnrichmentCollectionPending EnrichmentCollectionStatus = "PENDING"
)

func (c *internalClient) EnrichTransaction(enrichmentReq EnrichmentRequest) (enrichment *EnrichmentResponse, err error) {
	requestBody, err := json.Marshal(enrichmentReq)
	if err != nil {
		return
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.xyo.financial/v1/ai/finance/enrichment/transaction",
		bytes.NewReader(requestBody),
	)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))

	resp, err := c.httpClient.request(req)
	if err != nil {
		return
	}

	err = json.NewDecoder(resp.Body).Decode(enrichment)
	return
}

func (c *internalClient) EnrichTransactionCollection(enrichmentReq []EnrichmentRequest) (enrichment *EnrichTransactionCollectionResponse, err error) {
	requestBody, err := json.Marshal(enrichmentReq)
	if err != nil {
		return
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.xyo.financial/v1/ai/finance/enrichment/transactions",
		bytes.NewReader(requestBody),
	)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := c.httpClient.request(req)
	if err != nil {
		return
	}

	err = json.NewDecoder(resp.Body).Decode(enrichment)
	return
}

func (c *internalClient) EnrichTransactionCollectionStatus(ID string) (status EnrichmentCollectionStatus, err error) {
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

	var response struct {
		Status EnrichmentCollectionStatus `json:"status"`
	}

	return response.Status, json.NewDecoder(resp.Body).Decode(&response)
}
