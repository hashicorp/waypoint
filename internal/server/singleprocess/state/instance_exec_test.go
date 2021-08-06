package state

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/go-memdb"
	"github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestInstanceExecCreateByTargetedInstance_disabled(t *testing.T) {
	require := require.New(t)

	s := TestState(t)
	defer s.Close()

	// Create an instance
	instance := testInstance(t, nil)
	instance.DisableExec = true
	require.NoError(s.InstanceCreate(instance))

	{
		// Create an instance exec targetting the specific instance
		rec := &InstanceExec{
			InstanceId: instance.Id,
		}
		err := s.InstanceExecCreateByTargetedInstance(instance.Id, rec)
		require.Error(err)
		require.Equal(codes.PermissionDenied, status.Code(err))
	}

	// List them
	list, err := s.InstanceExecListByInstanceId(instance.Id, nil)
	require.NoError(err)
	require.Len(list, 0)
}

func TestInstanceExecCreateByDeploymentId_invalidDeployment(t *testing.T) {
	require := require.New(t)

	s := TestState(t)
	defer s.Close()

	// Create an instance
	rec := &InstanceExec{}
	err := s.InstanceExecCreateByDeployment("nope", rec)
	require.Error(err)
	require.Equal(codes.ResourceExhausted, status.Code(err))
}

func TestInstanceExecCreateByDeploymentId_allDisabled(t *testing.T) {
	require := require.New(t)

	s := TestState(t)
	defer s.Close()

	// Create an instance
	instance := testInstance(t, nil)
	instance.DisableExec = true
	require.NoError(s.InstanceCreate(instance))

	{
		// Create an instance exec targetting the specific instance
		rec := &InstanceExec{}
		err := s.InstanceExecCreateByDeployment(instance.DeploymentId, rec)
		require.Error(err)
		require.Equal(codes.ResourceExhausted, status.Code(err))
	}

	// List them
	list, err := s.InstanceExecListByInstanceId(instance.Id, nil)
	require.NoError(err)
	require.Len(list, 0)
}

func TestInstanceExecCreateByDeploymentId_valid(t *testing.T) {
	require := require.New(t)

	s := TestState(t)
	defer s.Close()

	// Create an instance
	instance := testInstance(t, nil)
	require.NoError(s.InstanceCreate(instance))

	{
		// Create an instance exec
		rec := &InstanceExec{}
		require.NoError(s.InstanceExecCreateByDeployment(instance.DeploymentId, rec))
		require.NotEmpty(rec.Id)
		require.Equal(instance.Id, rec.InstanceId)

		// Test single get
		found, err := s.InstanceExecById(rec.Id)
		require.NoError(err)
		require.Equal(rec, found)
	}

	// List them
	ws := memdb.NewWatchSet()
	list, err := s.InstanceExecListByInstanceId(instance.Id, ws)
	require.NoError(err)
	require.Len(list, 1)
	require.True(ws.Watch(time.After(50 * time.Millisecond)))

	{
		// Next one shuld get the same instance since its all we have
		rec := &InstanceExec{}
		require.NoError(s.InstanceExecCreateByDeployment(instance.DeploymentId, rec))
		require.NotEmpty(rec.Id)
		require.Equal(instance.Id, rec.InstanceId)

		// Should fire the watch
		require.False(ws.Watch(time.After(50 * time.Millisecond)))
	}

	list, err = s.InstanceExecListByInstanceId(instance.Id, nil)
	require.NoError(err)
	require.Len(list, 2)

	// Create another instance
	instance = testInstance(t, &Instance{Id: "B", DeploymentId: "A"})
	require.NoError(s.InstanceCreate(instance))

	{
		// Next one shuld get the B because its less loaded
		rec := &InstanceExec{}
		require.NoError(s.InstanceExecCreateByDeployment(instance.DeploymentId, rec))
		require.NotEmpty(rec.Id)
		require.Equal(instance.Id, rec.InstanceId)

		// Should be able to delete
		require.NoError(s.InstanceExecDelete(rec.Id))
	}

	// Create another instance, with exec disabled
	{
		instance := testInstance(t, &Instance{
			Id: "C", DeploymentId: "A", DisableExec: true})
		require.NoError(s.InstanceCreate(instance))

		// Should not get C
		rec := &InstanceExec{}
		require.NoError(s.InstanceExecCreateByDeployment(instance.DeploymentId, rec))
		require.NotEmpty(rec.Id)
		require.NotEqual("C", rec.InstanceId)

		// Should be able to delete
		require.NoError(s.InstanceExecDelete(rec.Id))
	}

	{
		// Create an instance exec targetting the specific instance
		rec := &InstanceExec{
			InstanceId: instance.Id,
		}
		require.NoError(s.InstanceExecCreateByTargetedInstance(instance.Id, rec))
		require.NotEmpty(rec.Id)
		require.Equal(instance.Id, rec.InstanceId)

		// Test single get
		found, err := s.InstanceExecById(rec.Id)
		require.NoError(err)
		require.Equal(rec, found)
	}
}

func TestInstanceExecCreateByDeploymentId_longrunningonly(t *testing.T) {
	require := require.New(t)

	s := TestState(t)
	defer s.Close()

	// Create an instance
	instance := testInstance(t, nil)
	require.NoError(s.InstanceCreate(instance))

	// Create a ondemand instance
	od := testInstance(t, nil)
	od.Id = "A2"
	od.Type = gen.Instance_ON_DEMAND
	require.NoError(s.InstanceCreate(od))

	// Create a virtual instance
	virt := testInstance(t, nil)
	virt.Id = "A3"
	virt.Type = gen.Instance_VIRTUAL
	require.NoError(s.InstanceCreate(virt))

	// We'll create 3 instance exec and make sure they all go to instance only
	for i := 0; i < 3; i++ {
		// Create an instance exec
		rec := &InstanceExec{}
		require.NoError(s.InstanceExecCreateByDeployment(instance.DeploymentId, rec))
		require.NotEmpty(rec.Id)
		require.Equal(instance.Id, rec.InstanceId)

		// Test single get
		found, err := s.InstanceExecById(rec.Id)
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

		require.NoError(s.InstanceCreate(instance))
	}()

	{
		// Create an instance exec
		rec := &InstanceExec{}
		require.NoError(s.InstanceExecCreateForVirtualInstance(ctx, instId, rec))
		require.NotEmpty(rec.Id)
		require.Equal(instId, rec.InstanceId)

		// Test single get
		found, err := s.InstanceExecById(rec.Id)
		require.NoError(err)
		require.Equal(rec, found)
	}

	// List them
	ws := memdb.NewWatchSet()
	list, err := s.InstanceExecListByInstanceId(instId, ws)
	require.NoError(err)
	require.Len(list, 1)
	require.True(ws.Watch(time.After(50 * time.Millisecond)))
}
