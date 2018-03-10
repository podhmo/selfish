package selfish

import (
	"github.com/google/go-github/github"
	"github.com/podhmo/commithistory"
)

type Commit = commithistory.Commit

// NewCommit creates and initializes a new Commit object.
func NewCommit(g *github.Gist, alias string, action string) *Commit {
	return &Commit{
		ID:        *g.ID,
		CreatedAt: *g.CreatedAt,
		Alias:     alias,
		Action:    action,
	}
}
