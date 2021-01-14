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

func TestServiceStartExecStream_badOpen(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	impl, err := New(WithDB(testDB(t)))
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
	impl, err := New(WithDB(testDB(t)))
	require.NoError(err)
	client := server.TestServer(t, impl)

	// Create an instance
	instanceId, deploymentId, closer := TestEntrypoint(t, client)
	defer closer()

	// Start exec with a bad starting message
	stream, err := client.StartExecStream(ctx)
	require.NoError(err)
	require.NoError(stream.Send(&pb.ExecStreamRequest{
		Event: &pb.ExecStreamRequest_Start_{
			Start: &pb.ExecStreamRequest_Start{
				Target: &pb.ExecStreamRequest_Start_DeploymentId{
					DeploymentId: deploymentId,
				},
				Args: []string{"foo", "bar"},
			},
		},
	}))

	// Should open
	{
		resp, err := stream.Recv()
		require.NoError(err)
		open, ok := resp.Event.(*pb.ExecStreamResponse_Open_)
		require.True(ok, "should be an open")
		require.NotNil(open)
	}

	// Get our instance exec value
	exec := testGetInstanceExec(t, impl, instanceId)
	require.Equal([]string{"foo", "bar"}, exec.Args)

	// Close send
	require.NoError(stream.CloseSend())

	// The above close send should trigger the stream to end.
	resp, err := stream.Recv()
	require.Error(err)
	require.Equal(io.EOF, err)
	require.Nil(resp)

	// The event channel should be closed
	_, active := <-exec.ClientEventCh
	require.False(active)
}

func TestServiceStartExecStream_eventExit(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(err)
	client := server.TestServer(t, impl)

	// Create an instance
	instanceId, deploymentId, closer := TestEntrypoint(t, client)
	defer closer()

	// Start stream
	stream, err := client.StartExecStream(ctx)
	require.NoError(err)
	require.NoError(stream.Send(&pb.ExecStreamRequest{
		Event: &pb.ExecStreamRequest_Start_{
			Start: &pb.ExecStreamRequest_Start{
				Target: &pb.ExecStreamRequest_Start_DeploymentId{
					DeploymentId: deploymentId,
				},
				Args: []string{"foo", "bar"},
			},
		},
	}))
	defer stream.CloseSend()

	// Should open
	{
		resp, err := stream.Recv()
		require.NoError(err)
		open, ok := resp.Event.(*pb.ExecStreamResponse_Open_)
		require.True(ok, "should be an open")
		require.NotNil(open)
	}

	// Get the record
	ws := memdb.NewWatchSet()
	list, err := testServiceImpl(impl).state.InstanceExecListByInstanceId(instanceId, ws)
	require.NoError(err)
	if len(list) == 0 {
		ws.Watch(time.After(1 * time.Second))
		list, err = impl.(*service).state.InstanceExecListByInstanceId(instanceId, ws)
		require.NoError(err)
	}
	require.Len(list, 1)
	exec := list[0]

	// Send an exit event
	exec.EntrypointEventCh <- &pb.EntrypointExecRequest{
		Event: &pb.EntrypointExecRequest_Exit_{
			Exit: &pb.EntrypointExecRequest_Exit{
				Code: 1,
			},
		},
	}

	// The above should trigger an exit event
	resp, err := stream.Recv()
	require.NoError(err)
	exitResp, ok := resp.Event.(*pb.ExecStreamResponse_Exit_)
	require.True(ok)
	require.Equal(int32(1), exitResp.Exit.Code)

	// Then we should get a close
	_, err = stream.Recv()
	require.Error(err)
	require.Equal(io.EOF, err)

	// The event channel should be closed
	_, active := <-exec.ClientEventCh
	require.False(active)
}

// When the InstanceExec EntrypointEventCh closes, we should exit.
func TestServiceStartExecStream_entrypointEventChClose(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(err)
	client := server.TestServer(t, impl)

	// Create an instance
	instanceId, deploymentId, closer := TestEntrypoint(t, client)
	defer closer()

	// Start stream
	stream, err := client.StartExecStream(ctx)
	require.NoError(err)
	require.NoError(stream.Send(&pb.ExecStreamRequest{
		Event: &pb.ExecStreamRequest_Start_{
			Start: &pb.ExecStreamRequest_Start{
				Target: &pb.ExecStreamRequest_Start_DeploymentId{
					DeploymentId: deploymentId,
				},
				Args: []string{"foo", "bar"},
			},
		},
	}))
	defer stream.CloseSend()

	// Should open
	{
		resp, err := stream.Recv()
		require.NoError(err)
		open, ok := resp.Event.(*pb.ExecStreamResponse_Open_)
		require.True(ok, "should be an open")
		require.NotNil(open)
	}

	// Get the record
	ws := memdb.NewWatchSet()
	list, err := testServiceImpl(impl).state.InstanceExecListByInstanceId(instanceId, ws)
	require.NoError(err)
	if len(list) == 0 {
		ws.Watch(time.After(1 * time.Second))
		list, err = impl.(*service).state.InstanceExecListByInstanceId(instanceId, ws)
		require.NoError(err)
	}
	require.Len(list, 1)
	exec := list[0]

	// Send an exit event
	close(exec.EntrypointEventCh)

	// Then we should get a close
	resp, err := stream.Recv()
	require.Error(err)
	require.Equal(io.EOF, err)
	require.Nil(resp)

	// The event channel should be closed
	_, active := <-exec.ClientEventCh
	require.False(active)
}

func testGetInstanceExec(t *testing.T, impl pb.WaypointServer, instanceId string) *state.InstanceExec {
	ws := memdb.NewWatchSet()
	list, err := testServiceImpl(impl).state.InstanceExecListByInstanceId(instanceId, ws)
	require.NoError(t, err)
	if len(list) == 0 {
		ws.Watch(time.After(1 * time.Second))
		list, err = impl.(*service).state.InstanceExecListByInstanceId(instanceId, ws)
		require.NoError(t, err)
	}
	require.Len(t, list, 1)

	return list[0]
}

func TestServiceStartExecStream_targeted(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(err)
	client := server.TestServer(t, impl)

	// Create an instance
	instanceId, _, closer := TestEntrypoint(t, client)
	defer closer()

	// Start exec with a bad starting message
	stream, err := client.StartExecStream(ctx)
	require.NoError(err)
	require.NoError(stream.Send(&pb.ExecStreamRequest{
		Event: &pb.ExecStreamRequest_Start_{
			Start: &pb.ExecStreamRequest_Start{
				Target: &pb.ExecStreamRequest_Start_InstanceId{
					InstanceId: instanceId,
				},
				Args: []string{"foo", "bar"},
			},
		},
	}))

	// Should open
	{
		resp, err := stream.Recv()
		require.NoError(err)
		open, ok := resp.Event.(*pb.ExecStreamResponse_Open_)
		require.True(ok, "should be an open")
		require.NotNil(open)
	}

	// Get our instance exec value
	exec := testGetInstanceExec(t, impl, instanceId)
	require.Equal([]string{"foo", "bar"}, exec.Args)

	// Close send
	require.NoError(stream.CloseSend())

	// The above close send should trigger the stream to end.
	resp, err := stream.Recv()
	require.Error(err)
	require.Equal(io.EOF, err)
	require.Nil(resp)

	// The event channel should be closed
	_, active := <-exec.ClientEventCh
	require.False(active)
}
