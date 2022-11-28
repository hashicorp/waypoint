package boltdbstate

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/go-memdb"
	"github.com/stretchr/testify/require"
)

func TestInstanceLogsCreate(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	s := TestState(t)
	defer s.Close()

	{
		// Create an instance exec
		rec := &InstanceLogs{}

		require.NoError(s.InstanceLogsCreate(ctx, "a", rec))
		require.NotEmpty(rec.Id)
		require.Equal("a", rec.InstanceId)

		// Test single get
		found, err := s.InstanceLogsById(ctx, rec.Id)
		require.NoError(err)
		require.Equal(rec, found)

		// Test single get
		found2, err := s.InstanceLogsByInstanceId(ctx, "a")
		require.NoError(err)
		require.Equal(rec, found2)
	}

	// List them
	ws := memdb.NewWatchSet()
	list, err := s.InstanceLogsListByInstanceId(ctx, "a", ws)
	require.NoError(err)
	require.Len(list, 1)
	require.True(ws.Watch(time.After(50 * time.Millisecond)))
}
