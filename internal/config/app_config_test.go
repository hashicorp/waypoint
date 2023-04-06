// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func TestConfigVars(t *testing.T) {
	cases := []struct {
		File string
		Err  string
		Func func(*testing.T, *Config)
	}{
		{
			"empty.hcl",
			"",
			func(t *testing.T, c *Config) {
				vars, err := c.Config.ConfigVars()
				require.NoError(t, err)
				require.Empty(t, vars)
			},
		},

		{
			"env_single.hcl",
			"",
			func(t *testing.T, c *Config) {
				require := require.New(t)

				vars, err := c.Config.ConfigVars()
				require.NoError(err)
				require.Len(vars, 1)

				v := vars[0]
				require.Equal("RAILS_ENV", v.Name)
				require.False(v.Internal)
				require.False(v.NameIsPath)
				require.Equal("production", v.Value.(*pb.ConfigVar_Static).Static)

				p, ok := v.Target.AppScope.(*pb.ConfigVar_Target_Project)
				require.True(ok)
				require.Equal("p", p.Project.Project)
			},
		},

		{
			"file_single.hcl",
			"",
			func(t *testing.T, c *Config) {
				require := require.New(t)

				vars, err := c.Config.ConfigVars()
				require.NoError(err)
				require.Len(vars, 1)

				v := vars[0]
				require.Equal("temp.yml", v.Name)
				require.True(v.NameIsPath)
				require.False(v.Internal)
				require.Equal("contents", v.Value.(*pb.ConfigVar_Static).Static)
			},
		},

		{
			"internal.hcl",
			"",
			func(t *testing.T, c *Config) {
				require := require.New(t)

				vars, err := c.Config.ConfigVars()
				require.NoError(err)
				require.Len(vars, 3)

				{
					v := vars[0]
					require.Equal("direct", v.Name)
					require.False(v.NameIsPath)
					require.False(v.Internal)
					require.Equal("V", v.Value.(*pb.ConfigVar_Static).Static)
				}
				{
					v := vars[1]
					require.Equal("interpolated", v.Name)
					require.False(v.NameIsPath)
					require.False(v.Internal)
					require.Equal("value: V", v.Value.(*pb.ConfigVar_Static).Static)
				}
				{
					v := vars[2]
					require.Equal("value", v.Name)
					require.False(v.NameIsPath)
					require.True(v.Internal)
					require.Equal("V", v.Value.(*pb.ConfigVar_Static).Static)
				}
			},
		},

		{
			"workspace.hcl",
			"",
			func(t *testing.T, c *Config) {
				require := require.New(t)

				vars, err := c.Config.ConfigVars()
				require.NoError(err)
				require.Len(vars, 2)

				{
					v := vars[0]
					require.Equal("bar", v.Name)
					require.False(v.NameIsPath)
					require.False(v.Internal)
					require.Equal("baz", v.Value.(*pb.ConfigVar_Static).Static)
					require.Equal("dev", v.Target.Workspace.Workspace)
				}
				{
					v := vars[1]
					require.Equal("foo", v.Name)
					require.False(v.NameIsPath)
					require.False(v.Internal)
					require.Equal("bar", v.Value.(*pb.ConfigVar_Static).Static)
					require.Nil(v.Target.Workspace)
				}
			},
		},

		{
			"app_labels.hcl",
			"",
			func(t *testing.T, c *Config) {
				require := require.New(t)

				// Root level vars should be empty
				{
					vars, err := c.Config.ConfigVars()
					require.NoError(err)
					require.Len(vars, 0)
				}

				// Get our app
				app, err := c.App("api", nil)
				require.NoError(err)

				vars, err := app.Config.ConfigVars()
				require.NoError(err)
				require.Len(vars, 2)

				{
					v := vars[0]
					require.Equal("bar", v.Name)
					require.False(v.NameIsPath)
					require.False(v.Internal)
					require.Equal("baz", v.Value.(*pb.ConfigVar_Static).Static)
					require.Equal("dev", v.Target.Workspace.Workspace)

					s, ok := v.Target.AppScope.(*pb.ConfigVar_Target_Application)
					require.True(ok)
					require.Equal("api", s.Application.Application)
				}
				{
					v := vars[1]
					require.Equal("foo", v.Name)
					require.False(v.NameIsPath)
					require.False(v.Internal)
					require.Equal("bar", v.Value.(*pb.ConfigVar_Static).Static)
					require.Nil(v.Target.Workspace)

					s, ok := v.Target.AppScope.(*pb.ConfigVar_Target_Application)
					require.True(ok)
					require.Equal("api", s.Application.Application)
				}
			},
		},

		{
			"runner.hcl",
			"",
			func(t *testing.T, c *Config) {
				require := require.New(t)

				vars, err := c.Config.ConfigVars()
				require.NoError(err)
				require.Len(vars, 2)

				{
					v := vars[0]
					require.Equal("bar", v.Name)
					require.False(v.NameIsPath)
					require.False(v.Internal)
					require.Equal("baz", v.Value.(*pb.ConfigVar_Static).Static)

					require.NotNil(v.Target.Runner)
					_, ok := v.Target.Runner.Target.(*pb.Ref_Runner_Any)
					require.True(ok)
				}
				{
					v := vars[1]
					require.Equal("foo", v.Name)
					require.False(v.NameIsPath)
					require.False(v.Internal)
					require.Equal("bar", v.Value.(*pb.ConfigVar_Static).Static)
					require.Nil(v.Target.Workspace)
					require.Nil(v.Target.Runner)
				}
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.File, func(t *testing.T) {
			require := require.New(t)

			cfg, err := Load(filepath.Join("testdata", "configvars", tt.File), nil)
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
