package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

func ppJSON(target interface{}) {
	b, err := json.Marshal(target)
	if err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer
	json.Indent(&out, b, " ", "    ")
	out.WriteTo(os.Stdout)
}

// CreateClient is factory of github client
func CreateClient(token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	return github.NewClient(tc)
}

// NewGist is shorthand of github.Gist object creation
func NewGist(filenames []string) (*github.Gist, error) {
	public := true
	files := make(map[github.GistFilename]github.GistFile)

	for _, filename := range filenames {
		gistfile, err := NewGistFile(filename)
		if err != nil {
			log.Printf("skip file=%s err=%v\n", filename, err)
			continue
		}
		k := github.GistFilename(path.Base(filename))
		files[k] = *gistfile
	}

	gist := github.Gist{
		Public: &public,
		Files:  files,
	}
	return &gist, nil
}

// NewGistFile is shorthand of github.GistFile object creation
func NewGistFile(filename string) (*github.GistFile, error) {
	basename := path.Base(filename)
	finfo, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}
	size := int(finfo.Size())

	byte, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	content := string(byte)

	gistfile := github.GistFile{
		Size:     &size,
		Filename: &basename,
		Content:  &content,
	}
	return &gistfile, nil
}

// config
const (
	defaultHistFile string = "selfish.history"
	defaultAlias    string = "head"
)

// persistent
type commit struct {
	ID        string
	CreatedAt time.Time
	Alias     string // optional
}

func loadCommit(filename string, alias string) (*commit, error) {
	if _, err := os.Stat(filename); err != nil {
		return nil, errors.Wrap(err, ":")
	}
	fp, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, ":")
	}
	defer fp.Close()

	sc := bufio.NewScanner(fp)
	for sc.Scan() {
		line := sc.Text()
		// id@alias@CreatedAt
		data := strings.SplitN(line, "@", 3)
		if data[1] == alias {
			createdAt, err := time.Parse(time.RubyDate, data[2])
			if err != nil {
				return nil, errors.Wrap(err, ":")
			}
			c := commit{ID: data[0], Alias: data[1], CreatedAt: createdAt}
			return &c, nil
		}
	}
	return nil, nil
}

func saveCommit(filename string, c commit) error {
	fp, err := ioutil.TempFile(".", filename)
	if err != nil {
		return errors.Wrap(err, ":")
	}
	w := bufio.NewWriter(fp)
	defer func() {
		tmpname := fp.Name()
		fp.Close()
		os.Rename(path.Join(".", tmpname), path.Join(".", filename))
	}()

	createdAt := c.CreatedAt.Format(time.RubyDate)
	// id@alias@CreatedAt
	fmt.Fprintf(w, "%s@%s@%s\n", c.ID, c.Alias, createdAt)

	if _, err := os.Stat(filename); err == nil {
		fp, err := os.Open(filename)
		if err != nil {
			return errors.Wrap(err, ":")
		}
		defer fp.Close()
		sc := bufio.NewScanner(fp)
		for sc.Scan() {
			buf := sc.Bytes()
			w.Write(buf)
		}
	}
	w.Flush()
	return nil
}

func newCommit(g *github.Gist, alias string) commit {
	c := commit{
		ID:        *g.ID,
		CreatedAt: *g.CreatedAt,
		Alias:     alias,
	}
	return c
}

// AppMain is main function of Application
func AppMain(client *github.Client, filenames []string) {
	gist, err := NewGist(filenames)
	if err != nil {
		fmt.Printf("%+v\n", err)
		log.Fatal(err)
	}

	g, response, err := client.Gists.Create(gist)
	if err != nil {
		fmt.Printf("%+v\n", err)
		log.Fatal(err)
	}

	c := newCommit(g, defaultAlias)
	err = saveCommit(defaultHistFile, c)
	if err != nil {
		fmt.Printf("%+v\n", err)
		log.Fatal(err)
	}

	fmt.Println("g ----------------------------------------")
	ppJSON(g)
	fmt.Println("response ----------------------------------------")
	ppJSON(response)
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Fprintf(os.Stderr, "selfish <token>\n")
		os.Exit(1)
	}
	token := os.Args[1]
	client := CreateClient(token)
	AppMain(client, os.Args[2:])
}
