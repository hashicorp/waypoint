package singleprocess

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// TestServer starts a singleprocess server and returns the connected client.
// We use t.Cleanup to ensure resources are automatically cleaned up.
func TestServer(t testing.T) pb.WaypointClient {
	impl, err := New(testDB(t))
	require.NoError(t, err)
	return server.TestServer(t, impl)
}

func TestEntrypoint(t testing.T, client pb.WaypointClient) (string, string, func()) {
	instanceId, err := server.Id()
	require.NoError(t, err)

	// Create the config
	stream, err := client.EntrypointConfig(context.Background(), &pb.EntrypointConfigRequest{
		InstanceId:   instanceId,
		DeploymentId: "A",
	})
	require.NoError(t, err)

	// Wait for the first config so that we know we're registered
	_, err = stream.Recv()
	require.NoError(t, err)

	return instanceId, "A", func() {
		stream.CloseSend()
	}
}

func testDB(t testing.T) *bolt.DB {
	t.Helper()

	// Temporary directory for the database
	td, err := ioutil.TempDir("", "test")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(td) })

	// Create the DB
	db, err := bolt.Open(filepath.Join(td, "test.db"), 0600, nil)
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	// Init
	require.NoError(t, dbInit(db))

	return db
}
