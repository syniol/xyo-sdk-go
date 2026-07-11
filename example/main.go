package main

import (
	"fmt"
	"log"

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
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Merchant:   ", resp.Merchant)
	fmt.Println("Description:", resp.Description)
	fmt.Println("Categories: ", resp.Categories)
	fmt.Println("Logo:       ", resp.Logo)
}
