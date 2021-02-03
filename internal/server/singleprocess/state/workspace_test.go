package state

import (
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

func TestWorkspace(t *testing.T) {
	t.Run("List is empty by default", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		result, err := s.WorkspaceList()
		require.NoError(err)
		require.Empty(result)
	})

	t.Run("List non-empty", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a build
		require.NoError(s.BuildPut(false, serverptypes.TestValidBuild(t, &pb.Build{
			Id: "1",
		})))
		require.NoError(s.BuildPut(false, serverptypes.TestValidBuild(t, &pb.Build{
			Id: "2",
		})))
		require.NoError(s.BuildPut(false, serverptypes.TestValidBuild(t, &pb.Build{
			Id: "3",
			Application: &pb.Ref_Application{
				Application: "B",
				Project:     "B",
			},
		})))

		// Create some other resources
		require.NoError(s.DeploymentPut(false, serverptypes.TestValidDeployment(t, &pb.Deployment{
			Id: "1",
		})))

		// Workspace list should only list one
		{
			result, err := s.WorkspaceList()
			require.NoError(err)
			require.Len(result, 1)

			ws := result[0]
			require.Len(ws.Projects, 2)
			require.Len(ws.Projects[0].Applications, 1)
			require.Len(ws.Projects[1].Applications, 1)
		}

		// Create a new workspace
		require.NoError(s.BuildPut(false, serverptypes.TestValidBuild(t, &pb.Build{
			Id: "4",
			Workspace: &pb.Ref_Workspace{
				Workspace: "2",
			},
		})))
		{
			result, err := s.WorkspaceList()
			require.NoError(err)
			require.Len(result, 2)
		}
	})
}
