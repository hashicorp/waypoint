package statetest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func init() {
	tests["workspace"] = []testFunc{
		TestWorkspace,
		TestWorkspaceProject,
		TestWorkspaceApp,
		TestWorkspacePut,
	}
}

func TestWorkspace(t *testing.T, factory Factory, restartF RestartFactory) {
	ctx := context.Background()
	t.Run("List is empty by default", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		result, err := s.WorkspaceList(ctx)
		require.NoError(err)
		require.Empty(result)
	})

	t.Run("List non-empty", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create a build
		require.NoError(s.BuildPut(ctx, false, serverptypes.TestValidBuild(t, &pb.Build{
			Id: "1",
		})))
		require.NoError(s.BuildPut(ctx, false, serverptypes.TestValidBuild(t, &pb.Build{
			Id: "2",
		})))
		require.NoError(s.BuildPut(ctx, false, serverptypes.TestValidBuild(t, &pb.Build{
			Id: "3",
			Application: &pb.Ref_Application{
				Application: "B",
				Project:     "B",
			},
		})))

		// Create some other resources
		require.NoError(s.DeploymentPut(ctx, false, serverptypes.TestValidDeployment(t, &pb.Deployment{
			Id: "1",
		})))

		// Workspace list should only list one
		{
			result, err := s.WorkspaceList(ctx)
			require.NoError(err)
			require.Len(result, 1)

			ws := result[0]
			require.Len(ws.Projects, 2)
			require.Len(ws.Projects[0].Applications, 1)
			require.Len(ws.Projects[1].Applications, 1)
		}

		// Create a new workspace
		require.NoError(s.BuildPut(ctx, false, serverptypes.TestValidBuild(t, &pb.Build{
			Id: "4",
			Workspace: &pb.Ref_Workspace{
				Workspace: "w2",
			},
		})))
		{
			result, err := s.WorkspaceList(ctx)
			require.NoError(err)
			require.Len(result, 2)
		}
	})
}

func TestWorkspacePut(t *testing.T, factory Factory, _ RestartFactory) {
	ctx := context.Background()
	t.Run("Default", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		{
			workspace, err := s.WorkspaceGet(ctx, "default")
			require.Equal(codes.NotFound, status.Code(err))
			require.Error(err)
			require.Nil(workspace)
		}

		// Put
		err := s.WorkspacePut(ctx, serverptypes.TestWorkspace(t, &pb.Workspace{
			Name: "default",
		}))
		require.NoError(err)

		{
			workspace, err := s.WorkspaceGet(ctx, "default")
			require.NoError(err)
			require.NotNil(workspace)
			require.Equal(workspace.Name, "default")
		}
	})

	t.Run("Only 1 workspace per name", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		{
			workspace, err := s.WorkspaceGet(ctx, "default")
			require.Equal(codes.NotFound, status.Code(err))
			require.Error(err)
			require.Nil(workspace)
		}

		// Put
		err := s.WorkspacePut(ctx, serverptypes.TestWorkspace(t, &pb.Workspace{
			Name: "default",
		}))
		require.NoError(err)

		{
			workspace, err := s.WorkspaceGet(ctx, "default")
			require.NoError(err)
			require.NotNil(workspace)
			require.Equal(workspace.Name, "default")
		}

		// Only 1
		{
			workspace, err := s.WorkspaceList(ctx)
			require.NoError(err)
			require.NotNil(workspace)
			require.Len(workspace, 1)
		}

		// Put
		err = s.WorkspacePut(ctx, serverptypes.TestWorkspace(t, &pb.Workspace{
			Name: "default",
		}))
		require.NoError(err)

		// Only 1
		{
			workspace, err := s.WorkspaceList(ctx)
			require.NoError(err)
			require.NotNil(workspace)
			require.Len(workspace, 1)
		}
	})

	t.Run("No spaces in name", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Put with a bad name
		err := s.WorkspacePut(ctx, serverptypes.TestWorkspace(t, &pb.Workspace{
			Name: "no spaces allowed",
		}))
		require.Error(err)
	})

	t.Run("Allow underscores and hyphens", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Underscores and hyphens are fine
		err := s.WorkspacePut(ctx, serverptypes.TestWorkspace(t, &pb.Workspace{
			Name: "special_and-allowed",
		}))
		require.NoError(err)
	})

	t.Run("Multi List", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Put default
		err := s.WorkspacePut(ctx, serverptypes.TestWorkspace(t, &pb.Workspace{
			Name: "default",
		}))
		require.NoError(err)

		// Put dev
		err = s.WorkspacePut(ctx, serverptypes.TestWorkspace(t, &pb.Workspace{
			Name: "dev",
		}))
		require.NoError(err)

		// Put staging
		err = s.WorkspacePut(ctx, serverptypes.TestWorkspace(t, &pb.Workspace{
			Name: "staging",
		}))
		require.NoError(err)

		{
			workspace, err := s.WorkspaceList(ctx)
			require.NoError(err)
			require.NotNil(workspace)
			require.Len(workspace, 3)
		}
	})

	t.Run("Preserves Projects", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Put a Workspace with a Project
		err := s.WorkspacePut(ctx, serverptypes.TestWorkspace(t, &pb.Workspace{
			Name: "staging",
			Projects: []*pb.Workspace_Project{
				{
					Project:   &pb.Ref_Project{Project: "projectA"},
					Workspace: &pb.Ref_Workspace{Workspace: "staging"},
				},
			},
		}))
		require.NoError(err)

		// Put again, without projects
		err = s.WorkspacePut(ctx, serverptypes.TestWorkspace(t, &pb.Workspace{
			Name: "staging",
		}))
		require.NoError(err)

		{
			workspace, err := s.WorkspaceGet(ctx, "staging")
			require.NoError(err)
			require.NotNil(workspace)
			require.Equal(workspace.Name, "staging")
			require.Len(workspace.Projects, 1)
		}
	})

	// Enforce that workspaces cannot start or end with either hyphens and
	// underscores, or contain spaces
	invalidNames := []string{
		"cannot contain spaces",
		" cannot start with spaces",
		"-starts_with-hyphen",
		"_starts-with_underscore",
		"_ends-with_underscore-_",
		"_ends-with_underscore-_",
	}

	for _, invalidName := range invalidNames {
		// hyphens and underscores are allowed, but names cannot start with them
		t.Run("Invalid_"+invalidName, func(t *testing.T) {
			require := require.New(t)

			s := factory(t)
			defer s.Close()

			// Workspace names cannot start with underscore or hyphens
			err := s.WorkspacePut(ctx, serverptypes.TestWorkspace(t, &pb.Workspace{
				Name: invalidName,
			}))
			require.Error(err)
		})
	}
}

