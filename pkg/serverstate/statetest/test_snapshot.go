package statetest

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func init() {
	tests["snapshot"] = []testFunc{
		TestSnapshotRestore,
		TestSnapshotRestore_corrupt,
	}
}

func TestSnapshotRestore(t *testing.T, factory Factory, factoryRestart RestartFactory) {
	require := require.New(t)

	s := factory(t)
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
	s = factoryRestart(t, s)

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
	s = factoryRestart(t, s)

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

func TestSnapshotRestore_corrupt(t *testing.T, factory Factory, factoryRestart RestartFactory) {
	require := require.New(t)

	s := factory(t)
	defer s.Close()

	// Restore with garbage data
	require.Error(s.StageRestoreSnapshot(strings.NewReader(
		"I am probably not a valid BoltDB file.")))

	// Reboot!
	factoryRestart(t, s)
}
