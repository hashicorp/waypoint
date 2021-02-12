package state

import (
	"testing"
	"time"

	"github.com/hashicorp/go-memdb"
	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func TestConfigSource(t *testing.T) {
	t.Run("basic put and get", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create
		require.NoError(s.ConfigSourceSet(&pb.ConfigSource{
			Scope: &pb.ConfigSource_Global{
				Global: &pb.Ref_Global{},
			},

			Type:   "vault",
			Config: map[string]string{},
		}))

		{
			// Get it exactly
			vs, err := s.ConfigSourceGet(&pb.GetConfigSourceRequest{
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
			vs, err := s.ConfigSourceGet(&pb.GetConfigSourceRequest{
				Scope: &pb.GetConfigSourceRequest_Global{
					Global: &pb.Ref_Global{},
				},
			})
			require.NoError(err)
			require.Len(vs, 1)
		}

		{
			// non-matching type
			vs, err := s.ConfigSourceGet(&pb.GetConfigSourceRequest{
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

		s := TestState(t)
		defer s.Close()

		// Create
		require.NoError(s.ConfigSourceSet(&pb.ConfigSource{
			Scope: &pb.ConfigSource_Global{
				Global: &pb.Ref_Global{},
			},

			Type:   "vault",
			Config: map[string]string{},
		}))

		{
			// Get it exactly
			vs, err := s.ConfigSourceGet(&pb.GetConfigSourceRequest{
				Scope: &pb.GetConfigSourceRequest_Global{
					Global: &pb.Ref_Global{},
				},

				Type: "vault",
			})
			require.NoError(err)
			require.Len(vs, 1)
		}

		// Create
		require.NoError(s.ConfigSourceSet(&pb.ConfigSource{
			Scope: &pb.ConfigSource_Global{
				Global: &pb.Ref_Global{},
			},

			Type: "vault",

			Delete: true,
		}))

		{
			// Get it exactly
			vs, err := s.ConfigSourceGet(&pb.GetConfigSourceRequest{
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

		s := TestState(t)
		defer s.Close()

		// Create
		require.NoError(s.ConfigSourceSet(&pb.ConfigSource{
			Scope: &pb.ConfigSource_Global{
				Global: &pb.Ref_Global{},
			},

			Type:   "vault",
			Config: map[string]string{},
		}))

		var hash uint64

		// Get it exactly
		{
			vs, err := s.ConfigSourceGet(&pb.GetConfigSourceRequest{
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
		require.NoError(s.ConfigSourceSet(&pb.ConfigSource{
			Scope: &pb.ConfigSource_Global{
				Global: &pb.Ref_Global{},
			},

			Type:   "vault",
			Config: map[string]string{},
		}))

		// Get it exactly
		{
			vs, err := s.ConfigSourceGet(&pb.GetConfigSourceRequest{
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
		require.NoError(s.ConfigSourceSet(&pb.ConfigSource{
			Scope: &pb.ConfigSource_Global{
				Global: &pb.Ref_Global{},
			},

			Type:   "vault",
			Config: map[string]string{"a": "b"},
		}))

		// Get it exactly
		{
			vs, err := s.ConfigSourceGet(&pb.GetConfigSourceRequest{
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
}

func TestConfigSourceWatch(t *testing.T) {
	t.Run("basic put and get", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		ws := memdb.NewWatchSet()

		// Get it with watch
		vs, err := s.ConfigSourceGetWatch(&pb.GetConfigSourceRequest{
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
		require.NoError(s.ConfigSourceSet(&pb.ConfigSource{
			Scope: &pb.ConfigSource_Global{
				Global: &pb.Ref_Global{},
			},

			Type:   "vault",
			Config: map[string]string{},
		}))

		require.False(ws.Watch(time.After(100 * time.Millisecond)))
	})
}
