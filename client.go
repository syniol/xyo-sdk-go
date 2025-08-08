package xyo

type ClientConfig struct {
	APIKey string
}

type Client struct {
	config *ClientConfig
}

func NewClient(opt *ClientConfig) *Client {
	return &Client{
		config: opt,
	}
}
