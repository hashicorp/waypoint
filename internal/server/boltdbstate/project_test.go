package boltdbstate

import (
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProject(t *testing.T) {
	t.Run("create and get and delete", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a project with one app
		require.NoError(s.ProjectPut(&pb.Project{
			Name: "test",
			Applications: []*pb.Application{
				{
					Project: &pb.Ref_Project{Project: "test"},
					Name:    "testApp",
				},
			},
		}))

		// Read it back
		projectBeforeDelete, err := s.ProjectGet(&pb.Ref_Project{Project: "test"})
		require.NoError(err)
		require.NotNil(projectBeforeDelete)

		// Create a build
		require.NoError(s.BuildPut(false, &pb.Build{
			Id: "testBuild",
			Application: &pb.Ref_Application{
				Application: "testApp",
				Project:     "test",
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		// Delete the project (this should also delete the build)
		err = s.ProjectDelete(&pb.Ref_Project{Project: "test"})
		require.NoError(err)

		// Attempt to get the project again (expected error)
		_, err = s.ProjectGet(&pb.Ref_Project{Project: "test"})
		require.Error(err)
	})
}
