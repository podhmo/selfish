package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/podhmo/selfish"
	"github.com/podhmo/structflag"
)

func main() {
	config := &selfish.Config{}
	b := structflag.NewBuilder()
	fs := b.Build(config)
	if err := fs.Parse(os.Args[1:]); err != nil {
		fs.Usage()
		os.Exit(2)
	}
	config.Files = fs.Args()

	if err := run(config); err != nil {
		log.Fatalf("!!%+v", err)
	}
}

func run(config *selfish.Config) error {
	ctx := context.Background()
	app, err := selfish.NewApp(config)
	if err != nil {
		if err == selfish.ErrAccessTokenNotfound {
			fmt.Fprintln(os.Stderr, "if config file is not found. then")
			fmt.Println(`
		mkdir -p ~/.config/selfish
		cat <<-EOS > ~/.config/selfish/config.json
		{
		  "access_token": "<your github access token>"
		}
		EOS
		`)
			os.Exit(1)
		}
		return err
	}

	var latestCommit *selfish.Commit
	if config.Alias != "" {
		var c selfish.Commit
		if err := app.CommitHistory.LoadCommit(app.Config.Profile.HistFile, config.Alias, &c); err != nil {
			if !app.CommitHistory.IsNotFound(err) {
				return errors.Wrap(err, "load commit")
			}
		} else {
			latestCommit = &c
		}
	}

	files := config.Files
	if app.IsDelete && config.Alias != "" {
		return app.Delete(ctx, latestCommit, config.Alias)
	} else if latestCommit == nil {
		return app.Create(ctx, latestCommit, config.Alias, files)
	} else {
		return app.Update(ctx, latestCommit, config.Alias, files)
	}
}
