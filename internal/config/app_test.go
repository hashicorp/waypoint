package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
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

		{
			"build.hcl",
			"test",
			func(t *testing.T, c *App) {
				b, err := c.Build(nil)
				require.NoError(t, err)
				require.Equal(t, "bar", b.Labels["foo"])

				r, err := c.Registry(nil)
				require.NoError(t, err)
				require.Nil(t, r)
			},
		},

		{
			"build_use.hcl",
			"test",
			func(t *testing.T, c *App) {
				b, err := c.Build(nil)
				require.NoError(t, err)

				op := b.Operation()
				require.NotNil(t, op)

				var p testPluginBuildConfig
				diag := op.Configure(&p, nil)
				if diag.HasErrors() {
					t.Fatal(diag.Error())
				}

				require.NotEmpty(t, p.config.Foo)
			},
		},

		{
			"build_registry.hcl",
			"test",
			func(t *testing.T, c *App) {
				r, err := c.Registry(nil)
				require.NoError(t, err)
				require.NotNil(t, r)
				require.Equal(t, "docker", r.Use.Type)
			},
		},

		{
			"config_env.hcl",
			"test",
			func(t *testing.T, c *App) {
				require := require.New(t)

				vars, err := c.Config.ConfigVars()
				require.NoError(err)

				// test the static value
				require.Len(vars, 1)
				static, ok := vars[0].Value.(*pb.ConfigVar_Static)
				require.True(ok)
				require.Equal("hello", static.Static)
			},
		},

		{
			"config_env_dynamic.hcl",
			"test",
			func(t *testing.T, c *App) {
				require := require.New(t)

				vars, err := c.Config.ConfigVars()
				require.NoError(err)

				// test the static value
				require.Len(vars, 1)
				val, ok := vars[0].Value.(*pb.ConfigVar_Dynamic)
				require.True(ok)
				require.Equal("DATABASE_URL", vars[0].Name)
				require.Equal("vault", val.Dynamic.From)
				require.Equal(map[string]string{
					"path": "foo/",
				}, val.Dynamic.Config)
			},
		},

		{
			"config_env_merge.hcl",
			"test",
			func(t *testing.T, c *App) {
				require := require.New(t)

				vars, err := c.ConfigVars()
				require.NoError(err)
				require.Len(vars, 2)
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

func TestAppValidate(t *testing.T) {
	cases := []struct {
		File string
		App  string
		Err  string
	}{
		{
			"app.hcl",
			"foo",
			"",
		},

		{
			"app.hcl",
			"relative_above_root",
			"must be a child",
		},

		{
			"app.hcl",
			"system_label",
			"reserved for system",
		},
	}

	for _, tt := range cases {
		t.Run(tt.File, func(t *testing.T) {
			require := require.New(t)

			cfg, err := Load(filepath.Join("testdata", "validate", tt.File), "")
			require.NoError(err)

			app, err := cfg.App(tt.App, nil)
			require.NoError(err)
			require.NotNil(app)

			err = app.Validate()
			if tt.Err == "" {
				require.NoError(err)
				return
			}

			require.Error(err)
			require.Contains(err.Error(), tt.Err)
		})
	}
}

// testPluginBuildConfig implements component.Configurable to test that we
// decode HCL properly.
type testPluginBuildConfig struct {
	config struct {
		Foo string `hcl:"foo,attr"`
	}
}

func (p *testPluginBuildConfig) Config() (interface{}, error) {
	return &p.config, nil
}
