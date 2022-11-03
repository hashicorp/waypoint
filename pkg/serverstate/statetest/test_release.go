package statetest

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

func init() {
	tests["release"] = []testFunc{
		TestRelease,
	}
}

func TestRelease(t *testing.T, factory Factory, restartF RestartFactory) {
	ctx := context.Background()
	t.Run("CRUD operations", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Write project
		ref := &pb.Ref_Project{Project: "foo"}
		require.NoError(s.ProjectPut(ctx, serverptypes.TestProject(t, &pb.Project{
			Name: ref.Project,
		})))

		// Has no apps
		{
			resp, err := s.ProjectGet(ctx, ref)
			require.NoError(err)
			require.NotNil(resp)
			require.Empty(resp.Applications)
		}

		app := &pb.Ref_Application{
			Project:     ref.Project,
			Application: "testapp",
		}

		ws := &pb.Ref_Workspace{
			Workspace: "default",
		}

		// Add
		err := s.ReleasePut(ctx, false, serverptypes.TestRelease(t, &pb.Release{
			Id:          "d1",
			Application: app,
			Workspace:   ws,
			Status: &pb.Status{
				State:     pb.Status_SUCCESS,
				StartTime: timestamppb.Now(),
			},
		}))
		require.NoError(err)

		// Can read
		{
			resp, err := s.ReleaseGet(ctx, &pb.Ref_Operation{
				Target: &pb.Ref_Operation_Id{
					Id: "d1",
				},
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Can read latest
		{
			resp, err := s.ReleaseLatest(ctx, app, &pb.Ref_Workspace{Workspace: "default"})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Update
		ts := timestamppb.Now()
		err = s.ReleasePut(ctx, true, serverptypes.TestRelease(t, &pb.Release{
			Id:          "d1",
			Application: app,
			Workspace:   ws,
			Status: &pb.Status{
				State:        pb.Status_SUCCESS,
				StartTime:    timestamppb.Now(),
				CompleteTime: ts,
			},
		}))
		require.NoError(err)

		{
			resp, err := s.ReleaseGet(ctx, &pb.Ref_Operation{
				Target: &pb.Ref_Operation_Id{
					Id: "d1",
				},
			})
			require.NoError(err)
			require.NotNil(resp)

			require.Equal(ts.AsTime(), resp.Status.CompleteTime.AsTime())
		}

		// Add another in another workspace, and see Latest change
		err = s.ReleasePut(ctx, false, serverptypes.TestRelease(t, &pb.Release{
			Id:          "d2",
			Application: app,
			Workspace:   &pb.Ref_Workspace{Workspace: "non-default"},
			Status: &pb.Status{
				State:        pb.Status_SUCCESS,
				StartTime:    timestamppb.New(ts.AsTime().Add(time.Second)),
				CompleteTime: timestamppb.New(ts.AsTime().Add(2 * time.Second)),
			},
		}))
		require.NoError(err)

		// Get by workspace
		{
			resp, err := s.ReleaseLatest(ctx, app, &pb.Ref_Workspace{Workspace: "non-default"})
			require.NoError(err)
			require.NotNil(resp)
			require.Equal("d2", resp.Id)
		}

		// Get with no workspace
		{
			resp, err := s.ReleaseLatest(ctx, app, nil)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal("d2", resp.Id)
		}

		// Add another release in another workspace, and ensure
		// ReleaseLatest (with no workspace filter) returns the latest one.

		{
			resp, err := s.ReleaseList(ctx, app)
			require.NoError(err)

			require.Len(resp, 2)
		}

		/*
				TODO: singleprocess/state's usage of Desc is broken.
			{
				resp, err := s.ReleaseList(app, serverstate.ListWithOrder(&pb.OperationOrder{
					Order: pb.OperationOrder_START_TIME,
					Desc:  false,
					Limit: 1,
				}))
				require.NoError(err)

				require.Len(resp, 1)

				require.Equal("d1", resp[0].Id)
			}
		*/

		{
			resp, err := s.ReleaseList(ctx, app, serverstate.ListWithOrder(&pb.OperationOrder{
				Order: pb.OperationOrder_START_TIME,
				Desc:  true,
				Limit: 1,
			}))
			require.NoError(err)

			require.Len(resp, 1)

			require.Equal("d2", resp[0].Id)
		}

		err = s.ReleasePut(ctx, false, serverptypes.TestRelease(t, &pb.Release{
			Id:          "d3",
			Application: app,
			Workspace:   ws,
			Status: &pb.Status{
				State:     pb.Status_ERROR,
				StartTime: timestamppb.Now(),
			},
		}))
		require.NoError(err)

		{
			resp, err := s.ReleaseList(ctx, app)
			require.NoError(err)

			require.Len(resp, 3)
		}

		{
			resp, err := s.ReleaseList(ctx, app,
				serverstate.ListWithOrder(&pb.OperationOrder{
					Order: pb.OperationOrder_START_TIME,
					Desc:  true,
				}),
				serverstate.ListWithStatusFilter(&pb.StatusFilter{
					Filters: []*pb.StatusFilter_Filter{
						{
							Filter: &pb.StatusFilter_Filter_State{
								State: pb.Status_ERROR,
							},
						},
					},
				}),
			)
			require.NoError(err)

			require.Len(resp, 1)

			require.Equal("d3", resp[0].Id)
		}

		{
			resp, err := s.ReleaseList(ctx, app,
				serverstate.ListWithOrder(&pb.OperationOrder{
					Order: pb.OperationOrder_START_TIME,
					Desc:  true,
				}),
				serverstate.ListWithStatusFilter(
					&pb.StatusFilter{
						Filters: []*pb.StatusFilter_Filter{
							{
								Filter: &pb.StatusFilter_Filter_State{
									State: pb.Status_ERROR,
								},
							},
						},
					},
					&pb.StatusFilter{
						Filters: []*pb.StatusFilter_Filter{
							{
								Filter: &pb.StatusFilter_Filter_State{
									State: pb.Status_SUCCESS,
								},
							},
						},
					},
				),
			)
			require.NoError(err)

			require.Len(resp, 3)
		}
	})

}
