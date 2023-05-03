package handlertest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func init() {
	tests["trigger"] = []testFunc{
		TestServiceTrigger,
		TestServiceTrigger_GetTrigger,
		TestServiceTrigger_ListTriggersSimple,
		TestServiceTrigger_ListTriggersWithFilters,
		TestServiceTrigger_DeleteTrigger,
		TestServiceTrigger_RunTrigger,
	}
}

func TestServiceTrigger(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	type Req = pb.UpsertTriggerRequest

	t.Run("create and update", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertTrigger(ctx, &pb.UpsertTriggerRequest{
			Trigger: serverptypes.TestValidTrigger(t, nil),
		})
		require.NoError(err)
		require.NotNil(resp)
		result := resp.Trigger
		require.NotEmpty(result.Id)

		// Let's write some data
		testName := "TestyTest"
		result.Name = testName
		resp, err = client.UpsertTrigger(ctx, &pb.UpsertTriggerRequest{
			Trigger: result,
		})
		require.NoError(err)
		require.NotNil(resp)
		result = resp.Trigger
		require.Equal(result.Name, testName)
	})

	t.Run("upsert fails if invalid Name for path given", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertTrigger(ctx, &pb.UpsertTriggerRequest{
			Trigger: serverptypes.TestValidTrigger(t, nil),
		})
		require.NoError(err)
		require.NotNil(resp)
		result := resp.Trigger
		require.NotEmpty(result.Id)

		// Let's write some bad data
		testName := "../TestyTest" // this is an invalid path
		result.Name = testName
		resp, err = client.UpsertTrigger(ctx, &pb.UpsertTriggerRequest{
			Trigger: result,
		})
		require.Error(err)
		require.Nil(resp)
	})

	t.Run("create uses default workspace if unset", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertTrigger(ctx, &pb.UpsertTriggerRequest{
			Trigger: &pb.Trigger{Project: &pb.Ref_Project{Project: "test_proj"}},
		})
		require.NoError(err)
		require.NotNil(resp)
		result := resp.Trigger
		require.NotEmpty(result.Id)
		require.NotEmpty(result.Workspace)
		require.Equal(result.Workspace.Workspace, "default")
	})

	t.Run("errors if no project defined", func(t *testing.T) {
		require := require.New(t)

		// create should error with no project defined
		resp, err := client.UpsertTrigger(ctx, &pb.UpsertTriggerRequest{
			Trigger: &pb.Trigger{},
		})
		require.Error(err)
		require.Nil(resp)
	})

	t.Run("update non-existent creates a new trigger", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertTrigger(ctx, &Req{
			Trigger: serverptypes.TestValidTrigger(t, &pb.Trigger{
				Id: "nope",
			}),
		})
		require.NoError(err)
		require.NotNil(resp)
		require.Equal(resp.Trigger.Id, "nope")
	})
}

func TestServiceTrigger_GetTrigger(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	resp, err := client.UpsertTrigger(ctx, &pb.UpsertTriggerRequest{
		Trigger: serverptypes.TestValidTrigger(t, nil),
	})
	require.NoError(t, err)
	triggerId := resp.Trigger.Id

	type Req = pb.UpsertTriggerRequest

	t.Run("get existing", func(t *testing.T) {
		require := require.New(t)

		// Get, should return a trigger
		resp, err := client.GetTrigger(ctx, &pb.GetTriggerRequest{
			Ref: &pb.Ref_Trigger{
				Id: resp.Trigger.Id,
			},
		})
		require.NoError(err)
		require.NotNil(resp.Trigger)
		require.NotEmpty(resp.Trigger.Id)
		require.Equal(triggerId, resp.Trigger.Id)
	})

	t.Run("get non-existing", func(t *testing.T) {
		require := require.New(t)

		// get, should fail
		resp, err := client.GetTrigger(ctx, &pb.GetTriggerRequest{
			Ref: &pb.Ref_Trigger{
				Id: "nope",
			},
		})
		require.Error(err)
		require.Nil(resp)
		st, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, st.Code())
	})
}

