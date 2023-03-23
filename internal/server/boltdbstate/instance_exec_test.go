// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package boltdbstate

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/go-memdb"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

func TestInstanceExecCreateByTargetedInstance_disabled(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	s := TestState(t)
	defer s.Close()

	// Create an instance
	instance := testInstance(t, nil)
	instance.DisableExec = true
	require.NoError(s.InstanceCreate(ctx, instance))

	{
		// Create an instance exec targetting the specific instance
		rec := &serverstate.InstanceExec{
			InstanceId: instance.Id,
		}
		err := s.InstanceExecCreateByTargetedInstance(ctx, instance.Id, rec)
		require.Error(err)
		require.Equal(codes.PermissionDenied, status.Code(err))
	}

	// List them
	list, err := s.InstanceExecListByInstanceId(ctx, instance.Id, nil)
	require.NoError(err)
	require.Len(list, 0)
}

func TestInstanceExecCreateByDeploymentId_invalidDeployment(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	s := TestState(t)
	defer s.Close()

	// Create an instance
	rec := &serverstate.InstanceExec{}
	err := s.InstanceExecCreateByDeployment(ctx, "nope", rec)
	require.Error(err)
	require.Equal(codes.ResourceExhausted, status.Code(err))
}

func TestInstanceExecCreateByDeploymentId_allDisabled(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	s := TestState(t)
	defer s.Close()

	// Create an instance
	instance := testInstance(t, nil)
	instance.DisableExec = true
	require.NoError(s.InstanceCreate(ctx, instance))

	{
		// Create an instance exec targetting the specific instance
		rec := &serverstate.InstanceExec{}
		err := s.InstanceExecCreateByDeployment(ctx, instance.DeploymentId, rec)
		require.Error(err)
		require.Equal(codes.ResourceExhausted, status.Code(err))
	}

	// List them
	list, err := s.InstanceExecListByInstanceId(ctx, instance.Id, nil)
	require.NoError(err)
	require.Len(list, 0)
}

func TestInstanceExecCreateByDeploymentId_valid(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	s := TestState(t)
	defer s.Close()

	// Create an instance
	instance := testInstance(t, nil)
	require.NoError(s.InstanceCreate(ctx, instance))

	{
		// Create an instance exec
		rec := &serverstate.InstanceExec{}
		require.NoError(s.InstanceExecCreateByDeployment(ctx, instance.DeploymentId, rec))
		require.NotEmpty(rec.Id)
		require.Equal(instance.Id, rec.InstanceId)

		// Test single get
		found, err := s.InstanceExecById(ctx, rec.Id)
		require.NoError(err)
		require.Equal(rec, found)
	}

	// List them
	ws := memdb.NewWatchSet()
	list, err := s.InstanceExecListByInstanceId(ctx, instance.Id, ws)
	require.NoError(err)
	require.Len(list, 1)
	require.True(ws.Watch(time.After(50 * time.Millisecond)))

	{
		// Next one shuld get the same instance since its all we have
		rec := &serverstate.InstanceExec{}
		require.NoError(s.InstanceExecCreateByDeployment(ctx, instance.DeploymentId, rec))
		require.NotEmpty(rec.Id)
		require.Equal(instance.Id, rec.InstanceId)

		// Should fire the watch
		require.False(ws.Watch(time.After(50 * time.Millisecond)))
	}

	list, err = s.InstanceExecListByInstanceId(ctx, instance.Id, nil)
	require.NoError(err)
	require.Len(list, 2)

	// Create another instance
	instance = testInstance(t, &serverstate.Instance{Id: "B", DeploymentId: "A"})
	require.NoError(s.InstanceCreate(ctx, instance))

	{
		// Next one shuld get the B because its less loaded
		rec := &serverstate.InstanceExec{}
		require.NoError(s.InstanceExecCreateByDeployment(ctx, instance.DeploymentId, rec))
		require.NotEmpty(rec.Id)
		require.Equal(instance.Id, rec.InstanceId)

		// Should be able to delete
		require.NoError(s.InstanceExecDelete(ctx, rec.Id))
	}

	// Create another instance, with exec disabled
	{
		instance := testInstance(t, &serverstate.Instance{
			Id: "C", DeploymentId: "A", DisableExec: true})
		require.NoError(s.InstanceCreate(ctx, instance))

		// Should not get C
		rec := &serverstate.InstanceExec{}
		require.NoError(s.InstanceExecCreateByDeployment(ctx, instance.DeploymentId, rec))
		require.NotEmpty(rec.Id)
		require.NotEqual("C", rec.InstanceId)

		// Should be able to delete
		require.NoError(s.InstanceExecDelete(ctx, rec.Id))
	}

	{
		// Create an instance exec targetting the specific instance
		rec := &serverstate.InstanceExec{
			InstanceId: instance.Id,
		}
		require.NoError(s.InstanceExecCreateByTargetedInstance(ctx, instance.Id, rec))
		require.NotEmpty(rec.Id)
		require.Equal(instance.Id, rec.InstanceId)

		// Test single get
		found, err := s.InstanceExecById(ctx, rec.Id)
		require.NoError(err)
		require.Equal(rec, found)
	}
}

func TestInstanceExecCreateByDeploymentId_longrunningonly(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	s := TestState(t)
	defer s.Close()

	// Create an instance
	instance := testInstance(t, nil)
	require.NoError(s.InstanceCreate(ctx, instance))

	// Create a ondemand instance
	od := testInstance(t, nil)
	od.Id = "A2"
	od.Type = gen.Instance_ON_DEMAND
	require.NoError(s.InstanceCreate(ctx, od))

	// Create a virtual instance
	virt := testInstance(t, nil)
	virt.Id = "A3"
	virt.Type = gen.Instance_VIRTUAL
	require.NoError(s.InstanceCreate(ctx, virt))

	// We'll create 3 instance exec and make sure they all go to instance only
	for i := 0; i < 3; i++ {
		// Create an instance exec
		rec := &serverstate.InstanceExec{}
		require.NoError(s.InstanceExecCreateByDeployment(ctx, instance.DeploymentId, rec))
		require.NotEmpty(rec.Id)
		require.Equal(instance.Id, rec.InstanceId)

		// Test single get
		found, err := s.InstanceExecById(ctx, rec.Id)
		require.NoError(err)
		require.Equal(rec, found)
	}
}

func TestInstanceExecCreateForVirtualInstance(t *testing.T) {
	require := require.New(t)

	s := TestState(t)
	defer s.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	instId := "instA"

	go func() {
		time.Sleep(time.Second)
		instance := testInstance(t, nil)
		instance.Id = instId

		require.NoError(s.InstanceCreate(ctx, instance))
	}()

	{
		// Create an instance exec
		rec := &serverstate.InstanceExec{}
		require.NoError(s.InstanceExecCreateForVirtualInstance(ctx, instId, rec))
		require.NotEmpty(rec.Id)
		require.Equal(instId, rec.InstanceId)

		// Test single get
		found, err := s.InstanceExecById(ctx, rec.Id)
		require.NoError(err)
		require.Equal(rec, found)
	}

	// List them
	ws := memdb.NewWatchSet()
	list, err := s.InstanceExecListByInstanceId(ctx, instId, ws)
	require.NoError(err)
	require.Len(list, 1)
	require.True(ws.Watch(time.After(50 * time.Millisecond)))
}
