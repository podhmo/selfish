package internal_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/podhmo/selfish/internal"
)

func TestNewGist(t *testing.T) {
	tests := []struct {
		name      string
		filenames []string
		public    bool
	}{
		{
			name:      "single valid file",
			filenames: []string{"testdata/file1.txt"},
			public:    true,
		},
		{
			name:      "multiple valid files",
			filenames: []string{"testdata/file1.txt", "testdata/file2.txt"},
			public:    true,
		},
		{
			name:      "file not found, ignored",
			filenames: []string{"testdata/nonexistent.txt"},
			public:    true,
		},
		{
			name:      "empty file, ignored",
			filenames: []string{"testdata/empty.txt"},
			public:    true,
		},
		{
			name:      "readme with title",
			filenames: []string{"testdata/readme.md"},
			public:    true,
		},
		{
			name:      "readme with title, case insensitive",
			filenames: []string{"testdata/README.md"},
			public:    true,
		},
		{
			name:      "secret gist",
			filenames: []string{"testdata/file1.txt"},
			public:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gist, err := internal.NewGist(tt.filenames, tt.public)
			if err != nil {
				t.Errorf("NewGist() unexpected error = %v", err)
				return
			}

			if gist.Public == nil {
				t.Errorf("NewGist() public is nil")
				return
			}
			if want, got := tt.public, *gist.Public; want != got {
				t.Errorf("NewGist() public = %v, want %v", got, want)
				return
			}

			if slices.ContainsFunc(tt.filenames, func(s string) bool { return strings.ToLower(s) == "testdata/readme.md" }) {
				if gist.Description == nil {
					t.Errorf("NewGist() description is nil")
					return
				}
				if want, got := "Readme", *gist.Description; want != got {
					t.Errorf("NewGist() description = %v, want %v", got, want)
					return
				}
			}
		})
	}
}