func TestServiceTrigger_ListTriggersSimple(t *testing.T, factory Factory) {
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	client, _ := factory(t)

	_, err := client.UpsertTrigger(ctx, &pb.UpsertTriggerRequest{
		Trigger: serverptypes.TestValidTrigger(t, nil),
	})
	require.NoError(err)

	_, err = client.UpsertTrigger(ctx, &pb.UpsertTriggerRequest{
		Trigger: serverptypes.TestValidTrigger(t, nil),
	})
	require.NoError(err)

	_, err = client.UpsertTrigger(ctx, &pb.UpsertTriggerRequest{
		Trigger: serverptypes.TestValidTrigger(t, nil),
	})
	require.NoError(err)

	t.Run("list", func(t *testing.T) {
		respList, err := client.ListTriggers(ctx, &pb.ListTriggerRequest{})
		require.NoError(err)
		require.Equal(len(respList.Triggers), 3)
	})
}

func TestServiceTrigger_ListTriggersWithFilters(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	_, err := client.UpsertTrigger(ctx, &pb.UpsertTriggerRequest{
		Trigger: serverptypes.TestValidTrigger(t, nil),
	})
	_, err = client.UpsertTrigger(ctx, &pb.UpsertTriggerRequest{
		Trigger: serverptypes.TestValidTrigger(t, nil),
	})
	_, err = client.UpsertTrigger(ctx, &pb.UpsertTriggerRequest{
		Trigger: serverptypes.TestValidTrigger(t, nil),
	})

	t.Run("list default workspace triggers", func(t *testing.T) {
		require := require.New(t)

		respList, err := client.ListTriggers(ctx, &pb.ListTriggerRequest{
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		})
		require.NoError(err)
		require.Equal(len(respList.Triggers), 3)
	})

	t.Run("list non-existent workspace triggers", func(t *testing.T) {
		require := require.New(t)

		respList, err := client.ListTriggers(ctx, &pb.ListTriggerRequest{
			Workspace: &pb.Ref_Workspace{Workspace: "fake"},
		})
		require.NoError(err)
		require.Equal(len(respList.Triggers), 0)
	})

	t.Run("list project triggers", func(t *testing.T) {
		require := require.New(t)

		respList, err := client.ListTriggers(ctx, &pb.ListTriggerRequest{
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
			Project:   &pb.Ref_Project{Project: "p_test"},
		})
		require.NoError(err)
		require.Equal(len(respList.Triggers), 3)
	})

	t.Run("list app triggers", func(t *testing.T) {
		require := require.New(t)

		respList, err := client.ListTriggers(ctx, &pb.ListTriggerRequest{
			Workspace:   &pb.Ref_Workspace{Workspace: "default"},
			Project:     &pb.Ref_Project{Project: "p_test"},
			Application: &pb.Ref_Application{Project: "p_test", Application: "a_test"},
		})
		require.NoError(err)
		require.Equal(len(respList.Triggers), 3)
	})

	t.Run("filter to one app", func(t *testing.T) {
		require := require.New(t)

		_, err = client.UpsertTrigger(ctx, &pb.UpsertTriggerRequest{
			Trigger: &pb.Trigger{
				Project: &pb.Ref_Project{
					Project: "secret_project",
				},
				Application: &pb.Ref_Application{
					Application: "another_one",
					Project:     "secret_project",
				},
				Workspace: &pb.Ref_Workspace{
					Workspace: "staging",
				},
			},
		})
		require.NoError(err)

		respList, err := client.ListTriggers(ctx, &pb.ListTriggerRequest{
			Workspace:   &pb.Ref_Workspace{Workspace: "staging"},
			Project:     &pb.Ref_Project{Project: "secret_project"},
			Application: &pb.Ref_Application{Project: "secret_project", Application: "another_one"},
		})
		require.NoError(err)
		require.Equal(1, len(respList.Triggers))
	})

	t.Run("filter on tags", func(t *testing.T) {
		require := require.New(t)

		_, err = client.UpsertTrigger(ctx, &pb.UpsertTriggerRequest{
			Trigger: &pb.Trigger{
				Project: &pb.Ref_Project{
					Project: "secret_project",
				},
				Application: &pb.Ref_Application{
					Application: "another_one",
					Project:     "secret_project",
				},
				Workspace: &pb.Ref_Workspace{
					Workspace: "staging",
				},
				Tags: []string{"prod", "test"},
			},
		})
		require.NoError(err)

		respList, err := client.ListTriggers(ctx, &pb.ListTriggerRequest{
			Workspace: &pb.Ref_Workspace{Workspace: "staging"},
			Tags:      []string{"prod"},
		})
		require.NoError(err)
		require.Equal(1, len(respList.Triggers))
	})

	t.Run("filter on missing tags returns nothing", func(t *testing.T) {
		require := require.New(t)

		_, err = client.UpsertTrigger(ctx, &pb.UpsertTriggerRequest{
			Trigger: &pb.Trigger{
				Project: &pb.Ref_Project{
					Project: "secret_project",
				},
				Application: &pb.Ref_Application{
					Application: "another_one",
					Project:     "secret_project",
				},
				Workspace: &pb.Ref_Workspace{
					Workspace: "staging",
				},
				Tags: []string{"prod", "test"},
			},
		})
		require.NoError(err)

		respList, err := client.ListTriggers(ctx, &pb.ListTriggerRequest{
			Workspace: &pb.Ref_Workspace{Workspace: "staging"},
			Tags:      []string{"pikachu"},
		})
		require.NoError(err)
		require.Equal(0, len(respList.Triggers))
	})

	t.Run("filter with missing workspace on app returns error", func(t *testing.T) {
		require := require.New(t)

		_, err := client.ListTriggers(ctx, &pb.ListTriggerRequest{
			Application: &pb.Ref_Application{Project: "secret_project", Application: "another_one"},
		})
		require.Error(err)
	})
}

