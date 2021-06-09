package core

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	componentmocks "github.com/hashicorp/waypoint-plugin-sdk/component/mocks"
	sdk "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	"github.com/hashicorp/waypoint/internal/config"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestAppDeploymentStatusReport(t *testing.T) {
	t.Run("with no status implementation", func(t *testing.T) {
		ctx := context.Background()
		require := require.New(t)

		// Our mock platform, which must also implement Status
		mock := &mockPlatformStatus{}

		// Make our factory for platforms
		factory := TestFactory(t, component.PlatformType)
		TestFactoryRegister(t, factory, "test", mock)

		// Make our app
		app := TestApp(t, TestProject(t,
			WithConfig(config.TestConfig(t, testPlatformConfig)),
			WithFactory(component.PlatformType, factory),
		), "test")
		client := app.client

		// We're using GetVersionInfoResponse here just because it is a proto message
		// that can be converted to an any.Any easily. We never use it, it's just to keep
		// the tests from blowing up with a nil reference.
		mockPluginArtifact := &pb.GetVersionInfoResponse{}

		anyval, err := ptypes.MarshalAny(mockPluginArtifact)
		require.NoError(err)

		aresp, err := client.UpsertPushedArtifact(ctx, &pb.UpsertPushedArtifactRequest{
			Artifact: serverptypes.TestValidArtifact(t, &pb.PushedArtifact{
				Artifact: &pb.Artifact{
					Artifact: anyval,
				},
			}),
		})
		require.NoError(err)

		resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
			Deployment: serverptypes.TestValidDeployment(t, &pb.Deployment{
				ArtifactId: aresp.Artifact.Id,
			}),
		})
		require.NoError(err)
		deploy := resp.Deployment

		// not implemented
		mock.Status.On("StatusFunc").Return(nil)

		// Status Report
		srResp, err := app.DeploymentStatusReport(context.Background(), deploy)
		require.NoError(err)
		require.Nil(srResp)

	})

	t.Run("with status implementation on deploy", func(t *testing.T) {
		ctx := context.Background()
		require := require.New(t)

		// Our mock platform, which must also implement Status
		mock := &mockPlatformStatus{}

		// Make our factory for platforms
		factory := TestFactory(t, component.PlatformType)
		TestFactoryRegister(t, factory, "test", mock)

		// Make our app
		app := TestApp(t, TestProject(t,
			WithConfig(config.TestConfig(t, testPlatformConfig)),
			WithFactory(component.PlatformType, factory),
		), "test")
		client := app.client

		// We're using GetVersionInfoResponse here just because it is a proto message
		// that can be converted to an any.Any easily. We never use it, it's just to keep
		// the tests from blowing up with a nil reference.
		mockPluginArtifact := &pb.GetVersionInfoResponse{}

		anyval, err := ptypes.MarshalAny(mockPluginArtifact)
		require.NoError(err)

		aresp, err := client.UpsertPushedArtifact(ctx, &pb.UpsertPushedArtifactRequest{
			Artifact: serverptypes.TestValidArtifact(t, &pb.PushedArtifact{
				Artifact: &pb.Artifact{
					Artifact: anyval,
				},
			}),
		})
		require.NoError(err)

		resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
			Deployment: serverptypes.TestValidDeployment(t, &pb.Deployment{
				ArtifactId: aresp.Artifact.Id,
			}),
		})
		require.NoError(err)
		deploy := resp.Deployment

		statusReportTs := timestamppb.Now()
		mock.Status.On("StatusFunc").Return(func(context.Context) (*sdk.StatusReport, error) {
			return &sdk.StatusReport{
				GeneratedTime: statusReportTs,
			}, nil
		})

		// Status Report
		srResp, err := app.DeploymentStatusReport(context.Background(), deploy)
		statusReport := &sdk.StatusReport{}
		anypb.UnmarshalTo(srResp.StatusReport, statusReport, proto.UnmarshalOptions{})
		require.NoError(err)
		require.NotNil(srResp.StatusReport)
		require.NotNil(statusReport.Health)

		// Verify that we have a Target of the right type with the right id
		require.IsType(srResp.TargetId, &pb.StatusReport_DeploymentId{})
		require.Equal(srResp.TargetId.(*pb.StatusReport_DeploymentId).DeploymentId, deploy.Id)

		// Verify that the status report timestamp made it into the server resp
		require.NotNil(srResp.GeneratedTime)
		require.True(srResp.GeneratedTime.AsTime().Equal(statusReportTs.AsTime()))
	})
}

