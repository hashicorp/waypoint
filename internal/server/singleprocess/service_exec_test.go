package singleprocess

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func TestServiceStartExecStream_badOpen(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	impl, err := New(testDB(t))
	require.NoError(err)
	client := server.TestServer(t, impl)

	// Start exec with a bad starting message
	stream, err := client.StartExecStream(ctx)
	require.NoError(err)
	require.NoError(stream.Send(&pb.ExecStreamRequest{
		Event: &pb.ExecStreamRequest_Input_{
			Input: &pb.ExecStreamRequest_Input{},
		},
	}))

	// Wait for data
	resp, err := stream.Recv()
	require.Error(err)
	require.Equal(codes.FailedPrecondition, status.Code(err))
	require.Nil(resp)
}

func TestServiceStartExecStream_start(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	impl, err := New(testDB(t))
	require.NoError(err)
	client := server.TestServer(t, impl)

	// Create an instance
	_, deploymentId, closer := TestEntrypoint(t, client)
	defer closer()

	// Start exec with a bad starting message
	stream, err := client.StartExecStream(ctx)
	require.NoError(err)
	require.NoError(stream.Send(&pb.ExecStreamRequest{
		Event: &pb.ExecStreamRequest_Start_{
			Start: &pb.ExecStreamRequest_Start{
				DeploymentId: deploymentId,
				Args:         []string{"foo", "bar"},
			},
		},
	}))

	// Close send
	require.NoError(stream.CloseSend())

	// The above close send should trigger the stream to end.
	resp, err := stream.Recv()
	require.Error(err)
	require.Equal(io.EOF, err)
	require.Nil(resp)
}
