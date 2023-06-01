package statetest

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/go-memdb"
	"github.com/stretchr/testify/require"
	empty "google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func init() {
	tests["config"] = []testFunc{TestConfig, TestConfigWatch}
}

func TestConfig(t *testing.T, factory Factory, restartF RestartFactory) {
	// NOTE(mitchellh): A lot of the tests below use the "UnusedScope"
	// field. This is done on purpose because I wanted to retain tests
	// from our old format to ensure that we have backwards compatibility.
	// New functionality uses the new format.
	ctx := context.Background()
	t.Run("basic put and get", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create a config
		require.NoError(s.ConfigSet(ctx, &pb.ConfigVar{
			UnusedScope: &pb.ConfigVar_Project{
				Project: &pb.Ref_Project{
					Project: "foo",
				},
			},

			Name:  "foo",
			Value: &pb.ConfigVar_Static{Static: "bar"},
		}))

		// Create a runner config, we should never get this
		require.NoError(s.ConfigSet(ctx, &pb.ConfigVar{
			Target: &pb.ConfigVar_Target{
				AppScope: &pb.ConfigVar_Target_Project{
					Project: &pb.Ref_Project{
						Project: "foo",
					},
				},

				Runner: &pb.Ref_Runner{
					Target: &pb.Ref_Runner_Any{
						Any: &pb.Ref_RunnerAny{},
					},
				},
			},

			Name:  "bar",
			Value: &pb.ConfigVar_Static{Static: "bar"},
		}))

		{
			// Get it exactly
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Project{
					Project: &pb.Ref_Project{Project: "foo"},
				},

				Prefix: "foo",
			})
			require.NoError(err)
			require.Len(vs, 1)
		}

		{
			// Get it via a prefix match
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Project{
					Project: &pb.Ref_Project{Project: "foo"},
				},

				Prefix: "",
			})
			require.NoError(err)
			require.Len(vs, 1)
		}

		{
			// non-matching prefix
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Project{
					Project: &pb.Ref_Project{Project: "foo"},
				},

				Prefix: "bar",
			})
			require.NoError(err)
			require.Empty(vs)
		}
	})

	t.Run("explicit delete", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create a config
		require.NoError(s.ConfigSet(ctx, &pb.ConfigVar{
			UnusedScope: &pb.ConfigVar_Project{
				Project: &pb.Ref_Project{
					Project: "foo",
				},
			},

			Name:  "foo",
			Value: &pb.ConfigVar_Static{Static: "bar"},
		}))

		require.NoError(s.ConfigSet(ctx, &pb.ConfigVar{
			UnusedScope: &pb.ConfigVar_Project{
				Project: &pb.Ref_Project{
					Project: "foo",
				},
			},

			Name:  "barbar",
			Value: &pb.ConfigVar_Static{Static: "barbar"},
		}))

		// Create a runner config, we should never get this
		require.NoError(s.ConfigSet(ctx, &pb.ConfigVar{
			Target: &pb.ConfigVar_Target{
				AppScope: &pb.ConfigVar_Target_Project{
					Project: &pb.Ref_Project{
						Project: "foo",
					},
				},

				Runner: &pb.Ref_Runner{
					Target: &pb.Ref_Runner_Any{
						Any: &pb.Ref_RunnerAny{},
					},
				},
			},

			Name:  "bar",
			Value: &pb.ConfigVar_Static{Static: "bar"},
		}))

		{
			// Get it exactly
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Project{
					Project: &pb.Ref_Project{Project: "foo"},
				},

				Prefix: "foo",
			})
			require.NoError(err)
			require.Len(vs, 1)
		}

		{
			// delete it
			var vars []*pb.ConfigVar
			vars = append(vars, &pb.ConfigVar{
				Target: &pb.ConfigVar_Target{
					AppScope: &pb.ConfigVar_Target_Project{
						Project: &pb.Ref_Project{
							Project: "foo",
						},
					},
				},

				Name: "foo",
			})
			err := s.ConfigDelete(ctx, vars...)
			require.NoError(err)

			// It's gone
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Project{
					Project: &pb.Ref_Project{Project: "foo"},
				},

				Prefix: "foo",
			})
			require.NoError(err)
			require.Len(vs, 0)
		}
	})

	t.Run("deletes before writes", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create a config
		require.NoError(s.ConfigSet(ctx, &pb.ConfigVar{
			UnusedScope: &pb.ConfigVar_Project{
				Project: &pb.Ref_Project{
					Project: "foo",
				},
			},

			Name:  "foo",
			Value: &pb.ConfigVar_Static{Static: "bar"},
		}, &pb.ConfigVar{
			UnusedScope: &pb.ConfigVar_Project{
				Project: &pb.Ref_Project{
					Project: "foo",
				},
			},

			Name: "foo",
			Value: &pb.ConfigVar_Unset{
				Unset: &empty.Empty{},
			},
		}))

		{
			// Get it exactly
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Project{
					Project: &pb.Ref_Project{Project: "foo"},
				},

				Prefix: "foo",
			})
			require.NoError(err)
			require.Len(vs, 1)

			require.Equal("bar", vs[0].Value.(*pb.ConfigVar_Static).Static)
		}

		{
			// Get it via a prefix match
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Project{
					Project: &pb.Ref_Project{Project: "foo"},
				},

				Prefix: "",
			})
			require.NoError(err)
			require.Len(vs, 1)
		}

		{
			// non-matching prefix
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Project{
					Project: &pb.Ref_Project{Project: "foo"},
				},

				Prefix: "bar",
			})
			require.NoError(err)
			require.Empty(vs)
		}
	})

	t.Run("merging", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create vars
		require.NoError(s.ConfigSet(
			ctx,
			&pb.ConfigVar{
				Target: &pb.ConfigVar_Target{
					AppScope: &pb.ConfigVar_Target_Global{
						Global: &pb.Ref_Global{},
					},
				},

				Name:  "global",
				Value: &pb.ConfigVar_Static{Static: "value"},
			},
			&pb.ConfigVar{
				UnusedScope: &pb.ConfigVar_Project{
					Project: &pb.Ref_Project{
						Project: "foo",
					},
				},

				Name:  "project",
				Value: &pb.ConfigVar_Static{Static: "value"},
			},
			&pb.ConfigVar{
				UnusedScope: &pb.ConfigVar_Project{
					Project: &pb.Ref_Project{
						Project: "foo",
					},
				},

				Name:  "hello",
				Value: &pb.ConfigVar_Static{Static: "project"},
			},
			&pb.ConfigVar{
				UnusedScope: &pb.ConfigVar_Application{
					Application: &pb.Ref_Application{
						Project:     "foo",
						Application: "bar",
					},
				},

				Name:  "hello",
				Value: &pb.ConfigVar_Static{Static: "app"},
			},
		))

		{
			// Get our merged variables
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Application{
					Application: &pb.Ref_Application{
						Project:     "foo",
						Application: "bar",
					},
				},
			})
			require.NoError(err)
			require.Len(vs, 3)

			// They are sorted, so check on them
			require.Equal("global", vs[0].Name)
			require.Equal("value", vs[0].Value.(*pb.ConfigVar_Static).Static)
			require.Equal("hello", vs[1].Name)
			require.Equal("app", vs[1].Value.(*pb.ConfigVar_Static).Static)
			require.Equal("project", vs[2].Name)
			require.Equal("value", vs[2].Value.(*pb.ConfigVar_Static).Static)
		}

		{
			// Get project scoped variables. This should return everything.
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Project{
					Project: &pb.Ref_Project{
						Project: "foo",
					},
				},
			})
			require.NoError(err)
			require.Len(vs, 4)
		}
	})

	t.Run("delete", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create a var
		require.NoError(s.ConfigSet(ctx, &pb.ConfigVar{
			UnusedScope: &pb.ConfigVar_Project{
				Project: &pb.Ref_Project{
					Project: "foo",
				},
			},

			Name:  "foo",
			Value: &pb.ConfigVar_Static{Static: "bar"},
		}))

		{
			// Get it exactly
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Project{
					Project: &pb.Ref_Project{Project: "foo"},
				},

				Prefix: "foo",
			})
			require.NoError(err)
			require.Len(vs, 1)
		}

		// Delete it
		require.NoError(s.ConfigSet(ctx, &pb.ConfigVar{
			UnusedScope: &pb.ConfigVar_Project{
				Project: &pb.Ref_Project{
					Project: "foo",
				},
			},

			Name: "foo",
		}))

		// Should not exist
		{
			// Get it exactly
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Project{
					Project: &pb.Ref_Project{Project: "foo"},
				},

				Prefix: "foo",
			})
			require.NoError(err)
			require.Len(vs, 0)
		}
	})

	t.Run("delete with unset", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create a var
		require.NoError(s.ConfigSet(ctx, &pb.ConfigVar{
			UnusedScope: &pb.ConfigVar_Project{
				Project: &pb.Ref_Project{
					Project: "foo",
				},
			},

			Name:  "foo",
			Value: &pb.ConfigVar_Static{Static: "bar"},
		}))

		{
			// Get it exactly
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Project{
					Project: &pb.Ref_Project{Project: "foo"},
				},

				Prefix: "foo",
			})
			require.NoError(err)
			require.Len(vs, 1)
		}

		// Delete it
		require.NoError(s.ConfigSet(ctx, &pb.ConfigVar{
			UnusedScope: &pb.ConfigVar_Project{
				Project: &pb.Ref_Project{
					Project: "foo",
				},
			},

			Name: "foo",
			Value: &pb.ConfigVar_Unset{
				Unset: &empty.Empty{},
			},
		}))

		// Should not exist
		{
			// Get it exactly
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Project{
					Project: &pb.Ref_Project{Project: "foo"},
				},

				Prefix: "foo",
			})
			require.NoError(err)
			require.Len(vs, 0)
		}
	})

	t.Run("delete with empty static value", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create a var
		require.NoError(s.ConfigSet(ctx, &pb.ConfigVar{
			UnusedScope: &pb.ConfigVar_Project{
				Project: &pb.Ref_Project{
					Project: "foo",
				},
			},

			Name:  "foo",
			Value: &pb.ConfigVar_Static{Static: "bar"},
		}))

		{
			// Get it exactly
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Project{
					Project: &pb.Ref_Project{Project: "foo"},
				},

				Prefix: "foo",
			})
			require.NoError(err)
			require.Len(vs, 1)
		}

		// Delete it
		require.NoError(s.ConfigSet(ctx, &pb.ConfigVar{
			UnusedScope: &pb.ConfigVar_Project{
				Project: &pb.Ref_Project{
					Project: "foo",
				},
			},

			Name:  "foo",
			Value: &pb.ConfigVar_Static{Static: ""},
		}))

		// Should not exist
		{
			// Get it exactly
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Project{
					Project: &pb.Ref_Project{Project: "foo"},
				},

				Prefix: "foo",
			})
			require.NoError(err)
			require.Len(vs, 0)
		}
	})

	t.Run("runner configs any", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create the config
		require.NoError(s.ConfigSet(ctx, &pb.ConfigVar{
			UnusedScope: &pb.ConfigVar_Runner{
				Runner: &pb.Ref_Runner{
					Target: &pb.Ref_Runner_Any{
						Any: &pb.Ref_RunnerAny{},
					},
				},
			},

			Name:  "foo",
			Value: &pb.ConfigVar_Static{Static: "bar"},
		}))

		// Create a var that shouldn't match
		require.NoError(s.ConfigSet(ctx, &pb.ConfigVar{
			UnusedScope: &pb.ConfigVar_Project{
				Project: &pb.Ref_Project{
					Project: "foo",
				},
			},

			Name:  "bar",
			Value: &pb.ConfigVar_Static{Static: "baz"},
		}))

		{
			// Get it exactly.
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Runner: &pb.Ref_RunnerId{Id: "R_A"},

				Prefix: "foo",
			})
			require.NoError(err)
			require.Len(vs, 1)
		}

		{
			// Get it via a prefix match
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Runner: &pb.Ref_RunnerId{Id: "R_A"},

				Prefix: "",
			})
			require.NoError(err)
			require.Len(vs, 1)
		}

		{
			// non-matching prefix
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Runner: &pb.Ref_RunnerId{Id: "R_A"},

				Prefix: "bar",
			})
			require.NoError(err)
			require.Empty(vs)
		}
	})

	t.Run("runner configs targeting ID", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create the config
		require.NoError(s.ConfigSet(ctx, &pb.ConfigVar{
			UnusedScope: &pb.ConfigVar_Runner{
				Runner: &pb.Ref_Runner{
					Target: &pb.Ref_Runner_Id{
						Id: &pb.Ref_RunnerId{
							Id: "R_A",
						},
					},
				},
			},

			Name:  "foo",
			Value: &pb.ConfigVar_Static{Static: "bar"},
		}))

		// Create a var that shouldn't match
		require.NoError(s.ConfigSet(ctx, &pb.ConfigVar{
			UnusedScope: &pb.ConfigVar_Project{
				Project: &pb.Ref_Project{
					Project: "foo",
				},
			},

			Name:  "bar",
			Value: &pb.ConfigVar_Static{Static: "baz"},
		}))

		{
			// Get it exactly.
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Runner: &pb.Ref_RunnerId{Id: "R_A"},

				Prefix: "foo",
			})
			require.NoError(err)
			require.Len(vs, 1)
		}

		{
			// Doesn't match
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Runner: &pb.Ref_RunnerId{Id: "R_B"},

				Prefix: "foo",
			})
			require.NoError(err)
			require.Len(vs, 0)
		}
	})

	t.Run("runner configs targeting any and ID", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create the config
		require.NoError(s.ConfigSet(ctx, &pb.ConfigVar{
			UnusedScope: &pb.ConfigVar_Runner{
				Runner: &pb.Ref_Runner{
					Target: &pb.Ref_Runner_Any{
						Any: &pb.Ref_RunnerAny{},
					},
				},
			},

			Name:  "foo",
			Value: &pb.ConfigVar_Static{Static: "bar"},
		}))

		require.NoError(s.ConfigSet(ctx, &pb.ConfigVar{
			UnusedScope: &pb.ConfigVar_Runner{
				Runner: &pb.Ref_Runner{
					Target: &pb.Ref_Runner_Id{
						Id: &pb.Ref_RunnerId{
							Id: "R_A",
						},
					},
				},
			},

			Name:  "foo",
			Value: &pb.ConfigVar_Static{Static: "baz"},
		}))

		{
			// Get it exactly.
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Runner: &pb.Ref_RunnerId{Id: "R_A"},

				Prefix: "foo",
			})
			require.NoError(err)
			require.Len(vs, 1)
			require.Equal("baz", vs[0].Value.(*pb.ConfigVar_Static).Static)
		}
	})

	t.Run("runner configs scoped to an app", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create the config
		require.NoError(s.ConfigSet(ctx, &pb.ConfigVar{
			Target: &pb.ConfigVar_Target{
				AppScope: &pb.ConfigVar_Target_Application{
					Application: &pb.Ref_Application{
						Project:     "foo",
						Application: "bar",
					},
				},

				Runner: &pb.Ref_Runner{
					Target: &pb.Ref_Runner_Any{
						Any: &pb.Ref_RunnerAny{},
					},
				},
			},

			Name:  "foo",
			Value: &pb.ConfigVar_Static{Static: "bar"},
		}))

		require.NoError(s.ConfigSet(ctx, &pb.ConfigVar{
			Target: &pb.ConfigVar_Target{
				AppScope: &pb.ConfigVar_Target_Global{
					Global: &pb.Ref_Global{},
				},

				Runner: &pb.Ref_Runner{
					Target: &pb.Ref_Runner_Id{
						Id: &pb.Ref_RunnerId{
							Id: "R_A",
						},
					},
				},
			},

			Name:  "bar",
			Value: &pb.ConfigVar_Static{Static: "baz"},
		}))

		{
			// Get it exactly.
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Application{
					Application: &pb.Ref_Application{
						Project:     "foo",
						Application: "bar",
					},
				},

				Runner: &pb.Ref_RunnerId{Id: "R_A"},
			})
			require.NoError(err)
			require.Len(vs, 2)
			require.Equal("bar", vs[0].Name)
			require.Equal("foo", vs[1].Name)
		}

		{
			// Get it for a global runner scope.
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Runner: &pb.Ref_RunnerId{Id: "R_A"},
			})
			require.NoError(err)
			require.Len(vs, 1)
			require.Equal("bar", vs[0].Name)
		}
	})

	t.Run("workspace matching", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create a config
		require.NoError(s.ConfigSet(ctx, &pb.ConfigVar{
			Target: &pb.ConfigVar_Target{
				AppScope: &pb.ConfigVar_Target_Application{
					Application: &pb.Ref_Application{
						Project:     "foo",
						Application: "bar",
					},
				},

				Workspace: &pb.Ref_Workspace{Workspace: "dev"},
			},

			Name:  "foo",
			Value: &pb.ConfigVar_Static{Static: "bar"},
		}))

		{
			// No workspace set
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Application{
					Application: &pb.Ref_Application{
						Project:     "foo",
						Application: "bar",
					},
				},

				Prefix: "foo",
			})
			require.NoError(err)
			require.Len(vs, 1)
		}

		{
			// Matching workspace
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Application{
					Application: &pb.Ref_Application{
						Project:     "foo",
						Application: "bar",
					},
				},

				Workspace: &pb.Ref_Workspace{Workspace: "dev"},
			})
			require.NoError(err)
			require.Len(vs, 1)
		}

		{
			// Non-Matching workspace
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Application{
					Application: &pb.Ref_Application{
						Project:     "foo",
						Application: "bar",
					},
				},

				Workspace: &pb.Ref_Workspace{Workspace: "devno"},
			})
			require.NoError(err)
			require.Len(vs, 0)
		}
	})

	t.Run("workspace set and not set", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create a config that overrides based on workspace
		require.NoError(s.ConfigSet(ctx, &pb.ConfigVar{
			Target: &pb.ConfigVar_Target{
				AppScope: &pb.ConfigVar_Target_Application{
					Application: &pb.Ref_Application{
						Project:     "foo",
						Application: "bar",
					},
				},
			},

			Name:  "foo",
			Value: &pb.ConfigVar_Static{Static: "one"},
		}))
		require.NoError(s.ConfigSet(ctx, &pb.ConfigVar{
			Target: &pb.ConfigVar_Target{
				AppScope: &pb.ConfigVar_Target_Application{
					Application: &pb.Ref_Application{
						Project:     "foo",
						Application: "bar",
					},
				},

				Workspace: &pb.Ref_Workspace{Workspace: "dev"},
			},

			Name:  "foo",
			Value: &pb.ConfigVar_Static{Static: "two"},
		}))

		{
			// Matching workspace
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Application{
					Application: &pb.Ref_Application{
						Project:     "foo",
						Application: "bar",
					},
				},

				Workspace: &pb.Ref_Workspace{Workspace: "dev"},
			})
			require.NoError(err)
			require.Len(vs, 1)
			require.Equal("two", vs[0].Value.(*pb.ConfigVar_Static).Static)
		}

		{
			// Non-Matching workspace
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Application{
					Application: &pb.Ref_Application{
						Project:     "foo",
						Application: "bar",
					},
				},

				Workspace: &pb.Ref_Workspace{Workspace: "devno"},
			})
			require.NoError(err)
			require.Len(vs, 1)
			require.Equal("one", vs[0].Value.(*pb.ConfigVar_Static).Static)
		}
	})

	t.Run("workspace conflict", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create a config that overrides based on workspace
		require.NoError(s.ConfigSet(ctx, &pb.ConfigVar{
			Target: &pb.ConfigVar_Target{
				AppScope: &pb.ConfigVar_Target_Application{
					Application: &pb.Ref_Application{
						Project:     "foo",
						Application: "bar",
					},
				},

				Workspace: &pb.Ref_Workspace{Workspace: "staging"},
			},

			Name:  "foo",
			Value: &pb.ConfigVar_Static{Static: "one"},
		}))
		require.NoError(s.ConfigSet(ctx, &pb.ConfigVar{
			Target: &pb.ConfigVar_Target{
				AppScope: &pb.ConfigVar_Target_Application{
					Application: &pb.Ref_Application{
						Project:     "foo",
						Application: "bar",
					},
				},

				Workspace: &pb.Ref_Workspace{Workspace: "dev"},
			},

			Name:  "foo",
			Value: &pb.ConfigVar_Static{Static: "two"},
		}))

		{
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Application{
					Application: &pb.Ref_Application{
						Project:     "foo",
						Application: "bar",
					},
				},

				Workspace: &pb.Ref_Workspace{Workspace: "dev"},
			})
			require.NoError(err)
			require.Len(vs, 1)
			require.Equal("two", vs[0].Value.(*pb.ConfigVar_Static).Static)
		}

		{
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Application{
					Application: &pb.Ref_Application{
						Project:     "foo",
						Application: "bar",
					},
				},

				Workspace: &pb.Ref_Workspace{Workspace: "staging"},
			})
			require.NoError(err)
			require.Len(vs, 1)
			require.Equal("one", vs[0].Value.(*pb.ConfigVar_Static).Static)
		}
	})

	t.Run("label matching", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create a build
		require.NoError(s.ConfigSet(ctx, &pb.ConfigVar{
			Target: &pb.ConfigVar_Target{
				AppScope: &pb.ConfigVar_Target_Application{
					Application: &pb.Ref_Application{
						Project:     "foo",
						Application: "bar",
					},
				},

				LabelSelector: "env == dev",
			},

			Name:  "foo",
			Value: &pb.ConfigVar_Static{Static: "bar"},
		}))

		{
			// No labels set
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Application{
					Application: &pb.Ref_Application{
						Project:     "foo",
						Application: "bar",
					},
				},

				Prefix: "foo",
			})
			require.NoError(err)
			require.Len(vs, 0)
		}

		{
			// Matching labels
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Application{
					Application: &pb.Ref_Application{
						Project:     "foo",
						Application: "bar",
					},
				},

				Labels: map[string]string{"env": "dev"},
			})
			require.NoError(err)
			require.Len(vs, 1)
		}

		{
			// Non-Matching workspace
			vs, err := s.ConfigGet(ctx, &pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Application{
					Application: &pb.Ref_Application{
						Project:     "foo",
						Application: "bar",
					},
				},

				Labels: map[string]string{"env": "devno"},
			})
			require.NoError(err)
			require.Len(vs, 0)
		}
	})
}

func TestConfigWatch(t *testing.T, factory Factory, restartF RestartFactory) {
	ctx := context.Background()
	t.Run("watches for new variables", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		ws := memdb.NewWatchSet()

		// Get it with watch
		vs, err := s.ConfigGetWatch(ctx, &pb.ConfigGetRequest{
			Scope: &pb.ConfigGetRequest_Project{
				Project: &pb.Ref_Project{Project: "foo"},
			},

			Prefix: "foo",
		}, ws)
		require.NoError(err)
		require.Len(vs, 0)

		// Watch should block
		require.True(ws.Watch(time.After(10 * time.Millisecond)))

		// Create a config
		require.NoError(s.ConfigSet(ctx, &pb.ConfigVar{
			UnusedScope: &pb.ConfigVar_Project{
				Project: &pb.Ref_Project{
					Project: "foo",
				},
			},

			Name:  "foo",
			Value: &pb.ConfigVar_Static{Static: "bar"},
		}))

		require.False(ws.Watch(time.After(3 * time.Second)))
	})
}
