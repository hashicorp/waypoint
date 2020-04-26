package state

import (
	"testing"
	"time"

	"github.com/hashicorp/go-memdb"
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
	instance := &Instance{Id: "A", DeploymentId: "A"}
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
	instance = &Instance{Id: "B", DeploymentId: "A"}
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
}
