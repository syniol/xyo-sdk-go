package xyo

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestServerAndClient(t *testing.T, expectedMethod, expectedPath string, statusCode int, responsePayload interface{}) (*httptest.Server, Client) {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != expectedMethod {
			t.Errorf("expected method %q, got %q", expectedMethod, r.Method)
		}
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %q, got %q", expectedPath, r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-api-key" {
			t.Errorf("missing or invalid Authorization header: %q", r.Header.Get("Authorization"))
		}
		if expectedMethod != http.MethodGet && r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("missing or invalid Content-Type header: %q", r.Header.Get("Content-Type"))
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("missing or invalid Accept header: %q", r.Header.Get("Accept"))
		}

		w.WriteHeader(statusCode)
		if responsePayload != nil {
			_ = json.NewEncoder(w).Encode(responsePayload)
		}
	}))

	client, err := NewClient(&ClientConfig{
		APIKey:  "test-api-key",
		BaseURL: ts.URL,
	})
	if err != nil {
		t.Fatalf("failed to create test client: %v", err)
	}

	return ts, client
}

func TestEnrichTransaction(t *testing.T) {
	reqPayload := &EnrichmentRequest{
		Content:     "Some Random Content",
		CountryCode: "GB",
	}

	t.Run("non-200 status returns error", func(t *testing.T) {
		errPayload := map[string]interface{}{
			"errors": []map[string]interface{}{
				{
					"type":     "Invalid API Key",
					"status":   http.StatusForbidden,
					"title":    "Invalid API Key",
					"instance": "InvalidClientAPIKeyException",
					"detail":   "Credits expired or an invalid API Key is given",
				},
			},
		}

		ts, client := newTestServerAndClient(t, http.MethodPost, "/v1/ai/finance/enrichment/transaction", http.StatusForbidden, errPayload)
		defer ts.Close()

		_, err := client.EnrichTransaction(context.Background(), reqPayload)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		var apiErr *ErrorResponse
		if !errors.As(err, &apiErr) {
			t.Fatalf("expected error to unwrap to *ErrorResponse, got: %v", err)
		}
		if apiErr.HTTPStatusCode != http.StatusForbidden {
			t.Errorf("expected HTTP status %d, got %d", http.StatusForbidden, apiErr.HTTPStatusCode)
		}
		if len(apiErr.Errors) != 1 {
			t.Fatalf("expected 1 error, got %d", len(apiErr.Errors))
		}
		if apiErr.Errors[0].Title != "Invalid API Key" {
			t.Errorf("expected title %q, got %q", "Invalid API Key", apiErr.Errors[0].Title)
		}
	})

	t.Run("200 OK decodes response", func(t *testing.T) {
		payload := map[string]interface{}{
			"merchant":    "Syniol Limited",
			"description": "Software and Cloud Platform Consultancy",
			"logo":        "base64/png;31233232dsdsdaaersdasjhdsfi",
			"categories":  []string{"Tech"},
			"location":    "United Kingdom, England",
			"address":     "London, O2",
		}
		ts, client := newTestServerAndClient(t, http.MethodPost, "/v1/ai/finance/enrichment/transaction", http.StatusOK, payload)
		defer ts.Close()

		resp, err := client.EnrichTransaction(context.Background(), reqPayload)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Merchant != "Syniol Limited" {
			t.Errorf("expected merchant %q, got %q", "Syniol Limited", resp.Merchant)
		}
	})
}

func TestEnrichTransactionCollection(t *testing.T) {
	requests := []*EnrichmentRequest{
		{Content: "Some Random Content", CountryCode: "GB"},
		{Content: "Some Random Content 2", CountryCode: "US"},
	}

	t.Run("non-200 status returns error", func(t *testing.T) {
		ts, client := newTestServerAndClient(t, http.MethodPost, "/v1/ai/finance/enrichment/transactions", http.StatusBadRequest, nil)
		defer ts.Close()

		_, err := client.EnrichTransactionCollection(context.Background(), requests)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("200 OK decodes response", func(t *testing.T) {
		payload := map[string]interface{}{
			"id":   "72c037df-d0d3-43ee-9470-323ff35a2e50",
			"link": "https://api.xyo.financial/ai/transactions/download/72c037df-d0d3-43ee-9470-323ff35a2e50.tar.gz",
		}
		ts, client := newTestServerAndClient(t, http.MethodPost, "/v1/ai/finance/enrichment/transactions", http.StatusOK, payload)
		defer ts.Close()

		resp, err := client.EnrichTransactionCollection(context.Background(), requests)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.ID != "72c037df-d0d3-43ee-9470-323ff35a2e50" {
			t.Errorf("expected ID %q, got %q", "72c037df-d0d3-43ee-9470-323ff35a2e50", resp.ID)
		}
	})
}

func TestEnrichTransactionCollectionStatus(t *testing.T) {
	t.Run("non-200 status returns error", func(t *testing.T) {
		ts, client := newTestServerAndClient(t, http.MethodGet, "/v1/ai/finance/enrichment/transactions/status/asdsd", http.StatusBadRequest, nil)
		defer ts.Close()

		_, err := client.EnrichTransactionCollectionStatus(context.Background(), "asdsd")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("200 OK returns correct status", func(t *testing.T) {
		payload := map[string]interface{}{
			"status": EnrichmentCollectionStatusReady,
		}
		ts, client := newTestServerAndClient(t, http.MethodGet, "/v1/ai/finance/enrichment/transactions/status/72c037df-d0d3-43ee-9470-323ff35a2e50", http.StatusOK, payload)
		defer ts.Close()

		actual, err := client.EnrichTransactionCollectionStatus(context.Background(), "72c037df-d0d3-43ee-9470-323ff35a2e50")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if actual != EnrichmentCollectionStatusReady {
			t.Errorf("expected status %q, got %q", EnrichmentCollectionStatusReady, actual)
		}
	})
}
