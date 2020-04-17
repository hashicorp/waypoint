package singleprocess

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/mitchellh/devflow/internal/server"
	pb "github.com/mitchellh/devflow/internal/server/gen"
)

func TestServiceGetLogStream(t *testing.T) {
	ctx := context.Background()

	// Create our server
	impl, err := New(testDB(t))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

	// Register our instances
	configClient, err := client.EntrypointConfig(ctx, &pb.EntrypointConfigRequest{
		DeploymentId: "d",
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
		DeploymentId: "d",
	})
	require.NoError(err)

	// Get a batch
	batch, err := logRecvClient.Recv()
	require.NoError(err)
	require.NotEmpty(batch.Lines)
	require.Len(batch.Lines, 25)
}
