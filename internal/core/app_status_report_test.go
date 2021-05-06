package core

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	componentmocks "github.com/hashicorp/waypoint-plugin-sdk/component/mocks"
	sdk "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	"github.com/hashicorp/waypoint/internal/config"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
	"github.com/stretchr/testify/require"
)

func TestAppStatusReport(t *testing.T) {
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
		srResp, statusReport, err := app.StatusReport(context.Background(), deploy, nil)
		require.NoError(err)
		require.Nil(statusReport)
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

		mock.Status.On("StatusFunc").Return(func(context.Context) (*sdk.StatusReport, error) {
			return &sdk.StatusReport{}, nil
		})

		// Status Report
		_, statusReport, err := app.StatusReport(context.Background(), deploy, nil)
		require.NoError(err)
		require.NotNil(statusReport)
		require.NotNil(statusReport.Health)

	})
}

type mockPlatformStatus struct {
	componentmocks.Platform
	componentmocks.Status
}
