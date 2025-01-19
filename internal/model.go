package internal

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/google/go-github/v68/github"
	"github.com/pkg/errors"
	"github.com/podhmo/selfish/pkg/commithistory"
)

type Commit = commithistory.Commit

// NewGist is shorthand of github.Gist object creation
func NewGist(filenames []string) (*github.Gist, error) {
	public := true
	title := ""
	files := make(map[github.GistFilename]github.GistFile)

	for _, filename := range filenames {
		gistfile, err := NewGistFile(filename)
		if err != nil {
			log.Printf("skip file=%s err=%v\n", filename, err)
			continue
		}

		// guess title, first heading of markdown. (This code using wasteful memory, but it's ok)
		if (strings.ToLower(filename)) == "readme.md" && gistfile.Content != nil {
			text := strings.TrimLeft(*gistfile.Content, "\n\t  ")
			for _, line := range strings.Split(text, "\n") {
				if strings.HasPrefix(line, "# ") {
					title = strings.TrimSpace(strings.TrimPrefix(line, "# "))
					break
				}
			}
		}

		k := github.GistFilename(path.Base(filename))
		files[k] = *gistfile
	}

	gist := github.Gist{
		Public:      &public,
		Files:       files,
		Description: &title,
	}
	return &gist, nil
}

// NewGistFile is shorthand of github.GistFile object creation
func NewGistFile(filename string) (*github.GistFile, error) {
	basename := path.Base(filename)
	finfo, err := os.Stat(filename)
	if err != nil {
		return nil, errors.Wrap(err, "stat")
	}
	size := int(finfo.Size())
	if size == 0 {
		return nil, fmt.Errorf("empty file")
	}

	byte, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "read all")
	}
	content := string(byte)

	gistfile := github.GistFile{
		Size:     &size,
		Filename: &basename,
		Content:  &content,
	}
	return &gistfile, nil
}
