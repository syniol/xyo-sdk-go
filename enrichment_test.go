package xyo

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
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

		if responsePayload != nil {
			w.Header().Set("Content-Type", "application/json")
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
		ts, client := newTestServerAndClient(t, http.MethodGet, "/v1/ai/finance/enrichment/status/asdsd", http.StatusBadRequest, nil)
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
		ts, client := newTestServerAndClient(t, http.MethodGet, "/v1/ai/finance/enrichment/status/72c037df-d0d3-43ee-9470-323ff35a2e50", http.StatusOK, payload)
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

func TestDownloadEnrichmentCollection(t *testing.T) {
	t.Run("non-200 status returns error", func(t *testing.T) {
		ts, client := newTestServerAndClient(t, http.MethodGet, "/downloads/123.tar.gz", http.StatusBadRequest, nil)
		defer ts.Close()

		_, err := client.DownloadEnrichmentCollection(context.Background(), ts.URL+"/downloads/123.tar.gz")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("200 OK streams and decodes tarball", func(t *testing.T) {
		// Create an in-memory .tar.gz containing one JSON file
		var buf bytes.Buffer
		gzw := gzip.NewWriter(&buf)
		tw := tar.NewWriter(gzw)

		payload := map[string]interface{}{
			"merchant":    "Syniol Limited",
			"description": "Bulk Test",
		}
		b, _ := json.Marshal(payload)

		_ = tw.WriteHeader(&tar.Header{
			Name:     "transaction_0.json",
			Mode:     0600,
			Size:     int64(len(b)),
			Typeflag: tar.TypeReg,
		})
		_, _ = tw.Write(b)
		_ = tw.Close()
		_ = gzw.Close()

		ts, client := newTestServerAndClient(t, http.MethodGet, "/downloads/123.tar.gz", http.StatusOK, nil)
		defer ts.Close()

		// Override the test server to return our custom tarball bytes instead of JSON
		ts.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(buf.Bytes())
		})

		results, err := client.DownloadEnrichmentCollection(context.Background(), ts.URL+"/downloads/123.tar.gz")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		if results[0].Merchant != "Syniol Limited" {
			t.Errorf("expected merchant %q, got %q", "Syniol Limited", results[0].Merchant)
		}
	})
}

func TestEnrichTransactions(t *testing.T) {
	requests := []*EnrichmentRequest{
		{Content: "Costa PICKUP", CountryCode: "GB"},
		{Content: "STRBUKS GREENWICH", CountryCode: "GB"},
	}

	t.Run("200 OK decodes response", func(t *testing.T) {
		payload := map[string]interface{}{
			"id":   "batch-987",
			"link": "https://api.xyo.financial/ai/transactions/download/batch-987.tar.gz",
		}
		ts, client := newTestServerAndClient(t, http.MethodPost, "/v1/ai/finance/enrichment/transactions", http.StatusOK, payload)
		defer ts.Close()

		resp, err := client.EnrichTransactions(context.Background(), requests)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.ID != "batch-987" {
			t.Errorf("expected ID %q, got %q", "batch-987", resp.ID)
		}
		if resp.Link != "https://api.xyo.financial/ai/transactions/download/batch-987.tar.gz" {
			t.Errorf("expected Link %q, got %q", "https://api.xyo.financial/ai/transactions/download/batch-987.tar.gz", resp.Link)
		}
	})
}

func TestGetEnrichmentStatus(t *testing.T) {
	t.Run("200 OK returns correct status", func(t *testing.T) {
		payload := map[string]interface{}{
			"status": EnrichmentCollectionStatusReady,
		}
		ts, client := newTestServerAndClient(t, http.MethodGet, "/v1/ai/finance/enrichment/status/batch-987", http.StatusOK, payload)
		defer ts.Close()

		actual, err := client.GetEnrichmentStatus(context.Background(), "batch-987")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if actual != EnrichmentCollectionStatusReady {
			t.Errorf("expected status %q, got %q", EnrichmentCollectionStatusReady, actual)
		}
	})
}

func TestEnrichTransaction_NilRequest(t *testing.T) {
	_, client := newTestServerAndClient(t, http.MethodPost, "/v1/ai/finance/enrichment/transaction", http.StatusOK, nil)
	_, err := client.EnrichTransaction(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil request, got nil")
	}
}

func TestEnrichTransactions_NilRequestInSlice(t *testing.T) {
	_, client := newTestServerAndClient(t, http.MethodPost, "/v1/ai/finance/enrichment/transactions", http.StatusOK, nil)
	reqs := []*EnrichmentRequest{
		{Content: "Valid", CountryCode: "GB"},
		nil,
	}
	_, err := client.EnrichTransactions(context.Background(), reqs)
	if err == nil {
		t.Fatal("expected error for slice containing nil request, got nil")
	}
}

func TestDownloadEnrichmentCollection_RFC7807Error(t *testing.T) {
	errPayload := map[string]interface{}{
		"errors": []map[string]interface{}{
			{
				"type":   "https://api.xyo.financial/errors/download-expired",
				"status": http.StatusGone,
				"title":  "Download Link Expired",
				"detail": "The requested batch download link has expired.",
			},
		},
	}
	ts, client := newTestServerAndClient(t, http.MethodGet, "/downloads/expired.tar.gz", http.StatusGone, errPayload)
	defer ts.Close()

	_, err := client.DownloadEnrichmentCollection(context.Background(), ts.URL+"/downloads/expired.tar.gz")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *ErrorResponse
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected error to unwrap to *ErrorResponse, got %v", err)
	}
	if apiErr.HTTPStatusCode != http.StatusGone {
		t.Errorf("expected HTTP status %d, got %d", http.StatusGone, apiErr.HTTPStatusCode)
	}
	if len(apiErr.Errors) != 1 || apiErr.Errors[0].Title != "Download Link Expired" {
		t.Errorf("unexpected error payload: %+v", apiErr.Errors)
	}
}

