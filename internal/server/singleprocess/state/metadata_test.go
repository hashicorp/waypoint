package state

import (
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

func TestMetadata(t *testing.T) {
	t.Run("sets metadata on the project", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		name := "abcde"
		// Set
		err := s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: name,
		}))
		require.NoError(err)

		err = s.MetadataSet(&pb.MetadataSetRequest{
			Scope: &pb.MetadataSetRequest_Project{
				Project: &pb.Ref_Project{Project: name},
			},
			Value: &pb.MetadataSetRequest_FileChangeSignal{
				FileChangeSignal: "HUP",
			},
		})

		require.NoError(err)

		proj, err := s.ProjectGet(&pb.Ref_Project{
			Project: name,
		})
		require.NoError(err)

		require.Equal("HUP", proj.FileChangeSignal)
	})

	t.Run("sets metadata on an application", func(t *testing.T) {
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

		err = s.MetadataSet(&pb.MetadataSetRequest{
			Scope: &pb.MetadataSetRequest_Application{
				Application: &pb.Ref_Application{
					Project:     name,
					Application: "app",
				},
			},
			Value: &pb.MetadataSetRequest_FileChangeSignal{
				FileChangeSignal: "HUP",
			},
		})

		require.NoError(err)

		proj, err := s.ProjectGet(&pb.Ref_Project{
			Project: name,
		})
		require.NoError(err)

		require.Equal("", proj.FileChangeSignal)

		app, err := s.AppGet(&pb.Ref_Application{
			Project:     name,
			Application: "app",
		})
		require.NoError(err)

		require.Equal("HUP", app.FileChangeSignal)
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

		err = s.MetadataSet(&pb.MetadataSetRequest{
			Scope: &pb.MetadataSetRequest_Project{
				Project: &pb.Ref_Project{Project: name},
			},
			Value: &pb.MetadataSetRequest_FileChangeSignal{
				FileChangeSignal: "HUP",
			},
		})

		require.NoError(err)

		sig, err := s.MetadataGetFileChangeSignal(&pb.Ref_Application{
			Project:     name,
			Application: "app",
		})
		require.NoError(err)

		require.Equal("HUP", sig)

		err = s.MetadataSet(&pb.MetadataSetRequest{
			Scope: &pb.MetadataSetRequest_Application{
				Application: &pb.Ref_Application{
					Project:     name,
					Application: "app",
				},
			},
			Value: &pb.MetadataSetRequest_FileChangeSignal{
				FileChangeSignal: "TERM",
			},
		})

		require.NoError(err)

		sig, err = s.MetadataGetFileChangeSignal(&pb.Ref_Application{
			Project:     name,
			Application: "app",
		})
		require.NoError(err)

		require.Equal("TERM", sig)
	})
}
