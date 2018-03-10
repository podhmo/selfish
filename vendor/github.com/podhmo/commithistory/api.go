package commithistory

import (
	"path/filepath"

	"github.com/podhmo/commithistory/config"
	"github.com/podhmo/commithistory/history"
)

// todo: --profile

// Config :
type Config struct {
	*config.Config
}

// New :
func New(name string) *Config {
	return &Config{Config: config.New(name)}
}

// LoadCommit :
func (c *Config) LoadCommit(filename, alias string, ob history.Parsable) error {
	dirpath, err := c.Dir(c.Name)
	if err != nil {
		return err
	}
	path := filepath.Join(dirpath, filename)
	return history.LoadFile(path, ob, alias)
}

// SaveCommit :
func (c *Config) SaveCommit(filename string, ob history.Unparsable) error {
	dirpath, err := c.Dir(c.Name)
	if err != nil {
		return err
	}
	path := filepath.Join(dirpath, filename)
	return history.SaveFile(path, ob)
}

// IsNotFound :
func (c *Config) IsNotFound(err error) bool {
	return history.IsNotFound(err)
}
