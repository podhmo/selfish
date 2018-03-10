package selfish

import (
	"context"

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

// Create :
func (c *Client) Create(ctx context.Context, filenames []string) (*github.Gist, error) {
	gist, err := NewGist(filenames)
	if err != nil {
		return nil, err
	}
	g, _, err := c.Gists.Create(ctx, gist)
	return g, err
}

// Update :
func (c *Client) Update(ctx context.Context, latestCommit *Commit, filenames []string) (*github.Gist, error) {
	gist, err := NewGist(filenames)
	if err != nil {
		return nil, err
	}
	g, _, err := c.Gists.Edit(ctx, latestCommit.ID, gist)
	return g, err
}
