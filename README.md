# XYO Financial SDK Go (Golang)
![workflow](https://github.com/syniol/xyo-sdk-go/actions/workflows/makefile.yml/badge.svg)

<p align="center">
    <a href="https://xyo.financial" target="blank"><img alt="Go (Golang) Gopher Mascot" width="50%" src="https://github.com/syniol/xyo-sdk-go/blob/main/docs/mascot.png?raw=true" /></a>
</p>

This is an official SDK of XYO Financial for Go (Golang) Programming Language. The minimum requirement is version: `1.18`.


## Quickstart Guide
Client is an entry point to use the SDK. You need a valid API Key obtainable from https://xyo.financial/dashboard

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
		APIKey: "YourAPIKeyFromXYO.FinancialDashboard",
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
