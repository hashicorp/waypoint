package handlertest

import (
	"context"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/stretchr/testify/require"
	"testing"
)

func init() {
	tests["event"] = []testFunc{
		TestEvent,
	}
}
func TestEvent(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	t.Run("paginate events, resource count less than page size", func(t *testing.T) {
		require := require.New(t)

		proj := &pb.Project{
			Name: "p_test",
		}
		refProj := &pb.Ref_Project{Project: "p_test"}
		_, err := client.UpsertProject(context.Background(), &pb.UpsertProjectRequest{
			Project: proj,
		})
		require.NoError(err)

		_, err = client.UpsertApplication(context.Background(), &pb.UpsertApplicationRequest{
			Project: refProj,
			Name:    "a_test",
		})
		require.NoError(err)

		// Create build, deployment, release
		createEvents(client, ctx, t)

		eventResp, err := client.UI_ListEvents(ctx, &pb.UI_ListEventsRequest{
			Application: &pb.Ref_Application{
				Application: "a_test",
				Project:     "p_test",
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
			Pagination: &pb.PaginationRequest{
				PageSize:          5,
				NextPageToken:     "",
				PreviousPageToken: "",
			},
			Sorting: &pb.SortingRequest{OrderBy: []string{"event_timestamp desc"}},
		})
		require.NoError(err)
		require.NotNil(eventResp)
		// Only have 3 resources, 1: build, 1: deployment, 1: release
		require.Len(eventResp.Events, 3)
	})

	t.Run("paginate events, resource count equal page size", func(t *testing.T) {
		require := require.New(t)

		// Initialize our app
		TestApp(t, client, &pb.Ref_Application{
			Application: "a_test",
			Project:     "p_test",
		})

		// Create build, deployment, release
		createEvents(client, ctx, t)

		eventResp, err := client.UI_ListEvents(ctx, &pb.UI_ListEventsRequest{
			Application: &pb.Ref_Application{
				Application: "a_test",
				Project:     "p_test",
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
			Pagination: &pb.PaginationRequest{
				PageSize:          3,
				NextPageToken:     "",
				PreviousPageToken: "",
			},
			Sorting: &pb.SortingRequest{OrderBy: []string{"event_timestamp desc"}},
		})
		require.NoError(err)
		require.NotNil(eventResp)
		// Only have 3 resources, 1: build, 1: deployment, 1: release
		require.Len(eventResp.Events, 3)
	})
	t.Run("paginate events, resource count more than page size", func(t *testing.T) {
		require := require.New(t)

		// Initialize our app
		TestApp(t, client, &pb.Ref_Application{
			Application: "a_test",
			Project:     "p_test",
		})

		// Create 2 of all: build, deployment, release
		createEvents(client, ctx, t)
		//time.Sleep(2 * time.Second)
		//createEvents(client, ctx, t)
		eventResp, err := client.UI_ListEvents(ctx, &pb.UI_ListEventsRequest{
			Application: &pb.Ref_Application{
				Application: "a_test",
				Project:     "p_test",
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
			Pagination: &pb.PaginationRequest{
				PageSize:          2,
				NextPageToken:     "",
				PreviousPageToken: "",
			},
			Sorting: &pb.SortingRequest{OrderBy: []string{"event_timestamp desc"}},
		})
		require.NoError(err)
		require.NotNil(eventResp)
		// Only have 6 resources, 2: build, 2: deployment, 2: release
		require.Len(eventResp.Events, 2)
	})

	//TODO: add more tests
}

// create simple build, deployment, release for eventListBundling to be cleaner
func createEvents(client pb.WaypointClient, ctx context.Context, t *testing.T) {
	// Create Build, should get an ID back
	buildResp, err := client.UpsertBuild(ctx, &pb.UpsertBuildRequest{
		Build: serverptypes.TestValidBuild(t, nil),
	})
	require.NoError(t, err)
	require.NotNil(t, buildResp)

	build := buildResp.Build

	artifact := serverptypes.TestValidArtifact(t, nil)
	artifact.BuildId = build.Id

	artifactResp, err := client.UpsertPushedArtifact(ctx, &pb.UpsertPushedArtifactRequest{
		Artifact: artifact,
	})
	require.NoError(t, err)
	require.NotNil(t, artifactResp)

	dep := serverptypes.TestValidDeployment(t, nil)
	dep.ArtifactId = artifactResp.Artifact.Id

	// Create Deployment, should get an ID back
	deploymentResp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment: serverptypes.TestValidDeployment(t, &pb.Deployment{
			//Component: &pb.Component{
			//	Name: "testapp",
			//},
			Application: &pb.Ref_Application{
				Application: "a_test",
				Project:     "p_test",
			},
			ArtifactId: artifactResp.Artifact.Id,
		}),
	})
	require.NoError(t, err)
	require.NotNil(t, deploymentResp)
	depResult := deploymentResp.Deployment
	require.NotEmpty(t, depResult.Id)

	// Create Release, should get an ID back
	relResp, err := client.UpsertRelease(ctx, &pb.UpsertReleaseRequest{
		Release: serverptypes.TestValidRelease(t, nil),
	})
	require.NoError(t, err)
	require.NotNil(t, relResp)
	relResult := relResp.Release
	require.NotEmpty(t, relResult.Id)
}