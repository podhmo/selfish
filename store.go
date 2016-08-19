package selfish

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

import (
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)


// persistent

type commit struct {
	ID        string
	CreatedAt time.Time
	Alias     string // optional
}

// LoadCommit loading uploading history
func LoadCommit(filename string, alias string) (*commit, error) {
	if _, err := os.Stat(filename); err != nil {
		return nil, errors.Wrap(err, ":")
	}
	fp, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, ":")
	}
	defer fp.Close()
	return loadCommit(fp, alias)
}

func loadCommit(r io.Reader, alias string) (*commit, error) {
	sc := bufio.NewScanner(r)
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
