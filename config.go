package selfish

import (
	"github.com/pkg/errors"
	"github.com/podhmo/selfish/internal/commithistory"
)

// Config is mapping object for application config
type Config struct {
	DefaultAlias string `json:"default_alias"`
	AccessToken  string `json:"access_token"`
	HistFile     string `json:"hist_file"`
}

const (
	defaultAlias    = "head"
	defaultHistFile = "selfish.history"
)

// ResolveAlias :
func (c *Config) ResolveAlias(alias string) string {
	if alias == "" {
		return c.DefaultAlias
	}
	return alias
}

// LoadConfig loads configuration file, if configuration file is not existed, then return default config.
func LoadConfig(c *commithistory.Config) (*Config, error) {
	var conf Config
	if err := c.Load("config.json", &conf); err != nil {
		return nil, errors.Wrap(err, "load config")
	}
	if conf.DefaultAlias == "" {
		conf.DefaultAlias = defaultAlias
	}
	if conf.HistFile == "" {
		conf.HistFile = defaultHistFile
	}
	return &conf, nil
}
