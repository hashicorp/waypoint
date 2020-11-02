package state

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

func TestSnapshotRestore(t *testing.T) {
	require := require.New(t)

	s := TestState(t)
	defer s.Close()

	// Create some data
	err := s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
		Name: "A",
	}))
	require.NoError(err)
	resp, err := s.ProjectGet(&pb.Ref_Project{
		Project: "A",
	})
	require.NoError(err)
	require.NotNil(resp)

	// Snapshot
	var buf bytes.Buffer
	require.NoError(s.CreateSnapshot(&buf))

	// Create more data that isn't in the snapshot
	err = s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
		Name: "B",
	}))
	require.NoError(err)

	// Restore
	require.NoError(s.StageRestoreSnapshot(bytes.NewReader(buf.Bytes())))

	// Reboot!
	s, err = TestStateRestart(t, s)
	require.NoError(err)

	// Should find first record and not the second
	{
		resp, err := s.ProjectGet(&pb.Ref_Project{
			Project: "A",
		})
		require.NoError(err)
		require.NotNil(resp)
	}
	{
		_, err := s.ProjectGet(&pb.Ref_Project{
			Project: "B",
		})
		require.Error(err)
		require.Equal(codes.NotFound, status.Code(err))
	}

	// Create more data
	err = s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
		Name: "B",
	}))
	require.NoError(err)

	// Reboot again, should not restore again
	s, err = TestStateRestart(t, s)
	require.NoError(err)

	// Should find both records
	{
		resp, err := s.ProjectGet(&pb.Ref_Project{
			Project: "A",
		})
		require.NoError(err)
		require.NotNil(resp)
	}
	{
		resp, err := s.ProjectGet(&pb.Ref_Project{
			Project: "B",
		})
		require.NoError(err)
		require.NotNil(resp)
	}
}

func TestSnapshotRestore_corrupt(t *testing.T) {
	require := require.New(t)

	s := TestState(t)
	defer s.Close()

	// Restore with garbage data
	require.Error(s.StageRestoreSnapshot(strings.NewReader(
		"I am probably not a valid BoltDB file.")))

	// Reboot!
	s, err := TestStateRestart(t, s)
	require.NoError(err)
}
