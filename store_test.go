package selfish

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

func TestSaveCommit(t *testing.T) {
	r := bytes.NewBufferString(strings.Trim(`
435048d99f77300e5c33dd7ab46dbca2@head@Thu Aug 18 12:33:45 +0000 2016@create
`, "\n"))
	now, err := time.Parse(time.RubyDate, "Fri Aug 19 04:50:21 +0000 2016")
	if err != nil {
		t.Errorf("time parse error: %s", err)
	}
	commit := Commit{
		ID:        "xxxx",
		CreatedAt: now,
		Alias:     "newItem",
		Action:    "create",
	}
	w := &bytes.Buffer{}

	err = saveCommit(w, r, commit)
	if err != nil {
		t.Errorf("saveCommit failured: %s", err)
	}

	text, err := ioutil.ReadAll(w)
	if err != nil {
		t.Errorf("invalid contents: %s", err)
	}
	result := strings.Split(string(text), "\n")[0]
	expected := "xxxx@newItem@Fri Aug 19 04:50:21 +0000 2016@create"
	if result != expected {
		t.Errorf("commit line is must be %q but %q", expected, result)
	}
}

func TestLoadCommit(t *testing.T) {
	r := bytes.NewBufferString(strings.Trim(`
68332035342c0cbd0ed6792e0869882c@head@Fri Aug 19 04:50:21 +0000 2016@update
435048d99f77300e5c33dd7ab46dbca2@head@Thu Aug 18 12:33:45 +0000 2016@create
`, "\n"))

	cases := []struct {
		alias      string
		expectedID string
		found      bool
	}{
		{alias: "head", expectedID: "68332035342c0cbd0ed6792e0869882c", found: true},
		{alias: "hmm", found: false},
	}

	for _, c := range cases {
		commit, err := loadCommit(r, c.alias)
		if err != nil {
			t.Error(err)
		}
		if c.found {
			if commit == nil {
				t.Errorf("%q should be found. but not found", c.alias)
			}
			if commit.ID != c.expectedID {
				t.Errorf("expected id is %q but found id is %q", c.expectedID, commit.ID)
			}
		}
	}
}
