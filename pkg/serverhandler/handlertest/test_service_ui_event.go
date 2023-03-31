package handlertest

import (
	"context"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
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

	t.Run("paginate events", func(t *testing.T) {
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

		time.Sleep(250 * time.Millisecond)

		// Create build, deployment, release
		createEvents(client, ctx, t)

		eventResp, err := client.UI_ListEvents(ctx, &pb.UI_ListEventsRequest{
			Application: &pb.Ref_Application{
				Application: "a_test",
				Project:     "p_test",
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
			Pagination: &pb.PaginationRequest{
				PageSize:          6,
				NextPageToken:     "",
				PreviousPageToken: "",
			},
			Sorting: &pb.SortingRequest{OrderBy: []string{"event_timestamp desc"}},
		})
		require.NoError(err)
		require.NotNil(eventResp)
		// Only have 3 resources, 1: build, 1: deployment, 1: release
		require.Len(eventResp.Events, 6)

		var nextPageToken string
		t.Run("paginate events, test get first page", func(t *testing.T) {
			//first page, and next page
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
			nextPageToken = eventResp.Pagination.NextPageToken

			require.NoError(err)
			require.NotNil(eventResp)
		})

		t.Run("paginate events, test next page token", func(t *testing.T) {
			//first page, and next page
			eventResp, err := client.UI_ListEvents(ctx, &pb.UI_ListEventsRequest{
				Application: &pb.Ref_Application{
					Application: "a_test",
					Project:     "p_test",
				},
				Workspace: &pb.Ref_Workspace{Workspace: "default"},
				Pagination: &pb.PaginationRequest{
					PageSize:          5,
					NextPageToken:     nextPageToken,
					PreviousPageToken: "",
				},
				Sorting: &pb.SortingRequest{OrderBy: []string{"event_timestamp desc"}},
			})
			require.NoError(err)
			require.NotNil(eventResp)
			require.Len(eventResp.Events, 1) //there is only 1 resource left
		})
	})
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