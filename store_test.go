package selfish

import (
	"bytes"
	"strings"
	"testing"
)

func TestLoadCommit(t *testing.T) {
	r := bytes.NewBufferString(strings.Trim(`
68332035342c0cbd0ed6792e0869882c@head@Fri Aug 19 04:50:21 +0000 2016
435048d99f77300e5c33dd7ab46dbca2@head@Thu Aug 18 12:33:45 +0000 2016
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
