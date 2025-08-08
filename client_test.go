package xyo_test

import (
	"testing"

	"github.com/syniol/xyo-sdk-go"
)

func TestNewClient(t *testing.T) {
	_, _ = xyo.
		NewClient(&xyo.ClientConfig{APIKey: "sss"}).
		EnrichmentStatus("Dummy Transaction")
}
