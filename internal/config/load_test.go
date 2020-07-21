package config

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindPath(t *testing.T) {
	cases := []struct {
		Start    string
		Expected string
	}{
		{
			"findpath-current",
			"findpath-current/FILENAME",
		},

		{
			"findpath-nested/a/b/c/d",
			"findpath-nested/a/b/FILENAME",
		},
	}

	for _, tt := range cases {
		t.Run(tt.Start, func(t *testing.T) {
			start := filepath.Join("testdata", tt.Start)
			filename := "target.hcl"

			result, err := FindPath(start, filename)
			require.NoError(t, err)

			expected := strings.ReplaceAll(tt.Expected, "FILENAME", filename)
			require.Equal(t, filepath.Join("testdata", expected), result)
		})
	}
}

func TestFindPath_notFound(t *testing.T) {
	filename := "this-should-never-exist-12344321.hcl"

	result, err := FindPath("", filename)
	require.NoError(t, err)
	require.Empty(t, result)
}
