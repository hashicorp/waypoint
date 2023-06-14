// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package handlertest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func init() {
	tests["deploy"] = []testFunc{
		TestServiceDeployment,
		TestServiceDeployment_GetDeployment,
		TestServiceDeployment_GetLatestDeployment,
		TestServiceDeployment_ListDeployments,
	}
}

func TestServiceDeployment(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	// Simplify writing tests
	type Req = pb.UpsertDeploymentRequest

	t.Run("create and update", func(t *testing.T) {
		require := require.New(t)

		artifact := serverptypes.TestValidArtifact(t, nil)

		artifactresp, err := client.UpsertPushedArtifact(ctx, &pb.UpsertPushedArtifactRequest{
			Artifact: artifact,
		})

		dep := serverptypes.TestValidDeployment(t, nil)
		dep.ArtifactId = artifactresp.Artifact.Id

		// Create, should get an ID back
		resp, err := client.UpsertDeployment(ctx, &Req{
			Deployment: dep,
		})
		require.NoError(err)
		require.NotNil(resp)
		result := resp.Deployment
		require.NotEmpty(result.Id)

		// Let's write some data
		result.Status = server.NewStatus(pb.Status_RUNNING)
		resp, err = client.UpsertDeployment(ctx, &Req{
			Deployment: result,
		})
		require.NoError(err)
		require.NotNil(resp)
		result = resp.Deployment
		require.NotNil(result.Status)
		require.Equal(pb.Status_RUNNING, result.Status.State)
	})

	t.Run("update non-existent", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertDeployment(ctx, &Req{
			Deployment: serverptypes.TestValidDeployment(t, &pb.Deployment{
				Id: "nope",
			}),
		})
		require.Error(err)
		require.Nil(resp)
		st, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, st.Code())
	})
}

func TestServiceDeployment_GetDeployment(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	artifact := serverptypes.TestValidArtifact(t, nil)

	artifactresp, err := client.UpsertPushedArtifact(ctx, &pb.UpsertPushedArtifactRequest{
		Artifact: artifact,
	})

	dep := serverptypes.TestValidDeployment(t, nil)
	dep.ArtifactId = artifactresp.Artifact.Id
	// Best way to mock for now is to make a request
	resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment: dep,
	})

	require.NoError(t, err)

	// Simplify writing tests
	type Req = pb.GetDeploymentRequest

	t.Run("get existing", func(t *testing.T) {
		require := require.New(t)

		// Get, should return a deployment
		deployment, err := client.GetDeployment(ctx, &Req{
			Ref: &pb.Ref_Operation{
				Target: &pb.Ref_Operation_Id{Id: resp.Deployment.Id},
			},
		})
		require.NoError(err)
		require.NotNil(deployment)
		require.NotEmpty(deployment.Id)
	})

	t.Run("get non-existing", func(t *testing.T) {
		require := require.New(t)

		// get, should fail
		resp, err := client.GetDeployment(ctx, &Req{
			Ref: &pb.Ref_Operation{
				Target: &pb.Ref_Operation_Id{Id: "nope"},
			},
		})
		require.Error(err)
		require.Nil(resp)
		st, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, st.Code())
	})
}

func TestServiceDeployment_GetLatestDeployment(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	var application *pb.Ref_Application
	var workspace *pb.Ref_Workspace

	var latestSuccessfulDepResp *pb.UpsertDeploymentResponse
	// Upsert some deployments
	depNum := []int{1, 2, 3}
	for index := range depNum {
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
		if index == len(depNum)-1 { // make latest deployment a failed one
			dep.Status.State = pb.Status_ERROR
			dep.Status.Error = status.New(codes.Internal, "test failed deployment").Proto()
		}

		// Best way to mock for now is to make a request
		resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
			Deployment: dep,
		})

		require.NoError(t, err)

		// set latest successful deployment response
		if resp.Deployment.Status.State == pb.Status_SUCCESS {
			latestSuccessfulDepResp = resp
		}
		application = resp.Deployment.Application
		workspace = resp.Deployment.Workspace
	}

	// Simplify writing tests
	type Req = pb.GetLatestDeploymentRequest

	t.Run("get latest", func(t *testing.T) {
		require := require.New(t)

		latestDeploymentResponse, err := client.GetLatestDeployment(ctx, &Req{
			Application: application,
			Workspace:   workspace,
		})

		// GetLatest, should return the latest successful deployment
		require.NoError(err)
		require.NotEmpty(latestDeploymentResponse)
		require.Equal(latestDeploymentResponse.Deployment.Id, latestSuccessfulDepResp.Deployment.Id)
	})
}

func TestServiceDeployment_ListDeployments(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

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

	// Best way to mock for now is to make a request
	resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment: dep,
	})

	require.NoError(t, err)

	// Simplify writing tests
	type Req = pb.ListDeploymentsRequest

	t.Run("list", func(t *testing.T) {
		require := require.New(t)

		// Get, should return a deployment
		deployments, err := client.ListDeployments(ctx, &Req{
			Application: resp.Deployment.Application,
		})
		require.NoError(err)
		require.NotEmpty(deployments)
		require.Equal(deployments.Deployments[0].Id, resp.Deployment.Id)
	})

	t.Run("list with artifact", func(t *testing.T) {
		require := require.New(t)

		// Get, should return a deployment
		deployments, err := client.ListDeployments(ctx, &Req{
			Application: resp.Deployment.Application,
			LoadDetails: pb.Deployment_ARTIFACT,
		})
		require.NoError(err)
		require.NotEmpty(deployments)
		require.Equal(deployments.Deployments[0].Id, resp.Deployment.Id)
		require.NotNil(deployments.Deployments[0].Preload.Artifact)
		require.Nil(deployments.Deployments[0].Preload.Build)
		require.Equal(deployments.Deployments[0].Preload.Artifact.Id, artifactresp.Artifact.Id)
	})

	t.Run("list with build", func(t *testing.T) {
		require := require.New(t)

		// Get, should return a deployment
		deployments, err := client.ListDeployments(ctx, &Req{
			Application: resp.Deployment.Application,
			LoadDetails: pb.Deployment_BUILD,
		})
		require.NoError(err)
		require.NotEmpty(deployments)
		require.Equal(deployments.Deployments[0].Id, resp.Deployment.Id)
		require.NotNil(deployments.Deployments[0].Preload.Artifact)
		require.NotNil(deployments.Deployments[0].Preload.Build)
		require.Equal(deployments.Deployments[0].Preload.Artifact.Id, artifactresp.Artifact.Id)
		require.Equal(deployments.Deployments[0].Preload.Build.Id, build.Id)
	})
}
