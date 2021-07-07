package singleprocess

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/grpcmetadata"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

func TestServiceGetLogStreamCases(t *testing.T) {
	// Simplify writing tests
	type Req = pb.UpsertDeploymentRequest

	type instLog struct {
		pb.Waypoint_EntrypointLogStreamClient

		id string
	}

	mkinsts := func(
		t *testing.T, ctx context.Context, client pb.WaypointClient, cnt int,
	) ([]*instLog, *pb.Deployment) {
		// pb.Waypoint_EntrypointLogStreamClient, string) {
		// Register our instances
		resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
			Deployment: serverptypes.TestValidDeployment(t, &pb.Deployment{
				Component: &pb.Component{
					Name: "testapp",
				},
			}),
		})

		require.NoError(t, err)
		dep := resp.Deployment

		var out []*instLog

		for i := 0; i < cnt; i++ {
			id := string('1' + rune(i))

			configClient, err := client.EntrypointConfig(ctx, &pb.EntrypointConfigRequest{
				DeploymentId: dep.Id,
				InstanceId:   id,
			})

			require.NoError(t, err)
			_, err = configClient.Recv()
			require.NoError(t, err)

			logSendClient, err := client.EntrypointLogStream(ctx)
			require.NoError(t, err)

			out = append(out, &instLog{
				Waypoint_EntrypointLogStreamClient: logSendClient,
				id:                                 id,
			})
		}

		return out, dep
	}

	t.Run("one instance with all known lines", func(t *testing.T) {
		ctx := context.Background()

		// Create our server
		impl, err := New(WithDB(testDB(t)))
		require.NoError(t, err)
		client := server.TestServer(t, impl)

		insts, dep := mkinsts(t, ctx, client, 1)

		require := require.New(t)

		lsc := insts[0]
		id := insts[0].id

		// Create the stream and send some log messages
		for i := 0; i < 5; i++ {
			var entries []*pb.LogBatch_Entry
			for j := 0; j < 5; j++ {
				entries = append(entries, &pb.LogBatch_Entry{
					Line: strconv.Itoa(5*i + j),
				})
			}

			lsc.Send(&pb.EntrypointLogBatch{
				InstanceId: id,
				Lines:      entries,
			})
		}
		time.Sleep(100 * time.Millisecond)

		// Connect to the stream and download the logs
		logRecvClient, err := client.GetLogStream(ctx, &pb.GetLogStreamRequest{
			Scope: &pb.GetLogStreamRequest_DeploymentId{
				DeploymentId: dep.Id,
			},
		})
		require.NoError(err)

		// Get a batch
		batch, err := logRecvClient.Recv()
		require.NoError(err)
		require.NotEmpty(batch.Lines)
		require.Len(batch.Lines, 25)
	})

	t.Run("two instances with all known lines interleaved", func(t *testing.T) {
		ctx := context.Background()

		// Create our server
		impl, err := New(WithDB(testDB(t)))
		require.NoError(t, err)
		client := server.TestServer(t, impl)

		insts, dep := mkinsts(t, ctx, client, 2)

		require := require.New(t)

		start := time.Now()

		// Create the stream and send some log messages
		for i := 0; i < 5; i++ {
			var entries []*pb.LogBatch_Entry
			for j := 0; j < 5; j++ {
				entries = append(entries, &pb.LogBatch_Entry{
					Line:      strconv.Itoa(5*i + j),
					Timestamp: timestamppb.New(start.Add(time.Duration(i) * time.Minute)),
				})
			}

			for _, lsc := range insts {
				lsc.Send(&pb.EntrypointLogBatch{
					InstanceId: lsc.id,
					Lines:      entries,
				})

				time.Sleep(100 * time.Millisecond)
			}
		}
		time.Sleep(100 * time.Millisecond)

		// Connect to the stream and download the logs
		logRecvClient, err := client.GetLogStream(ctx, &pb.GetLogStreamRequest{
			Scope: &pb.GetLogStreamRequest_DeploymentId{
				DeploymentId: dep.Id,
			},
		})
		require.NoError(err)

		var ts time.Time

		// Get a batch
		for i := 0; i < 10; i++ {
			batch, err := logRecvClient.Recv()
			require.NoError(err)

			require.Len(batch.Lines, 5)

			if i%2 == 0 {
				require.Equal(insts[0].id, batch.InstanceId)
			} else {
				require.Equal(insts[1].id, batch.InstanceId)
			}

			if ts.IsZero() {
				ts = batch.Lines[0].Timestamp.AsTime()
			} else {
				bt := batch.Lines[0].Timestamp.AsTime()

				require.True(ts.Equal(bt) || ts.Before(bt))

				ts = bt
			}
		}
	})

	t.Run("one instance with some known lines and then a new instance", func(t *testing.T) {
		ctx := context.Background()

		// Create our server
		impl, err := New(WithDB(testDB(t)))
		require.NoError(t, err)
		client := server.TestServer(t, impl)

		insts, dep := mkinsts(t, ctx, client, 1)

		require := require.New(t)

		lsc := insts[0]
		id := insts[0].id

		// Create the stream and send some log messages
		for i := 0; i < 5; i++ {
			var entries []*pb.LogBatch_Entry
			for j := 0; j < 5; j++ {
				entries = append(entries, &pb.LogBatch_Entry{
					Line: strconv.Itoa(5*i + j),
				})
			}

			lsc.Send(&pb.EntrypointLogBatch{
				InstanceId: id,
				Lines:      entries,
			})
		}
		time.Sleep(100 * time.Millisecond)

		// Connect to the stream and download the logs
		logRecvClient, err := client.GetLogStream(ctx, &pb.GetLogStreamRequest{
			Scope: &pb.GetLogStreamRequest_DeploymentId{
				DeploymentId: dep.Id,
			},
		})
		require.NoError(err)

		// Get a batch
		batch, err := logRecvClient.Recv()
		require.NoError(err)
		require.NotEmpty(batch.Lines)
		require.Len(batch.Lines, 25)

		go func() {
			time.Sleep(100 * time.Millisecond)

			configClient, err := client.EntrypointConfig(ctx, &pb.EntrypointConfigRequest{
				DeploymentId: dep.Id,
				InstanceId:   "x",
			})

			require.NoError(err)
			_, err = configClient.Recv()
			require.NoError(err)

			logSendClient, err := client.EntrypointLogStream(ctx)
			require.NoError(err)

			logSendClient.Send(&pb.EntrypointLogBatch{
				InstanceId: "x",
				Lines: []*pb.LogBatch_Entry{
					{
						Line: "starting",
					},
					{
						Line: "finished",
					},
				},
			})
		}()

		// Get another batch
		batch, err = logRecvClient.Recv()
		require.NoError(err)
		require.NotEmpty(batch.Lines)
		require.Len(batch.Lines, 2)

		require.Equal("x", batch.InstanceId)

		require.Equal("starting", batch.Lines[0].Line)
		require.Equal("finished", batch.Lines[1].Line)
	})

	t.Run("one instance with some known lines and some on-demand", func(t *testing.T) {
		ctx := context.Background()

		// Create our server
		impl, err := New(WithDB(testDB(t)))
		require.NoError(t, err)
		client := server.TestServer(t, impl)

		insts, dep := mkinsts(t, ctx, client, 1)

		require := require.New(t)

		lsc := insts[0]
		id := insts[0].id

		// Create the stream and send some log messages
		for i := 0; i < 5; i++ {
			var entries []*pb.LogBatch_Entry
			for j := 0; j < 5; j++ {
				entries = append(entries, &pb.LogBatch_Entry{
					Line: strconv.Itoa(5*i + j),
				})
			}

			lsc.Send(&pb.EntrypointLogBatch{
				InstanceId: id,
				Lines:      entries,
			})
		}
		time.Sleep(100 * time.Millisecond)

		// Connect to the stream and download the logs
		logRecvClient, err := client.GetLogStream(ctx, &pb.GetLogStreamRequest{
			Scope: &pb.GetLogStreamRequest_DeploymentId{
				DeploymentId: dep.Id,
			},
		})
		require.NoError(err)

		// Get a batch
		batch, err := logRecvClient.Recv()
		require.NoError(err)
		require.NotEmpty(batch.Lines)
		require.Len(batch.Lines, 25)

		go func() {
			time.Sleep(100 * time.Millisecond)

			lsc.Send(&pb.EntrypointLogBatch{
				InstanceId: id,
				Lines: []*pb.LogBatch_Entry{
					{
						Line: "starting",
					},
					{
						Line: "finished",
					},
				},
			})
		}()

		// Get another batch
		batch, err = logRecvClient.Recv()
		require.NoError(err)
		require.NotEmpty(batch.Lines)
		require.Len(batch.Lines, 2)

		require.Equal("starting", batch.Lines[0].Line)
		require.Equal("finished", batch.Lines[1].Line)
	})
}

