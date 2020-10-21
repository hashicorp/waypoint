package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigApp_compare(t *testing.T) {
	cases := []struct {
		File string
		App  string
		Func func(*testing.T, *App)
	}{
		{
			"app.hcl",
			"dontexist",
			func(t *testing.T, c *App) {
				require.Nil(t, c)
			},
		},

		{
			"app.hcl",
			"foo",
			func(t *testing.T, c *App) {
				expected, err := filepath.Abs(filepath.Join("testdata", "compare"))
				require.NoError(t, err)

				require := require.New(t)
				require.Equal("foo", c.Name)
				require.Equal(expected, c.Path)
			},
		},

		{
			"app_path_relative.hcl",
			"foo",
			func(t *testing.T, c *App) {
				expected, err := filepath.Abs(filepath.Join("testdata", "compare", "bar"))
				require.NoError(t, err)

				require := require.New(t)
				require.Equal("foo", c.Name)
				require.Equal(expected, c.Path)
			},
		},

		{
			"app_labels.hcl",
			"bar",
			func(t *testing.T, c *App) {
				expected, err := filepath.Abs(filepath.Join("testdata", "compare"))
				require.NoError(t, err)

				require := require.New(t)
				require.Equal("bar", c.Name)
				require.NotEmpty(c.Labels["pwd"])
				require.Equal(expected, c.Labels["project"])
				require.Equal(filepath.Join(expected, "bar"), c.Labels["app"])
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.File, func(t *testing.T) {
			require := require.New(t)

			cfg, err := Load(filepath.Join("testdata", "compare", tt.File), "")
			require.NoError(err)

			app, err := cfg.App(tt.App, nil)
			require.NoError(err)

			tt.Func(t, app)
		})
	}
}
