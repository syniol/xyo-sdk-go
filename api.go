package xyo

import (
	"encoding/json"
	"net/http"
)

func EnrichTransaction(enrichmentReq EnrichmentRequest) (enrichment interface{}, err error) {
	req, err := http.NewRequest(
		"POST",
		"https://xyo.financial/v1/ai/transaction",
		nil,
	)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&enrichment)

	return
}