func TestServiceGetLogStream_depPlugin(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

	// Register our instances
	resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment: serverptypes.TestValidDeployment(t, &pb.Deployment{
			Component: &pb.Component{
				Name: "testapp",
			},
			HasLogsPlugin: true,
		}),
	})
	require.NoError(t, err)

	fakeRunner, err := server.Id()
	require.NoError(t, err)

	ctx = grpcmetadata.AddRunner(ctx, fakeRunner)

	dep := resp.Deployment

	// Connect to the stream and download the logs
	logRecvClient, err := client.GetLogStream(ctx, &pb.GetLogStreamRequest{
		Scope: &pb.GetLogStreamRequest_DeploymentId{
			DeploymentId: dep.Id,
		},
	})
	require.NoError(t, err)

	// Observe that a job to start the exec plugin has been queued
	time.Sleep(time.Second)

	jobs, err := testServiceImpl(impl).state.JobList()
	require.NoError(t, err)

	require.True(t, len(jobs) == 1)

	TestRunner(t, client, &pb.Runner{Id: fakeRunner})

	job := jobs[0]
	require.Equal(t, pb.Job_QUEUED, job.State)
	require.Equal(t, fakeRunner, job.TargetRunner.Target.(*pb.Ref_Runner_Id).Id.Id)
	require.Equal(t, resp.Deployment.Application, job.Application)

	// We force the job forward so that the server side moves forward not.

	rs, err := client.RunnerJobStream(ctx)
	require.NoError(t, err)
	require.NoError(t, rs.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Request_{
			Request: &pb.RunnerJobStreamRequest_Request{
				RunnerId: fakeRunner,
			},
		},
	}))

	jobResp, err := rs.Recv()
	require.NoError(t, err)
	assignment, ok := jobResp.Event.(*pb.RunnerJobStreamResponse_Assignment)
	require.True(t, ok, "should be an assignment")
	require.NotNil(t, assignment)
	require.Equal(t, job.Id, assignment.Assignment.Job.Id)

	require.NoError(t, rs.Send(&pb.RunnerJobStreamRequest{
		Event: &pb.RunnerJobStreamRequest_Ack_{
			Ack: &pb.RunnerJobStreamRequest_Ack{},
		},
	}))

	instanceId := assignment.Assignment.Job.Operation.(*pb.Job_Logs).Logs.InstanceId

	// Create the stream and send some log messages
	logSendClient, err := client.EntrypointLogStream(ctx)
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		var entries []*pb.LogBatch_Entry
		for j := 0; j < 5; j++ {
			entries = append(entries, &pb.LogBatch_Entry{
				Line: strconv.Itoa(5*i + j),
			})
		}

		logSendClient.Send(&pb.EntrypointLogBatch{
			InstanceId: instanceId,
			Lines:      entries,
		})
	}
	time.Sleep(100 * time.Millisecond)

	// Get a batch
	batch, err := logRecvClient.Recv()
	require.NoError(t, err)
	require.NotEmpty(t, batch.Lines)
	require.Len(t, batch.Lines, 25)
}

