package state

import (
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

func TestMetadata(t *testing.T) {
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
