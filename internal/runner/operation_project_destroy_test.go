// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package runner

import (
	"context"
	"github.com/hashicorp/go-hclog"
	configinternal "github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/core"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/singleprocess"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProjectDestroyOp(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()
	client := singleprocess.TestServer(t)

	log := hclog.New(&hclog.LoggerOptions{
		Name:            "test-runner",
		Level:           hclog.Debug,
		IncludeLocation: true,
	})

	// Initialize our runner
	runner, err := New(
		WithClient(client),
		WithLogger(log),
	)
	require.NoError(err)
	defer runner.Close()

	// Start it
	require.NoError(runner.Start(ctx))

	const projectName string = "projectDestroyTestProject"
	const appName string = "projectDestroyTestApp"

	// Create a project for us to destroy
	projectResp, err := client.UpsertProject(ctx, &pb.UpsertProjectRequest{
		Project: &pb.Project{
			Name: projectName,
			Applications: []*pb.Application{
				{
					Project: &pb.Ref_Project{Project: projectName},
					Name:    appName,
				},
			},
		},
	})
	require.NoError(err)
	require.NotNil(projectResp)

	// Create a deployment for us to destroy
	deploymentResp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment: &pb.Deployment{
			Application: &pb.Ref_Application{
				Application: appName,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}})
	require.NoError(err)
	require.NotNil(deploymentResp)

	job := &pb.Job{
		Operation: &pb.Job_DestroyProject{
			DestroyProject: &pb.Job_DestroyProjectOp{
				Project: &pb.Ref_Project{
					Project: projectName,
				},
				SkipDestroyResources: false,
			},
		},
	}
	project := core.TestProject(t, core.WithClient(client))
	require.NoError(err)
	require.NotNil(project)

	cfg := configinternal.TestConfig(t, core.TestProjectConfig)

	res, err := runner.executeDestroyProjectOp(ctx, runner.logger, job, project, cfg)
	require.NoError(err)
	require.NotNil(t, res.ProjectDestroy)
	require.NotNil(t, res.ProjectDestroy.JobId)

	// Verify that we can't get the project we deleted
	getProjectResp, err := client.GetProject(ctx, &pb.GetProjectRequest{Project: &pb.Ref_Project{Project: projectName}})
	require.Error(err)
	require.Nil(getProjectResp)

	// Verify that we can't get the deployment destroyed
	getDeploymentResp, err := client.GetDeployment(ctx, &pb.GetDeploymentRequest{
		Ref: &pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: deploymentResp.Deployment.Id}},
	})
	require.Error(err)
	require.Nil(getDeploymentResp)
}

// NOTE: This test is very similar to the previous test, with the exception of the
// skip destroy resources option being set, because in either case, after the project
// destroy op is complete, the deployment will be removed from the DB, and should return
// an error when we try to get it. There is still value in this test though, because it
// verifies that there was no error even if we skip destroying resources
func TestProjectDestroyOp_SkipDestroyResources(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()
	client := singleprocess.TestServer(t)

	log := hclog.New(&hclog.LoggerOptions{
		Name:            "test-runner",
		Level:           hclog.Debug,
		IncludeLocation: true,
	})

	// Initialize our runner
	runner, err := New(
		WithClient(client),
		WithLogger(log),
	)
	require.NoError(err)
	defer runner.Close()

	// Start it
	require.NoError(runner.Start(ctx))

	// Create a project for us to destroy
	projectResp, err := client.UpsertProject(ctx, &pb.UpsertProjectRequest{
		Project: &pb.Project{
			Name: "testProject",
			Applications: []*pb.Application{
				{
					Project: &pb.Ref_Project{Project: "testProject"},
					Name:    "testApp",
				},
			},
		},
	})
	require.NoError(err)
	require.NotNil(projectResp)

	// Create a deployment for us to destroy
	deploymentResp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment: &pb.Deployment{
			Application: &pb.Ref_Application{
				Application: "testApp",
				Project:     "testProject",
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}})
	require.NoError(err)
	require.NotNil(deploymentResp)

	job := &pb.Job{
		Operation: &pb.Job_DestroyProject{
			DestroyProject: &pb.Job_DestroyProjectOp{
				Project: &pb.Ref_Project{
					Project: "testProject",
				},
				SkipDestroyResources: true,
			},
		},
	}
	project := core.TestProject(t, core.WithClient(client))
	require.NoError(err)
	require.NotNil(project)

	cfg := configinternal.TestConfig(t, core.TestProjectConfig)
	res, err := runner.executeDestroyProjectOp(ctx, runner.logger, job, project, cfg)
	require.NoError(err)
	require.NotNil(t, res.ProjectDestroy)
	require.NotNil(t, res.ProjectDestroy.JobId)

	// Verify that we can't get the project we deleted
	getProjectResp, err := client.GetProject(ctx, &pb.GetProjectRequest{Project: &pb.Ref_Project{Project: "testProject"}})
	require.Error(err)
	require.Nil(getProjectResp)

	// Verify that we can't get the deployment destroyed
	getDeploymentResp, err := client.GetDeployment(ctx, &pb.GetDeploymentRequest{
		Ref: &pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: deploymentResp.Deployment.Id}},
	})
	require.Error(err)
	require.Nil(getDeploymentResp)
}
