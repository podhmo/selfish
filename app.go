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

	"github.com/pkg/errors"
	"github.com/podhmo/selfish/model"
	"github.com/podhmo/selfish/pkg/commithistory"
	"github.com/toqueteos/webbrowser"
)

type Commit = model.Commit // TODO: remove

// App :
type App struct {
	CommitHistory *commithistory.API
	Client        *Client
	Config        *Config

	IsSilent bool
	IsDelete bool
	Alias    string
}

// Delete :
func (app *App) Delete(ctx context.Context, latestCommit *model.Commit, alias string) error {
	if latestCommit == nil {
		return errors.Errorf("alias=%q is not found", alias)
	}

	gistID := latestCommit.ID
	if _, err := app.Client.Gists.Delete(ctx, gistID); err != nil {
		return errors.Wrapf(err, "gist api delete")
	}

	c := model.Commit{ID: gistID, Alias: alias, CreatedAt: time.Now(), Action: "delete"}
	if err := app.CommitHistory.SaveCommit(app.Config.HistFile, &c); err != nil {
		return errors.Wrap(err, "save commit")
	}
	fmt.Fprintf(os.Stderr, "deleted. (id=%q)\n", gistID)
	return nil
}

// Create :
func (app *App) Create(ctx context.Context, latestCommit *model.Commit, alias string, filenames []string) error {
	action := "create"
	g, err := app.Client.Create(ctx, filenames)
	if err != nil {
		return errors.Wrapf(err, "gist api %s", action)
	}

	c := model.NewCommit(g, app.Config.ResolveAlias(alias), action)
	if err := app.CommitHistory.SaveCommit(app.Config.HistFile, c); err != nil {
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
	g, err := app.Client.Update(ctx, latestCommit, filenames)
	if err != nil {
		return errors.Wrapf(err, "gist api %s", action)
	}

	c := model.NewCommit(g, app.Config.ResolveAlias(alias), action)
	if err := app.CommitHistory.SaveCommit(app.Config.HistFile, c); err != nil {
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