func TestServiceTrigger_DeleteTrigger(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	resp, err := client.UpsertTrigger(ctx, &pb.UpsertTriggerRequest{
		Trigger: serverptypes.TestValidTrigger(t, nil),
	})
	require.NoError(t, err)
	triggerId := resp.Trigger.Id

	type Req = pb.UpsertTriggerRequest

	t.Run("get existing then delete", func(t *testing.T) {
		require := require.New(t)

		// Get, should return a trigger
		resp, err := client.GetTrigger(ctx, &pb.GetTriggerRequest{
			Ref: &pb.Ref_Trigger{
				Id: resp.Trigger.Id,
			},
		})
		require.NoError(err)
		require.NotNil(resp.Trigger)
		require.NotEmpty(resp.Trigger.Id)
		require.Equal(triggerId, resp.Trigger.Id)

		_, err = client.DeleteTrigger(ctx, &pb.DeleteTriggerRequest{
			Ref: &pb.Ref_Trigger{
				Id: triggerId,
			},
		})
		require.NoError(err)

		// get, should fail
		resp, err = client.GetTrigger(ctx, &pb.GetTriggerRequest{
			Ref: &pb.Ref_Trigger{
				Id: triggerId,
			},
		})
		require.Error(err)
		require.Nil(resp)
		st, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, st.Code())
	})

	t.Run("delete non-existing", func(t *testing.T) {
		require := require.New(t)

		resp, err := client.DeleteTrigger(ctx, &pb.DeleteTriggerRequest{
			Ref: &pb.Ref_Trigger{
				Id: "nope",
			},
		})
		require.NoError(err)
		require.NotNil(resp)
	})
}

