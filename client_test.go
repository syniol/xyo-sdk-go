package xyo

import (
	"net/http"
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Run("healthy", func(t *testing.T) {
		client := &internalClient{
			httpClient: &httpClient{
				Request: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
					}, nil
				}},
		}

		if client.Health() != nil {
			t.Errorf("was not expecting an error but got: %s", client.Health().Error())
		}
	})

	t.Run("unhealthy", func(t *testing.T) {
		client := &internalClient{
			httpClient: &httpClient{
				Request: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusInternalServerError,
					}, nil
				}},
		}

		err := client.Health()
		if err == nil {
			t.Error("was expecting an error but got nothing back")
		}

		if !strings.Contains(err.Error(), "failed with status code 500") {
			t.Error("was expecting failed with status code 500 in an error message")
		}
	})
}
