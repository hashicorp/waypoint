package server

import (
	"context"
	"net"

	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	pb "github.com/mitchellh/devflow/internal/server/gen"
)

// TestServer starts a server and returns a gRPC client to that server.
// We use t.Cleanup to ensure resources are automatically cleaned up.
func TestServer(t testing.T, impl pb.DevflowServer) pb.DevflowClient {
	require := require.New(t)

	// Listen on a random port
	ln, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(err)
	t.Cleanup(func() { ln.Close() })

	// Create the server
	ctx, cancel := context.WithCancel(context.Background())
	go Run(
		WithContext(ctx),
		WithGRPC(ln),
		WithImpl(impl),
	)
	t.Cleanup(func() { cancel() })

	// Connect, this should retry in the case Run is not going yet
	conn, err := grpc.DialContext(ctx, ln.Addr().String(),
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	require.NoError(err)
	t.Cleanup(func() { conn.Close() })

	return pb.NewDevflowClient(conn)
}
