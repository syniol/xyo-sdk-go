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

type EnrichmentCollectionStatus string

const (
	EnrichmentCollectionSuccess EnrichmentCollectionStatus = "READY"
	EnrichmentCollectionFailure EnrichmentCollectionStatus = "FAILED"
	EnrichmentCollectionPending EnrichmentCollectionStatus = "PENDING"
)

type enrichmentRequester interface {
	EnrichTransaction(enrichmentReq EnrichmentRequest) (interface{}, error)
	EnrichTransactionCollection(enrichmentReq []EnrichmentRequest) ([]interface{}, error)
	EnrichmentStatus(ID string) string
}

func (c *Client) EnrichTransaction(enrichmentReq EnrichmentRequest) (enrichment interface{}, err error) {
	requestBody, err := json.Marshal(enrichmentReq)
	if err != nil {
		return
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://xyo.financial/v1/ai/transaction",
		bytes.NewReader(requestBody),
	)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&enrichment)

	return
}

func (c *Client) EnrichTransactionCollection(enrichmentReq []EnrichmentRequest) (enrichment []interface{}, err error) {
	requestBody, err := json.Marshal(enrichmentReq)
	if err != nil {
		return
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://xyo.financial/v1/ai/transactions",
		bytes.NewReader(requestBody),
	)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&enrichment)

	return
}

func (c *Client) EnrichmentStatus(ID string) (EnrichmentCollectionStatus, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("https://xyo.financial/v1/ai/enrichments/status/%s", ID),
		nil,
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	var response struct {
		Status EnrichmentCollectionStatus `json:"status"`
	}
	err = json.NewDecoder(resp.Body).Decode(&response)

	return response.Status, err
}
