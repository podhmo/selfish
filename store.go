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
		return nil, errors.Wrap(err, "stat")
	}
	fp, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "open")
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
				return nil, errors.Wrap(err, "time.parse")
			}
			c := Commit{ID: data[0], Alias: data[1], CreatedAt: createdAt}
			return &c, nil
		}
	}
	return nil, nil
}

// SaveCommit saving uploading history
func SaveCommit(filename string, c Commit) error {
	fp, err := ioutil.TempFile("", path.Base(filename))
	if err != nil {
		return errors.Wrap(err, "tempfile")
	}
	w := bufio.NewWriter(fp)
	defer func() {
		w.Flush()
		tmpname := fp.Name()
		fp.Close()
		os.Rename(tmpname, filename)
	}()
	if _, err := os.Stat(filename); err != nil {
		return saveCommit(w, nil, c)
	}

	rp, err := os.Open(filename)
	if err != nil {
		return errors.Wrap(err, "open")
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
		newline := []byte("\n")
		for sc.Scan() {
			buf := sc.Bytes()
			w.Write(buf)
			w.Write(newline)
		}
	}
	return nil
}

// NewCommit creates and initializes a new Commit object.
func NewCommit(g *github.Gist, alias string) Commit {
	c := Commit{
		ID:        *g.ID,
		CreatedAt: *g.CreatedAt,
		Alias:     alias,
	}
	return c
}
