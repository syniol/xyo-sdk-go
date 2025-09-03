package xyo

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestEnrichTransaction(t *testing.T) {
	t.Run("when status code is not 200 (OK)", func(t *testing.T) {
		client := &client{
			config: &ClientConfig{
				APIKey: "LWMzMGE0NmQ1MmNkNQo2OWNhNWVlMy0MWItYWIyZi1hMTc3ZTFkMDA0NDM=",
				httpClient: &httpClient{
					request: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusBadRequest,
						}, nil
					},
				},
			},
		}

		_, err := client.EnrichTransaction(&EnrichmentRequest{
			Content:     "Some Random Content",
			CountryCode: "GB",
		})
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("when status code is 200 (OK)", func(t *testing.T) {
		mockedEnrichmentResponse := map[string]interface{}{
			"merchant":    "Syniol Limited",
			"description": "Software and Cloud Platform Consultancy",
			"logo":        "base64/png;31233232dsdsdaaersdasjhdsfi",
			"categories":  []string{"Tech"},
		}

		jsonMockedEnrichmentResponse, _ := json.Marshal(mockedEnrichmentResponse)
		stringReadCloser := io.NopCloser(bytes.NewReader(jsonMockedEnrichmentResponse))

		client := &client{
			config: &ClientConfig{
				APIKey: "LWMzMGE0NmQ1MmNkNQo2OWNhNWVlMy0MWItYWIyZi1hMTc3ZTFkMDA0NDM=",
				httpClient: &httpClient{
					request: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							Body:       stringReadCloser,
							StatusCode: http.StatusOK,
						}, nil
					},
				},
			},
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
	t.Run("when status code is not 200 (OK)", func(t *testing.T) {
		client := &client{
			config: &ClientConfig{
				APIKey: "LWMzMGE0NmQ1MmNkNQo2OWNhNWVlMy0MWItYWIyZi1hMTc3ZTFkMDA0NDM=",
				httpClient: &httpClient{
					request: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusBadRequest,
						}, nil
					},
				},
			},
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

	t.Run("when status code is 200 (OK)", func(t *testing.T) {
		payloadMap := map[string]interface{}{
			"id":   "72c037df-d0d3-43ee-9470-323ff35a2e50",
			"link": "https://api.xyo.financial/ai/transactions/download/72c037df-d0d3-43ee-9470-323ff35a2e50.tar.gz",
		}
		serialisedPayload, _ := json.Marshal(payloadMap)

		client := &client{
			config: &ClientConfig{
				APIKey: "LWMzMGE0NmQ1MmNkNQo2OWNhNWVlMy0MWItYWIyZi1hMTc3ZTFkMDA0NDM=",
				httpClient: &httpClient{
					request: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							Body:       io.NopCloser(bytes.NewReader(serialisedPayload)),
							StatusCode: http.StatusOK,
						}, nil
					},
				},
			},
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
	t.Run("when status code is not 200 (OK)", func(t *testing.T) {
		client := &client{
			config: &ClientConfig{
				APIKey: "LWMzMGE0NmQ1MmNkNQo2OWNhNWVlMy0MWItYWIyZi1hMTc3ZTFkMDA0NDM=",
				httpClient: &httpClient{
					request: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusBadRequest,
						}, nil
					},
				},
			},
		}

		_, err := client.EnrichTransactionCollectionStatus("asdsd")
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("when status code is 200 (OK)", func(t *testing.T) {
		payloadMap := map[string]interface{}{
			"status": EnrichmentCollectionStatusReady,
		}
		serialisedPayload, _ := json.Marshal(payloadMap)

		client := &client{
			config: &ClientConfig{
				APIKey: "LWMzMGE0NmQ1MmNkNQo2OWNhNWVlMy0MWItYWIyZi1hMTc3ZTFkMDA0NDM=",
				httpClient: &httpClient{
					request: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							Body:       io.NopCloser(bytes.NewReader(serialisedPayload)),
							StatusCode: http.StatusOK,
						}, nil
					},
				},
			},
		}

		actual, err := client.EnrichTransactionCollectionStatus("72c037df-d0d3-43ee-9470-323ff35a2e50")
		if err != nil {
			t.Error("error", err)
		}

		if actual != EnrichmentCollectionStatusReady {
			t.Errorf("expected a status: '%s'", EnrichmentCollectionStatusReady)
		}
	})
}
