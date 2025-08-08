package xyo_test

import (
	"net/http"
	"testing"

	"github.com/syniol/xyo-sdk-go"
)

var server http.Server

func TestNewClient(t *testing.T) {
	client := xyo.
		NewClient(&xyo.ClientConfig{
			APIKey: "YourAPIKeyFromWebDashboard",
		})

	t.Run("health", func(t *testing.T) {
		err := client.Health()
		if err != nil {
			t.Error(err)
		}
	})
}
