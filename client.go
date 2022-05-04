package selfish

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

type Client interface {
	Delete(ctx context.Context, gistID string) error
	Create(ctx context.Context, filenames []string) (*CreateResult, error)
	Update(ctx context.Context, gistID string, filenames []string) (*UpdateResult, error)
}

type client struct {
	Github *github.Client
}

type CreateResult struct {
	HTMLURL string
	raw     *github.Gist
}

func (c *client) Create(ctx context.Context, filenames []string) (*CreateResult, error) {
	gist, err := NewGist(filenames)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid file name is found. %q", filenames)
	}
	g, _, err := c.Github.Gists.Create(ctx, gist)
	if err != nil {
		return nil, err
	}

	var htmlURL string
	if g.HTMLURL != nil {
		htmlURL = *g.HTMLURL
	}
	return &CreateResult{
		raw:     g,
		HTMLURL: htmlURL,
	}, nil
}

type UpdateResult struct {
	HTMLURL string
	raw     *github.Gist
}

func (c *client) Update(ctx context.Context, gistID string, filenames []string) (*UpdateResult, error) {
	gist, err := NewGist(filenames)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid file name is found. %q", filenames)
	}
	g, _, err := c.Github.Gists.Edit(ctx, gistID, gist)
	if err != nil {
		return nil, err
	}
	var htmlURL string
	if g.HTMLURL != nil {
		htmlURL = *g.HTMLURL
	}
	return &UpdateResult{
		raw:     g,
		HTMLURL: htmlURL,
	}, nil
}

func (c *client) Delete(ctx context.Context, gistID string) error {
	if _, err := c.Github.Gists.Delete(ctx, gistID); err != nil {
		return err
	}
	return nil
}
