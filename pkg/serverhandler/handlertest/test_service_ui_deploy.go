package handlertest

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func init() {
	tests["ui_deploy"] = []testFunc{
		TestServiceUI_Deployment_ListDeployments,
	}
}

func TestServiceUI_Deployment_ListDeployments(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	// Create a project with an application
	respProj, err := client.UpsertProject(ctx, &pb.UpsertProjectRequest{
		Project: serverptypes.TestProject(t, &pb.Project{
			Name: "Example",
			DataSource: &pb.Job_DataSource{
				Source: &pb.Job_DataSource_Local{
					Local: &pb.Job_Local{},
				},
			},
			Applications: []*pb.Application{
				{
					Project: &pb.Ref_Project{Project: "Example"},
					Name:    "apple-app",
				},
			},
		}),
	})
	require.NoError(t, err)
	require.NotNil(t, respProj)

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

	deployResp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment: serverptypes.TestValidDeployment(t, &pb.Deployment{
			Component: &pb.Component{
				Name: "testapp",
			},
			Application: &pb.Ref_Application{
				Application: "apple-app",
				Project:     "Example",
			},
			ArtifactId: artifactresp.Artifact.Id,
		}),
	})
	require.NoError(t, err)
	require.NotNil(t, deployResp)

	sr1, err := client.UpsertStatusReport(ctx, &pb.UpsertStatusReportRequest{
		StatusReport: serverptypes.TestValidStatusReport(t, &pb.StatusReport{
			TargetId: &pb.StatusReport_DeploymentId{
				DeploymentId: deployResp.Deployment.Id,
			},
			Application: &pb.Ref_Application{
				Application: "apple-app",
				Project:     "Example",
			},
			GeneratedTime: timestamppb.New(time.Now().Add(-1 * time.Minute)),
		}),
	})
	require.NoError(t, err)

	type Req = pb.UI_ListDeploymentsRequest

	t.Run("list", func(t *testing.T) {
		require := require.New(t)
		deployments, err := client.UI_ListDeployments(ctx, &Req{
			Application: deployResp.Deployment.Application,
			Workspace:   deployResp.Deployment.Workspace,
		})
		require.NoError(err)
		require.NotNil(deployments)
		require.NotNil(deployments.Deployments)
		require.Equal(len(deployments.Deployments), 1)

		// Operation exists and is what we expect
		deployment := deployments.Deployments[0]
		require.NotNil(deployment.Deployment)
		require.Equal(deployment.Deployment.Id, deployResp.Deployment.Id)

		// Status report exists and matches our operation
		require.NotNil(deployment.LatestStatusReport)
		require.IsType(deployment.LatestStatusReport.TargetId, &pb.StatusReport_DeploymentId{})
		require.Equal(deployment.LatestStatusReport.TargetId.(*pb.StatusReport_DeploymentId).DeploymentId, deployment.Deployment.Id)

		// Latest status report is what we inserted
		require.Equal(deployment.LatestStatusReport.Id, sr1.StatusReport.Id)
	})

	t.Run("list shows newest status report", func(t *testing.T) {
		// Insert another status report for the deployment, with a newer time
		sr2, err := client.UpsertStatusReport(ctx, &pb.UpsertStatusReportRequest{
			StatusReport: serverptypes.TestValidStatusReport(t, &pb.StatusReport{
				TargetId: &pb.StatusReport_DeploymentId{
					DeploymentId: deployResp.Deployment.Id,
				},
				Application: &pb.Ref_Application{
					Application: "apple-app",
					Project:     "Example",
				},
				GeneratedTime: timestamppb.Now(),
			}),
		})
		require.NoError(t, err)

		require := require.New(t)
		deployments, err := client.UI_ListDeployments(ctx, &Req{
			Application: deployResp.Deployment.Application,
			Workspace:   deployResp.Deployment.Workspace,
		})
		require.NoError(err)
		require.NotEmpty(deployments)
		require.NotEmpty(deployments.Deployments)
		require.Equal(len(deployments.Deployments), 1)

		// Operation exists and is what we expect
		deployment := deployments.Deployments[0]

		// Is the most recent status report for the deployment
		require.NotEmpty(deployment.LatestStatusReport)
		require.Equal(deployment.LatestStatusReport.Id, sr2.StatusReport.Id)
	})
}
