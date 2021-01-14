package client

import "github.com/javorszky/form3takehome/pkg/config"

type Client struct {
	baseURL string
}

// New returns a configured Client struct.
func New(cfg config.Config) Client {
	return Client{
		baseURL: cfg.AccountsAPIURL,
	}
}

func (c Client) List() {
}
