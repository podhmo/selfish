package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/podhmo/selfish"
	"github.com/podhmo/selfish/cmd/selfish/internal"
	"github.com/podhmo/selfish/pkg/commithistory"
	"github.com/podhmo/structflag"
)

// Option ...
type Option struct {
	Alias  string   `flag:"alias" help:"alias name of uploaded gists"`
	Delete bool     `flag:"delete" help:"delete uploaded gists"`
	Silent bool     `flag:"silent" help:"don't open gist pages with browser, after uploading"`
	Args   []string `flag:"-"`
}

func main() {
	opt := &Option{}
	b := structflag.NewBuilder()
	fs := b.Build(opt)
	if err := fs.Parse(os.Args[1:]); err != nil {
		fs.Usage()
		os.Exit(2)
	}
	opt.Args = fs.Args()

	if err := run(opt); err != nil {
		log.Fatalf("!!%+v", err)
	}
}

func run(opt *Option) error {
	ch := commithistory.New("selfish")
	c, err := selfish.LoadConfig(ch)
	if err != nil {
		return err
	}
	if c.AccessToken == "" {
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

	ctx := context.Background()
	app := &internal.App{
		CommitHistory: ch,
		Client:        selfish.NewClient(c),
		Config:        c,
		IsSilent:      opt.Silent,
		IsDelete:      opt.Delete,
		Alias:         opt.Alias,
	}

	var latestCommit *selfish.Commit
	if app.Alias != "" {
		var c selfish.Commit
		if err := app.CommitHistory.LoadCommit(app.Config.HistFile, opt.Alias, &c); err != nil {
			if !app.CommitHistory.IsNotFound(err) {
				return errors.Wrap(err, "load commit")
			}
		} else {
			latestCommit = &c
		}
	}

	files := opt.Args
	if app.IsDelete && app.Alias != "" {
		return app.Delete(ctx, latestCommit, app.Alias)
	} else if latestCommit == nil {
		return app.Create(ctx, latestCommit, app.Alias, files)
	} else {
		return app.Update(ctx, latestCommit, app.Alias, files)
	}
}
