package xyo

import (
	"fmt"
	"strings"
)

// APIError represents an RFC 7807 inspired error returned by the XYO API.
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

	var msgs []string
	for _, err := range e.Errors {
		msgs = append(msgs, fmt.Sprintf("%s: %s", err.Title, err.Detail))
	}
	return fmt.Sprintf("status %d: %s", e.HTTPStatusCode, strings.Join(msgs, ", "))
}
