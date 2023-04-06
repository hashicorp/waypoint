// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package handlertest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func init() {
	tests["ui_project"] = []testFunc{
		TestServiceUI_ListProjects,
		TestServiceUI_GetProject,
	}
}

func TestServiceUI_ListProjects(t *testing.T, factory Factory) {
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
	result, err := client.UI_ListProjects(ctx, &pb.UI_ListProjectsRequest{
		Pagination: &pb.PaginationRequest{
			PageSize: 10,
		},
	})
	require.NoError(err)
	require.Len(result.ProjectBundles, 2)
	require.EqualValues(2, result.TotalCount)
}

func TestServiceUI_GetProject(t *testing.T, factory Factory) {
	require := require.New(t)
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	project := serverptypes.TestProject(t, &pb.Project{
		Name: "example",
	})

	// Create a project
	_, err := client.UpsertProject(ctx, &pb.UpsertProjectRequest{
		Project: project,
	})
	require.NoError(err)

	// Queue an older InitOp job
	_, err = queueTestInitJob(t, ctx, client, project)
	require.NoError(err)

	// Queue a newer InitOp job
	queueJobResp, err := queueTestInitJob(t, ctx, client, project)
	require.NoError(err)
	require.NotEmpty(queueJobResp.JobId)

	// Get the project using UI_GetProject
	getProjectResp, err := client.UI_GetProject(ctx, &pb.UI_GetProjectRequest{
		Project: &pb.Ref_Project{
			Project: "example",
		},
	})

	require.NoError(err)
	require.NotNil(getProjectResp)
	require.NotNil(getProjectResp.Project, "should load a project")
	require.Equal(getProjectResp.Project.Name, "example", "should load the correct project")
	require.NotNil(getProjectResp.LatestInitJob, "should sideload an InitJob")
	require.Equal(
		getProjectResp.LatestInitJob.Id,
		queueJobResp.JobId,
		"should sideload the latest InitJob",
	)
}

func queueTestInitJob(
	t *testing.T,
	ctx context.Context,
	client pb.WaypointClient,
	project *pb.Project,
) (*pb.QueueJobResponse, error) {
	return client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, &pb.Job{
			Application: &pb.Ref_Application{
				Project: project.Name,
			},
			Workspace: &pb.Ref_Workspace{
				Workspace: "default",
			},
			Operation: &pb.Job_Init{},
		}),
	})
}
