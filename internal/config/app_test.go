// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package config

import (
	"path/filepath"
	"sort"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
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
			"build_scoped.hcl",
			"test",
			func(t *testing.T, c *App) {
				// Default
				b, err := c.Build(nil)
				require.NoError(t, err)
				require.Equal(t, "A", b.Use.Type)

				// Production workspace
				b, err = c.Build(labelsCtx(map[string]string{
					"waypoint/workspace": "production",
				}))
				require.NoError(t, err)
				require.Equal(t, "B", b.Use.Type)

				// Non-matching workspace
				b, err = c.Build(labelsCtx(map[string]string{
					"waypoint/workspace": "dev",
				}))
				require.NoError(t, err)
				require.Equal(t, "A", b.Use.Type)

				// No workspace
				b, err = c.Build(labelsCtx(map[string]string{
					"whatisthis": "dev",
				}))
				require.NoError(t, err)
				require.Equal(t, "A", b.Use.Type)

				// Labels
				b, err = c.Build(labelsCtx(map[string]string{
					"waypoint/workspace": "staging",
				}))
				require.NoError(t, err)
				require.Equal(t, "C", b.Use.Type)
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
			"build_registry_scoped.hcl",
			"test",
			func(t *testing.T, c *App) {
				// Default
				b, err := c.Registry(nil)
				require.NoError(t, err)
				require.Equal(t, "A", b.Use.Type)

				// Production workspace
				b, err = c.Registry(labelsCtx(map[string]string{
					"waypoint/workspace": "production",
				}))
				require.NoError(t, err)
				require.Equal(t, "B", b.Use.Type)
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

		{
			"config_internal.hcl",
			"test",
			func(t *testing.T, c *App) {
				require := require.New(t)

				vars, err := c.Config.ConfigVars()
				require.NoError(err)

				// test the static value
				require.Len(vars, 2)
				static, ok := vars[0].Value.(*pb.ConfigVar_Static)
				require.True(ok)
				require.Equal("hello", static.Static)

				static, ok = vars[1].Value.(*pb.ConfigVar_Static)
				require.True(ok)
				require.Equal("hello", static.Static)
			},
		},

		{
			"config_internal_dynamic.hcl",
			"test",
			func(t *testing.T, c *App) {
				require := require.New(t)

				vars, err := c.Config.ConfigVars()
				require.NoError(err)

				// test the static value
				require.Len(vars, 4)
				static, ok := vars[0].Value.(*pb.ConfigVar_Static)
				require.True(ok)
				require.Equal("extra: ${config.env.static}", static.Static)

				static, ok = vars[2].Value.(*pb.ConfigVar_Static)
				require.True(ok)
				require.Equal("${config.internal.greeting} ok?", static.Static)
			},
		},

		{
			"config_internal_partial.hcl",
			"test",
			func(t *testing.T, c *App) {
				require := require.New(t)

				vars, err := c.Config.ConfigVars()
				require.NoError(err)

				// test the static value
				require.Len(vars, 3)
				static, ok := vars[2].Value.(*pb.ConfigVar_Static)
				require.True(ok)
				require.Equal(`lower(config.internal.greeting, "FOO")`, static.Static)
			},
		},

		{
			"config_internal_escape.hcl",
			"test",
			func(t *testing.T, c *App) {
				require := require.New(t)

				vars, err := c.Config.ConfigVars()
				require.NoError(err)

				sort.Slice(vars, func(i, j int) bool {
					return vars[i].Name < vars[j].Name
				})

				// test the static value
				require.Len(vars, 3)
				static, ok := vars[0].Value.(*pb.ConfigVar_Static)
				require.True(ok)
				require.Equal(`templatestring("hostname = $${get_hostname()}\n", {"pass" = config.internal.pass})`, static.Static)
				static, ok = vars[2].Value.(*pb.ConfigVar_Static)
				require.True(ok)
				require.Equal("hostname = $${get_hostname()}\n", static.Static)
			},
		},

		{
			"config_reference_loop.hcl",
			"test",
			func(t *testing.T, c *App) {
				require := require.New(t)

				_, err := c.Config.ConfigVars()
				require.Error(err)

				vle := err.(*VariableLoopError)

				require.Equal([]string{"config.env.v1", "config.env.v2", "config.env.v3"}, vle.LoopVars)
			},
		},

		{
			"config_file.hcl",
			"test",
			func(t *testing.T, c *App) {
				require := require.New(t)

				vars, err := c.Config.ConfigVars()
				require.NoError(err)

				// test the static value
				require.Len(vars, 3)

				sort.Slice(vars, func(i, j int) bool {
					return vars[i].Name < vars[j].Name
				})

				require.Equal("blah.yml", vars[0].Name)

				static, ok := vars[0].Value.(*pb.ConfigVar_Static)
				require.True(ok)
				require.Equal("greeting: hello\n", static.Static)

				require.True(vars[0].NameIsPath)

				require.Equal("foo.yml", vars[1].Name)

				static, ok = vars[1].Value.(*pb.ConfigVar_Static)
				require.True(ok)
				require.Equal("foo: hello", static.Static)

				require.True(vars[1].NameIsPath)
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.File, func(t *testing.T) {
			require := require.New(t)

			cfg, err := Load(filepath.Join("testdata", "compare", tt.File), &LoadOptions{
				Workspace: "default",
			})
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

			cfg, err := Load(filepath.Join("testdata", "validate", tt.File), &LoadOptions{
				Workspace: "default",
			})
			require.NoError(err)

			app, err := cfg.App(tt.App, nil)
			require.NoError(err)
			require.NotNil(app)

			_, err = app.Validate()
			if tt.Err == "" {
				require.NoError(err)
				return
			}

			require.Error(err)
			require.Contains(err.Error(), tt.Err)
		})
	}
}

// labelsCtx creates an EvalContext with the labels set.
func labelsCtx(vs map[string]string) *hcl.EvalContext {
	mapValues := map[string]cty.Value{}
	for k, v := range vs {
		mapValues[k] = cty.StringVal(v)
	}

	return &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"labels": cty.MapVal(mapValues),
		},
	}
}

// testPluginBuildConfig implements component.Configurable to test that we
// decode HCL properly.
type testPluginBuildConfig struct {
	config struct {
		Foo string `hcl:"foo,attr"`
		Bar string `hcl:"bar,optional"`
	}
}

func (p *testPluginBuildConfig) Config() (interface{}, error) {
	return &p.config, nil
}
