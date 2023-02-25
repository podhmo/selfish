package selfish

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/v50/github"
	"github.com/pkg/errors"
	"github.com/podhmo/selfish/pkg/commithistory"
	"github.com/toqueteos/webbrowser"
	"golang.org/x/oauth2"
)

const (
	defaultAlias               = "head"
	defaultHistFile            = "selfish.history"
	defaultRepositoryDirectory = "gists"
)

// App :
type App struct {
	CommitHistory *commithistory.API
	Client        Client
	*Config       // TODO(podhmo): stop embedded
}

var ErrAccessTokenNotfound = fmt.Errorf("access token is not found")

func NewApp(c *Config) (*App, error) {
	// if c.Debug {
	// 	fmt.Fprintln(os.Stderr, "config loaded")
	// 	fprintJSON(os.Stderr, c)
	// }

	var chOptions []func(*commithistory.API)
	if c.ClientType == ClientTypeFake {
		chOptions = append(chOptions, commithistory.WithDryrun())
	}
	ch := commithistory.New("selfish", chOptions...)
	{
		if err := ch.Load("config.json", &c.Profile); err != nil {
			return nil, errors.Wrap(err, "load config")
		}
		if c.Profile.AccessToken == "" {
			return nil, ErrAccessTokenNotfound
		}

		if c.Profile.DefaultAlias == "" {
			c.Profile.DefaultAlias = defaultAlias
		}
		if c.Profile.HistFile == "" {
			c.Profile.HistFile = defaultHistFile
		}
		if c.Profile.RepositoryDirectory == "" {
			c.Profile.RepositoryDirectory = defaultRepositoryDirectory
		}
	}

	var gh Client
	switch c.ClientType {
	case ClientTypeFake:
		gh = &fakeClient{W: os.Stderr}
	case ClientTypeGithub:
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: c.Profile.AccessToken},
		)
		tc := oauth2.NewClient(oauth2.NoContext, ts)
		gh = &client{Github: github.NewClient(tc)}
	default:
		return nil, fmt.Errorf("unexpected client type %q", c.ClientType)
	}

	return &App{
		CommitHistory: ch,
		Client:        gh,
		Config:        c,
	}, nil
}

// Delete :
func (app *App) Delete(ctx context.Context, latestCommit *Commit) error {
	alias := app.Alias
	if latestCommit == nil {
		return errors.Errorf("alias=%q is not found", alias)
	}

	gistID := latestCommit.ID
	if err := app.Client.Delete(ctx, gistID); err != nil {
		return errors.Wrapf(err, "gist api delete")
	}

	c := Commit{ID: gistID, Alias: alias, CreatedAt: time.Now(), Action: "delete"}
	if err := app.CommitHistory.SaveCommit(app.Config.Profile.HistFile, &c); err != nil {
		return errors.Wrap(err, "save commit")
	}
	fmt.Fprintf(os.Stderr, "deleted. (id=%q)\n", gistID)
	return nil
}

// Create :
func (app *App) Create(ctx context.Context, latestCommit *Commit, filenames []string) (*Commit, error) {
	action := "create"
	alias := app.Alias
	g, err := app.Client.Create(ctx, filenames)
	if err != nil {
		return nil, errors.Wrapf(err, "gist api %s", action)
	}

	commit := &Commit{
		ID:        g.GistID,
		CreatedAt: g.CreatedAt,
		Alias:     app.Config.ResolveAlias(alias),
		Action:    action,
	}
	if err := app.CommitHistory.SaveCommit(app.Config.Profile.HistFile, commit); err != nil {
		return commit, errors.Wrap(err, "save commit")
	}

	fmt.Fprintf(os.Stderr, "%s success. (id=%q)\n", action, commit.ID)
	if !app.IsSilent && g.HTMLURL != "" {
		fmt.Fprintf(os.Stderr, "opening.. %q\n", g.HTMLURL)
		webbrowser.Open(g.HTMLURL)
	}
	// PrintJSON(g)
	return commit, nil
}

// Update :
func (app *App) Update(ctx context.Context, latestCommit *Commit, filenames []string) (*Commit, error) {
	action := "update"
	alias := app.Alias

	g, err := app.Client.Update(ctx, latestCommit.ID, filenames)
	if err != nil {
		return nil, errors.Wrapf(err, "gist api %s", action)
	}

	commit := &Commit{
		ID:        g.GistID,
		CreatedAt: g.CreatedAt,
		Alias:     app.Config.ResolveAlias(alias),
		Action:    action,
	}
	if err := app.CommitHistory.SaveCommit(app.Config.Profile.HistFile, commit); err != nil {
		return commit, err
	}

	fmt.Fprintf(os.Stderr, "%s success. (id=%q)\n", action, commit.ID)
	if !app.IsSilent {
		fmt.Fprintf(os.Stderr, "opening.. %q\n", g.HTMLURL)
		webbrowser.Open(g.HTMLURL)
	}

	if app.Config.Debug {
		fprintJSON(os.Stderr, g)
	}
	return commit, nil
}
