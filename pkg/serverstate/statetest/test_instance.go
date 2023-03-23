// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statetest

import (
	"context"
	"testing"
	"time"

	"github.com/imdario/mergo"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/hashicorp/go-memdb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

func init() {
	tests["instance"] = []testFunc{
		TestInstance,
		TestInstanceByDeployment,
	}
}

func TestInstance(t *testing.T, factory Factory, restartF RestartFactory) {
	ctx := context.Background()
	testInstance := func(t *testing.T, v *serverstate.Instance) *serverstate.Instance {
		if v == nil {
			v = &serverstate.Instance{}
		}

		require.NoError(t, mergo.Merge(v, &serverstate.Instance{
			Id:           "A",
			DeploymentId: "B",
			Project:      "C",
			Application:  "D",
			Workspace:    "E",
		}))

		return v
	}

	t.Run("crud", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		ref := &pb.Ref_Project{Project: "foo"}
		require.NoError(s.ProjectPut(ctx, serverptypes.TestProject(t, &pb.Project{
			Name: ref.Project,
		})))

		app := &pb.Ref_Application{
			Project:     ref.Project,
			Application: "testapp",
		}

		ws := &pb.Ref_Workspace{
			Workspace: "default",
		}

		// Add
		err := s.DeploymentPut(ctx, false, serverptypes.TestDeployment(t, &pb.Deployment{
			Id:          "B",
			Application: app,
			Workspace:   ws,
			Status: &pb.Status{
				State:     pb.Status_SUCCESS,
				StartTime: timestamppb.Now(),
			},
		}))
		require.NoError(err)

		// Create an instance
		rec := &serverstate.Instance{
			Id:           "A",
			DeploymentId: "B",
			Project:      ref.Project,
			Application:  app.Application,
			Workspace:    ws.Workspace,
		}

		require.NoError(s.InstanceCreate(ctx, rec))

		// We should be able to find it
		found, err := s.InstanceById(ctx, rec.Id)
		require.NoError(err)
		require.Equal(rec, found)

		// Delete that instance
		require.NoError(s.InstanceDelete(ctx, rec.Id))

		// Delete again should be fine
		require.NoError(s.InstanceDelete(ctx, rec.Id))
	})

	t.Run("not found", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// We should be able to find it
		found, err := s.InstanceById(ctx, "nope")
		require.Error(err)
		require.Nil(found)
		require.Equal(codes.NotFound, status.Code(err))
	})

	t.Run("by app", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		ref := &pb.Ref_Project{Project: "foo"}
		require.NoError(s.ProjectPut(ctx, serverptypes.TestProject(t, &pb.Project{
			Name: ref.Project,
		})))

		app := &pb.Ref_Application{
			Project:     ref.Project,
			Application: "testapp",
		}

		wsRef := &pb.Ref_Workspace{
			Workspace: "default",
		}

		// Add
		err := s.DeploymentPut(ctx, false, serverptypes.TestDeployment(t, &pb.Deployment{
			Id:          "B",
			Application: app,
			Workspace:   wsRef,
			Status: &pb.Status{
				State:     pb.Status_SUCCESS,
				StartTime: timestamppb.Now(),
			},
		}))
		require.NoError(err)

		// Empty with nothing
		ws := memdb.NewWatchSet()
		list, err := s.InstancesByApp(ctx, app, nil, ws)
		require.NoError(err)
		require.Empty(list)

		// Watch should block
		require.True(ws.Watch(time.After(10 * time.Millisecond)))

		// Create an instance
		rec := testInstance(t, &serverstate.Instance{Project: ref.Project, Application: app.Application})
		require.NoError(s.InstanceCreate(ctx, rec))

		// Should be triggered
		require.False(ws.Watch(time.After(3 * time.Second)))

		// Should have values
		list, err = s.InstancesByApp(ctx, app, nil, nil)
		require.NoError(err)
		require.Len(list, 1)

		// Should not for other app
		//nolint:govet,copylocks
		ref2 := *app
		ref2.Application = "NO"
		list, err = s.InstancesByApp(ctx, &ref2, nil, nil)
		require.NoError(err)
		require.Empty(list)
	})

	t.Run("by app workspace", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		ref := &pb.Ref_Project{Project: "foo"}
		require.NoError(s.ProjectPut(ctx, serverptypes.TestProject(t, &pb.Project{
			Name: ref.Project,
		})))

		app := &pb.Ref_Application{
			Project:     ref.Project,
			Application: "testapp",
		}

		wsRef := &pb.Ref_Workspace{
			Workspace: "default",
		}

		// Add
		err := s.DeploymentPut(ctx, false, serverptypes.TestDeployment(t, &pb.Deployment{
			Id:          "B",
			Application: app,
			Workspace:   wsRef,
			Status: &pb.Status{
				State:     pb.Status_SUCCESS,
				StartTime: timestamppb.Now(),
			},
		}))
		require.NoError(err)

		// Empty with nothing
		ws := memdb.NewWatchSet()
		list, err := s.InstancesByApp(ctx, app, wsRef, ws)
		require.NoError(err)
		require.Empty(list)

		// Watch should block
		require.True(ws.Watch(time.After(10 * time.Millisecond)))

		// Create an instance
		rec := testInstance(t, &serverstate.Instance{
			Project: ref.Project, Application: app.Application, Workspace: wsRef.Workspace})
		require.NoError(s.InstanceCreate(ctx, rec))

		// Should be triggered
		require.False(ws.Watch(time.After(3 * time.Second)))

		// Should have values
		list, err = s.InstancesByApp(ctx, app, wsRef, nil)
		require.NoError(err)
		require.Len(list, 1)

		// Should not for other app
		//nolint:govet,copylocks
		ref2 := *wsRef
		ref2.Workspace = "NO"
		list, err = s.InstancesByApp(ctx, app, &ref2, nil)
		require.NoError(err)
		require.Empty(list)
	})
}

