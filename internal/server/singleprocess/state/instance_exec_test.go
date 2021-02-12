package state

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/go-memdb"
	"github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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

func TestCalculateInstanceExecByDeployment(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	s := TestState(t)
	defer s.Close()

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

	// Should get an error because there are no long running instances
	_, err := s.CalculateInstanceExecByDeployment(od.DeploymentId)
	require.Error(err)

	// Create an instance
	instance := testInstance(t, nil)
	require.NoError(s.InstanceCreate(instance))

	// Run it 3 times and make sure we only see the long running instance
	for i := 0; i < 3; i++ {
		inst, err := s.CalculateInstanceExecByDeployment(instance.DeploymentId)
		require.NoError(err)
		require.Equal(instance.Id, inst.Id)
	}

	// Create another long running instance
	lr := testInstance(t, nil)
	lr.Id = "A4"
	lr.Type = gen.Instance_LONG_RUNNING
	require.NoError(s.InstanceCreate(lr))

	// Get an instance, it should be one of the long running ones
	reserve, err := s.CalculateInstanceExecByDeployment(instance.DeploymentId)
	require.NoError(err)

	var exec InstanceExec
	require.NoError(s.InstanceExecCreateByTargetedInstance(reserve.Id, &exec))

	// ok, now see that on the next time, we get the other long running instance
	reserve2, err := s.CalculateInstanceExecByDeployment(instance.DeploymentId)
	require.NoError(err)
	assert.NotEqual(reserve.Id, reserve2.Id)
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
