package plugin

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/config"
)

func TestDiscover(t *testing.T) {
	cases := []struct {
		Name   string
		Paths  []string
		Plugin *config.Plugin
		Err    string
		Result *exec.Cmd
	}{
		{
			"No paths",
			nil,
			&config.Plugin{Name: "foo"},
			"",
			nil,
		},

		{
			"Does not exist",
			[]string{
				filepath.Join("testdata", "pathA"),
				filepath.Join("testdata", "pathB"),
			},
			&config.Plugin{Name: "foo"},
			"",
			nil,
		},

		{
			"In one path",
			[]string{
				filepath.Join("testdata", "pathA"),
				filepath.Join("testdata", "pathB"),
			},
			&config.Plugin{Name: "b"},
			"",
			&exec.Cmd{
				Path: filepath.Join("testdata", "pathB", "waypoint-plugin-b"),
				Args: []string{filepath.Join("testdata", "pathB", "waypoint-plugin-b")},
			},
		},

		{
			"In two paths",
			[]string{
				filepath.Join("testdata", "pathA"),
				filepath.Join("testdata", "pathB"),
			},
			&config.Plugin{Name: "a"},
			"",
			&exec.Cmd{
				Path: filepath.Join("testdata", "pathA", "waypoint-plugin-a"),
				Args: []string{filepath.Join("testdata", "pathA", "waypoint-plugin-a")},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)

			result, err := Discover(tt.Plugin, tt.Paths)
			if tt.Err != "" {
				require.Error(err)
				require.Contains(err.Error(), tt.Err)
				return
			}
			require.NoError(err)

			if tt.Result == nil {
				require.Nil(result)
				return
			}

			// Things to compare
			result = &exec.Cmd{
				Path: result.Path,
				Args: result.Args,
				Env:  result.Env,
				Dir:  result.Dir,
			}

			require.Equal(tt.Result, result)
		})
	}
}
