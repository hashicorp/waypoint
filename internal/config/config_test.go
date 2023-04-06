// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad_compare(t *testing.T) {
	cases := []struct {
		File string
		Err  string
		Func func(*testing.T, *Config)
	}{
		{
			"project.hcl",
			"",
			func(t *testing.T, c *Config) {
				require.Equal(t, "hello", c.Project)
			},
		},

		{
			"project_pwd.hcl",
			"",
			func(t *testing.T, c *Config) {
				require.NotEmpty(t, c.Project)
			},
		},

		{
			"project_path_project.hcl",
			"",
			func(t *testing.T, c *Config) {
				expected, err := filepath.Abs(filepath.Join("testdata", "compare"))
				require.NoError(t, err)
				require.Equal(t, expected, c.Project)
			},
		},

		{
			"project_function.hcl",
			"",
			func(t *testing.T, c *Config) {
				require.Equal(t, "HELLO", c.Project)
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.File, func(t *testing.T) {
			require := require.New(t)

			cfg, err := Load(filepath.Join("testdata", "compare", tt.File), nil)
			if tt.Err != "" {
				require.Error(err)
				require.Contains(err.Error(), tt.Err)
				return
			}
			require.NoError(err)

			tt.Func(t, cfg)
		})
	}
}

func TestConfig_variableDecode(t *testing.T) {
	cases := []struct {
		file string
		err  string
	}{
		{
			"valid.hcl",
			"",
		},
		{
			"invalid_type.hcl",
			"Invalid type specification",
		},

		{
			"invalid_def.hcl",
			"Invalid default value for variable",
		},
		{
			"duplicate_def.hcl",
			"Duplicate variable",
		},
	}

	for _, tt := range cases {
		t.Run(tt.file, func(t *testing.T) {
			require := require.New(t)

			_, err := Load(filepath.Join("testdata", "validate", tt.file), nil)

			if tt.err == "" {
				require.NoError(err)
				return
			}

			require.Error(err)
			require.Contains(err.Error(), tt.err)
		})
	}
}
