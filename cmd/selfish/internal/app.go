package internal

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/podhmo/selfish"
	"github.com/podhmo/selfish/cmd/selfish/internal/pputil"
	"github.com/podhmo/selfish/internal/commithistory"
	"github.com/toqueteos/webbrowser"
)

// App :
type App struct {
	CommitHistory *commithistory.API
	Client        *selfish.Client
	Config        *selfish.Config

	IsSilent bool
	IsDelete bool
	Alias    string
}

// NewApp :
func NewApp(
	commitHistory *commithistory.API,
	Client *selfish.Client,
	Config *selfish.Config,

	IsSilent bool,
	IsDelete bool,
	Alias string,
) *App {
	return &App{
		CommitHistory: commitHistory,
		Client:        Client,
		Config:        Config,

		IsSilent: IsSilent,
		IsDelete: IsDelete,
		Alias:    Alias,
	}
}

// Run :
func (app *App) Run(ctx context.Context, files []string) (err error) {
	var latestCommit *selfish.Commit
	if app.Alias != "" {
		latestCommit, err = app.FindLatestCommit(app.Config.HistFile, app.Alias)
		if err != nil {
			return err
		}
	}

	if app.IsDelete && app.Alias != "" {
		return app.Delete(ctx, latestCommit, app.Alias)
	} else if latestCommit == nil {
		return app.Create(ctx, latestCommit, app.Alias, files)
	} else {
		return app.Update(ctx, latestCommit, app.Alias, files)
	}
}

// FindLatestCommit :
func (app *App) FindLatestCommit(filename, alias string) (*selfish.Commit, error) {
	var c selfish.Commit
	if err := app.CommitHistory.LoadCommit(filename, alias, &c); err != nil {
		if app.CommitHistory.IsNotFound(err) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "load commit")
	}
	return &c, nil
}

// Delete :
func (app *App) Delete(ctx context.Context, latestCommit *selfish.Commit, alias string) error {
	if latestCommit == nil {
		return errors.Errorf("alias=%q is not found", alias)
	}

	gistID := latestCommit.ID
	if _, err := app.Client.Gists.Delete(ctx, gistID); err != nil {
		return errors.Wrapf(err, "gist api delete")
	}

	c := selfish.Commit{ID: gistID, Alias: alias, CreatedAt: time.Now(), Action: "delete"}
	if err := app.CommitHistory.SaveCommit(app.Config.HistFile, &c); err != nil {
		return errors.Wrap(err, "save commit")
	}
	fmt.Fprintf(os.Stderr, "deleted. (id=%q)\n", gistID)
	return nil
}

// Create :
func (app *App) Create(ctx context.Context, latestCommit *selfish.Commit, alias string, filenames []string) error {
	action := "create"
	g, err := app.Client.Create(ctx, filenames)
	if err != nil {
		return errors.Wrapf(err, "gist api %s", action)
	}

	c := selfish.NewCommit(g, app.Config.ResolveAlias(alias), action)
	if err := app.CommitHistory.SaveCommit(app.Config.HistFile, c); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "%s success. (id=%q)\n", action, c.ID)
	if !app.IsSilent {
		fmt.Fprintf(os.Stderr, "opening.. %q\n", *g.HTMLURL)
		webbrowser.Open(*g.HTMLURL)
	}
	// pputil.PrintJSON(g)
	return nil
}

// Update :
func (app *App) Update(ctx context.Context, latestCommit *selfish.Commit, alias string, filenames []string) error {
	action := "update"
	g, err := app.Client.Update(ctx, latestCommit, filenames)
	if err != nil {
		return errors.Wrapf(err, "gist api %s", action)
	}

	c := selfish.NewCommit(g, app.Config.ResolveAlias(alias), action)
	if err := app.CommitHistory.SaveCommit(app.Config.HistFile, c); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "%s success. (id=%q)\n", action, c.ID)
	if !app.IsSilent {
		fmt.Fprintf(os.Stderr, "opening.. %q\n", *g.HTMLURL)
		webbrowser.Open(*g.HTMLURL)
	}

	if ok, _ := strconv.ParseBool(os.Getenv("DEBUG")); ok {
		pputil.FprintJSON(os.Stderr, g)
	}
	return nil
}
