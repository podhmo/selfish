package commithistory

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

// Config :
type Config struct {
	Name    string
	Profile string // optional
	Dryrun  bool

	JoinPath func(profile string, paths ...string) string
	Default  func(path string) error
	Dir      func(name string) (string, error)
}

func DefaultConfig() *Config {
	c := &Config{}
	if c.Dir == nil {
		c.Dir = DefaultConfigDir
	}
	if c.Default == nil {
		c.Default = DefaultWriteFile
	}
	if c.JoinPath == nil {
		c.JoinPath = DefaultJoinPath
	}
	return c
}

// DefaultJoinPath :
func DefaultJoinPath(profile string, paths ...string) string {
	if profile == "" {
		return filepath.Join(paths...)
	}
	paths[len(paths)-1] = fmt.Sprintf("%s.%s", profile, paths[len(paths)-1])
	return filepath.Join(paths...)
}

// DefaultWriteFile :
func DefaultWriteFile(path string) error {
	log.Printf("create. %q\n", path)
	fp, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fp.Close()
	_, rerr := fmt.Fprintln(fp, "{}")
	return rerr
}

// DefaultConfigDir :
func DefaultConfigDir(name string) (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(u.HomeDir, ".config", name), nil
}

func (c *Config) FilePath(name string) (string, error) {
	d, err := c.Dir(c.Name)
	if err != nil {
		return "", err
	}

	path := c.JoinPath(c.Profile, d, name)
	return path, nil
}

// Load :
func (c *Config) Load(name string, ob interface{}) error {
	path, err := c.FilePath(name)
	if err != nil {
		return err
	}
	log.Printf("load. %q\n", path)

	var fp io.ReadCloser
	fp, err = os.Open(path)
	if err != nil {
		log.Printf("not found. %q\n", path)
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			return err
		}
		if err := c.Default(path); err != nil {
			return err
		}
		fp, err = os.Open(path)
		if err != nil {
			return err
		}
	}
	defer fp.Close()

	decoder := json.NewDecoder(fp)
	return decoder.Decode(ob)
}

// Save :
func (c *Config) Save(name string, ob interface{}) error {
	path, err := c.FilePath(name)
	if err != nil {
		return err
	}
	log.Printf("save. %q\n", path)

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	fp, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fp.Close()
	encoder := json.NewEncoder(fp)
	return encoder.Encode(ob)
}
