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

	t.Run("list events, count less than page size", func(t *testing.T){
		require := require.New(t)


		// Initialize our app
		TestApp(t, client, &pb.Ref_Application{
			Application: "a_test",
			Project:     "p_test",
		})

		//create build, deployment, release
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
			Sorting: &pb.SortingRequest{OrderBy: []string{"event_type", "event_timestamp desc"}},
		})
		require.NoError(err)
		require.NotNil(eventResp)
		require.Len(eventResp, 3)
	})




}
// create simple build, deployment, release for eventListBundling to be cleaner
func createEvents(client pb.WaypointClient, ctx context.Context, t *testing.T) {

	// Create Build, should get an ID back
	buildresp, err := client.UpsertBuild(ctx, &pb.UpsertBuildRequest{
		Build: serverptypes.TestValidBuild(t, nil),
	})
	require.NoError(t, err)
	require.NotNil(t, buildresp)

	build := buildresp.Build

	artifact := serverptypes.TestValidArtifact(t, nil)
	artifact.BuildId = build.Id

	artifactresp, err := client.UpsertPushedArtifact(ctx, &pb.UpsertPushedArtifactRequest{
		Artifact: artifact,
	})
	require.NoError(t, err)
	require.NotNil(t, artifactresp)

	dep := serverptypes.TestValidDeployment(t, nil)
	dep.ArtifactId = artifactresp.Artifact.Id

	// Create Deployment, should get an ID back
	depResp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment: serverptypes.TestValidDeployment(t, nil),
	})
	require.NoError(t, err)
	require.NotNil(t, depResp)
	depResult := depResp.Deployment
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