func TestDownloadEnrichmentCollection_ExternalHostNoAuthLeak(t *testing.T) {
	externalServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if auth := r.Header.Get("Authorization"); auth != "" {
			t.Errorf("expected NO Authorization header for external host, got %q", auth)
		}
		w.WriteHeader(http.StatusOK)
		// Return valid empty tarball
		var buf bytes.Buffer
		gzw := gzip.NewWriter(&buf)
		tw := tar.NewWriter(gzw)
		_ = tw.Close()
		_ = gzw.Close()
		_, _ = w.Write(buf.Bytes())
	}))
	defer externalServer.Close()

	// Client is configured with a different base URL
	client, err := NewClient(&ClientConfig{
		APIKey:  "secret-token",
		BaseURL: "https://api.xyo.financial",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	results, err := client.DownloadEnrichmentCollection(context.Background(), externalServer.URL+"/presigned.tar.gz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestNewClient_ConfigAlias(t *testing.T) {
	c, err := NewClient(&Config{
		APIKey: "test-key",
	})
	if err != nil {
		t.Fatalf("failed to create client with Config alias: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestGetEnrichmentStatus_EmptyID(t *testing.T) {
	_, client := newTestServerAndClient(t, http.MethodGet, "/v1/ai/finance/enrichment/status/test", http.StatusOK, nil)
	_, err := client.GetEnrichmentStatus(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty id, got nil")
	}
}

func TestDownloadEnrichmentCollection_EmptyURL(t *testing.T) {
	_, client := newTestServerAndClient(t, http.MethodGet, "/downloads/123.tar.gz", http.StatusOK, nil)
	_, err := client.DownloadEnrichmentCollection(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty download URL, got nil")
	}
}

func TestDownloadEnrichmentCollection_ExceedsMaxEntrySize(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		gzw := gzip.NewWriter(&buf)
		tw := tar.NewWriter(gzw)

		// Create header claiming size exceeds DefaultMaxEntryBytes
		_ = tw.WriteHeader(&tar.Header{
			Name:     "bomb.json",
			Mode:     0600,
			Size:     DefaultMaxEntryBytes + 1024,
			Typeflag: tar.TypeReg,
		})
		_ = tw.Close()
		_ = gzw.Close()

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(buf.Bytes())
	}))
	defer ts.Close()

	client, err := NewClient(&ClientConfig{
		APIKey:  "test-api-key",
		BaseURL: ts.URL,
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.DownloadEnrichmentCollection(context.Background(), ts.URL+"/bomb.tar.gz")
	if err == nil {
		t.Fatal("expected error when entry size exceeds limit, got nil")
	}
}

func TestVersionConstants(t *testing.T) {
	if Version == "" {
		t.Error("expected non-empty Version")
	}
	if DefaultUserAgent != "xyo-sdk-go/"+Version {
		t.Errorf("expected DefaultUserAgent 'xyo-sdk-go/%s', got %q", Version, DefaultUserAgent)
	}
}
