package internal

import (
	"fmt"
	"net/http"
)

func MandatoryAPIHeaders(req *http.Request, apiKey string) {
	if req.Method == http.MethodPost || req.Method == http.MethodPut || req.Method == http.MethodPatch {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("User-Agent", "xyo-sdk-go/1.0.0")
}
