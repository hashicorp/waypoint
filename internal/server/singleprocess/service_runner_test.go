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
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

// Complete happy path job stream
func TestServiceRunnerJobStream_complete(t *testing.T) {
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
	id, _ := TestRunner(t, client, nil)

	// Start a job request
	stream, err := client.RunnerJobStream(ctx)
	require.NoError(err)
	require.NoError(stream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Request_{
			Request: &pb.RunnerJobStreamRequest_Request{
				RunnerId: id,
			},
		},
	}))

	// Wait for assignment and ack
	{
		resp, err := stream.Recv()
		require.NoError(err)
		assignment, ok := resp.Event.(*pb.RunnerJobStreamResponse_Assignment)
		require.True(ok, "should be an assignment")
		require.NotNil(assignment)
		require.Equal(queueResp.JobId, assignment.Assignment.Job.Id)

		require.NoError(stream.Send(&pb.RunnerJobStreamRequest{
			Event: &pb.RunnerJobStreamRequest_Ack_{
				Ack: &pb.RunnerJobStreamRequest_Ack{},
			},
		}))
	}

	// Complete the job
	require.NoError(stream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Complete_{
			Complete: &pb.RunnerJobStreamRequest_Complete{},
		},
	}))

	// Should be done
	_, err = stream.Recv()
	require.Error(err)
	require.Equal(io.EOF, err)

	// Query our job and it should be done
	job, err := testServiceImpl(impl).state.JobById(queueResp.JobId, nil)
	require.NoError(err)
	require.Equal(pb.Job_SUCCESS, job.State)
}

func TestServiceRunnerJobStream_badOpen(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	impl, err := New(testDB(t))
	require.NoError(err)
	client := server.TestServer(t, impl)

	// Start exec with a bad starting message
	stream, err := client.RunnerJobStream(ctx)
	require.NoError(err)
	require.NoError(stream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Ack_{
			Ack: &pb.RunnerJobStreamRequest_Ack{},
		},
	}))

	// Wait for data
	resp, err := stream.Recv()
	require.Error(err)
	require.Equal(codes.FailedPrecondition, status.Code(err))
	require.Nil(resp)
}

func TestServiceRunnerJobStream_errorBeforeAck(t *testing.T) {
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
	id, _ := TestRunner(t, client, nil)

	// Start a job request
	stream, err := client.RunnerJobStream(ctx)
	require.NoError(err)
	require.NoError(stream.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Request_{
			Request: &pb.RunnerJobStreamRequest_Request{
				RunnerId: id,
			},
		},
	}))

	// Wait for assignment and DONT ack, send an error instead
	{
		resp, err := stream.Recv()
		require.NoError(err)
		assignment, ok := resp.Event.(*pb.RunnerJobStreamResponse_Assignment)
		require.True(ok, "should be an assignment")
		require.NotNil(assignment)
		require.Equal(queueResp.JobId, assignment.Assignment.Job.Id)

		require.NoError(stream.Send(&pb.RunnerJobStreamRequest{
			Event: &pb.RunnerJobStreamRequest_Error_{
				Error: &pb.RunnerJobStreamRequest_Error{
					Error: status.Newf(codes.Unknown, "error").Proto(),
				},
			},
		}))
	}

	// Should be done
	_, err = stream.Recv()
	require.Error(err)
	require.Equal(io.EOF, err)

	// Query our job and it should be queued again
	job, err := testServiceImpl(impl).state.JobById(queueResp.JobId, nil)
	require.NoError(err)
	require.Equal(pb.Job_QUEUED, job.State)
}

func TestServiceRunnerGetDeploymentConfig(t *testing.T) {
	ctx := context.Background()

	t.Run("with no server config", func(t *testing.T) {
		require := require.New(t)

		// Create our server
		impl, err := New(testDB(t))
		require.NoError(err)
		client := server.TestServer(t, impl)

		// Request deployment config
		_, err = client.RunnerGetDeploymentConfig(ctx, &pb.RunnerGetDeploymentConfigRequest{})
		require.Error(err)
		require.Equal(codes.Aborted, status.Code(err))
	})

	t.Run("with server config", func(t *testing.T) {
		require := require.New(t)

		// Create our server
		impl, err := New(testDB(t))
		require.NoError(err)
		client := server.TestServer(t, impl)

		// Set some config
		_, err = client.SetServerConfig(ctx, &pb.SetServerConfigRequest{
			Config: serverptypes.TestServerConfig(t, nil),
		})
		require.NoError(err)

		// Request deployment config
		resp, err := client.RunnerGetDeploymentConfig(ctx, &pb.RunnerGetDeploymentConfigRequest{})
		require.NoError(err)
		require.NotEmpty(resp.ServerAddr)
	})
}
