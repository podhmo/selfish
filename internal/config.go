package internal

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/pflag"
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
	ClientType ClientType `flag:"client"`
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

// for flagstruct.HasHelpText
func (v ClientType) HelpText() string {
	return "- client is the option for debug {github, fake}"
}

// for pflag.Value
func (v *ClientType) String() string {
	if v == nil {
		return "<nil>"
	}
	return string(*v)
}

// for pflag.Value
func (v *ClientType) Set(value string) error {
	if v == nil {
		*v = ClientTypeGithub
	} else {
		*v = ClientType(strings.ToLower(value))
	}
	return v.Validate()
}

// for pflag.Value
func (v *ClientType) Type() string {
	return "ClientType"
}

var _ pflag.Value = (*ClientType)(nil)
