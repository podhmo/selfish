package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/podhmo/flagstruct"
	"github.com/podhmo/selfish/internal"
)

func main() {
	config := &internal.Config{ClientType: internal.ClientTypeGithub}
	b := flagstruct.NewBuilder()
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

func run(config *internal.Config) error {
	ctx := context.Background()
	app, err := internal.NewApp(config)
	if err != nil {
		if err == internal.ErrAccessTokenNotfound {
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

	var latestCommit *internal.Commit
	if config.Alias != "" {
		var c internal.Commit
		if err := app.CommitHistory.LoadCommit(app.Config.Profile.HistFile, config.Alias, &c); err != nil {
			if !app.CommitHistory.IsNotFound(err) {
				return errors.Wrap(err, "load commit")
			}
		} else {
			latestCommit = &c
		}
	}

	scanResult := ScanFiles(config.Files)
	if app.Config.Debug {
		json.NewEncoder(os.Stderr).Encode(scanResult)
	}

	var commit *internal.Commit
	if app.IsDelete && config.Alias != "" {
		return app.Delete(ctx, latestCommit)
	} else if latestCommit == nil {
		commit, err = app.Create(ctx, latestCommit, scanResult.TextFiles)
	} else if len(scanResult.TextFiles) == 0 {
		log.Println("empty, skipped")
		commit = latestCommit
	} else {
		commit, err = app.Update(ctx, latestCommit, scanResult.TextFiles)
	}

	if err != nil {
		return err
	}

	// handling image files (see: https://gist.github.com/mroderick/1afdd71aa69f6b29601d335751a1a9be)
	if len(scanResult.BinaryFiles) > 0 && commit != nil {
		if app.Config.Debug {
			fmt.Fprintln(os.Stderr, "----------------------------------------")
			fmt.Fprintln(os.Stderr, "binary files are detected", scanResult.BinaryFiles)
			fmt.Fprintln(os.Stderr, "----------------------------------------")
		}

		// TODO(podhmo): refactoring
		rootDir, err := app.CommitHistory.FilePath(app.Config.Profile.RepositoryDirectory)
		if err != nil {
			log.Printf("WARN: handling binaries is failed. %+v\nignored.", err)
			return nil
		}
		if err := os.MkdirAll(rootDir, 0744); err != nil {
			log.Printf("WARN: handling binaries is failed. %+v\nignored.", err)
			return nil
		}
		repoDir := filepath.Join(rootDir, commit.ID)
		repoURL := fmt.Sprintf("git@github.com:%s.git", commit.ID)

		if _, err := os.Stat(repoDir); err != nil {
			if !os.IsNotExist(err) {
				log.Printf("WARN: handling binaries is failed. %+v\nignored.", err)
				return nil
			}

			log.Println("clone repository; git clone", repoURL)
			cmd := exec.Command("git", "clone", repoURL)
			cmd.Dir = rootDir
			if err := cmd.Run(); err != nil {
				log.Printf("WARN: handling binaries is failed. in git clone. %+v\nignored.", err)
				return nil
			}
		}

		log.Printf("sync repository; cd %s && git pull", repoDir)
		cmd := exec.Command("git", "pull", "--rebase")
		cmd.Dir = repoDir
		if err := cmd.Run(); err != nil {
			log.Printf("git pull is failed. %+v (but continued)", err)
		}

		copied := make([]string, 0, len(scanResult.BinaryFiles))
		for _, fname := range scanResult.BinaryFiles {
			if err := func(fname string) error {
				rf, err := os.Open(fname)
				if err != nil {
					return err
				}
				defer rf.Close()
				wf, err := os.Create(filepath.Join(repoDir, filepath.Base(fname)))
				if err != nil {
					return err
				}
				defer wf.Close()
				if _, err := io.Copy(wf, rf); err != nil {
					return err
				}
				copied = append(copied, filepath.Base(fname))
				return nil
			}(fname); err != nil {
				log.Printf("ignored for %+v, in copy file (%q)", err, fname)
				continue
			}
		}

		if len(copied) > 0 {
			log.Printf("git add && git push")
			{
				cmd := exec.Command("git", append([]string{"add"}, copied...)...)
				cmd.Dir = repoDir
				if err := cmd.Run(); err != nil {
					log.Printf("WARN: handling binaries is failed. in git add. %+v\nignored.", err)
					return nil
				}
			}
			{
				cmd := exec.Command("git", "commit", "-m", "from selfish")
				cmd.Dir = repoDir
				if err := cmd.Run(); err != nil {
					log.Printf("WARN: handling binaries is failed. in git commit. %+v\nignored.", err)
					return nil
				}
			}
			{
				cmd := exec.Command("git", "push")
				cmd.Dir = repoDir
				if err := cmd.Run(); err != nil {
					log.Printf("WARN: handling binaries is failed. in git push. %+v\nignored.", err)
					return nil
				}
			}
		}
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

func ScanFiles(files []string) ScanResult {
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
	return r
}
