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

	g, _, err := client.Gists.Create(gist)
	if err != nil {
		return errors.Wrap(err, "gist api create")
	}

	c := selfish.NewCommit(g, defaultAlias)
	err = selfish.SaveCommit(defaultHistFile, c)
	if err != nil {
		return err
	}

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