func TestServiceTrigger_RunTrigger(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	t.Run("running a missing trigger", func(t *testing.T) {
		require := require.New(t)

		// try to run a trigger, require error
		resp, err := client.RunTrigger(ctx, &pb.RunTriggerRequest{
			Ref: &pb.Ref_Trigger{
				Id: "nope",
			},
		})

		require.Error(err)
		require.Nil(resp)
	})

	t.Run("running a trigger when project is missing", func(t *testing.T) {
		require := require.New(t)

		respTrigger, err := client.UpsertTrigger(ctx, &pb.UpsertTriggerRequest{
			Trigger: &pb.Trigger{
				Id:   "test",
				Name: "build-trigger",
				Operation: &pb.Trigger_Build{
					Build: &pb.Job_BuildOp{
						DisablePush: false,
					},
				},
				Project: &pb.Ref_Project{
					Project: "secret_project",
				},
				Application: &pb.Ref_Application{
					Application: "another_one",
					Project:     "secret_project",
				},
				Workspace: &pb.Ref_Workspace{
					Workspace: "staging",
				},
			},
		})
		require.NoError(err)
		require.NotNil(respTrigger)

		// try to run a trigger, require error
		resp, err := client.RunTrigger(ctx, &pb.RunTriggerRequest{
			Ref: &pb.Ref_Trigger{
				Id: respTrigger.Trigger.Id,
			},
		})

		require.Error(err)
		require.Nil(resp)
		st, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, st.Code())
	})

	t.Run("queues a single registered trigger URL on an application", func(t *testing.T) {
		require := require.New(t)

		// Create a project with an application
		respProj, err := client.UpsertProject(ctx, &pb.UpsertProjectRequest{
			Project: serverptypes.TestProject(t, &pb.Project{
				Name: "secret_project",
				DataSource: &pb.Job_DataSource{
					Source: &pb.Job_DataSource_Local{
						Local: &pb.Job_Local{},
					},
				},
				Applications: []*pb.Application{
					{
						Project: &pb.Ref_Project{Project: "secret_project"},
						Name:    "another_one",
					},
				},
			}),
		})
		require.NoError(err)
		require.NotNil(respProj)

		respTrigger, err := client.UpsertTrigger(ctx, &pb.UpsertTriggerRequest{
			Trigger: &pb.Trigger{
				Id:   "test",
				Name: "build-trigger",
				Operation: &pb.Trigger_Build{
					Build: &pb.Job_BuildOp{
						DisablePush: false,
					},
				},
				Project: &pb.Ref_Project{
					Project: "secret_project",
				},
				Application: &pb.Ref_Application{
					Application: "another_one",
					Project:     "secret_project",
				},
				Workspace: &pb.Ref_Workspace{
					Workspace: "staging",
				},
			},
		})
		require.NoError(err)
		require.NotNil(respTrigger)

		resp, err := client.RunTrigger(ctx, &pb.RunTriggerRequest{
			Ref: &pb.Ref_Trigger{
				Id: respTrigger.Trigger.Id,
			},
		})
		require.NoError(err)
		require.NotNil(resp)
		require.Len(resp.JobIds, 1)
	})

	t.Run("queues multiple registered trigger URLs for all apps in a project-scoped request", func(t *testing.T) {
		require := require.New(t)

		// Create a project with an application
		respProj, err := client.UpsertProject(ctx, &pb.UpsertProjectRequest{
			Project: serverptypes.TestProject(t, &pb.Project{
				Name: "multi-app",
				DataSource: &pb.Job_DataSource{
					Source: &pb.Job_DataSource_Local{
						Local: &pb.Job_Local{},
					},
				},
			}),
		})
		require.NoError(err)
		require.NotNil(respProj)

		_, err = client.UpsertApplication(ctx, &pb.UpsertApplicationRequest{
			Project: &pb.Ref_Project{Project: "multi-app"},
			Name:    "app-one",
		})
		require.NoError(err)
		_, err = client.UpsertApplication(ctx, &pb.UpsertApplicationRequest{
			Project: &pb.Ref_Project{Project: "multi-app"},
			Name:    "app-two",
		})
		require.NoError(err)
		_, err = client.UpsertApplication(ctx, &pb.UpsertApplicationRequest{
			Project: &pb.Ref_Project{Project: "multi-app"},
			Name:    "app-three",
		})
		require.NoError(err)
		_, err = client.UpsertApplication(ctx, &pb.UpsertApplicationRequest{
			Project: &pb.Ref_Project{Project: "multi-app"},
			Name:    "app-four",
		})
		require.NoError(err)

		respTrigger, err := client.UpsertTrigger(ctx, &pb.UpsertTriggerRequest{
			Trigger: &pb.Trigger{
				Id:   "test",
				Name: "build-trigger",
				Operation: &pb.Trigger_Build{
					Build: &pb.Job_BuildOp{
						DisablePush: false,
					},
				},
				Project: &pb.Ref_Project{
					Project: "multi-app",
				},
				Workspace: &pb.Ref_Workspace{
					Workspace: "staging",
				},
			},
		})
		require.NoError(err)
		require.NotNil(respTrigger)

		err = nil
		resp, err := client.RunTrigger(ctx, &pb.RunTriggerRequest{
			Ref: &pb.Ref_Trigger{
				Id: respTrigger.Trigger.Id,
			},
		})
		require.NoError(err)
		require.NotNil(resp)
		require.Len(resp.JobIds, 4)
	})
}
