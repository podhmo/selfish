package selfish

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Client object
type Client struct {
	Config *Config
	*github.Client
}

// CreateClient is factory of github client
func CreateClient() (*Client, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	gclient := createGithubClient(config.AccessToken)
	client := Client{Client: gclient, Config: config}
	return &client, nil
}

func createGithubClient(token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	return github.NewClient(tc)
}
