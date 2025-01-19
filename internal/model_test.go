package internal_test

import (
	"slices"
	"testing"

	"github.com/podhmo/selfish/internal"
)

func TestNewGist(t *testing.T) {
	tests := []struct {
		name      string
		filenames []string
	}{
		{
			name:      "single valid file",
			filenames: []string{"testdata/file1.txt"},
		},
		{
			name:      "multiple valid files",
			filenames: []string{"testdata/file1.txt", "testdata/file2.txt"},
		},
		{
			name:      "file not found, ignored",
			filenames: []string{"testdata/nonexistent.txt"},
		},
		{
			name:      "empty file, ignored",
			filenames: []string{"testdata/empty.txt"},
		},
		{
			name:      "readme with title",
			filenames: []string{"testdata/Readme.md"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gist, err := internal.NewGist(tt.filenames)
			if err != nil {
				t.Errorf("NewGist() unexpected error = %v", err)
				return
			}

			if slices.Contains(tt.filenames, "Readme.md") {
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
