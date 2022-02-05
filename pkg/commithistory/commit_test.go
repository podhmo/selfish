package commithistory

import (
	"bytes"
	"strings"
	"testing"

	"github.com/podhmo/selfish/pkg/commithistory/timeutil"
)

func TestUnparse(t *testing.T) {
	c := Commit{
		ID:        "435048d99f77300e5c33dd7ab46dbca2",
		CreatedAt: timeutil.MustParse("2018-03-10T16:23:57+09:00"),
		Alias:     "head",
		Action:    "create",
	}
	var b bytes.Buffer
	if err := c.Unparse(&b); err != nil {
		t.Fatal(err)
	}
	got := strings.TrimRight(b.String(), "\n")

	expected := "435048d99f77300e5c33dd7ab46dbca2,head,2018-03-10T16:23:57+09:00,create"
	if got != expected {
		t.Errorf("expected %q, but %q", expected, got)
	}
}

func TestParse(t *testing.T) {
	source := []string{"435048d99f77300e5c33dd7ab46dbca2", "head", "2018-03-10T16:23:57+09:00", "create"}
	var got Commit
	if err := got.Parse(source); err != nil {
		t.Fatal(err)
	}
	expected := Commit{
		ID:        "435048d99f77300e5c33dd7ab46dbca2",
		CreatedAt: timeutil.MustParse("2018-03-10T16:23:57+09:00"),
		Alias:     "head",
		Action:    "create",
	}

	if expected.ID != got.ID {
		t.Errorf("expected ID %q, but got %q\n", expected.ID, got.ID)
	}
	if expected.Alias != got.Alias {
		t.Errorf("expected Alias %q, but got %q\n", expected.Alias, got.Alias)
	}
	if expected.CreatedAt != got.CreatedAt {
		t.Errorf("expected CreatedAt %q, but got %q\n", expected.CreatedAt, got.CreatedAt)
	}
	if expected.Action != got.Action {
		t.Errorf("expected Action %q, but got %q\n", expected.Action, got.Action)
	}
}
