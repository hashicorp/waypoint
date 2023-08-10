// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

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
	tests["ui_release"] = []testFunc{
		TestServiceUI_Release_ListReleases,
	}
}

func TestServiceUI_Release_ListReleases(t *testing.T, factory Factory) {
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

	releaseResp, err := client.UpsertRelease(ctx, &pb.UpsertReleaseRequest{
		Release: serverptypes.TestValidRelease(t, &pb.Release{
			DeploymentId: deployResp.Deployment.Id,
			Application: &pb.Ref_Application{
				Application: "apple-app",
				Project:     "Example",
			},
		}),
	})
	require.NoError(t, err)

	sr1, err := client.UpsertStatusReport(ctx, &pb.UpsertStatusReportRequest{
		StatusReport: serverptypes.TestValidStatusReport(t, &pb.StatusReport{
			TargetId: &pb.StatusReport_ReleaseId{
				ReleaseId: releaseResp.Release.Id,
			},
			Application: &pb.Ref_Application{
				Application: "apple-app",
				Project:     "Example",
			},
			GeneratedTime: timestamppb.New(time.Now().Add(-1 * time.Minute)),
		}),
	})
	require.NoError(t, err)

	type Req = pb.UI_ListReleasesRequest

	t.Run("list", func(t *testing.T) {
		require := require.New(t)
		releases, err := client.UI_ListReleases(ctx, &Req{
			Application: releaseResp.Release.Application,
			Workspace:   releaseResp.Release.Workspace,
		})
		require.NoError(err)
		require.NotNil(releases)
		require.NotNil(releases.Releases)
		require.Equal(len(releases.Releases), 1)

		// Operation exists and is what we expect
		release := releases.Releases[0]
		require.NotNil(release.Release)
		require.Equal(release.Release.Id, releaseResp.Release.Id)

		// Status report exists and matches our operation
		require.NotEmpty(release.LatestStatusReport)
		require.IsType(release.LatestStatusReport.TargetId, &pb.StatusReport_ReleaseId{})
		require.Equal(release.LatestStatusReport.TargetId.(*pb.StatusReport_ReleaseId).ReleaseId, release.Release.Id)

		// Latest status report is what we inserted
		require.Equal(release.LatestStatusReport.Id, sr1.StatusReport.Id)
	})

	t.Run("list shows newest status report", func(t *testing.T) {
		// Insert another status report for the release, with a newer time
		sr2, err := client.UpsertStatusReport(ctx, &pb.UpsertStatusReportRequest{
			StatusReport: serverptypes.TestValidStatusReport(t, &pb.StatusReport{
				TargetId: &pb.StatusReport_ReleaseId{
					ReleaseId: releaseResp.Release.Id,
				},
				Application: &pb.Ref_Application{
					Application: "apple-app",
					Project:     "Example",
				},
				GeneratedTime: timestamppb.New(time.Now()),
			}),
		})
		require.NoError(t, err)

		require := require.New(t)
		releases, err := client.UI_ListReleases(ctx, &Req{
			Application: releaseResp.Release.Application,
			Workspace:   releaseResp.Release.Workspace,
		})
		require.NoError(err)
		require.NotNil(releases)
		require.NotNil(releases.Releases)
		require.Equal(len(releases.Releases), 1)

		// Operation exists and is what we expect
		release := releases.Releases[0]

		// Is the most recent status report for the release
		require.NotEmpty(release.LatestStatusReport)
		require.Equal(release.LatestStatusReport.Id, sr2.StatusReport.Id)
	})
}
