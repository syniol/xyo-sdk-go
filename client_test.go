package xyo_test

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"testing"

	"github.com/syniol/xyo-sdk-go"
)

var server http.Server

func init() {
	// Spinning up a mock TLS Server for test.
	// /etc/hosts is populated in Docker with 127.0.0.1 api.xyo.financial
	server := http.Server{
		Addr: ":80",
	}

	http.HandleFunc("/v1/ai/finance/enrichment", func(wr http.ResponseWriter, rq *http.Request) {
		response, _ := json.MarshalIndent(map[string]bool{"healthy": true}, "", "\t")
		_, _ = wr.Write(response)
	})

	http.HandleFunc("/v1/ai/finance/enrichments", func(wr http.ResponseWriter, rq *http.Request) {
		response, _ := json.MarshalIndent(map[string]bool{"healthy": true}, "", "\t")
		_, _ = wr.Write(response)
	})

	http.HandleFunc("/v1/ai/finance/enrichments/status", func(wr http.ResponseWriter, rq *http.Request) {
		response, _ := json.MarshalIndent(map[string]bool{"healthy": true}, "", "\t")
		_, _ = wr.Write(response)
	})

	go func() {
		println("🚀 RESTful TCP Server is running at: 127.0.0.1:80")
		log.Fatal(server.ListenAndServe())
	}()
}

func TestNewClient(t *testing.T) {
	_ = xyo.
		NewClient(&xyo.ClientConfig{
			APIKey: "YourAPIKeyFromWebDashboard",
		})

	t.Run("health", func(t *testing.T) {
		t.Log("API Test call to Health Status API")
	})

	t.Run("Enrichment", func(t *testing.T) {
		t.Log("API Test call to Enrichment API")
		// todo: API Test call to dummy API
	})

	t.Run("Enrichments", func(t *testing.T) {
		t.Log("API Test call to Bulk Enrichment Collection API")
		// todo: API Test call to dummy API
	})

	t.Run("Enrichments Status", func(t *testing.T) {
		t.Log("API Test call to Bulk Enrichment Collection Status API")
		// todo: API Test call to dummy API
	})

	// Shutting down test server
	t.Cleanup(func() {
		if err := server.Shutdown(context.TODO()); err != nil {
			t.Fatal(err)
		}
	})

}
