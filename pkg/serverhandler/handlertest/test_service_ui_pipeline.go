package handlertest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func init() {
	tests["ui_pipeline"] = []testFunc{
		TestServiceUI_ListPipelines,
	}
}

func TestServiceUI_ListPipelines(t *testing.T, factory Factory) {
	require := require.New(t)
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	// Create projects
	for _, name := range []string{"alpha", "beta"} {
		_, err := client.UpsertProject(ctx, &pb.UpsertProjectRequest{
			Project: serverptypes.TestProject(t, &pb.Project{
				Name: name,
			}),
		})
		require.NoError(err)
	}

	// Call the method
	result, err := client.UI_ListPipelines(ctx, &pb.UI_ListPipelinesRequest{
		Pagination: &pb.PaginationRequest{
			PageSize: 10,
		},
	})
	require.NoError(err)
	require.Len(result.PipelineBundles, 2)
	require.EqualValues(2, result.TotalCount)
}
