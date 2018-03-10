package selfish

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Client object
type Client struct {
	*github.Client
}

// NewClient is factory of github client
func NewClient(c *Config) *Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.AccessToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	return &Client{Client: github.NewClient(tc)}
}