func TestWorkspaceProject(t *testing.T, factory Factory, restartF RestartFactory) {
	ctx := context.Background()
	t.Run("List non-empty", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create a build
		require.NoError(s.BuildPut(ctx, false, serverptypes.TestValidBuild(t, &pb.Build{
			Id: "1",
		})))
		require.NoError(s.BuildPut(ctx, false, serverptypes.TestValidBuild(t, &pb.Build{
			Id: "2",
		})))
		require.NoError(s.BuildPut(ctx, false, serverptypes.TestValidBuild(t, &pb.Build{
			Id: "3",
			Application: &pb.Ref_Application{
				Application: "B",
				Project:     "B",
			},
			Workspace: &pb.Ref_Workspace{
				Workspace: "w1",
			},
		})))

		// Workspace list should return only 1 for B
		{
			result, err := s.WorkspaceListByProject(ctx, &pb.Ref_Project{
				Project: "B",
			})
			require.NoError(err)
			require.Len(result, 1)

			ws := result[0]
			require.Equal("w1", ws.Name)
			require.Len(ws.Projects, 1)
		}

		// Create a new workspace
		require.NoError(s.BuildPut(ctx, false, serverptypes.TestValidBuild(t, &pb.Build{
			Id: "4",
			Application: &pb.Ref_Application{
				Application: "B",
				Project:     "B",
			},
			Workspace: &pb.Ref_Workspace{
				Workspace: "w2",
			},
		})))
		{
			result, err := s.WorkspaceListByProject(ctx, &pb.Ref_Project{
				Project: "B",
			})
			require.NoError(err)
			require.Len(result, 2)
		}
	})
}

func TestWorkspaceApp(t *testing.T, factory Factory, restartF RestartFactory) {
	ctx := context.Background()
	t.Run("List non-empty", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create a build
		require.NoError(s.BuildPut(ctx, false, serverptypes.TestValidBuild(t, &pb.Build{
			Id: "1",
		})))
		require.NoError(s.BuildPut(ctx, false, serverptypes.TestValidBuild(t, &pb.Build{
			Id: "2",
			Application: &pb.Ref_Application{
				Application: "B",
				Project:     "A",
			},
		})))
		require.NoError(s.BuildPut(ctx, false, serverptypes.TestValidBuild(t, &pb.Build{
			Id: "3",
			Application: &pb.Ref_Application{
				Application: "B",
				Project:     "B",
			},
			Workspace: &pb.Ref_Workspace{
				Workspace: "w1",
			},
		})))

		// Workspace list should return only 1 for B,B
		{
			result, err := s.WorkspaceListByApp(ctx, &pb.Ref_Application{
				Application: "B",
				Project:     "B",
			})
			require.NoError(err)
			require.Len(result, 1)

			ws := result[0]
			require.Equal("w1", ws.Name)
			require.Len(ws.Projects, 1)
		}
	})
}
