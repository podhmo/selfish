package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/podhmo/commithistory"
	"github.com/podhmo/selfish"
	"github.com/podhmo/selfish/cmd/selfish/internal"
)

var aliasFlag = flag.String("alias", "", "alias name of uploaded gists")
var deleteFlag = flag.Bool("delete", false, "delete uploaded gists")
var silentFlag = flag.Bool("silent", false, "deactivate webbrowser open, after gists uploading")

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()
	c := commithistory.New("selfish")
	config, err := selfish.LoadConfig(c)

	if config.AccessToken == "" {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
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

	client := selfish.NewClient(config)

	app := &internal.App{Config: config, Client: client, C: c, IsSilent: *silentFlag}
	ctx := context.Background()

	var latestCommit *selfish.Commit
	if *aliasFlag != "" {
		latestCommit, err = app.FindLatestCommit(app.Config.HistFile, *aliasFlag)
		if err != nil {
			return err
		}
	}
	if *deleteFlag && *aliasFlag != "" {
		return app.Delete(ctx, latestCommit, *aliasFlag)
	} else if latestCommit == nil {
		return app.Create(ctx, latestCommit, *aliasFlag, flag.Args())
	} else {
		return app.Update(ctx, latestCommit, *aliasFlag, flag.Args())
	}
}
