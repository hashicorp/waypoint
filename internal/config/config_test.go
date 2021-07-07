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

		{
			"scoped_settings.hcl",
			"",
			func(t *testing.T, c *Config) {
				require.Equal(t, 3, len(c.Runner.ScopedSettings.Workspaces))
				require.Equal(t, "develop", c.Runner.ScopedSettings.Workspaces[0].Workspace)
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.File, func(t *testing.T) {
			req := require.New(t)

			cfg, err := Load(filepath.Join("testdata", "compare", tt.File), nil)
			if tt.Err != "" {
				req.Error(err)
				req.Contains(err.Error(), tt.Err)
				return
			}
			req.NoError(err)

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
			req := require.New(t)

			_, err := Load(filepath.Join("testdata", "validate", tt.file), nil)

			if tt.err == "" {
				req.NoError(err)
				return
			}

			req.Error(err)
			req.Contains(err.Error(), tt.err)
		})
	}
}
