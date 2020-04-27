package singleprocess

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/hashicorp/go-memdb"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/singleprocess/state"
)

func TestServiceEntrypointExecStream_badOpen(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	impl, err := New(testDB(t))
	require.NoError(err)
	client := server.TestServer(t, impl)

	// Start exec with a bad starting message
	stream, err := client.EntrypointExecStream(ctx)
	require.NoError(err)
	require.NoError(stream.Send(&pb.EntrypointExecRequest{
		Event: &pb.EntrypointExecRequest_Output_{
			Output: &pb.EntrypointExecRequest_Output{},
		},
	}))

	// Wait for data
	resp, err := stream.Recv()
	require.Error(err)
	require.Equal(codes.FailedPrecondition, status.Code(err))
	require.Nil(resp)
}

func TestServiceEntrypointExecStream_invalidInstanceId(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	impl, err := New(testDB(t))
	require.NoError(err)
	client := server.TestServer(t, impl)

	// Start exec with a bad starting message
	stream, err := client.EntrypointExecStream(ctx)
	require.NoError(err)
	require.NoError(stream.Send(&pb.EntrypointExecRequest{
		Event: &pb.EntrypointExecRequest_Open_{
			Open: &pb.EntrypointExecRequest_Open{
				InstanceId: "nope",
				Index:      0,
			},
		},
	}))

	// Wait for data
	resp, err := stream.Recv()
	require.Error(err)
	require.Equal(codes.NotFound, status.Code(err))
	require.Nil(resp)
}

func TestServiceEntrypointExecStream_invalidSessionId(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	impl, err := New(testDB(t))
	require.NoError(err)
	client := server.TestServer(t, impl)

	exec, closer := testRegisterExec(t, client, impl)
	defer closer()

	// Start exec with a bad starting message
	stream, err := client.EntrypointExecStream(ctx)
	require.NoError(err)
	require.NoError(stream.Send(&pb.EntrypointExecRequest{
		Event: &pb.EntrypointExecRequest_Open_{
			Open: &pb.EntrypointExecRequest_Open{
				InstanceId: exec.InstanceId,
				Index:      exec.Id + 4,
			},
		},
	}))

	// Wait for data
	resp, err := stream.Recv()
	require.Error(err)
	require.Equal(codes.NotFound, status.Code(err))
	require.Nil(resp)
}

func TestServiceEntrypointExecStream_closeSend(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	impl, err := New(testDB(t))
	require.NoError(err)
	client := server.TestServer(t, impl)

	exec, closer := testRegisterExec(t, client, impl)
	defer closer()

	// Start exec with a bad starting message
	stream, err := client.EntrypointExecStream(ctx)
	require.NoError(err)
	require.NoError(stream.Send(&pb.EntrypointExecRequest{
		Event: &pb.EntrypointExecRequest_Open_{
			Open: &pb.EntrypointExecRequest_Open{
				InstanceId: exec.InstanceId,
				Index:      exec.Id,
			},
		},
	}))

	// Close our sending side
	require.NoError(stream.CloseSend())

	// Wait for data
	resp, err := stream.Recv()
	require.Error(err)
	require.Equal(io.EOF, err)
	require.Nil(resp)
}

func TestServiceEntrypointExecStream_doubleStart(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	impl, err := New(testDB(t))
	require.NoError(err)
	client := server.TestServer(t, impl)

	exec, closer := testRegisterExec(t, client, impl)
	defer closer()

	// Start exec
	stream, err := client.EntrypointExecStream(ctx)
	require.NoError(err)
	require.NoError(stream.Send(&pb.EntrypointExecRequest{
		Event: &pb.EntrypointExecRequest_Open_{
			Open: &pb.EntrypointExecRequest_Open{
				InstanceId: exec.InstanceId,
				Index:      exec.Id,
			},
		},
	}))
	defer stream.CloseSend()

	// Start a second exec
	stream2, err := client.EntrypointExecStream(ctx)
	require.NoError(err)
	require.NoError(stream2.Send(&pb.EntrypointExecRequest{
		Event: &pb.EntrypointExecRequest_Open_{
			Open: &pb.EntrypointExecRequest_Open{
				InstanceId: exec.InstanceId,
				Index:      exec.Id,
			},
		},
	}))
	defer stream2.CloseSend()

	// Wait for data
	resp, err := stream2.Recv()
	require.Error(err)
	require.Equal(codes.FailedPrecondition, status.Code(err))
	require.Nil(resp)
}

func testRegisterExec(t *testing.T, client pb.WaypointClient, impl pb.WaypointServer) (*state.InstanceExec, func()) {
	// Create an instance
	instanceId, deploymentId, closer := TestEntrypoint(t, client)
	defer closer()

	// Start exec
	stream, err := client.StartExecStream(context.Background())
	require.NoError(t, err)
	require.NoError(t, stream.Send(&pb.ExecStreamRequest{
		Event: &pb.ExecStreamRequest_Start_{
			Start: &pb.ExecStreamRequest_Start{
				DeploymentId: deploymentId,
				Args:         []string{"foo", "bar"},
			},
		},
	}))

	// Wait for the registered exec
	ws := memdb.NewWatchSet()
	list, err := testServiceImpl(impl).state.InstanceExecListByInstanceId(instanceId, ws)
	require.NoError(t, err)
	if len(list) == 0 {
		ws.Watch(time.After(1 * time.Second))
		list, err = impl.(*service).state.InstanceExecListByInstanceId(instanceId, ws)
		require.NoError(t, err)
	}
	require.Len(t, list, 1)

	return list[0], func() {
		stream.CloseSend()
	}
}
