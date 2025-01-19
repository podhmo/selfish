package internal

import (
	"context"
	"log"
	"time"

	"github.com/google/go-github/v68/github"
	"github.com/pkg/errors"
)

type Client struct {
	Github *github.Client
}

type CreateResult struct {
	GistID    string
	CreatedAt time.Time
	HTMLURL   string
	raw       *github.Gist
}

func (c *Client) Create(ctx context.Context, filenames []string) (*CreateResult, error) {
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
	var createdAt time.Time
	if g.CreatedAt != nil {
		createdAt = g.CreatedAt.Time
	}
	var gistID string
	if g.ID != nil {
		gistID = *g.ID
	}
	return &CreateResult{
		GistID:    gistID,
		CreatedAt: createdAt,
		HTMLURL:   htmlURL,
		raw:       g,
	}, nil
}

type UpdateResult struct {
	GistID    string
	CreatedAt time.Time
	HTMLURL   string
	raw       *github.Gist
}

func (c *Client) Update(ctx context.Context, gistID string, filenames []string) (*UpdateResult, error) {
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
	var createdAt time.Time
	if g.CreatedAt != nil {
		createdAt = g.CreatedAt.Time
	}
	if g.ID != nil {
		if *g.ID != gistID {
			log.Printf("WARN: gistId is mismatch %q != %q", gistID, *g.ID)
		}
	}
	return &UpdateResult{
		GistID:    gistID,
		CreatedAt: createdAt,
		HTMLURL:   htmlURL,
		raw:       g,
	}, nil
}

func (c *Client) Delete(ctx context.Context, gistID string) error {
	if _, err := c.Github.Gists.Delete(ctx, gistID); err != nil {
		return err
	}
	return nil
}
