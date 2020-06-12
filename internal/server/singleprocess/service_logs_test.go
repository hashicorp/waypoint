package singleprocess

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

func TestServiceGetLogStream(t *testing.T) {
	ctx := context.Background()

	// Create our server
	impl, err := New(testDB(t))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

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
		DeploymentId: dep.Id,
	})
	require.NoError(err)

	// Get a batch
	batch, err := logRecvClient.Recv()
	require.NoError(err)
	require.NotEmpty(batch.Lines)
	require.Len(batch.Lines, 25)
}
