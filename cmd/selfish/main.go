package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

import (
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"github.com/toqueteos/webbrowser"
)

import (
	"github.com/podhmo/selfish"
)

func appDelete(client *selfish.Client, alias string) error {
	config := client.Config

	// W: ignore err
	var latestCommit *selfish.Commit
	latestCommit, err := selfish.LoadCommit(config.HistFile, alias)
	if err != nil {
		return err
	}
	if latestCommit == nil {
		return errors.Errorf("alias=%q is not found", alias)
	}

	gistID := latestCommit.ID
	_, err = client.Gists.Delete(gistID)

	if err != nil {
		return errors.Wrapf(err, "gist api delete")
	}

	c := selfish.Commit{ID: gistID, Alias: alias, CreatedAt: time.Now(), Action: "delete"}
	err = selfish.SaveCommit(config.HistFile, c)
	fmt.Fprintf(os.Stderr, "deleted. (id=%q)\n", gistID)
	return nil
}

func appMain(client *selfish.Client, alias string, filenames []string) error {
	gist, err := selfish.NewGist(filenames)
	if err != nil {
		return err
	}

	config := client.Config

	// W: ignore err
	var latestCommit *selfish.Commit
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

	var saveAlias string
	if alias == "" {
		saveAlias = config.DefaultAlias
	} else {
		saveAlias = alias
	}

	c := selfish.NewCommit(g, saveAlias, action)
	err = selfish.SaveCommit(config.HistFile, c)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "%s success. (id=%q)\n", action, c.ID)
	if !*silentFlag {
		fmt.Fprintf(os.Stderr, "opening.. %q\n", *g.HTMLURL)
		webbrowser.Open(*g.HTMLURL)
	}
	// selfish.PrintJSON(g)
	return nil
}

var aliasFlag = flag.String("alias", "", "alias name of uploaded gists")
var deleteFlag = flag.Bool("delete", false, "delete uploaded gists")
var silentFlag = flag.Bool("silent", false, "deactivate webbrowser open, after gists uploading")

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
	if *deleteFlag && *aliasFlag != "" {
		err = appDelete(client, *aliasFlag)
	} else {
		err = appMain(client, *aliasFlag, flag.Args())
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
