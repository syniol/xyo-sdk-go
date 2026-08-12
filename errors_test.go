package xyo

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorResponse_Error(t *testing.T) {
	t.Run("empty errors slice", func(t *testing.T) {
		err := &ErrorResponse{
			HTTPStatusCode: http.StatusBadGateway,
		}
		if err.Error() != "status 502" {
			t.Errorf("expected 'status 502', got %q", err.Error())
		}
	})

	t.Run("single error", func(t *testing.T) {
		err := &ErrorResponse{
			HTTPStatusCode: http.StatusBadRequest,
			Errors: []*APIError{
				{
					Title:  "Invalid Parameter",
					Detail: "Content length exceeds limit",
				},
			},
		}
		expected := "status 400: Invalid Parameter: Content length exceeds limit"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("multiple errors", func(t *testing.T) {
		err := &ErrorResponse{
			HTTPStatusCode: http.StatusUnprocessableEntity,
			Errors: []*APIError{
				{
					Title:  "Missing Field",
					Detail: "CountryCode is required",
				},
				{
					Title:  "Invalid Field",
					Detail: "Content cannot be empty",
				},
			},
		}
		expected := "status 422: Missing Field: CountryCode is required, Invalid Field: Content cannot be empty"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})
}

func TestParseOpenAPIError_NilAndGenericErrors(t *testing.T) {
	t.Run("nil error returns nil", func(t *testing.T) {
		if err := parseOpenAPIError(nil, "op", nil); err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	})

	t.Run("generic error with non-zero http response status", func(t *testing.T) {
		stdErr := errors.New("underlying network failure")
		httpResp := &http.Response{StatusCode: 503}

		parsed := parseOpenAPIError(stdErr, "TestOp", httpResp)
		if parsed == nil {
			t.Fatal("expected error, got nil")
		}
		expected := "xyo: TestOp: status 503: underlying network failure"
		if parsed.Error() != expected {
			t.Errorf("expected %q, got %q", expected, parsed.Error())
		}
	})

	t.Run("generic error without http response", func(t *testing.T) {
		stdErr := fmt.Errorf("context deadline exceeded")

		parsed := parseOpenAPIError(stdErr, "TestOp", nil)
		if parsed == nil {
			t.Fatal("expected error, got nil")
		}
		expected := "xyo: TestOp: context deadline exceeded"
		if parsed.Error() != expected {
			t.Errorf("expected %q, got %q", expected, parsed.Error())
		}
	})
}

func TestParseOpenAPIError_ThroughClientExecution(t *testing.T) {
	t.Run("structured error response model", func(t *testing.T) {
		errPayload := map[string]interface{}{
			"errors": []map[string]interface{}{
				{
					"type":     "https://api.xyo.financial/errors/unauthorized",
					"status":   401,
					"title":    "Unauthorized Access",
					"detail":   "Invalid API key provided",
					"instance": "AuthException",
				},
			},
		}
		ts, client := newTestServerAndClient(t, http.MethodPost, "/v1/ai/finance/enrichment/transaction", http.StatusUnauthorized, errPayload)
		defer ts.Close()

		_, err := client.EnrichTransaction(context.Background(), &EnrichmentRequest{
			Content:     "Coffee",
			CountryCode: "GB",
		})
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		var apiErr *ErrorResponse
		if !errors.As(err, &apiErr) {
			t.Fatalf("expected error to unwrap to *ErrorResponse, got: %v", err)
		}
		if apiErr.HTTPStatusCode != 401 {
			t.Errorf("expected status 401, got %d", apiErr.HTTPStatusCode)
		}
		if len(apiErr.Errors) != 1 || apiErr.Errors[0].Title != "Unauthorized Access" {
			t.Errorf("unexpected error details: %+v", apiErr.Errors)
		}
	})

	t.Run("raw JSON fallback error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain") // Content-Type is text/plain so model decode fails, raw body parsed
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"errors":[{"title":"Raw Title","detail":"Raw Detail","status":400}]}`))
		}))
		defer ts.Close()

		client, err := NewClient(&ClientConfig{
			APIKey:  "test-api-key",
			BaseURL: ts.URL,
		})
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}

		_, err = client.EnrichTransaction(context.Background(), &EnrichmentRequest{
			Content:     "Coffee",
			CountryCode: "GB",
		})
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		var apiErr *ErrorResponse
		if !errors.As(err, &apiErr) {
			t.Fatalf("expected error to unwrap to *ErrorResponse, got: %v", err)
		}
		if len(apiErr.Errors) != 1 || apiErr.Errors[0].Title != "Raw Title" {
			t.Errorf("unexpected error response: %+v", apiErr.Errors)
		}
	})

	t.Run("plain non-JSON HTTP status fallback", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte(`<html><body>502 Bad Gateway</body></html>`))
		}))
		defer ts.Close()

		client, err := NewClient(&ClientConfig{
			APIKey:  "test-api-key",
			BaseURL: ts.URL,
		})
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}

		_, err = client.EnrichTransaction(context.Background(), &EnrichmentRequest{
			Content:     "Coffee",
			CountryCode: "GB",
		})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		expected := "xyo: EnrichTransaction: status 502"
		if err.Error() != expected {
			t.Errorf("expected error %q, got %q", expected, err.Error())
		}
	})
}
