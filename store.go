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

// Commit is a tiny expression of gist uploading history.
type Commit struct {
	ID        string
	CreatedAt time.Time
	Alias     string // optional
}

// LoadCommit loading uploading history
func LoadCommit(filename string, alias string) (*Commit, error) {
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

func loadCommit(r io.Reader, alias string) (*Commit, error) {
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
			c := Commit{ID: data[0], Alias: data[1], CreatedAt: createdAt}
			return &c, nil
		}
	}
	return nil, nil
}

// SaveCommit saving uploading history
func SaveCommit(filename string, c Commit) error {
	fp, err := ioutil.TempFile(".", filename)
	if err != nil {
		return errors.Wrap(err, ":")
	}
	w := bufio.NewWriter(fp)
	defer func() {
		w.Flush()
		tmpname := fp.Name()
		fp.Close()
		os.Rename(path.Join(".", tmpname), path.Join(".", filename))
	}()
	if _, err := os.Stat(filename); err != nil {
		return saveCommit(w, nil, c)
	}

	rp, err := os.Open(filename)
	if err != nil {
		return errors.Wrap(err, ":")
	}
	defer rp.Close()
	return saveCommit(w, rp, c)
}

func saveCommit(w io.Writer, r io.Reader, c Commit) error {
	createdAt := c.CreatedAt.Format(time.RubyDate)
	// id@alias@CreatedAt
	fmt.Fprintf(w, "%s@%s@%s\n", c.ID, c.Alias, createdAt)

	if r != nil {
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			buf := sc.Bytes()
			w.Write(buf)
		}
	}
	return nil
}

func newCommit(g *github.Gist, alias string) Commit {
	c := Commit{
		ID:        *g.ID,
		CreatedAt: *g.CreatedAt,
		Alias:     alias,
	}
	return c
}
