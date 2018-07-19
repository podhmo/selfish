package internal

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/podhmo/commithistory"
	"github.com/podhmo/selfish"
	"github.com/toqueteos/webbrowser"
)

// App :
type App struct {
	C        *commithistory.Config // todo: rename
	Client   *selfish.Client
	Config   *selfish.Config
	IsSilent bool
}

// FindLatestCommit :
func (app *App) FindLatestCommit(filename, alias string) (*selfish.Commit, error) {
	var c selfish.Commit
	if err := app.C.LoadCommit(filename, alias, &c); err != nil {
		if app.C.IsNotFound(err) {
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
	if err := app.C.SaveCommit(app.Config.HistFile, &c); err != nil {
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
	if err := app.C.SaveCommit(app.Config.HistFile, c); err != nil {
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
	if err := app.C.SaveCommit(app.Config.HistFile, c); err != nil {
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
