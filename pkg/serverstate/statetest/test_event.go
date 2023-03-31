package statetest

import (
	"context"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
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
		require.NoError(s.EventPut(ctx, build))

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
		require.NoError(s.EventPut(ctx, dep))

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
		require.NoError(s.EventPut(ctx, release))

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
