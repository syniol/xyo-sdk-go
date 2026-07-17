package main

import (
	"fmt"

	"github.com/syniol/xyo-sdk-go"
)

func main() {
	// Verify the SDK can be imported and the client can be instantiated.
	// For a full usage example see the README.
	_, err := xyo.NewClient(&xyo.ClientConfig{
		APIKey: "example-api-key",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully imported and instantiated the XYO Client")
}
