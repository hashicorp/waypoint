package state

import (
	"testing"

	"github.com/stretchr/testify/require"
	/*
		"google.golang.org/grpc/codes"
		"google.golang.org/grpc/status"
	*/

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

func TestApplication(t *testing.T) {
	t.Run("Put adds a new application", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Write project
		ref := &pb.Ref_Project{Project: "foo"}
		require.NoError(s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: ref.Project,
		})))

		// Has no apps
		{
			resp, err := s.ProjectGet(ref)
			require.NoError(err)
			require.NotNil(resp)
			require.Empty(resp.Applications)
		}

		// Add
		app, err := s.AppPut(serverptypes.TestApplication(t, &pb.Application{
			Project: ref,
		}))
		require.NoError(err)

		// Can read
		{
			resp, err := s.AppGet(&pb.Ref_Application{
				Project:     ref.Project,
				Application: app.Name,
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Has apps
		{
			resp, err := s.ProjectGet(ref)
			require.NoError(err)
			require.NotNil(resp)
			require.Len(resp.Applications, 1)
		}
	})

	t.Run("Put non-existent project", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Write project
		ref := &pb.Ref_Project{Project: "foo"}

		// Add
		app, err := s.AppPut(serverptypes.TestApplication(t, &pb.Application{
			Project: ref,
		}))
		require.NoError(err)

		// Can read
		{
			resp, err := s.AppGet(&pb.Ref_Application{
				Project:     ref.Project,
				Application: app.Name,
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Has project
		{
			resp, err := s.ProjectGet(ref)
			require.NoError(err)
			require.NotNil(resp)
			require.Len(resp.Applications, 1)
		}
	})

	t.Run("Put appends to existing list of applications", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Write project
		ref := &pb.Ref_Project{Project: "foo"}
		require.NoError(s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: ref.Project,
			Applications: []*pb.Application{
				serverptypes.TestApplication(t, nil),
			},
		})))

		// Add
		_, err := s.AppPut(serverptypes.TestApplication(t, &pb.Application{
			Project: ref,
			Name:    "next",
		}))
		require.NoError(err)

		// Has apps
		{
			resp, err := s.ProjectGet(ref)
			require.NoError(err)
			require.NotNil(resp)
			require.Len(resp.Applications, 2)
		}
	})

	t.Run("Put updates an existing application", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Write project
		ref := &pb.Ref_Project{Project: "foo"}
		require.NoError(s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: ref.Project,
			Applications: []*pb.Application{
				serverptypes.TestApplication(t, &pb.Application{
					Name: "foo",
				}),
			},
		})))

		// Add
		_, err := s.AppPut(serverptypes.TestApplication(t, &pb.Application{
			Project: ref,
			Name:    "foo",
		}))
		require.NoError(err)

		// Has apps
		{
			resp, err := s.ProjectGet(ref)
			require.NoError(err)
			require.NotNil(resp)
			require.Len(resp.Applications, 1)
		}
	})

	t.Run("reads file change signal upward", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		name := "abcde"
		// Set
		err := s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: name,
		}))
		require.NoError(err)

		_, err = s.AppPut(&pb.Application{
			Project: &pb.Ref_Project{Project: name},
			Name:    "app",
		})
		require.NoError(err)

		err = s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name:             name,
			FileChangeSignal: "HUP",
		}))
		require.NoError(err)

		sig, err := s.GetFileChangeSignal(&pb.Ref_Application{
			Project:     name,
			Application: "app",
		})
		require.NoError(err)

		require.Equal("HUP", sig)

		_, err = s.AppPut(&pb.Application{
			Project:          &pb.Ref_Project{Project: name},
			Name:             "app",
			FileChangeSignal: "TERM",
		})
		require.NoError(err)

		sig, err = s.GetFileChangeSignal(&pb.Ref_Application{
			Project:     name,
			Application: "app",
		})
		require.NoError(err)

		require.Equal("TERM", sig)
	})
}
