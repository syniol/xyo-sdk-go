package xyo

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

// newTestClient constructs a client with a stubbed HTTP transport.
// statusCode is returned for every request. body may be nil for error-path tests.
func newTestClient(_ *testing.T, statusCode int, body io.ReadCloser) *client {
	if body == nil {
		body = http.NoBody
	}
	return &client{
		apiKey:     "test-api-key",
		apiBaseURL: defaultAPIBaseURL,
		http: &httpClient{
			request: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: statusCode,
					Body:       body,
				}, nil
			},
		},
	}
}

func TestEnrichTransaction(t *testing.T) {
	t.Run("non-200 status returns error", func(t *testing.T) {
		c := newTestClient(t, http.StatusBadRequest, nil)

		_, err := c.EnrichTransaction(&EnrichmentRequest{
			Content:     "Some Random Content",
			CountryCode: "GB",
		})
		if err == nil {
			t.Fatal("expected error, got nil")
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
		body := marshalBody(t, payload)

		c := newTestClient(t, http.StatusOK, body)
		resp, err := c.EnrichTransaction(&EnrichmentRequest{
			Content:     "Some Random Content",
			CountryCode: "GB",
		})
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
		c := newTestClient(t, http.StatusBadRequest, nil)

		_, err := c.EnrichTransactionCollection(requests)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("200 OK decodes response", func(t *testing.T) {
		payload := map[string]interface{}{
			"id":   "72c037df-d0d3-43ee-9470-323ff35a2e50",
			"link": "https://api.xyo.financial/ai/transactions/download/72c037df-d0d3-43ee-9470-323ff35a2e50.tar.gz",
		}
		body := marshalBody(t, payload)

		c := newTestClient(t, http.StatusOK, body)
		resp, err := c.EnrichTransactionCollection(requests)
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
		c := newTestClient(t, http.StatusBadRequest, nil)

		_, err := c.EnrichTransactionCollectionStatus("asdsd")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("200 OK returns correct status", func(t *testing.T) {
		payload := map[string]interface{}{
			"status": EnrichmentCollectionStatusReady,
		}
		body := marshalBody(t, payload)

		c := newTestClient(t, http.StatusOK, body)
		actual, err := c.EnrichTransactionCollectionStatus("72c037df-d0d3-43ee-9470-323ff35a2e50")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if actual != EnrichmentCollectionStatusReady {
			t.Errorf("expected status %q, got %q", EnrichmentCollectionStatusReady, actual)
		}
	})
}

// marshalBody is a test helper that serialises v to JSON and returns a ReadCloser.
func marshalBody(t *testing.T, v interface{}) io.ReadCloser {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshalBody: %v", err)
	}
	return io.NopCloser(bytes.NewReader(b))
}
