package state

import (
	"testing"
	"time"

	"github.com/hashicorp/go-memdb"
	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func TestConfig(t *testing.T) {
	t.Run("basic put and get", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a build
		require.NoError(s.ConfigSet(&pb.ConfigVar{
			Scope: &pb.ConfigVar_Project{
				Project: &pb.Ref_Project{
					Project: "foo",
				},
			},

			Name:  "foo",
			Value: "bar",
		}))

		{
			// Get it exactly
			vs, err := s.ConfigGet(&pb.ConfigGetRequest{
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
			vs, err := s.ConfigGet(&pb.ConfigGetRequest{
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
			vs, err := s.ConfigGet(&pb.ConfigGetRequest{
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

		s := TestState(t)
		defer s.Close()

		// Create a build
		require.NoError(s.ConfigSet(
			&pb.ConfigVar{
				Scope: &pb.ConfigVar_Project{
					Project: &pb.Ref_Project{
						Project: "foo",
					},
				},

				Name:  "global",
				Value: "value",
			},
			&pb.ConfigVar{
				Scope: &pb.ConfigVar_Project{
					Project: &pb.Ref_Project{
						Project: "foo",
					},
				},

				Name:  "hello",
				Value: "project",
			},
			&pb.ConfigVar{
				Scope: &pb.ConfigVar_Application{
					Application: &pb.Ref_Application{
						Project:     "foo",
						Application: "bar",
					},
				},

				Name:  "hello",
				Value: "app",
			},
		))

		{
			// Get our merged variables
			vs, err := s.ConfigGet(&pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Application{
					Application: &pb.Ref_Application{
						Project:     "foo",
						Application: "bar",
					},
				},
			})
			require.NoError(err)
			require.Len(vs, 2)

			// They are sorted, so check on them
			require.Equal("global", vs[0].Name)
			require.Equal("value", vs[0].Value)
			require.Equal("hello", vs[1].Name)
			require.Equal("app", vs[1].Value)
		}

		{
			// Get project scoped variables. This should return everything.
			vs, err := s.ConfigGet(&pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Project{
					Project: &pb.Ref_Project{
						Project: "foo",
					},
				},
			})
			require.NoError(err)
			require.Len(vs, 3)
		}
	})

	t.Run("delete", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a var
		require.NoError(s.ConfigSet(&pb.ConfigVar{
			Scope: &pb.ConfigVar_Project{
				Project: &pb.Ref_Project{
					Project: "foo",
				},
			},

			Name:  "foo",
			Value: "bar",
		}))

		{
			// Get it exactly
			vs, err := s.ConfigGet(&pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Project{
					Project: &pb.Ref_Project{Project: "foo"},
				},

				Prefix: "foo",
			})
			require.NoError(err)
			require.Len(vs, 1)
		}

		// Delete it
		require.NoError(s.ConfigSet(&pb.ConfigVar{
			Scope: &pb.ConfigVar_Project{
				Project: &pb.Ref_Project{
					Project: "foo",
				},
			},

			Name: "foo",
		}))

		// Should not exist
		{
			// Get it exactly
			vs, err := s.ConfigGet(&pb.ConfigGetRequest{
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

		s := TestState(t)
		defer s.Close()

		// Create the config
		require.NoError(s.ConfigSet(&pb.ConfigVar{
			Scope: &pb.ConfigVar_Runner{
				Runner: &pb.Ref_Runner{
					Target: &pb.Ref_Runner_Any{
						Any: &pb.Ref_RunnerAny{},
					},
				},
			},

			Name:  "foo",
			Value: "bar",
		}))

		// Create a var that shouldn't match
		require.NoError(s.ConfigSet(&pb.ConfigVar{
			Scope: &pb.ConfigVar_Project{
				Project: &pb.Ref_Project{
					Project: "foo",
				},
			},

			Name:  "bar",
			Value: "baz",
		}))

		{
			// Get it exactly.
			vs, err := s.ConfigGet(&pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Runner{
					Runner: &pb.Ref_RunnerId{Id: "R_A"},
				},

				Prefix: "foo",
			})
			require.NoError(err)
			require.Len(vs, 1)
		}

		{
			// Get it via a prefix match
			vs, err := s.ConfigGet(&pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Runner{
					Runner: &pb.Ref_RunnerId{Id: "R_A"},
				},

				Prefix: "",
			})
			require.NoError(err)
			require.Len(vs, 1)
		}

		{
			// non-matching prefix
			vs, err := s.ConfigGet(&pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Runner{
					Runner: &pb.Ref_RunnerId{Id: "R_A"},
				},

				Prefix: "bar",
			})
			require.NoError(err)
			require.Empty(vs)
		}
	})

	t.Run("runner configs targeting ID", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create the config
		require.NoError(s.ConfigSet(&pb.ConfigVar{
			Scope: &pb.ConfigVar_Runner{
				Runner: &pb.Ref_Runner{
					Target: &pb.Ref_Runner_Id{
						Id: &pb.Ref_RunnerId{
							Id: "R_A",
						},
					},
				},
			},

			Name:  "foo",
			Value: "bar",
		}))

		// Create a var that shouldn't match
		require.NoError(s.ConfigSet(&pb.ConfigVar{
			Scope: &pb.ConfigVar_Project{
				Project: &pb.Ref_Project{
					Project: "foo",
				},
			},

			Name:  "bar",
			Value: "baz",
		}))

		{
			// Get it exactly.
			vs, err := s.ConfigGet(&pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Runner{
					Runner: &pb.Ref_RunnerId{Id: "R_A"},
				},

				Prefix: "foo",
			})
			require.NoError(err)
			require.Len(vs, 1)
		}

		{
			// Doesn't match
			vs, err := s.ConfigGet(&pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Runner{
					Runner: &pb.Ref_RunnerId{Id: "R_B"},
				},

				Prefix: "foo",
			})
			require.NoError(err)
			require.Len(vs, 0)
		}
	})

	t.Run("runner configs targeting any and ID", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create the config
		require.NoError(s.ConfigSet(&pb.ConfigVar{
			Scope: &pb.ConfigVar_Runner{
				Runner: &pb.Ref_Runner{
					Target: &pb.Ref_Runner_Any{
						Any: &pb.Ref_RunnerAny{},
					},
				},
			},

			Name:  "foo",
			Value: "bar",
		}))

		require.NoError(s.ConfigSet(&pb.ConfigVar{
			Scope: &pb.ConfigVar_Runner{
				Runner: &pb.Ref_Runner{
					Target: &pb.Ref_Runner_Id{
						Id: &pb.Ref_RunnerId{
							Id: "R_A",
						},
					},
				},
			},

			Name:  "foo",
			Value: "baz",
		}))

		{
			// Get it exactly.
			vs, err := s.ConfigGet(&pb.ConfigGetRequest{
				Scope: &pb.ConfigGetRequest_Runner{
					Runner: &pb.Ref_RunnerId{Id: "R_A"},
				},

				Prefix: "foo",
			})
			require.NoError(err)
			require.Len(vs, 1)
			require.Equal("baz", vs[0].Value)
		}
	})
}

func TestConfigWatch(t *testing.T) {
	t.Run("basic put and get", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		ws := memdb.NewWatchSet()

		// Get it with watch
		vs, err := s.ConfigGetWatch(&pb.ConfigGetRequest{
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
		require.NoError(s.ConfigSet(&pb.ConfigVar{
			Scope: &pb.ConfigVar_Project{
				Project: &pb.Ref_Project{
					Project: "foo",
				},
			},

			Name:  "foo",
			Value: "bar",
		}))

		require.False(ws.Watch(time.After(100 * time.Millisecond)))
	})
}
