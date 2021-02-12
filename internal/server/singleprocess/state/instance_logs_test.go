package state

import (
	"testing"
	"time"

	"github.com/hashicorp/go-memdb"
	"github.com/stretchr/testify/require"
)

func TestInstanceLogsCreate(t *testing.T) {
	require := require.New(t)

	s := TestState(t)
	defer s.Close()

	{
		// Create an instance exec
		rec := &InstanceLogs{}

		require.NoError(s.InstanceLogsCreate("a", rec))
		require.NotEmpty(rec.Id)
		require.Equal("a", rec.InstanceId)

		// Test single get
		found, err := s.InstanceLogsById(rec.Id)
		require.NoError(err)
		require.Equal(rec, found)

		// Test single get
		found2, err := s.InstanceLogsByInstanceId("a")
		require.NoError(err)
		require.Equal(rec, found2)
	}

	// List them
	ws := memdb.NewWatchSet()
	list, err := s.InstanceLogsListByInstanceId("a", ws)
	require.NoError(err)
	require.Len(list, 1)
	require.True(ws.Watch(time.After(50 * time.Millisecond)))
}
