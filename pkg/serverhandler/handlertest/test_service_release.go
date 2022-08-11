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
	tests["release"] = []testFunc{
		TestServiceRelease,
		TestServiceRelease_GetRelease,
		TestServiceRelease_ListReleases,
	}
}

func TestServiceRelease(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	// Simplify writing tests
	type Req = pb.UpsertReleaseRequest

	t.Run("create and update", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertRelease(ctx, &Req{
			Release: serverptypes.TestValidRelease(t, nil),
		})
		require.NoError(err)
		require.NotNil(resp)
		result := resp.Release
		require.NotEmpty(result.Id)

		// Let's write some data
		result.Status = server.NewStatus(pb.Status_RUNNING)
		resp, err = client.UpsertRelease(ctx, &Req{
			Release: result,
		})
		require.NoError(err)
		require.NotNil(resp)
		result = resp.Release
		require.NotNil(result.Status)
		require.Equal(pb.Status_RUNNING, result.Status.State)
	})

	t.Run("update non-existent", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertRelease(ctx, &Req{
			Release: serverptypes.TestValidRelease(t, &pb.Release{
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

func TestServiceRelease_GetRelease(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	// Best way to mock for now is to make a request
	resp, err := client.UpsertRelease(ctx, &pb.UpsertReleaseRequest{
		Release: serverptypes.TestValidRelease(t, nil),
	})

	require.NoError(t, err)

	// Simplify writing tests
	type Req = pb.GetReleaseRequest

	t.Run("get existing", func(t *testing.T) {
		require := require.New(t)

		// Get, should return a release
		release, err := client.GetRelease(ctx, &Req{
			Ref: &pb.Ref_Operation{
				Target: &pb.Ref_Operation_Id{Id: resp.Release.Id},
			},
		})
		require.NoError(err)
		require.NotNil(release)
		require.NotEmpty(release.Id)
	})

	t.Run("get non-existing", func(t *testing.T) {
		require := require.New(t)

		// get, should fail
		resp, err := client.GetRelease(ctx, &Req{
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

func TestServiceRelease_ListReleases(t *testing.T, factory Factory) {
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

	depresp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment: dep,
	})
	require.NoError(t, err)
	require.NotNil(t, depresp)

	release := serverptypes.TestValidRelease(t, nil)
	release.DeploymentId = depresp.Deployment.Id

	// Best way to mock for now is to make a request
	resp, err := client.UpsertRelease(ctx, &pb.UpsertReleaseRequest{
		Release: release,
	})

	require.NoError(t, err)

	// Simplify writing tests
	type Req = pb.ListReleasesRequest

	t.Run("list", func(t *testing.T) {
		require := require.New(t)

		// Get, should return a release
		releases, err := client.ListReleases(ctx, &Req{
			Application: resp.Release.Application,
		})
		require.NoError(err)
		require.NotEmpty(releases)
		require.Equal(releases.Releases[0].Id, resp.Release.Id)
	})

	t.Run("list with artifact", func(t *testing.T) {
		require := require.New(t)

		// Get, should return a release
		releases, err := client.ListReleases(ctx, &Req{
			Application: resp.Release.Application,
			LoadDetails: pb.Release_ARTIFACT,
		})
		require.NoError(err)
		require.NotEmpty(releases)
		require.Equal(releases.Releases[0].Id, resp.Release.Id)
		require.NotNil(releases.Releases[0].Preload.Artifact)
		require.Nil(releases.Releases[0].Preload.Build)
		require.Equal(releases.Releases[0].Preload.Artifact.Id, artifactresp.Artifact.Id)
	})

	t.Run("list with build", func(t *testing.T) {
		require := require.New(t)

		// Get, should return a release
		releases, err := client.ListReleases(ctx, &Req{
			Application: resp.Release.Application,
			LoadDetails: pb.Release_BUILD,
		})
		require.NoError(err)
		require.NotEmpty(releases)
		require.Equal(releases.Releases[0].Id, resp.Release.Id)
		require.NotNil(releases.Releases[0].Preload.Artifact)
		require.NotNil(releases.Releases[0].Preload.Build)
		require.Equal(releases.Releases[0].Preload.Artifact.Id, artifactresp.Artifact.Id)
		require.Equal(releases.Releases[0].Preload.Build.Id, build.Id)
	})
}
