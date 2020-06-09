package state

import (
	"testing"

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
	})
}
