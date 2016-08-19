package main

import (
	"flag"
	"fmt"
	"os"
)

import (
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"github.com/toqueteos/webbrowser"
)

import (
	"github.com/podhmo/selfish"
)

// AppMain is main function of Application
func AppMain(client *selfish.Client, alias string, filenames []string) error {
	gist, err := selfish.NewGist(filenames)
	if err != nil {
		return err
	}

	config := client.Config

	// W: ignore err
	var latestCommit *selfish.Commit
    if alias == "" {
        panic(alias)
    }
	if alias != "" {
		latestCommit, err = selfish.LoadCommit(config.HistFile, alias)
		if err != nil {
			return err
		}
	}

	var g *github.Gist
	var action string
	if latestCommit == nil {
		g, _, err = client.Gists.Create(gist)
		action = "create"
	} else {
		gistID := latestCommit.ID
		g, _, err = client.Gists.Edit(gistID, gist)
		action = "update"

	}

	if err != nil {
		return errors.Wrapf(err, "gist api %s", action)
	}

	c := selfish.NewCommit(g, config.DefaultAlias)
	err = selfish.SaveCommit(config.HistFile, c)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "%s success. (id=%q)\n", action, c.ID)
	fmt.Fprintf(os.Stderr, "redirect to %q\n", *g.HTMLURL)
	webbrowser.Open(*g.HTMLURL)
	// selfish.PrintJSON(g)
	return nil
}

var alias = flag.String("alias", "", "alias name of uploaded gist")

func main() {
	flag.Parse()
	client, err := selfish.CreateClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		fmt.Fprintln(os.Stderr, "if config file is not found. then")
		fmt.Println(`
mkdir -p ~/.selfish
cat <<-EOS > ~/.selfish/config.json
{
  "access_token": "<your github access token>"
}
EOS
`)
		os.Exit(1)
	}
	err = AppMain(client, *alias, flag.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
