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

	t.Run("put and get workspace scope", func(t *testing.T) {
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
