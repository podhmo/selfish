package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/podhmo/selfish"
	"github.com/podhmo/structflag"
)

func main() {
	config := &selfish.Config{}
	b := structflag.NewBuilder()
	fs := b.Build(config)
	if err := fs.Parse(os.Args[1:]); err != nil {
		fs.Usage()
		os.Exit(2)
	}
	config.Files = fs.Args()

	if err := run(config); err != nil {
		log.Fatalf("!!%+v", err)
	}
}

func run(config *selfish.Config) error {
	ctx := context.Background()
	app, err := selfish.NewApp(config)
	if err != nil {
		if err == selfish.ErrAccessTokenNotfound {
			fmt.Fprintln(os.Stderr, "if config file is not found. then")
			fmt.Println(`
		mkdir -p ~/.config/selfish
		cat <<-EOS > ~/.config/selfish/config.json
		{
		  "access_token": "<your github access token>"
		}
		EOS
		`)
			os.Exit(1)
		}
		return err
	}

	var latestCommit *selfish.Commit
	if config.Alias != "" {
		var c selfish.Commit
		if err := app.CommitHistory.LoadCommit(app.Config.Profile.HistFile, config.Alias, &c); err != nil {
			if !app.CommitHistory.IsNotFound(err) {
				return errors.Wrap(err, "load commit")
			}
		} else {
			latestCommit = &c
		}
	}

	scanResult := Scan(config.Files)

	var commit *selfish.Commit
	if app.IsDelete && config.Alias != "" {
		return app.Delete(ctx, latestCommit)
	} else if latestCommit == nil {
		commit, err = app.Create(ctx, latestCommit, scanResult.TextFiles)
	} else {
		commit, err = app.Update(ctx, latestCommit, scanResult.TextFiles)
	}

	if len(scanResult.BinaryFiles) > 0 && commit != nil {
		// TODO:
		fmt.Printf("git clone git@github.com:%s.git\n", commit.ID)
	}
	return nil
}

type ScanResult struct {
	TextFiles   []string
	BinaryFiles []string
}

const (
	TooLargeFileSize = 5 * (1024 * 1024) // 5Mb
)

func Scan(files []string) ScanResult {
	// TODO(podhmo): use io/fs
	textFiles := make([]string, 0, len(files))
	binaryFiles := make([]string, 0, len(files))

	for _, fname := range files {
		finfo, err := os.Stat(fname)
		if err != nil {
			log.Printf("ignored for %+v (%q)", err, fname)
			continue
		}

		if finfo.Size() > TooLargeFileSize {
			binaryFiles = append(binaryFiles, fname)
			continue
		}

		if err := func(fname string) error {
			b, err := os.ReadFile(fname)
			if err != nil {
				return err
			}
			contentType := http.DetectContentType(b)
			if strings.HasPrefix(contentType, "text/") {
				textFiles = append(textFiles, fname)
			} else {
				binaryFiles = append(binaryFiles, fname)
			}
			return nil
		}(fname); err != nil {
			log.Printf("ignored for %+v, in detect file type (%q)", err, fname)

		}
	}

	r := ScanResult{
		TextFiles:   textFiles,
		BinaryFiles: binaryFiles,
	}
	json.NewEncoder(os.Stderr).Encode(r)
	return r
}
