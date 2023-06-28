// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statetest

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/go-memdb"
	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func init() {
	tests["config_source"] = []testFunc{TestConfigSource, TestConfigSourceWatch}
}

func TestConfigSource(t *testing.T, factory Factory, restartF RestartFactory) {
	ctx := context.Background()
	t.Run("basic put and get", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create
		require.NoError(s.ConfigSourceSet(ctx, &pb.ConfigSource{
			Scope: &pb.ConfigSource_Global{
				Global: &pb.Ref_Global{},
			},

			Type:   "vault",
			Config: map[string]string{},
		}))

		{
			// Get it exactly
			vs, err := s.ConfigSourceGet(ctx, &pb.GetConfigSourceRequest{
				Scope: &pb.GetConfigSourceRequest_Global{
					Global: &pb.Ref_Global{},
				},

				Type: "vault",
			})
			require.NoError(err)
			require.Len(vs, 1)
		}

		{
			// Get it via any type
			vs, err := s.ConfigSourceGet(ctx, &pb.GetConfigSourceRequest{
				Scope: &pb.GetConfigSourceRequest_Global{
					Global: &pb.Ref_Global{},
				},
			})
			require.NoError(err)
			require.Len(vs, 1)
		}

		{
			// non-matching type
			vs, err := s.ConfigSourceGet(ctx, &pb.GetConfigSourceRequest{
				Scope: &pb.GetConfigSourceRequest_Global{
					Global: &pb.Ref_Global{},
				},

				Type: "vault-wrong",
			})
			require.NoError(err)
			require.Empty(vs)
		}
	})

	t.Run("delete", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create
		require.NoError(s.ConfigSourceSet(ctx, &pb.ConfigSource{
			Scope: &pb.ConfigSource_Global{
				Global: &pb.Ref_Global{},
			},

			Type:   "vault",
			Config: map[string]string{},
		}))

		{
			// Get it exactly
			vs, err := s.ConfigSourceGet(ctx, &pb.GetConfigSourceRequest{
				Scope: &pb.GetConfigSourceRequest_Global{
					Global: &pb.Ref_Global{},
				},

				Type: "vault",
			})
			require.NoError(err)
			require.Len(vs, 1)
		}

		// Create
		require.NoError(s.ConfigSourceSet(ctx, &pb.ConfigSource{
			Scope: &pb.ConfigSource_Global{
				Global: &pb.Ref_Global{},
			},

			Type: "vault",

			Delete: true,
		}))

		{
			// Get it exactly
			vs, err := s.ConfigSourceGet(ctx, &pb.GetConfigSourceRequest{
				Scope: &pb.GetConfigSourceRequest_Global{
					Global: &pb.Ref_Global{},
				},

				Type: "vault",
			})
			require.NoError(err)
			require.Len(vs, 0)
		}

		// Create
		require.NoError(s.ConfigSourceSet(ctx, &pb.ConfigSource{
			Scope: &pb.ConfigSource_Global{
				Global: &pb.Ref_Global{},
			},

			Type:   "vault",
			Config: map[string]string{},
		}))

		{
			// Get it exactly, then explicit delete
			vs, err := s.ConfigSourceGet(ctx, &pb.GetConfigSourceRequest{
				Scope: &pb.GetConfigSourceRequest_Global{
					Global: &pb.Ref_Global{},
				},

				Type: "vault",
			})
			require.NoError(err)
			require.Len(vs, 1)

			err = s.ConfigSourceDelete(ctx, vs...)

			// Get it exactly again, should be gone
			vs, err = s.ConfigSourceGet(ctx, &pb.GetConfigSourceRequest{
				Scope: &pb.GetConfigSourceRequest_Global{
					Global: &pb.Ref_Global{},
				},

				Type: "vault",
			})
			require.NoError(err)
			require.Len(vs, 0)
		}
	})

	t.Run("hash", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create
		require.NoError(s.ConfigSourceSet(ctx, &pb.ConfigSource{
			Scope: &pb.ConfigSource_Global{
				Global: &pb.Ref_Global{},
			},

			Type:   "vault",
			Config: map[string]string{},
		}))

		var hash uint64

		// Get it exactly
		{
			vs, err := s.ConfigSourceGet(ctx, &pb.GetConfigSourceRequest{
				Scope: &pb.GetConfigSourceRequest_Global{
					Global: &pb.Ref_Global{},
				},

				Type: "vault",
			})
			require.NoError(err)
			require.Len(vs, 1)

			hash = vs[0].Hash
			require.NotEmpty(hash)
		}

		// Modify without change
		require.NoError(s.ConfigSourceSet(ctx, &pb.ConfigSource{
			Scope: &pb.ConfigSource_Global{
				Global: &pb.Ref_Global{},
			},

			Type:   "vault",
			Config: map[string]string{},
		}))

		// Get it exactly
		{
			vs, err := s.ConfigSourceGet(ctx, &pb.GetConfigSourceRequest{
				Scope: &pb.GetConfigSourceRequest_Global{
					Global: &pb.Ref_Global{},
				},

				Type: "vault",
			})
			require.NoError(err)
			require.Len(vs, 1)

			// Hash should NOT change
			require.Equal(hash, vs[0].Hash)
		}

		// Modify
		require.NoError(s.ConfigSourceSet(ctx, &pb.ConfigSource{
			Scope: &pb.ConfigSource_Global{
				Global: &pb.Ref_Global{},
			},

			Type:   "vault",
			Config: map[string]string{"a": "b"},
		}))

		// Get it exactly
		{
			vs, err := s.ConfigSourceGet(ctx, &pb.GetConfigSourceRequest{
				Scope: &pb.GetConfigSourceRequest_Global{
					Global: &pb.Ref_Global{},
				},

				Type: "vault",
			})
			require.NoError(err)
			require.Len(vs, 1)

			// Hash should change
			require.NotEqual(hash, vs[0].Hash)
		}
	})

	t.Run("put and get workspace scoped global config sources", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		require.NoError(s.WorkspacePut(ctx, &pb.Workspace{
			Name: "dev",
		}))

		require.NoError(s.WorkspacePut(ctx, &pb.Workspace{
			Name: "prod",
		}))

		// Create dev Vault config source
		require.NoError(s.ConfigSourceSet(ctx, &pb.ConfigSource{
			Scope: &pb.ConfigSource_Global{
				Global: &pb.Ref_Global{},
			},

			Type: "vault",
			Config: map[string]string{
				"token": "abc",
			},
			Workspace: &pb.Ref_Workspace{
				Workspace: "dev",
			},
		}))

		// Create prod Vault config source
		require.NoError(s.ConfigSourceSet(ctx, &pb.ConfigSource{
			Scope: &pb.ConfigSource_Global{
				Global: &pb.Ref_Global{},
			},

			Type: "vault",
			Config: map[string]string{
				"token": "123",
			},
			Workspace: &pb.Ref_Workspace{
				Workspace: "prod",
			},
		}))

		// Get the dev Vault config source
		dcs, err := s.ConfigSourceGet(ctx, &pb.GetConfigSourceRequest{
			Scope: &pb.GetConfigSourceRequest_Global{
				Global: &pb.Ref_Global{},
			},
			Workspace: &pb.Ref_Workspace{Workspace: "dev"},
			Type:      "vault",
		})

		// Verify that we got back the expected config
		require.NoError(err)
		require.NotNil(dcs)
		require.Equal(dcs[0].Config["token"], "abc")

		// Get the prod Vault config source
		pcs, err := s.ConfigSourceGet(ctx, &pb.GetConfigSourceRequest{
			Scope: &pb.GetConfigSourceRequest_Global{
				Global: &pb.Ref_Global{},
			},
			Workspace: &pb.Ref_Workspace{Workspace: "prod"},
			Type:      "vault",
		})

		// Verify that we got back the expected config
		require.NoError(err)
		require.NotNil(pcs)
		require.Equal(pcs[0].Config["token"], "123")
	})

	t.Run("put and get workspace scoped project config sources", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		require.NoError(s.ProjectPut(ctx, &pb.Project{
			Name: "testProject",
		}))

		require.NoError(s.WorkspacePut(ctx, &pb.Workspace{
			Name: "dev",
		}))

		require.NoError(s.WorkspacePut(ctx, &pb.Workspace{
			Name: "prod",
		}))

		// Create dev Vault config source for the test project
		require.NoError(s.ConfigSourceSet(ctx, &pb.ConfigSource{
			Scope: &pb.ConfigSource_Project{
				Project: &pb.Ref_Project{Project: "testProject"},
			},

			Type: "vault",
			Config: map[string]string{
				"token": "abc",
			},
			Workspace: &pb.Ref_Workspace{
				Workspace: "dev",
			},
		}))

		// Create prod Vault config source for the test project
		require.NoError(s.ConfigSourceSet(ctx, &pb.ConfigSource{
			Scope: &pb.ConfigSource_Project{
				Project: &pb.Ref_Project{Project: "testProject"},
			},

			Type: "vault",
			Config: map[string]string{
				"token": "123",
			},
			Workspace: &pb.Ref_Workspace{
				Workspace: "prod",
			},
		}))

		// Get the dev Vault config source
		dcs, err := s.ConfigSourceGet(ctx, &pb.GetConfigSourceRequest{
			Scope: &pb.GetConfigSourceRequest_Project{
				Project: &pb.Ref_Project{
					Project: "testProject",
				},
			},
			Workspace: &pb.Ref_Workspace{Workspace: "dev"},
			Type:      "vault",
		})

		// Verify that we got back the expected config
		require.NoError(err)
		require.NotNil(dcs)
		require.Equal(dcs[0].Config["token"], "abc")

		// Get the prod Vault config source
		pcs, err := s.ConfigSourceGet(ctx, &pb.GetConfigSourceRequest{
			Scope: &pb.GetConfigSourceRequest_Project{
				Project: &pb.Ref_Project{
					Project: "testProject",
				},
			},
			Workspace: &pb.Ref_Workspace{Workspace: "prod"},
			Type:      "vault",
		})

		// Verify that we got back the expected config
		require.NoError(err)
		require.NotNil(pcs)
		require.Equal(pcs[0].Config["token"], "123")
	})

	// this test verifies that if there is no workspace-project scoped config
	// source but there is a global one, then we get the global one back when
	// requesting a workspace-project scoped config source
	t.Run("put global config source and get workspace-project scoped config source", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		require.NoError(s.ProjectPut(ctx, &pb.Project{
			Name: "testProject",
		}))

		require.NoError(s.WorkspacePut(ctx, &pb.Workspace{
			Name: "dev",
		}))

		// Create a config source at the global scope, no workspace
		require.NoError(s.ConfigSourceSet(ctx, &pb.ConfigSource{
			Scope: &pb.ConfigSource_Global{
				Global: &pb.Ref_Global{},
			},

			Type:   "vault",
			Config: map[string]string{},
		}))

		// Use a workspace and project scoped config source get request
		// There is no config source with that scope set, but we should still
		// get the global config source back
		resp, err := s.ConfigSourceGet(ctx, &pb.GetConfigSourceRequest{
			Scope: &pb.GetConfigSourceRequest_Project{
				Project: &pb.Ref_Project{
					Project: "testProject",
				},
			},
			Workspace: &pb.Ref_Workspace{
				Workspace: "dev",
			},
			Type: "vault",
		})
		require.NoError(err)
		require.NotNil(resp)
		require.Equal(1, len(resp))
	})

	// this test verifies that if we have a global config source not scoped to
	// any workspace, and we request a config source for a given project in a
	// specific workspace, then we get the more-tightly scoped config source
	t.Run("put global & workspace-project scoped config source and get workspace-project config source", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		require.NoError(s.ProjectPut(ctx, &pb.Project{
			Name: "testProject",
		}))

		require.NoError(s.WorkspacePut(ctx, &pb.Workspace{
			Name: "dev",
		}))

		// Create a config source at the global scope, no workspace
		require.NoError(s.ConfigSourceSet(ctx, &pb.ConfigSource{
			Scope: &pb.ConfigSource_Global{
				Global: &pb.Ref_Global{},
			},

			Type: "vault",
			Config: map[string]string{
				"token": "global",
			},
		}))

		// Create a config source at the project scope in dev workspace
		require.NoError(s.ConfigSourceSet(ctx, &pb.ConfigSource{
			Scope: &pb.ConfigSource_Project{
				Project: &pb.Ref_Project{
					Project: "testProject",
				},
			},
			Workspace: &pb.Ref_Workspace{
				Workspace: "dev",
			},

			Type: "vault",
			Config: map[string]string{
				"token": "dev",
			},
		}))

		// Use a workspace and project scoped config source get request
		// We should get back the workspace-project scoped config source as
		// the first result in our slice of config sources
		resp, err := s.ConfigSourceGet(ctx, &pb.GetConfigSourceRequest{
			Scope: &pb.GetConfigSourceRequest_Project{
				Project: &pb.Ref_Project{
					Project: "testProject",
				},
			},
			Workspace: &pb.Ref_Workspace{
				Workspace: "dev",
			},
			Type: "vault",
		})
		require.NoError(err)
		require.NotNil(resp)
		require.Equal(2, len(resp))
		require.Equal("global", resp[0].Config["token"])
		require.Equal("dev", resp[1].Config["token"])
	})

	// this test verifies that if we have a global config source not scoped to
	// any workspace, and we request a config source for a given app in a
	// specific workspace, then we get the more-tightly scoped config source
	t.Run("put global & workspace-project scoped config source and get workspace-app config source", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		require.NoError(s.ProjectPut(ctx, &pb.Project{
			Name: "testProject",
		}))

		_, err := s.AppPut(ctx, &pb.Application{
			Project: &pb.Ref_Project{Project: "testProject"},
			Name:    "testApp",
		})
		require.NoError(err)

		require.NoError(s.WorkspacePut(ctx, &pb.Workspace{
			Name: "dev",
		}))

		// Create a config source at the global scope, no workspace
		require.NoError(s.ConfigSourceSet(ctx, &pb.ConfigSource{
			Scope: &pb.ConfigSource_Global{
				Global: &pb.Ref_Global{},
			},

			Type: "vault",
			Config: map[string]string{
				"token": "abc",
			},
		}))

		// Create a config source at the app scope in dev workspace
		require.NoError(s.ConfigSourceSet(ctx, &pb.ConfigSource{
			Scope: &pb.ConfigSource_Application{Application: &pb.Ref_Application{
				Application: "testApp",
				Project:     "testProject",
			}},
			Workspace: &pb.Ref_Workspace{
				Workspace: "dev",
			},

			Type: "vault",
			Config: map[string]string{
				"token": "123",
			},
		}))

		// Use a workspace and app scoped config source get request
		// We should get back the workspace-app scoped config source as
		// the first result in our slice of config sources
		resp, err := s.ConfigSourceGet(ctx, &pb.GetConfigSourceRequest{
			Scope: &pb.GetConfigSourceRequest_Application{
				Application: &pb.Ref_Application{
					Application: "testApp",
					Project:     "testProject",
				},
			},
			Workspace: &pb.Ref_Workspace{
				Workspace: "dev",
			},
			Type: "vault",
		})
		require.NoError(err)
		require.NotNil(resp)
		require.Equal(2, len(resp))
		require.Equal("abc", resp[0].Config["token"])
		require.Equal("123", resp[1].Config["token"])
	})

	t.Run("get all config sources", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		projectName := "testProject"
		appName := "testApp"

		require.NoError(s.ProjectPut(ctx, &pb.Project{
			Name: projectName,
		}))

		_, err := s.AppPut(ctx, &pb.Application{
			Project: &pb.Ref_Project{Project: projectName},
			Name:    appName,
		})
		require.NoError(err)

		require.NoError(s.WorkspacePut(ctx, &pb.Workspace{
			Name: "dev",
		}))

		// Create a config source at the global scope, no workspace
		require.NoError(s.ConfigSourceSet(ctx, &pb.ConfigSource{
			Scope: &pb.ConfigSource_Global{
				Global: &pb.Ref_Global{},
			},

			Type: "vault",
			Config: map[string]string{
				"token": "vault-global",
			},
		}))

		// Create a config source at the project scope, no workspace
		require.NoError(s.ConfigSourceSet(ctx, &pb.ConfigSource{
			Scope: &pb.ConfigSource_Project{
				Project: &pb.Ref_Project{
					Project: projectName,
				},
			},

			Type: "terraform-cloud",
			Config: map[string]string{
				"token": "terraform-cloud-project",
			},
		}))

		// Create a config source at the app scope, no workspace
		require.NoError(s.ConfigSourceSet(ctx, &pb.ConfigSource{
			Scope: &pb.ConfigSource_Application{Application: &pb.Ref_Application{
				Application: appName,
				Project:     projectName,
			}},

			Type: "packer",
			Config: map[string]string{
				"token": "packer-app",
			},
		}))

		// Create a config source at the app scope in dev workspace
		require.NoError(s.ConfigSourceSet(ctx, &pb.ConfigSource{
			Scope: &pb.ConfigSource_Application{Application: &pb.Ref_Application{
				Application: appName,
				Project:     projectName,
			}},
			Workspace: &pb.Ref_Workspace{
				Workspace: "dev",
			},

			Type: "consul",
			Config: map[string]string{
				"token": "consul-app-dev-workspace",
			},
		}))

		// get all of the config sources we made
		sources, err := s.ConfigSourceGet(ctx, &pb.GetConfigSourceRequest{
			Scope: &pb.GetConfigSourceRequest_All{
				All: true,
			},
		})

		require.NoError(err)
		require.Equal(4, len(sources))
		require.Equal("vault-global", sources[0].Config["token"])
		require.Equal("terraform-cloud-project", sources[1].Config["token"])
		require.Equal("packer-app", sources[2].Config["token"])
		require.Equal("consul-app-dev-workspace", sources[3].Config["token"])
	})
}

func TestConfigSourceWatch(t *testing.T, factory Factory, restartF RestartFactory) {
	ctx := context.Background()
	t.Run("basic put and get", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		ws := memdb.NewWatchSet()

		// Get it with watch
		vs, err := s.ConfigSourceGetWatch(ctx, &pb.GetConfigSourceRequest{
			Scope: &pb.GetConfigSourceRequest_Global{
				Global: &pb.Ref_Global{},
			},

			Type: "vault",
		}, ws)
		require.NoError(err)
		require.Len(vs, 0)

		// Watch should block
		require.True(ws.Watch(time.After(10 * time.Millisecond)))

		// Create
		require.NoError(s.ConfigSourceSet(ctx, &pb.ConfigSource{
			Scope: &pb.ConfigSource_Global{
				Global: &pb.Ref_Global{},
			},

			Type:   "vault",
			Config: map[string]string{},
		}))

		require.False(ws.Watch(time.After(3 * time.Second)))
	})
}
