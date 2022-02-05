package history_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/podhmo/selfish/pkg/commithistory"
	"github.com/podhmo/selfish/pkg/commithistory/history"
	"github.com/podhmo/selfish/pkg/commithistory/timeutil"
)

func TestLoad(t *testing.T) {
	t.Run("file is not found, return EOF", func(t *testing.T) {
		var ob commithistory.Commit
		err := history.LoadFile("./notfound.history", &ob, "head")
		if err == nil {
			t.Fatal("must be error. but not error")
		}
	})
	t.Run("file is not found, be enable to check by IsNotfound()", func(t *testing.T) {
		var ob commithistory.Commit
		err := history.LoadFile("./notfound.history", &ob, "head")
		if !history.IsNotFound(err) {
			t.Errorf("must be not found error. but %s\n", err)
		}
	})

	t.Run("finding", func(t *testing.T) {
		type C struct {
			msg      string
			source   []string
			alias    string
			expected *commithistory.Commit
		}

		candidates := []C{
			{
				msg: "found",
				source: []string{
					"435048d99f77300e5c33dd7ab46dbca2,head,2018-03-10T16:23:57+09:00,create",
				},
				alias: "head",
				expected: &commithistory.Commit{
					ID:        "435048d99f77300e5c33dd7ab46dbca2",
					Alias:     "head",
					CreatedAt: timeutil.MustParse("2018-03-10T16:23:57+09:00"),
					Action:    "create",
				},
			},
			{
				msg: "not found",
				source: []string{
					"435048d99f77300e5c33dd7ab46dbca2,todo,2018-03-10T16:23:57+09:00,create",
				},
				alias:    "head",
				expected: nil,
			},
			{
				msg: "found, another",
				source: []string{
					"435048d99f77300e5c33dd7ab46dbca3,todo,2018-03-10T16:23:57+09:00,create",
					"435048d99f77300e5c33dd7ab46dbca2,head,2018-03-10T16:23:57+09:00,create",
				},
				alias: "head",
				expected: &commithistory.Commit{
					ID:        "435048d99f77300e5c33dd7ab46dbca2",
					Alias:     "head",
					CreatedAt: timeutil.MustParse("2018-03-10T16:23:57+09:00"),
					Action:    "create",
				},
			},
		}

		for _, c := range candidates {
			c := c
			t.Run(c.msg, func(t *testing.T) {
				fp, err := ioutil.TempFile(".", "")
				if err != nil {
					t.Fatal(err)
				}
				for _, line := range c.source {
					fmt.Fprintln(fp, line)
				}
				fp.Close()
				defer os.Remove(fp.Name())

				var got commithistory.Commit

				err = history.LoadFile(fp.Name(), &got, c.alias)
				if err != nil {
					if !history.IsNotFound(err) {
						t.Fatalf("unexpected error %s", err)
					}

					if c.expected == nil {
						return
					}
					t.Fatalf("invalid, must not be NOT FOUND")
				}

				if c.expected.ID != got.ID {
					t.Errorf("expected ID %q, but got %q\n", c.expected.ID, got.ID)
				}
				if c.expected.Alias != got.Alias {
					t.Errorf("expected Alias %q, but got %q\n", c.expected.Alias, got.Alias)
				}
				if c.expected.CreatedAt != got.CreatedAt {
					t.Errorf("expected CreatedAt %q, but got %q\n", c.expected.CreatedAt, got.CreatedAt)
				}
				if c.expected.Action != got.Action {
					t.Errorf("expected Action %q, but got %q\n", c.expected.Action, got.Action)
				}
			})
		}
	})
}