func TestInstanceByDeployment(t *testing.T, factory Factory, _ RestartFactory) {
	ctx := context.Background()
	require := require.New(t)

	s := factory(t)
	defer s.Close()

	ref := &pb.Ref_Project{Project: "foo"}
	require.NoError(s.ProjectPut(ctx, serverptypes.TestProject(t, &pb.Project{
		Name: ref.Project,
	})))

	app := &pb.Ref_Application{
		Project:     ref.Project,
		Application: "testapp",
	}

	ws := &pb.Ref_Workspace{
		Workspace: "default",
	}

	// Add two deployments
	require.NoError(s.DeploymentPut(ctx, false, serverptypes.TestDeployment(t, &pb.Deployment{
		Id:          "A",
		Application: app,
		Workspace:   ws,
		Status: &pb.Status{
			State:     pb.Status_SUCCESS,
			StartTime: timestamppb.Now(),
		},
	})))

	require.NoError(s.DeploymentPut(ctx, false, serverptypes.TestDeployment(t, &pb.Deployment{
		Id:          "B",
		Application: app,
		Workspace:   ws,
		Status: &pb.Status{
			State:     pb.Status_SUCCESS,
			StartTime: timestamppb.Now(),
		},
	})))

	// Create three instances, for deployment A, two for B

	initialAInst := &serverstate.Instance{
		Id:           "A",
		DeploymentId: "A",
		Project:      ref.Project,
		Application:  app.Application,
		Workspace:    ws.Workspace,
		DisableExec:  true,
	}
	require.NoError(s.InstanceCreate(ctx, initialAInst))

	require.NoError(s.InstanceCreate(ctx, &serverstate.Instance{
		Id:           "B1",
		DeploymentId: "B",
		Project:      ref.Project,
		Application:  app.Application,
		Workspace:    ws.Workspace,
	}))

	require.NoError(s.InstanceCreate(ctx, &serverstate.Instance{
		Id:           "B2",
		DeploymentId: "B",
		Project:      ref.Project,
		Application:  app.Application,
		Workspace:    ws.Workspace,
	}))

	t.Run("can get deployment A's instance", func(t *testing.T) {
		inst, err := s.InstancesByDeployment(ctx, "A", nil)
		require.NoError(err)
		require.Len(inst, 1)

		// Ensure all the fields have been set
		require.Equal(inst[0].Id, initialAInst.Id)
		require.Equal(inst[0].DeploymentId, initialAInst.DeploymentId)
		require.Equal(inst[0].Application, initialAInst.Application)
		require.Equal(inst[0].Project, initialAInst.Project)
		require.Equal(inst[0].DisableExec, initialAInst.DisableExec)
	})

	t.Run("can get deployment B's instances", func(t *testing.T) {
		inst, err := s.InstancesByDeployment(ctx, "B", nil)
		require.NoError(err)
		require.Len(inst, 2)

		// Ensure we got both of B's instances (but ignore order)
		require.True(inst[0].Id == "B1" || inst[0].Id == "B2")
		require.True(inst[1].Id == "B1" || inst[1].Id == "B2")
		require.True(inst[0].Id != inst[1].Id)
	})

}
