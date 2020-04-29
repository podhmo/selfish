package commithistory

import (
	"github.com/podhmo/selfish/internal/commithistory/config"
	"github.com/podhmo/selfish/internal/commithistory/history"
)

// Config :
type Config struct {
	*config.Config
}

// New :
func New(name string, ops ...func(*Config)) *Config {
	c := &Config{Config: config.New(name)}
	for _, op := range ops {
		op(c)
	}
	return c
}

// WithProfile :
func WithProfile(profile string) func(*Config) {
	return func(c *Config) {
		c.Config.Profile = profile
	}
}

// Load :
func (c *Config) Load(name string, ob interface{}) error {
	return c.Config.Load(name, ob)
}

// Save :
func (c *Config) Save(name string, ob interface{}) error {
	return c.Config.Save(name, ob)
}

// LoadCommit :
func (c *Config) LoadCommit(filename, alias string, ob history.Parsable) error {
	dirpath, err := c.Dir(c.Name)
	if err != nil {
		return err
	}
	path := c.JoinPath(c.Profile, dirpath, filename)
	return history.LoadFile(path, ob, alias)
}

// SaveCommit :
func (c *Config) SaveCommit(filename string, ob history.Unparsable) error {
	dirpath, err := c.Dir(c.Name)
	if err != nil {
		return err
	}
	path := c.JoinPath(c.Profile, dirpath, filename)
	return history.SaveFile(path, ob)
}

// IsNotFound :
func (c *Config) IsNotFound(err error) bool {
	return history.IsNotFound(err)
}
