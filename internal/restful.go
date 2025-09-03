package internal

import (
	"fmt"
	"net/http"
)

func MandatoryAPIHeaders(req *http.Request, apiKey string) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
}
