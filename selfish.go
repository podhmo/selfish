package selfish

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"github.com/podhmo/selfish/pkg/commithistory"
	"github.com/toqueteos/webbrowser"
	"golang.org/x/oauth2"
)

const (
	defaultAlias    = "head"
	defaultHistFile = "selfish.history"
)

// App :
type App struct {
	CommitHistory *commithistory.API
	Client        Client
	*Config       // TODO(podhmo): stop embedded
}

var ErrAccessTokenNotfound = fmt.Errorf("access token is not found")

func NewApp(c *Config) (*App, error) {
	ch := commithistory.New("selfish")
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
	}

	// if c.Debug {
	// 	fmt.Fprintln(os.Stderr, "config loaded")
	// 	fprintJSON(os.Stderr, c)
	// }

	var gh *github.Client
	{
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: c.Profile.AccessToken},
		)
		tc := oauth2.NewClient(oauth2.NoContext, ts)
		gh = github.NewClient(tc)
	}

	return &App{
		CommitHistory: ch,
		Client:        &client{Github: gh},
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
func (app *App) Create(ctx context.Context, latestCommit *Commit, filenames []string) error {
	action := "create"
	alias := app.Alias

	g, err := app.Client.Create(ctx, filenames)
	if err != nil {
		return errors.Wrapf(err, "gist api %s", action)
	}

	c := NewCommit(g.raw, app.Config.ResolveAlias(alias), action) // TODO(podhmo): remove *github.Gist
	if err := app.CommitHistory.SaveCommit(app.Config.Profile.HistFile, c); err != nil {
		return errors.Wrap(err, "save commit")
	}

	fmt.Fprintf(os.Stderr, "%s success. (id=%q)\n", action, c.ID)
	if !app.IsSilent && g.HTMLURL != "" {
		fmt.Fprintf(os.Stderr, "opening.. %q\n", g.HTMLURL)
		webbrowser.Open(g.HTMLURL)
	}
	// PrintJSON(g)
	return nil
}

// Update :
func (app *App) Update(ctx context.Context, latestCommit *Commit, filenames []string) error {
	action := "update"
	alias := app.Alias

	g, err := app.Client.Update(ctx, latestCommit.ID, filenames)
	if err != nil {
		return errors.Wrapf(err, "gist api %s", action)
	}

	c := NewCommit(g.raw, app.Config.ResolveAlias(alias), action) // TODO(podhmo): remove *github.Gist
	if err := app.CommitHistory.SaveCommit(app.Config.Profile.HistFile, c); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "%s success. (id=%q)\n", action, c.ID)
	if !app.IsSilent {
		fmt.Fprintf(os.Stderr, "opening.. %q\n", g.HTMLURL)
		webbrowser.Open(g.HTMLURL)
	}

	if app.Config.Debug {
		fprintJSON(os.Stderr, g)
	}
	return nil
}

// fprintJSON is pretty printed json output shorthand.
func fprintJSON(w io.Writer, data interface{}) {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		log.Fatal(err)
	}
}

// PrintJSON is similar that a relation about fmt.Printf and fmt.Fprintf.
func PrintJSON(data interface{}) {
	fprintJSON(os.Stdout, data)
}
