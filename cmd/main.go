package main

import (
	"fmt"
	"os"
	"path"
)

import (
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"github.com/toqueteos/webbrowser"
)

import (
	"github.com/podhmo/selfish"
)

// config
var (
	defaultHistFile string
	defaultAlias    string
)

func init() {
	defaultHistFile = path.Join(os.Getenv("HOME"), ".selfish.history")
	defaultAlias = "head"
	// fmt.Printf("history: %q, alias: %q\n", defaultHistFile, defaultAlias)
}

// AppMain is main function of Application
func AppMain(client *github.Client, filenames []string) error {
	gist, err := selfish.NewGist(filenames)
	if err != nil {
		return err
	}

	latestCommit, err := selfish.LoadCommit(defaultHistFile, "head")
	if err != nil {
		return err
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

	c := selfish.NewCommit(g, defaultAlias)
	err = selfish.SaveCommit(defaultHistFile, c)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "%s success. (id=%q)\n", action, c.ID)
	fmt.Fprintf(os.Stderr, "redirect to %q\n", *g.HTMLURL)
	webbrowser.Open(*g.HTMLURL)
	// selfish.PrintJSON(g)
	return nil
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Fprintf(os.Stderr, "selfish <token>\n")
		os.Exit(1)
	}
	token := os.Args[1]
	client := selfish.CreateClient(token)
	err := AppMain(client, os.Args[2:])
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}
