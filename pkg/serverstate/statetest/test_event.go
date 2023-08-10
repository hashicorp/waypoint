// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package statetest

import (
	"context"
	"github.com/hashicorp/waypoint/internal/pkg/jsonpb"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/hashicorp/waypoint/pkg/serverstate"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func init() {
	tests["event"] = []testFunc{
		TestEvent,
	}
}

func TestEvent(t *testing.T, factory Factory, restartF RestartFactory) {
	ctx := context.Background()

	t.Run("Basic put all types and check pagination", func(t *testing.T) {
		require := require.New(t)
		s := factory(t)
		defer s.Close()

		ws := &pb.Ref_Workspace{
			Workspace: "default",
		}

		// Write project
		refProj := &pb.Ref_Project{Project: "test"}
		require.NoError(s.ProjectPut(ctx, serverptypes.TestProject(t, &pb.Project{
			Name: refProj.Project,
		})))
		refApp := &pb.Ref_Application{
			Application: "test",
			Project:     refProj.Project,
		}
		_, err := s.AppPut(ctx, serverptypes.TestApplication(t, &pb.Application{
			Name:    refApp.Application,
			Project: refProj,
		}))
		require.NoError(err)

		// Put Build
		build := serverptypes.TestBuild(t, &pb.Build{
			Id:          "test",
			Application: refApp,
			Workspace:   ws,
			Status: &pb.Status{
				State:     pb.Status_SUCCESS,
				StartTime: timestamppb.Now(),
			},
		})
		require.NoError(s.BuildPut(ctx, false, build))

		var commit string
		if build.Preload.JobDataSourceRef != nil {
			commit = build.Preload.JobDataSourceRef.Ref.(*pb.Job_DataSource_Ref_Git).Git.Commit
		}
		buildDataSubset := &pb.UI_EventBuild{
			BuildId:   build.Id,
			Sequence:  build.Sequence,
			Component: build.Component,
			Workspace: build.Workspace,
			Status:    build.Status,
			Commit:    commit,
		}

		var buildBytes []byte
		buildBytes, err = jsonpb.Marshal(buildDataSubset)
		require.NoError(err)

		buildEvent := &serverstate.Event{
			EventType:      "build",
			Application:    refApp,
			EventData:      buildBytes,
			EventTimestamp: time.Now(),
		}

		require.NoError(s.EventPut(ctx, buildEvent))

		pt := timestamppb.Now()

		s.ArtifactPut(ctx, false, &pb.PushedArtifact{
			Application: refApp,
			Workspace:   ws,
			Sequence:    0,
			Id:          "test",
			Status: &pb.Status{
				State:        pb.Status_SUCCESS,
				StartTime:    pt,
				CompleteTime: pt,
			},
			BuildId:      "test",
			Labels:       nil,
			TemplateData: nil,
			Build:        build,
		})
		// Put Deployment
		dep := &pb.Deployment{
			Id:          "test",
			Application: refApp,
			Workspace:   ws,
			Status: &pb.Status{
				State:        pb.Status_SUCCESS,
				StartTime:    pt,
				CompleteTime: pt,
			},
			ArtifactId: "test",
		}
		require.NoError(s.DeploymentPut(ctx, false, dep))

		depDataSubset := &pb.UI_EventDeployment{
			DeploymentId:  dep.Id,
			Sequence:      dep.Sequence,
			Component:     dep.Component,
			Workspace:     dep.Workspace,
			BuildSequence: build.Sequence,
			Status:        dep.Status,
		}

		var depBytes []byte
		depBytes, err = jsonpb.Marshal(depDataSubset)
		require.NoError(err)

		depEvent := &serverstate.Event{
			EventType:      "deployment",
			Application:    refApp,
			EventData:      depBytes,
			EventTimestamp: time.Now(),
		}

		require.NoError(s.EventPut(ctx, depEvent))

		// Put Release
		release := &pb.Release{
			Id:          "test",
			Application: refApp,
			Workspace:   ws,
			Status: &pb.Status{
				State:        pb.Status_SUCCESS,
				StartTime:    pt,
				CompleteTime: pt,
			},
			DeploymentId: dep.Id,
		}
		require.NoError(s.ReleasePut(ctx, false, release))

		releaseDataSubset := &pb.UI_EventRelease{
			ReleaseId:          release.Id,
			Sequence:           release.Sequence,
			Component:          release.Component,
			Workspace:          release.Workspace,
			Status:             release.Status,
			DeploymentSequence: dep.Sequence,
		}

		var releaseBytes []byte
		releaseBytes, err = jsonpb.Marshal(releaseDataSubset)
		require.NoError(err)

		releaseEvent := &serverstate.Event{
			EventType:      "release",
			Application:    refApp,
			EventData:      releaseBytes,
			EventTimestamp: time.Now(),
		}

		require.NoError(s.EventPut(ctx, releaseEvent))

		// check
		resp, _, err := s.EventListBundles(ctx, &pb.UI_ListEventsRequest{
			Application: refApp,
			Workspace:   ws,
			Pagination: &pb.PaginationRequest{
				PageSize:          3,
				NextPageToken:     "",
				PreviousPageToken: "",
			},
		})
		require.NoError(err)
		require.Len(resp, 3)
	})
}
