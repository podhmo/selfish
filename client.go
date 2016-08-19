package selfish

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// CreateClient is factory of github client
func CreateClient(token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	return github.NewClient(tc)
}
