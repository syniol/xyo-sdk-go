package main

import "github.com/syniol/xyo-sdk-go"

func main() {
	xyo.NewClient(&xyo.ClientConfig{APIKey: "RandomBase64EncodedStringApiKey"})

	println("Successfully imported and instantiated the XYO Client")
}
