# XYO Financial SDK Go (Golang)
![workflow](https://github.com/syniol/xyo-sdk-go/actions/workflows/makefile.yml/badge.svg)


## Quick Start
You will need to have your API Key to create a client.

```go
package main

import (
	"encoding/json"
	"log"
	"fmt"

	"github.com/syniol/xyo-sdk-go"
)

func main() {
	client := xyo.NewClient(&xyo.ClientConfig{
		APIKey: "ghjdasd7321312ghjhgdsahjdf/dasdasuit34324e3274gdsa",
	})

	resp, err := client.EnrichTransaction(&xyo.EnrichmentRequest{
		Content:     "COSTA PICKUP",
		CountryCode: "GB",
	})
	log.Fatal(err)

	jsonResp, err := json.Marshal(resp)
	log.Fatal(err)
	fmt.Println(jsonResp)
}
```


#### Credits
Copyright &copy; 2025 Syniol Limited. All rights reserved.
