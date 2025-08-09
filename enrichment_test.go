package xyo

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestEnrichTransaction(t *testing.T) {
	t.Run("negative", func(t *testing.T) {
		client := &internalClient{
			httpClient: &httpClient{
				request: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusBadRequest,
					}, nil
				}},
			config: &ClientConfig{APIKey: "xsadsdsadada"},
		}

		_, err := client.EnrichTransaction(&EnrichmentRequest{
			Content:     "Some Random Content",
			CountryCode: "GB",
		})
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("positive", func(t *testing.T) {
		sss := map[string]interface{}{
			"merchant":    "ssdsdsa",
			"description": "FUCK O",
			"logo":        "cdafdafa",
			"categories":  []string{"ssdsdsa"},
		}

		sxxss, _ := json.Marshal(sss)
		stringReadCloser := io.NopCloser(bytes.NewReader(sxxss))

		client := &internalClient{
			httpClient: &httpClient{
				request: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						Body:       stringReadCloser,
						StatusCode: http.StatusOK,
					}, nil
				}},
			config: &ClientConfig{APIKey: "xsadsdsadada"},
		}

		_, err := client.EnrichTransaction(&EnrichmentRequest{
			Content:     "Some Random Content",
			CountryCode: "GB",
		})
		if err != nil {
			t.Error("error", err)
		}
	})
}

func TestEnrichTransactionCollection(t *testing.T) {
	t.Run("negative", func(t *testing.T) {
		client := &internalClient{
			httpClient: &httpClient{
				request: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusBadRequest,
					}, nil
				}},
			config: &ClientConfig{APIKey: "xsadsdsadada"},
		}

		_, err := client.EnrichTransactionCollection([]*EnrichmentRequest{
			{
				Content:     "Some Random Content",
				CountryCode: "GB",
			},
			{
				Content:     "Some Random Content 2",
				CountryCode: "US",
			},
		})
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("positive", func(t *testing.T) {
		payloadMap := map[string]interface{}{
			"id":   "72c037df-d0d3-43ee-9470-323ff35a2e50",
			"link": "https://api.xyo.financial/ai/transactions/download/72c037df-d0d3-43ee-9470-323ff35a2e50.tar.gz",
		}
		serialisedPayload, _ := json.Marshal(payloadMap)

		client := &internalClient{
			httpClient: &httpClient{
				request: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						Body:       io.NopCloser(bytes.NewReader(serialisedPayload)),
						StatusCode: http.StatusOK,
					}, nil
				}},
			config: &ClientConfig{APIKey: "xsadsdsadada"},
		}

		_, err := client.EnrichTransactionCollection([]*EnrichmentRequest{
			{
				Content:     "Some Random Content",
				CountryCode: "GB",
			},
			{
				Content:     "Some Random Content 2",
				CountryCode: "US",
			},
		})
		if err != nil {
			t.Error("error", err)
		}
	})
}

func TestEnrichTransactionCollectionStatus(t *testing.T) {
	t.Run("negative", func(t *testing.T) {
		client := &internalClient{
			httpClient: &httpClient{
				request: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusBadRequest,
					}, nil
				}},
			config: &ClientConfig{APIKey: "xsadsdsadada"},
		}

		_, err := client.EnrichTransactionCollectionStatus("asdsd")
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("positive", func(t *testing.T) {
		payloadMap := map[string]interface{}{
			"status": EnrichmentCollectionStatusReady,
		}
		serialisedPayload, _ := json.Marshal(payloadMap)

		client := &internalClient{
			httpClient: &httpClient{
				request: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						Body:       io.NopCloser(bytes.NewReader(serialisedPayload)),
						StatusCode: http.StatusOK,
					}, nil
				}},
			config: &ClientConfig{APIKey: "xsadsdsadada"},
		}

		actual, err := client.EnrichTransactionCollectionStatus("72c037df-d0d3-43ee-9470-323ff35a2e50")
		if err != nil {
			t.Error("error", err)
		}

		if actual != EnrichmentCollectionStatusReady {
			t.Error("expected a status: 'READY'")
		}
	})
}
