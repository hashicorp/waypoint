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
	tests["artifact"] = []testFunc{
		TestServiceArtifact,
		TestServiceArtifact_List,
	}
}

func TestServiceArtifact(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	// Simplify writing tests
	type Req = pb.UpsertPushedArtifactRequest

	t.Run("create and update", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertPushedArtifact(ctx, &Req{
			Artifact: serverptypes.TestValidArtifact(t, nil),
		})
		require.NoError(err)
		require.NotNil(resp)
		result := resp.Artifact
		require.NotEmpty(result.Id)

		// Let's write some data
		result.Status = server.NewStatus(pb.Status_RUNNING)
		resp, err = client.UpsertPushedArtifact(ctx, &Req{
			Artifact: result,
		})
		require.NoError(err)
		require.NotNil(resp)
		result = resp.Artifact
		require.NotNil(result.Status)
		require.Equal(pb.Status_RUNNING, result.Status.State)
	})

	t.Run("update non-existent", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertPushedArtifact(ctx, &Req{
			Artifact: serverptypes.TestValidArtifact(t, &pb.PushedArtifact{
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

func TestServiceArtifact_List(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	// Simplify writing tests
	type Req = pb.ListPushedArtifactsRequest

	t.Run("list with build", func(t *testing.T) {
		require := require.New(t)

		buildresp, err := client.UpsertBuild(ctx, &pb.UpsertBuildRequest{
			Build: serverptypes.TestValidBuild(t, nil),
		})
		require.NoError(err)
		require.NotNil(buildresp)

		build := buildresp.Build

		artifact := serverptypes.TestValidArtifact(t, nil)
		artifact.BuildId = build.Id

		resp, err := client.UpsertPushedArtifact(ctx, &pb.UpsertPushedArtifactRequest{
			Artifact: artifact,
		})
		require.NoError(err)
		require.NotNil(resp)

		artifact = resp.Artifact

		// Create, should get an ID back
		listresp, err := client.ListPushedArtifacts(ctx, &Req{
			Application:  artifact.Application,
			IncludeBuild: true,
		})
		require.NoError(err)
		require.NotNil(listresp)

		require.NotNil(listresp.Artifacts[0].Build)
	})

}
