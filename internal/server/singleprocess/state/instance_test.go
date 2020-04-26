package state

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestInstance_crud(t *testing.T) {
	require := require.New(t)

	s := TestState(t)
	defer s.Close()

	// Create an instance
	rec := &Instance{Id: "A", DeploymentId: "B"}
	require.NoError(s.InstanceCreate(rec))

	// We should be able to find it
	found, err := s.InstanceById(rec.Id)
	require.NoError(err)
	require.Equal(rec, found)

	// Delete that instance
	require.NoError(s.InstanceDelete(rec.Id))

	// Delete again should be fine
	require.NoError(s.InstanceDelete(rec.Id))
}

func TestInstanceById_notFound(t *testing.T) {
	require := require.New(t)

	s := TestState(t)
	defer s.Close()

	// We should be able to find it
	found, err := s.InstanceById("nope")
	require.Error(err)
	require.Nil(found)
	require.Equal(codes.NotFound, status.Code(err))
}
