// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiscover(t *testing.T) {
	cases := []struct {
		Name   string
		Paths  []string
		Plugin *Config
		Err    string
		Result *exec.Cmd
	}{
		{
			"No paths",
			nil,
			&Config{Name: "foo"},
			"",
			nil,
		},

		{
			"Does not exist",
			[]string{
				filepath.Join("testdata", "pathA"),
				filepath.Join("testdata", "pathB"),
			},
			&Config{Name: "foo"},
			"",
			nil,
		},

		{
			"In one path",
			[]string{
				filepath.Join("testdata", "pathA"),
				filepath.Join("testdata", "pathB"),
			},
			&Config{Name: "b"},
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
			&Config{Name: "a"},
			"",
			&exec.Cmd{
				Path: filepath.Join("testdata", "pathA", "waypoint-plugin-a"),
				Args: []string{filepath.Join("testdata", "pathA", "waypoint-plugin-a")},
			},
		},

		{
			"Matching checksum",
			[]string{
				filepath.Join("testdata", "pathA"),
				filepath.Join("testdata", "pathB"),
			},
			&Config{
				Name:     "b",
				Checksum: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			},
			"",
			&exec.Cmd{
				Path: filepath.Join("testdata", "pathB", "waypoint-plugin-b"),
				Args: []string{filepath.Join("testdata", "pathB", "waypoint-plugin-b")},
			},
		},

		{
			"Checksum mismatch",
			[]string{
				filepath.Join("testdata", "pathA"),
				filepath.Join("testdata", "pathB"),
			},
			&Config{
				Name:     "b",
				Checksum: "f3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			},
			"checksum",
			nil,
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
