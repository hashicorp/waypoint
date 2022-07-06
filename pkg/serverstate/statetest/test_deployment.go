package statetest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

func init() {
	tests["deployment"] = []testFunc{
		TestDeployment,
		TestDeploymentListFilter,
	}
}

func TestDeployment(t *testing.T, factory Factory, restartF RestartFactory) {
	t.Run("CRUD operations", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Write project
		ref := &pb.Ref_Project{Project: "foo"}
		require.NoError(s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: ref.Project,
		})))

		// Has no apps
		{
			resp, err := s.ProjectGet(ref)
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
		err := s.DeploymentPut(false, serverptypes.TestDeployment(t, &pb.Deployment{
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
			resp, err := s.DeploymentGet(&pb.Ref_Operation{
				Target: &pb.Ref_Operation_Id{
					Id: "d1",
				},
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Can read latest
		{
			resp, err := s.DeploymentLatest(app, &pb.Ref_Workspace{Workspace: "default"})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Update
		ts := timestamppb.Now()
		err = s.DeploymentPut(true, serverptypes.TestDeployment(t, &pb.Deployment{
			Id:          "d1",
			Application: app,
			Workspace:   ws,
			Status: &pb.Status{
				State:        pb.Status_SUCCESS,
				StartTime:    ts,
				CompleteTime: ts,
			},
		}))
		require.NoError(err)

		{
			resp, err := s.DeploymentGet(&pb.Ref_Operation{
				Target: &pb.Ref_Operation_Id{
					Id: "d1",
				},
			})
			require.NoError(err)
			require.NotNil(resp)

			require.Equal(ts.AsTime(), resp.Status.CompleteTime.AsTime())
		}

		// Add another and see Latset change
		// Add
		err = s.DeploymentPut(false, serverptypes.TestDeployment(t, &pb.Deployment{
			Id:          "d2",
			Application: app,
			Workspace:   ws,
			Status: &pb.Status{
				State:        pb.Status_SUCCESS,
				StartTime:    timestamppb.New(ts.AsTime().Add(time.Second)),
				CompleteTime: timestamppb.New(ts.AsTime().Add(2 * time.Second)),
			},
		}))
		require.NoError(err)

		{
			resp, err := s.DeploymentLatest(app, &pb.Ref_Workspace{Workspace: "default"})
			require.NoError(err)
			require.NotNil(resp)
			require.Equal("d2", resp.Id)
		}

		{
			resp, err := s.DeploymentList(app)
			require.NoError(err)

			require.Len(resp, 2)
		}

		/*
				TODO: singleprocess/state's usage of Desc is broken.
			{
				resp, err := s.DeploymentList(app, serverstate.ListWithOrder(&pb.OperationOrder{
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
			resp, err := s.DeploymentList(app, serverstate.ListWithOrder(&pb.OperationOrder{
				Order: pb.OperationOrder_START_TIME,
				Desc:  true,
				Limit: 1,
			}))
			require.NoError(err)

			require.Len(resp, 1)

			require.Equal("d2", resp[0].Id)
		}

		err = s.DeploymentPut(false, serverptypes.TestDeployment(t, &pb.Deployment{
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
			resp, err := s.DeploymentList(app)
			require.NoError(err)

			require.Len(resp, 3)
		}

		{
			resp, err := s.DeploymentList(app,
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
			resp, err := s.DeploymentList(app,
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

func TestDeploymentListFilter(t *testing.T, f Factory, rf RestartFactory) {
	require := require.New(t)

	s := f(t)
	defer s.Close()

	// Write project
	ref := &pb.Ref_Project{Project: "foo"}
	require.NoError(s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
		Name: ref.Project,
	})))

	app := &pb.Ref_Application{
		Project:     ref.Project,
		Application: "testapp",
	}

	ws := &pb.Ref_Workspace{
		Workspace: "default",
	}

	// Add a destroyed deployment
	require.NoError(
		s.DeploymentPut(false, serverptypes.TestDeployment(t, &pb.Deployment{
			Id:          "destroyed",
			Application: app,
			Workspace:   ws,
			State:       pb.Operation_DESTROYED,
			Status: &pb.Status{
				StartTime: timestamppb.Now(),
			},
		})),
	)

	require.NoError(
		s.DeploymentPut(false, serverptypes.TestDeployment(t, &pb.Deployment{
			Id:          "created",
			Application: app,
			Workspace:   ws,
			State:       pb.Operation_CREATED,
			Status: &pb.Status{
				StartTime: timestamppb.Now(),
			},
		})),
	)

	resp, err := s.DeploymentList(app, serverstate.ListWithPhysicalState(pb.Operation_CREATED))
	require.NoError(err)
	require.Len(resp, 1)
	require.Equal("created", resp[0].Id)

	resp, err = s.DeploymentList(app, serverstate.ListWithPhysicalState(pb.Operation_DESTROYED))
	require.NoError(err)
	require.Len(resp, 1)
	require.Equal("destroyed", resp[0].Id)

	resp, err = s.DeploymentList(app, serverstate.ListWithPhysicalState(pb.Operation_PENDING))
	require.NoError(err)
	require.Len(resp, 0)
}
