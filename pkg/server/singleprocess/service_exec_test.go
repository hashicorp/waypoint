// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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

	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/grpcmetadata"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/hashicorp/waypoint/pkg/serverstate"
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

	// Start an exec stream properly
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
	list, err := testStateInmem(impl).InstanceExecListByInstanceId(ctx, instanceId, ws)
	require.NoError(err)
	if len(list) == 0 {
		ws.Watch(time.After(1 * time.Second))
		list, err = testStateInmem(impl).InstanceExecListByInstanceId(ctx, instanceId, ws)
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

func TestServiceStartExecStream_eventError(t *testing.T) {
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
	list, err := testStateInmem(impl).InstanceExecListByInstanceId(ctx, instanceId, ws)
	require.NoError(err)
	if len(list) == 0 {
		ws.Watch(time.After(1 * time.Second))
		list, err = testStateInmem(impl).InstanceExecListByInstanceId(ctx, instanceId, ws)
		require.NoError(err)
	}
	require.Len(list, 1)
	exec := list[0]

	// Send an exit event
	exec.EntrypointEventCh <- &pb.EntrypointExecRequest{
		Event: &pb.EntrypointExecRequest_Error_{
			Error: &pb.EntrypointExecRequest_Error{
				Error: status.New(codes.DataLoss, "this is a bad thing").Proto(),
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
	list, err := testStateInmem(impl).InstanceExecListByInstanceId(ctx, instanceId, ws)
	require.NoError(err)
	if len(list) == 0 {
		ws.Watch(time.After(1 * time.Second))
		list, err = testStateInmem(impl).InstanceExecListByInstanceId(ctx, instanceId, ws)
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

func testGetInstanceExec(t *testing.T, impl pb.WaypointServer, instanceId string) *serverstate.InstanceExec {
	ctx := context.Background()
	ws := memdb.NewWatchSet()
	list, err := testStateInmem(impl).InstanceExecListByInstanceId(ctx, instanceId, ws)
	require.NoError(t, err)
	if len(list) == 0 {
		ws.Watch(time.After(1 * time.Second))
		list, err = testStateInmem(impl).InstanceExecListByInstanceId(ctx, instanceId, ws)
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

func TestServiceStartExecStream_startPlugin(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(err)
	client := server.TestServer(t, impl)

	// Create an instance
	resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment: serverptypes.TestValidDeployment(t, &pb.Deployment{
			Component: &pb.Component{
				Name: "testapp",
			},

			HasExecPlugin: true,
		}),
	})
	require.NoError(err)

	deploymentId := resp.Deployment.Id

	fakeRunner, err := server.Id()
	require.NoError(err)

	ctx = grpcmetadata.AddRunner(ctx, fakeRunner)

	// Start an exec stream
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

	// Observe that a job to start the exec plugin has been queued
	time.Sleep(time.Second)

	jobs, _, err := testServiceImpl(impl).state(ctx).JobList(ctx, &pb.ListJobsRequest{})
	require.NoError(err)

	require.True(len(jobs) == 1)

	job := jobs[0]
	require.Equal(pb.Job_QUEUED, job.State)
	require.Equal(fakeRunner, job.TargetRunner.Target.(*pb.Ref_Runner_Id).Id.Id)
	require.Equal(resp.Deployment.Application, job.Application)
}

func TestService_waitOnJobStarted(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	require := require.New(t)

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(err)
	client := server.TestServer(t, impl)

	// Create an instance
	resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment: serverptypes.TestValidDeployment(t, &pb.Deployment{
			Component: &pb.Component{
				Name: "testapp",
			},

			HasExecPlugin: true,
		}),
	})
	require.NoError(err)

	job := &pb.Job{
		Id:          "aabbcc",
		Workspace:   resp.Deployment.Workspace,
		Application: resp.Deployment.Application,
		Operation: &pb.Job_Noop_{
			Noop: &pb.Job_Noop{},
		},
		DataSource: &pb.Job_DataSource{
			Source: &pb.Job_DataSource_Local{
				Local: &pb.Job_Local{},
			},
		},
		TargetRunner: &pb.Ref_Runner{
			Target: &pb.Ref_Runner_Any{
				Any: &pb.Ref_RunnerAny{},
			},
		},
	}

	s := testServiceImpl(impl)

	// Queue the job
	err = s.state(ctx).JobCreate(ctx, job)
	require.NoError(err)

	go func() {
		time.Sleep(time.Second)
		s.state(ctx).JobExpire(ctx, job.Id)
	}()

	ts := time.Now()
	js, err := s.waitOnJobStarted(ctx, job.Id)
	require.NoError(err)
	dur := time.Since(ts)

	require.InDelta(1.0, dur.Seconds(), 1.0)

	require.Equal(pb.Job_ERROR, js)
}
