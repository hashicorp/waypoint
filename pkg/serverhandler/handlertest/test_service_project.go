package handlertest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	empty "google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func init() {
	tests["project"] = []testFunc{
		TestServiceProject,
		TestServiceProject_GetApplication,
		TestServiceProject_UpsertApplication,
		TestServiceProject_InvalidName,
	}
}

func TestServiceProject(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)
	project := ptypes.TestProject(t, &pb.Project{
		Name: "example",
	})

	t.Run("create and update", func(t *testing.T) {
		require := require.New(t)

		// Creates a project
		{
			resp, err := client.UpsertProject(ctx, &pb.UpsertProjectRequest{
				Project: project,
			})

			require.NoError(err)
			require.NotNil(resp)
			require.Len(resp.Project.Applications, 0)
			require.False(resp.Project.RemoteEnabled)
		}

		// Updates a project by making project remote
		{
			project.RemoteEnabled = true
			resp, err := client.UpsertProject(ctx, &pb.UpsertProjectRequest{
				Project: project,
			})
			require.NoError(err)
			require.NotNil(resp)
			require.True(resp.Project.RemoteEnabled)
		}
	})

	t.Run("get", func(t *testing.T) {
		require := require.New(t)

		// Returns an error for a missing project
		{
			resp, err := client.GetProject(ctx, &pb.GetProjectRequest{
				Project: &pb.Ref_Project{Project: "not-found"},
			})
			require.Error(err)
			require.Nil(resp)
		}

		// Returns a response for a project that exists
		{
			resp, err := client.GetProject(ctx, &pb.GetProjectRequest{
				Project: &pb.Ref_Project{Project: "example"},
			})
			require.NoError(err)
			require.NotNil(resp)
			require.Equal(resp.Project.Name, "example")
		}
	})

	t.Run("list", func(t *testing.T) {
		require := require.New(t)

		// Lists all projects
		{
			resp, err := client.ListProjects(ctx, &empty.Empty{})
			require.NoError(err)
			require.NotNil(resp)
			require.Len(resp.Projects, 1)
		}
	})

	t.Run("destroy", func(t *testing.T) {
		require := require.New(t)

		// Destroys the specified project
		{
			_, err := client.DestroyProject(ctx, &pb.DestroyProjectRequest{
				Project: &pb.Ref_Project{Project: "example"},
			})
			require.NoError(err)

			resp, err := client.GetProject(ctx, &pb.GetProjectRequest{
				Project: &pb.Ref_Project{Project: "not-found"},
			})
			require.Error(err)
			require.Nil(resp)
		}
	})
}

func TestServiceProject_GetApplication(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)
	project := ptypes.TestProject(t, &pb.Project{
		Name: "example",
	})

	t.Run("get", func(t *testing.T) {
		require := require.New(t)

		// Returns an error if the application doesn't exist
		{
			resp, err := client.GetApplication(ctx, &pb.GetApplicationRequest{
				Application: &pb.Ref_Application{Application: "doesnt-exist"},
			})
			require.Error(err)
			require.Nil(resp)
		}

		// Create a project
		resp, err := client.UpsertProject(ctx, &pb.UpsertProjectRequest{
			Project: project,
		})
		require.NoError(err)
		require.NotNil(resp)

		// Gets an application inside a project
		{
			resp, err := client.UpsertApplication(ctx, &pb.UpsertApplicationRequest{
				Project: &pb.Ref_Project{Project: "example"},
				Name:    "Apple",
			})
			require.NoError(err)
			require.NotNil(resp)

			respApp, err := client.GetApplication(ctx, &pb.GetApplicationRequest{
				Application: &pb.Ref_Application{
					Application: "Apple",
					Project:     "example",
				},
			})
			require.NoError(err)
			require.NotNil(respApp)
			require.Equal(respApp.Application.Name, "Apple")
		}
	})
}

func TestServiceProject_UpsertApplication(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)
	project := ptypes.TestProject(t, &pb.Project{
		Name: "example",
	})

	t.Run("create and update", func(t *testing.T) {
		require := require.New(t)

		// Returns an error if the project doesn't exist
		{
			resp, err := client.UpsertApplication(ctx, &pb.UpsertApplicationRequest{
				Project: &pb.Ref_Project{Project: "doesnt-exist"},
				Name:    "nope",
			})
			require.Error(err)
			require.Nil(resp)
		}

		//create a project
		resp, err := client.UpsertProject(ctx, &pb.UpsertProjectRequest{
			Project: project,
		})
		require.NoError(err)
		require.NotNil(resp)

		// creates an application inside a project
		{
			resp, err := client.UpsertApplication(ctx, &pb.UpsertApplicationRequest{
				Project: &pb.Ref_Project{Project: "example"},
				Name:    "Apple",
			})
			require.NoError(err)
			require.NotNil(resp)

			resp, err = client.UpsertApplication(ctx, &pb.UpsertApplicationRequest{
				Project: &pb.Ref_Project{Project: "example"},
				Name:    "Orange",
			})
			require.NoError(err)
			require.NotNil(resp)
			require.Equal(resp.Application.FileChangeSignal, "")
		}

		// updates a file change signal for the app
		{
			resp, err := client.UpsertApplication(ctx, &pb.UpsertApplicationRequest{
				Project:          &pb.Ref_Project{Project: "example"},
				Name:             "Orange",
				FileChangeSignal: "SIGINT",
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// adds runner profile if defined
		{
			resp, err := client.UpsertApplication(ctx, &pb.UpsertApplicationRequest{
				Project: &pb.Ref_Project{Project: "example"},
				Name:    "Apple",
			})
			require.NoError(err)
			require.NotNil(resp)
		}
	})
}

func TestServiceProject_InvalidName(t *testing.T, factory Factory) {
	ctx := context.Background()
	require := require.New(t)
	client, _ := factory(t)

	// GRPC Gateway interprets ../ as a path traversal, and therefore we cannot allow
	// '../' in any fields we use as path tokens.
	project := ptypes.TestProject(t, &pb.Project{
		Name: "../../",
	})

	// Create a project
	_, err := client.UpsertProject(ctx, &pb.UpsertProjectRequest{
		Project: project,
	})
	require.Error(err)
}