func TestServiceGetLogStream_byApp(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log := hclog.New(&hclog.LoggerOptions{
		Level: hclog.Trace,
	})

	// Create our server
	impl, err := New(WithDB(testDB(t)), WithLogger(log))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

	// Setup our references
	refApp := &pb.Ref_Application{
		Project:     "test",
		Application: "app",
	}
	refWs := &pb.Ref_Workspace{
		Workspace: "ws",
	}

	// Register our instances
	resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment: serverptypes.TestValidDeployment(t, &pb.Deployment{
			Application: refApp,
			Workspace:   refWs,
			Component: &pb.Component{
				Name: "testapp",
			},
			State: pb.Operation_CREATED,
			Status: &pb.Status{
				State: pb.Status_SUCCESS,
			},
		}),
	})

	require.NoError(t, err)

	dep := resp.Deployment
	configClient, err := client.EntrypointConfig(ctx, &pb.EntrypointConfigRequest{
		DeploymentId: dep.Id,
		InstanceId:   "1",
	})
	require.NoError(t, err)
	_, err = configClient.Recv()
	require.NoError(t, err)

	// Simplify writing tests
	type Req = pb.UpsertDeploymentRequest

	require := require.New(t)

	// Create the stream and send some log messages
	logSendClient, err := client.EntrypointLogStream(ctx)
	require.NoError(err)
	for i := 0; i < 5; i++ {
		var entries []*pb.LogBatch_Entry
		for j := 0; j < 5; j++ {
			entries = append(entries, &pb.LogBatch_Entry{
				Line: strconv.Itoa(5*i + j),
			})
		}

		logSendClient.Send(&pb.EntrypointLogBatch{
			InstanceId: "1",
			Lines:      entries,
		})
	}
	time.Sleep(100 * time.Millisecond)

	// Connect to the stream and download the logs
	logRecvClient, err := client.GetLogStream(ctx, &pb.GetLogStreamRequest{
		Scope: &pb.GetLogStreamRequest_Application_{
			Application: &pb.GetLogStreamRequest_Application{
				Application: refApp,
				Workspace:   refWs,
			},
		},
	})
	require.NoError(err)

	// Get a batch
	batch, err := logRecvClient.Recv()
	require.NoError(err)
	require.NotEmpty(batch.Lines)
	require.Len(batch.Lines, 25)
}
