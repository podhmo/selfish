package commithistory

import (
	"github.com/podhmo/selfish/pkg/commithistory/history"
)

type API struct {
	*Config
}

func New(name string, ops ...func(*API)) *API {
	c := &API{Config: DefaultConfig()}
	c.Config.Name = name
	for _, op := range ops {
		op(c)
	}
	return c
}

func WithProfile(profile string) func(*API) {
	return func(c *API) {
		c.Config.Profile = profile
	}
}

func (c *API) LoadCommit(filename, alias string, ob history.Parsable) error {
	dirpath, err := c.Dir(c.Name)
	if err != nil {
		return err
	}
	path := c.JoinPath(c.Profile, dirpath, filename)
	return history.LoadFile(path, ob, alias)
}

func (c *API) SaveCommit(filename string, ob history.Unparsable) error {
	dirpath, err := c.Dir(c.Name)
	if err != nil {
		return err
	}
	path := c.JoinPath(c.Profile, dirpath, filename)
	return history.SaveFile(path, ob)
}

func (c *API) IsNotFound(err error) bool {
	return history.IsNotFound(err)
}
