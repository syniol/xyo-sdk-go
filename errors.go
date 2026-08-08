package xyo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/xyo-financial/sdk-go/v2/openapi"
)

// APIError represents an RFC 7807-inspired error returned by the XYO API.
type APIError struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status,omitempty"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

// ErrorResponse represents the payload of an API error response.
// The XYO API wraps errors in an "errors" array.
type ErrorResponse struct {
	// HTTPStatusCode contains the HTTP response status code that triggered the error.
	HTTPStatusCode int `json:"-"`

	// Errors contains the list of API exceptions returned by the server.
	Errors []*APIError `json:"errors"`
}

// Error implements the standard error interface.
func (e *ErrorResponse) Error() string {
	if len(e.Errors) == 0 {
		return fmt.Sprintf("status %d", e.HTTPStatusCode)
	}

	if len(e.Errors) == 1 {
		err := e.Errors[0]
		return fmt.Sprintf("status %d: %s: %s", e.HTTPStatusCode, err.Title, err.Detail)
	}

	msgs := make([]string, 0, len(e.Errors))
	for _, err := range e.Errors {
		msgs = append(msgs, fmt.Sprintf("%s: %s", err.Title, err.Detail))
	}
	return fmt.Sprintf("status %d: %s", e.HTTPStatusCode, strings.Join(msgs, ", "))
}

// parseOpenAPIError converts a generated-client error into an SDK-level *ErrorResponse,
// preserving HTTP status codes and structured error fields from the API response body.
func parseOpenAPIError(err error, op string, httpResp *http.Response) error {
	if err == nil {
		return nil
	}

	statusCode := 0
	if httpResp != nil {
		statusCode = httpResp.StatusCode
	}

	var openapiErr *openapi.GenericOpenAPIError
	if errors.As(err, &openapiErr) {
		// Try the already-decoded model first.
		if model := openapiErr.Model(); model != nil {
			var errResp *openapi.ErrorResponse
			switch m := model.(type) {
			case openapi.ErrorResponse:
				errResp = &m
			case *openapi.ErrorResponse:
				errResp = m
			}
			if errResp != nil && len(errResp.Errors) > 0 {
				return fmt.Errorf("xyo: %s: %w", op, mapErrorResponse(errResp, statusCode))
			}
		}

		// Fall back to raw body parsing.
		if len(openapiErr.Body()) > 0 {
			var raw openapi.ErrorResponse
			if jsonErr := json.Unmarshal(openapiErr.Body(), &raw); jsonErr == nil && len(raw.Errors) > 0 {
				return fmt.Errorf("xyo: %s: %w", op, mapErrorResponse(&raw, statusCode))
			}
		}

		if statusCode != 0 {
			return fmt.Errorf("xyo: %s: status %d", op, statusCode)
		}
	}

	if httpResp != nil && httpResp.StatusCode != 0 {
		return fmt.Errorf("xyo: %s: status %d: %w", op, httpResp.StatusCode, err)
	}

	return fmt.Errorf("xyo: %s: %w", op, err)
}

func mapErrorResponse(src *openapi.ErrorResponse, statusCode int) *ErrorResponse {
	out := &ErrorResponse{HTTPStatusCode: statusCode}
	for _, e := range src.Errors {
		out.Errors = append(out.Errors, &APIError{
			Type:     e.GetType(),
			Title:    e.GetTitle(),
			Status:   int(e.GetStatus()),
			Detail:   e.GetDetail(),
			Instance: e.GetInstance(),
		})
	}
	if out.HTTPStatusCode == 0 && len(out.Errors) > 0 {
		out.HTTPStatusCode = out.Errors[0].Status
	}
	return out
}
