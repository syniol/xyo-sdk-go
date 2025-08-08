package xyo

type EnrichmentRequest struct {
	Content     string `json:"content"`
	CountryCode string `json:"countryCode"`
}

type enrichmentRequester interface {
	EnrichTransaction(enrichmentReq EnrichmentRequest)
	EnrichTransactionCollection(enrichmentReq []EnrichmentRequest)
	EnrichmentStatus(ID string)
}

func (c *Client) EnrichTransaction(enrichmentReq EnrichmentRequest) {}

func (c *Client) EnrichTransactionCollection(enrichmentReq []EnrichmentRequest) {}

func (c *Client) EnrichmentStatus(ID string) {}
