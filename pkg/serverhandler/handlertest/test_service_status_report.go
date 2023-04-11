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
	tests["status_report"] = []testFunc{
		TestServiceStatusReport,
		TestServiceStatusReport_GetStatusReport,
		TestServiceStatusReport_ListStatusReports,
		TestServiceStatusReport_ExpediteStatusReport,
	}
}

func TestServiceStatusReport(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	type Req = pb.UpsertStatusReportRequest

	t.Run("create and update", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertStatusReport(ctx, &pb.UpsertStatusReportRequest{
			StatusReport: serverptypes.TestValidStatusReport(t, nil),
		})
		require.NoError(err)
		require.NotNil(resp)
		result := resp.StatusReport
		require.NotEmpty(result.Id)

		// Let's write some data
		result.Status = server.NewStatus(pb.Status_RUNNING)
		resp, err = client.UpsertStatusReport(ctx, &Req{
			StatusReport: result,
		})
		require.NoError(err)
		require.NotNil(resp)
		result = resp.StatusReport
		require.NotNil(result.Status)
		require.Equal(pb.Status_RUNNING, result.Status.State)
	})

	t.Run("update non-existent", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertStatusReport(ctx, &Req{
			StatusReport: serverptypes.TestValidStatusReport(t, &pb.StatusReport{
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

func TestServiceStatusReport_GetStatusReport(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	statusReportResp, err := client.UpsertStatusReport(ctx, &pb.UpsertStatusReportRequest{
		StatusReport: serverptypes.TestValidStatusReport(t, nil),
	})
	require.NoError(t, err)

	type Req = pb.GetStatusReportRequest

	t.Run("get existing", func(t *testing.T) {
		require := require.New(t)

		// Get, should return a status report
		sp, err := client.GetStatusReport(ctx, &Req{
			Ref: &pb.Ref_Operation{
				Target: &pb.Ref_Operation_Id{Id: statusReportResp.StatusReport.Id},
			},
		})
		require.NoError(err)
		require.NotNil(sp)
		require.NotEmpty(sp.Id)
	})

	t.Run("get non-existing", func(t *testing.T) {
		require := require.New(t)

		// get, should fail
		resp, err := client.GetStatusReport(ctx, &Req{
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

func TestServiceStatusReport_ListStatusReports(t *testing.T, factory Factory) {
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

	artifact := serverptypes.TestValidArtifact(t, nil)

	artifactresp, err := client.UpsertPushedArtifact(ctx, &pb.UpsertPushedArtifactRequest{
		Artifact: artifact,
	})
	require.NoError(t, err)
	require.NotNil(t, artifactresp)

	deployment := &pb.Deployment{
		Component: &pb.Component{
			Name: "testapp",
		},
		Application: &pb.Ref_Application{
			Application: "apple-app",
			Project:     "Example",
		},
		ArtifactId: artifactresp.Artifact.Id,
	}

	deployResp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment: serverptypes.TestValidDeployment(t, deployment),
	})
	require.NoError(t, err)
	require.NotNil(t, deployResp)

	resp, err := client.UpsertStatusReport(ctx, &pb.UpsertStatusReportRequest{
		StatusReport: serverptypes.TestValidStatusReport(t, &pb.StatusReport{
			TargetId: &pb.StatusReport_DeploymentId{
				DeploymentId: deployResp.Deployment.Id,
			},
		}),
	})
	require.NoError(t, err)

	releaseResp, err := client.UpsertRelease(ctx, &pb.UpsertReleaseRequest{
		Release: serverptypes.TestValidRelease(t, &pb.Release{
			DeploymentId: deployResp.Deployment.Id,
		}),
	})
	require.NoError(t, err)

	releaseStatusResp, err := client.UpsertStatusReport(ctx, &pb.UpsertStatusReportRequest{
		StatusReport: serverptypes.TestValidStatusReport(t, &pb.StatusReport{
			TargetId: &pb.StatusReport_ReleaseId{
				ReleaseId: releaseResp.Release.Id,
			},
		}),
	})
	require.NoError(t, err)

	type Req = pb.ListStatusReportsRequest

	t.Run("list", func(t *testing.T) {
		require := require.New(t)

		// Get, should return a status report
		sr, err := client.ListStatusReports(ctx, &Req{
			Application: resp.StatusReport.Application,
		})
		require.NoError(err)
		require.NotEmpty(sr)
		require.Equal(len(sr.StatusReports), 2)

		// ensure each returned report matches the generated report id for both types
		for _, report := range sr.StatusReports {
			switch report.TargetId.(type) {
			case *pb.StatusReport_DeploymentId:
				require.Equal(report.Id, resp.StatusReport.Id)
			case *pb.StatusReport_ReleaseId:
				require.Equal(report.Id, releaseStatusResp.StatusReport.Id)
			}
		}
	})

	t.Run("list only deployment reports", func(t *testing.T) {
		require := require.New(t)

		// Get, should return a status report
		sr, err := client.ListStatusReports(ctx, &Req{
			Application: resp.StatusReport.Application,
			Target: &pb.ListStatusReportsRequest_Deployment{
				Deployment: &pb.Ref_Operation{
					Target: &pb.Ref_Operation_Id{
						Id: deployResp.Deployment.Id,
					},
				},
			},
		})
		require.NoError(err)
		require.NotEmpty(sr)
		require.Equal(len(sr.StatusReports), 1)
		require.Equal(sr.StatusReports[0].Id, resp.StatusReport.Id)
		require.Equal(sr.StatusReports[0].TargetId.(*pb.StatusReport_DeploymentId).DeploymentId, deployResp.Deployment.Id)
	})

	t.Run("list only release reports", func(t *testing.T) {
		require := require.New(t)

		// Get, should return a status report
		sr, err := client.ListStatusReports(ctx, &Req{
			Application: resp.StatusReport.Application,
			Target: &pb.ListStatusReportsRequest_Release{
				Release: &pb.Ref_Operation{
					Target: &pb.Ref_Operation_Id{
						Id: releaseResp.Release.Id,
					},
				},
			},
		})
		require.NoError(err)
		require.NotEmpty(sr)
		require.Equal(len(sr.StatusReports), 1)
		require.Equal(sr.StatusReports[0].Id, releaseStatusResp.StatusReport.Id)
		require.Equal(sr.StatusReports[0].TargetId.(*pb.StatusReport_ReleaseId).ReleaseId, releaseResp.Release.Id)
	})
}

func TestServiceStatusReport_ExpediteStatusReport(t *testing.T, factory Factory) {
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

	artifact := serverptypes.TestValidArtifact(t, nil)

	artifactresp, err := client.UpsertPushedArtifact(ctx, &pb.UpsertPushedArtifactRequest{
		Artifact: artifact,
	})
	require.NoError(t, err)
	require.NotNil(t, artifactresp)

	deployment := &pb.Deployment{
		Component: &pb.Component{
			Name: "testapp",
		},
		Application: &pb.Ref_Application{
			Application: "apple-app",
			Project:     "Example",
		},
		ArtifactId: artifactresp.Artifact.Id,
	}

	resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment: serverptypes.TestValidDeployment(t, deployment),
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	t.Run("Expedite Status Report", func(t *testing.T) {
		require := require.New(t)

		jobResp, err := client.ExpediteStatusReport(ctx, &pb.ExpediteStatusReportRequest{
			Target: &pb.ExpediteStatusReportRequest_Deployment{
				Deployment: &pb.Ref_Operation{
					Target: &pb.Ref_Operation_Id{Id: resp.Deployment.Id},
				},
			},
			Workspace: &pb.Ref_Workspace{
				Workspace: "default",
			},
		})
		require.NoError(err)
		require.NotEmpty(t, jobResp)
		require.NotNil(t, jobResp.JobId)
	})

	t.Run("Expedite Status Report with no workspace uses default and doesn't error", func(t *testing.T) {
		require := require.New(t)

		jobResp, err := client.ExpediteStatusReport(ctx, &pb.ExpediteStatusReportRequest{
			Target: &pb.ExpediteStatusReportRequest_Deployment{
				Deployment: &pb.Ref_Operation{
					Target: &pb.Ref_Operation_Id{Id: resp.Deployment.Id},
				},
			},
		})
		require.NoError(err)
		require.NotEmpty(t, jobResp)
		require.NotNil(t, jobResp.JobId)
	})

	require.NoError(t, err)
}
