package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

// Config :
type Config struct {
	Name     string
	Profile  string // optional
	M        MarshalUnmarshaller
	JoinPath func(profile string, paths ...string) string
	Default  func(path string) error
	Dir      func(name string) (string, error)
}

// New :
func New(name string, ops ...func(*Config)) *Config {
	c := &Config{Name: name}
	for _, op := range ops {
		op(c)
	}
	if c.Dir == nil {
		c.Dir = DefaultConfigDir
	}
	if c.Default == nil {
		c.Default = DefaultWriteFile
	}
	if c.JoinPath == nil {
		c.JoinPath = DefaultJoinPath
	}
	if c.M == nil {
		M := &JSONModule{}
		c.M = M
	}
	return c
}

// WithProfile :
func WithProfile(profile string) func(*Config) {
	return func(c *Config) {
		c.Profile = profile
	}
}

// WithMarshalUnmarshaller :
func WithMarshalUnmarshaller(m MarshalUnmarshaller) func(*Config) {
	return func(c *Config) {
		c.M = m
	}
}

// WithDirFunction :
func WithDirFunction(dir func(name string) (string, error)) func(*Config) {
	return func(c *Config) {
		c.Dir = dir
	}
}

// WithDefaultFunction :
func WithDefaultFunction(writefile func(path string) error) func(*Config) {
	return func(c *Config) {
		c.Default = writefile
	}
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

// Load :
func (c *Config) Load(name string, ob interface{}) error {
	d, err := c.Dir(c.Name)
	if err != nil {
		return err
	}

	path := c.JoinPath(c.Profile, d, name)
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
	return c.M.Unmarshal(fp, ob)
}

// Save :
func (c *Config) Save(name string, ob interface{}) error {
	d, err := c.Dir(c.Name)
	if err != nil {
		return err
	}

	path := c.JoinPath(c.Profile, d, name)
	log.Printf("save. %q\n", path)

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	fp, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fp.Close()
	return c.M.Marshal(fp, ob)
}
