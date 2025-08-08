package xyo_test

import (
	"testing"

	"github.com/syniol/xyo-sdk-go"
)

func TestNewClient(t *testing.T) {
	xyo.
		NewClient(&xyo.ClientConfig{APIKey: "sss"}).
		EnrichmentStatus("Dummy Transaction")
}
