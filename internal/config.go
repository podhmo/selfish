package internal

import (
	"fmt"
	"reflect"
)

// Config is mapping object for application config
type Config struct {
	Profile struct {
		DefaultAlias        string `json:"default_alias"`
		AccessToken         string `json:"access_token"`
		HistFile            string `json:"hist_file"`
		RepositoryDirectory string `json:"repositoryDirectory"`
	} `flag:"-"`

	Alias      string     `flag:"alias" help:"alias name of uploaded gists"`
	IsDelete   bool       `flag:"delete" help:"delete uploaded gists"`
	IsSilent   bool       `flag:"silent" help:"don't open gist pages with browser, after uploading"`
	ClientType ClientType `flag:"client" help:"if =fake, doesn't request {github, fake}"`
	Files      []string   `flag:"-"`

	Debug bool `flag:"debug"`
}

// ResolveAlias :
func (c *Config) ResolveAlias(alias string) string {
	if alias == "" {
		return c.Profile.DefaultAlias
	}
	return alias
}

////////////////////////////////////////

type ClientType string

const (
	ClientTypeGithub ClientType = "github"
	ClientTypeFake   ClientType = "fake"
)

func (v ClientType) Validate() error {
	switch v {
	case ClientTypeGithub, ClientTypeFake:
		return nil
	default:
		return fmt.Errorf("%v is an invalid value for %v", v, reflect.TypeOf(v))
	}
}

// for flags.TextVar

func (v *ClientType) UnmarshalText(data []byte) error {
	*v = ClientType(data)
	return v.Validate()
}

func (v ClientType) MarshalText() ([]byte, error) {
	return []byte(v), nil
}
