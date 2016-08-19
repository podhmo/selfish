package selfish

import (
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"os"
	"path"
)

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
		return nil, errors.Wrap(err, "stat")
	}
	size := int(finfo.Size())

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