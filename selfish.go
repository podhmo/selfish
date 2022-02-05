package selfish

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"github.com/podhmo/selfish/pkg/commithistory"
	"github.com/toqueteos/webbrowser"
	"golang.org/x/oauth2"
)

// Config is mapping object for application config
type Config struct {
	Profile struct {
		DefaultAlias string `json:"default_alias"`
		AccessToken  string `json:"access_token"`
		HistFile     string `json:"hist_file"`
	} `flag:"-"`

	Alias    string   `flag:"alias" help:"alias name of uploaded gists"`
	IsDelete bool     `flag:"delete" help:"delete uploaded gists"`
	IsSilent bool     `flag:"silent" help:"don't open gist pages with browser, after uploading"`
	Files    []string `flag:"-"`
}

// ResolveAlias :
func (c *Config) ResolveAlias(alias string) string {
	if alias == "" {
		return c.Profile.DefaultAlias
	}
	return alias
}

const (
	defaultAlias    = "head"
	defaultHistFile = "selfish.history"
)

// App :
type App struct {
	CommitHistory *commithistory.API
	Client        *github.Client

	*Config
}

var ErrAccessTokenNotfound = fmt.Errorf("access token is not found")

func NewApp(c *Config) (*App, error) {
	ch := commithistory.New("selfish")
	{
		profile := c.Profile
		if err := ch.Load("config.json", &profile); err != nil {
			return nil, errors.Wrap(err, "load config")
		}
		if profile.AccessToken == "" {
			return nil, ErrAccessTokenNotfound
		}

		if profile.DefaultAlias == "" {
			profile.DefaultAlias = defaultAlias
		}
		if profile.HistFile == "" {
			profile.HistFile = defaultHistFile
		}
	}

	var client *github.Client
	{
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: c.Profile.AccessToken},
		)
		tc := oauth2.NewClient(oauth2.NoContext, ts)
		client = github.NewClient(tc)
	}

	return &App{
		CommitHistory: ch,
		Client:        client,
		Config:        c,
	}, nil
}

// Delete :
func (app *App) Delete(ctx context.Context, latestCommit *Commit, alias string) error {
	if latestCommit == nil {
		return errors.Errorf("alias=%q is not found", alias)
	}

	gistID := latestCommit.ID
	if _, err := app.Client.Gists.Delete(ctx, gistID); err != nil {
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
func (app *App) Create(ctx context.Context, latestCommit *Commit, alias string, filenames []string) error {
	action := "create"

	gist, err := NewGist(filenames)
	if err != nil {
		return err
	}
	g, _, err := app.Client.Gists.Create(ctx, gist)
	if err != nil {
		return errors.Wrapf(err, "gist api %s", action)
	}

	c := NewCommit(g, app.Config.ResolveAlias(alias), action)
	if err := app.CommitHistory.SaveCommit(app.Config.Profile.HistFile, c); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "%s success. (id=%q)\n", action, c.ID)
	if !app.IsSilent {
		fmt.Fprintf(os.Stderr, "opening.. %q\n", *g.HTMLURL)
		webbrowser.Open(*g.HTMLURL)
	}
	// PrintJSON(g)
	return nil
}

// Update :
func (app *App) Update(ctx context.Context, latestCommit *Commit, alias string, filenames []string) error {
	action := "update"

	gist, err := NewGist(filenames)
	if err != nil {
		return err
	}
	g, _, err := app.Client.Gists.Edit(ctx, latestCommit.ID, gist)
	if err != nil {
		return errors.Wrapf(err, "gist api %s", action)
	}

	c := NewCommit(g, app.Config.ResolveAlias(alias), action)
	if err := app.CommitHistory.SaveCommit(app.Config.Profile.HistFile, c); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "%s success. (id=%q)\n", action, c.ID)
	if !app.IsSilent {
		fmt.Fprintf(os.Stderr, "opening.. %q\n", *g.HTMLURL)
		webbrowser.Open(*g.HTMLURL)
	}

	if ok, _ := strconv.ParseBool(os.Getenv("DEBUG")); ok {
		FprintJSON(os.Stderr, g)
	}
	return nil
}

// FprintJSON is pretty printed json output shorthand.
func FprintJSON(w io.Writer, data interface{}) {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		log.Fatal(err)
	}
}

// PrintJSON is similar that a relation about fmt.Printf and fmt.Fprintf.
func PrintJSON(data interface{}) {
	FprintJSON(os.Stdout, data)
}
