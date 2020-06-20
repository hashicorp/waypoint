package singleprocess

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

func TestServiceQueueJob(t *testing.T) {
	ctx := context.Background()

	// Create our server
	impl, err := New(testDB(t))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

	// Simplify writing tests
	type Req = pb.QueueJobRequest

	t.Run("create success", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.QueueJob(ctx, &Req{
			Job: serverptypes.TestJobNew(t, nil),
		})
		require.NoError(err)
		require.NotNil(resp)
		require.NotEmpty(resp.JobId)

		// Job should exist and be queued
		job, err := testServiceImpl(impl).state.JobById(resp.JobId, nil)
		require.NoError(err)
		require.Equal(pb.Job_QUEUED, job.State)
	})
}

func TestServiceGetJobStream_complete(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	impl, err := New(testDB(t))
	require.NoError(err)
	client := server.TestServer(t, impl)

	// Create a job
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{Job: serverptypes.TestJobNew(t, nil)})
	require.NoError(err)
	require.NotNil(queueResp)
	require.NotEmpty(queueResp.JobId)

	// Register our runner
	configStream, err := client.RunnerConfig(ctx, &pb.RunnerConfigRequest{
		Id: "R_A",
	})
	require.NoError(err)
	defer configStream.CloseSend()
	_, err = configStream.Recv()
	require.NoError(err)

	// Start a job request
	runnerStream, err := client.RunnerJobStream(ctx)
	require.NoError(err)
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Request_{
			Request: &pb.RunnerJobStreamRequest_Request{
				RunnerId: "R_A",
			},
		},
	}))

	// Wait for assignment and ack
	{
		resp, err := runnerStream.Recv()
		require.NoError(err)
		assignment, ok := resp.Event.(*pb.RunnerJobStreamResponse_Assignment)
		require.True(ok, "should be an assignment")
		require.NotNil(assignment)
		require.Equal(queueResp.JobId, assignment.Assignment.Job.Id)

		require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
			Event: &pb.RunnerJobStreamRequest_Ack_{
				Ack: &pb.RunnerJobStreamRequest_Ack{},
			},
		}))
	}

	// Get our job stream and verify we open
	stream, err := client.GetJobStream(ctx, &pb.GetJobStreamRequest{JobId: queueResp.JobId})
	require.NoError(err)
	{
		resp, err := stream.Recv()
		require.NoError(err)
		open, ok := resp.Event.(*pb.GetJobStreamResponse_Open_)
		require.True(ok, "should be an open")
		require.NotNil(open)

	}

	// We need to give the stream time to initialize the output readers
	// so our output below doesn't become buffered. This isn't really that
	// brittle and 100ms should be more than enough.
	time.Sleep(100 * time.Millisecond)

	// Send some output
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Terminal{
			Terminal: &pb.GetJobStreamResponse_Terminal{
				Lines: []*pb.GetJobStreamResponse_Terminal_Line{
					{Raw: "hello"},
					{Raw: "world"},
				},
			},
		},
	}))

	// Wait for output
	{
		resp, err := stream.Recv()
		require.NoError(err)
		event, ok := resp.Event.(*pb.GetJobStreamResponse_Terminal_)
		require.True(ok, "should be terminal data")
		require.NotNil(event)
		require.False(event.Terminal.Buffered, 2)
		require.Len(event.Terminal.Lines, 2)

	}

	// Complete the job
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Complete_{
			Complete: &pb.RunnerJobStreamRequest_Complete{},
		},
	}))

	// Wait for completion
	{
		resp, err := stream.Recv()
		require.NoError(err)
		event, ok := resp.Event.(*pb.GetJobStreamResponse_Complete_)
		require.True(ok, "should be terminal data")
		require.NotNil(event)
		require.Nil(event.Complete.Error)
	}
}

func TestServiceGetJobStream_bufferedData(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	impl, err := New(testDB(t))
	require.NoError(err)
	client := server.TestServer(t, impl)

	// Create a job
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{Job: serverptypes.TestJobNew(t, nil)})
	require.NoError(err)
	require.NotNil(queueResp)
	require.NotEmpty(queueResp.JobId)

	// Register our runner
	configStream, err := client.RunnerConfig(ctx, &pb.RunnerConfigRequest{
		Id: "R_A",
	})
	require.NoError(err)
	defer configStream.CloseSend()
	_, err = configStream.Recv()
	require.NoError(err)

	// Start a job request
	runnerStream, err := client.RunnerJobStream(ctx)
	require.NoError(err)
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Request_{
			Request: &pb.RunnerJobStreamRequest_Request{
				RunnerId: "R_A",
			},
		},
	}))

	// Wait for assignment and ack
	{
		resp, err := runnerStream.Recv()
		require.NoError(err)
		assignment, ok := resp.Event.(*pb.RunnerJobStreamResponse_Assignment)
		require.True(ok, "should be an assignment")
		require.NotNil(assignment)
		require.Equal(queueResp.JobId, assignment.Assignment.Job.Id)

		require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
			Event: &pb.RunnerJobStreamRequest_Ack_{
				Ack: &pb.RunnerJobStreamRequest_Ack{},
			},
		}))
	}

	// Send some output
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Terminal{
			Terminal: &pb.GetJobStreamResponse_Terminal{
				Lines: []*pb.GetJobStreamResponse_Terminal_Line{
					{Raw: "hello"},
					{Raw: "world"},
				},
			},
		},
	}))

	// Complete the job
	require.NoError(runnerStream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Complete_{
			Complete: &pb.RunnerJobStreamRequest_Complete{},
		},
	}))

	// Should be done
	_, err = runnerStream.Recv()
	require.Error(err)
	require.Equal(io.EOF, err)

	// Get our job stream and verify we open
	stream, err := client.GetJobStream(ctx, &pb.GetJobStreamRequest{JobId: queueResp.JobId})
	require.NoError(err)
	{
		resp, err := stream.Recv()
		require.NoError(err)
		open, ok := resp.Event.(*pb.GetJobStreamResponse_Open_)
		require.True(ok, "should be an open")
		require.NotNil(open)

	}

	// Wait for output
	{
		resp, err := stream.Recv()
		require.NoError(err)
		event, ok := resp.Event.(*pb.GetJobStreamResponse_Terminal_)
		require.True(ok, "should be terminal data")
		require.NotNil(event)
		require.True(event.Terminal.Buffered)
		require.Len(event.Terminal.Lines, 2)

	}

	// Wait for completion
	{
		resp, err := stream.Recv()
		require.NoError(err)
		event, ok := resp.Event.(*pb.GetJobStreamResponse_Complete_)
		require.True(ok, "should be terminal data")
		require.NotNil(event)
		require.Nil(event.Complete.Error)
	}
}
