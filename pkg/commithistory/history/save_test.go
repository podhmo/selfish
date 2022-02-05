package history_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/podhmo/selfish/pkg/commithistory"
	"github.com/podhmo/selfish/pkg/commithistory/history"
)

func must(t time.Time, err error) time.Time {
	if err != nil {
		panic(err)
	}
	return t
}

func TestSave(t *testing.T) {
	fp, err := ioutil.TempFile(".", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(fp.Name())

	fmt.Fprintln(fp, "435048d99f77300e5c33dd7ab46dbca2,head,2018-03-10T16:23:57+09:00,create")
	fp.Close()

	c := commithistory.Commit{
		ID:        "435048d99f77300e5c33dd7ab46dbca2",
		CreatedAt: must(time.Parse(time.RFC3339, "2018-03-10T16:23:57+09:00")),
		Alias:     "head",
		Action:    "create",
	}
	if err := history.SaveFile(fp.Name(), &c); err != nil {
		t.Fatal(err)
	}

	rp, err := os.Open(fp.Name())
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{
		"435048d99f77300e5c33dd7ab46dbca2,head,2018-03-10T16:23:57+09:00,create",
		"435048d99f77300e5c33dd7ab46dbca2,head,2018-03-10T16:23:57+09:00,create",
	}

	b, err := ioutil.ReadAll(rp)
	if err != nil {
		t.Fatal(err)
	}
	for i, line := range strings.Split(strings.TrimRight(string(b), "\n"), "\n") {
		if line != expected[i] {
			t.Errorf("line %d, expected=%q, got=%q\n", i, expected[i], line)
		}
	}
}