func TestAppReleaseStatusReport(t *testing.T) {
	t.Run("with no status implementation", func(t *testing.T) {
		ctx := context.Background()
		require := require.New(t)

		// Our mock platform, which must also implement Status
		mock := &mockReleaseStatus{}

		// Make our factory for platforms
		factory := TestFactory(t, component.ReleaseManagerType)
		TestFactoryRegister(t, factory, "test", mock)

		// Make our app
		app := TestApp(t, TestProject(t,
			WithConfig(config.TestConfig(t, testReleaseManagerConfig)),
			WithFactory(component.ReleaseManagerType, factory),
		), "test")
		client := app.client

		// We're using GetVersionInfoResponse here just because it is a proto message
		// that can be converted to an any.Any easily. We never use it, it's just to keep
		// the tests from blowing up with a nil reference.
		mockPluginArtifact := &pb.GetVersionInfoResponse{}

		anyval, err := ptypes.MarshalAny(mockPluginArtifact)
		require.NoError(err)

		aresp, err := client.UpsertPushedArtifact(ctx, &pb.UpsertPushedArtifactRequest{
			Artifact: serverptypes.TestValidArtifact(t, &pb.PushedArtifact{
				Artifact: &pb.Artifact{
					Artifact: anyval,
				},
			}),
		})
		require.NoError(err)

		resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
			Deployment: serverptypes.TestValidDeployment(t, &pb.Deployment{
				ArtifactId: aresp.Artifact.Id,
			}),
		})
		require.NoError(err)
		deploy := resp.Deployment

		releaseResp, err := client.UpsertRelease(ctx, &pb.UpsertReleaseRequest{
			Release: serverptypes.TestValidRelease(t, &pb.Release{
				DeploymentId: deploy.Id,
			}),
		})
		require.NoError(err)
		release := releaseResp.Release

		// not implemented
		mock.Status.On("StatusFunc").Return(nil)

		// Status Report
		srResp, err := app.ReleaseStatusReport(context.Background(), release)
		require.NoError(err)
		require.Nil(srResp)

	})
	t.Run("with status implementation on release", func(t *testing.T) {
		ctx := context.Background()
		require := require.New(t)

		// Our mock platform, which must also implement Status
		mock := &mockReleaseStatus{}

		// Make our factory for platforms
		factory := TestFactory(t, component.ReleaseManagerType)
		TestFactoryRegister(t, factory, "test", mock)

		// Make our app
		app := TestApp(t, TestProject(t,
			WithConfig(config.TestConfig(t, testReleaseManagerConfig)),
			WithFactory(component.ReleaseManagerType, factory),
		), "test")
		client := app.client

		// We're using GetVersionInfoResponse here just because it is a proto message
		// that can be converted to an any.Any easily. We never use it, it's just to keep
		// the tests from blowing up with a nil reference.
		mockPluginArtifact := &pb.GetVersionInfoResponse{}

		anyval, err := ptypes.MarshalAny(mockPluginArtifact)
		require.NoError(err)

		aresp, err := client.UpsertPushedArtifact(ctx, &pb.UpsertPushedArtifactRequest{
			Artifact: serverptypes.TestValidArtifact(t, &pb.PushedArtifact{
				Artifact: &pb.Artifact{
					Artifact: anyval,
				},
			}),
		})
		require.NoError(err)

		resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
			Deployment: serverptypes.TestValidDeployment(t, &pb.Deployment{
				ArtifactId: aresp.Artifact.Id,
			}),
		})
		require.NoError(err)
		deploy := resp.Deployment

		releaseResp, err := client.UpsertRelease(ctx, &pb.UpsertReleaseRequest{
			Release: serverptypes.TestValidRelease(t, &pb.Release{
				DeploymentId: deploy.Id,
			}),
		})
		require.NoError(err)
		release := releaseResp.Release

		statusReportTs := timestamppb.Now()
		mock.Status.On("StatusFunc").Return(func(context.Context) (*sdk.StatusReport, error) {
			return &sdk.StatusReport{
				GeneratedTime: statusReportTs,
			}, nil
		})

		// Status Report
		srResp, err := app.ReleaseStatusReport(context.Background(), release)
		statusReport := &sdk.StatusReport{}
		anypb.UnmarshalTo(srResp.StatusReport, statusReport, proto.UnmarshalOptions{})
		require.NoError(err)
		require.NotNil(srResp.StatusReport)
		require.NotNil(statusReport.Health)

		// Verify that we have a Target of the right type with the right id
		require.IsType(srResp.TargetId, &pb.StatusReport_ReleaseId{})
		require.Equal(srResp.TargetId.(*pb.StatusReport_ReleaseId).ReleaseId, release.Id)

		// Verify that the status report timestamp made it into the server resp
		require.NotNil(srResp.GeneratedTime)
		require.True(srResp.GeneratedTime.AsTime().Equal(statusReportTs.AsTime()))
	})
}

type mockPlatformStatus struct {
	componentmocks.Platform
	componentmocks.Status
}

type mockReleaseStatus struct {
	componentmocks.ReleaseManager
	componentmocks.Status
}
