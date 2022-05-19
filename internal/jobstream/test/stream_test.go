package jobstream

import (
	"context"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/jobstream"
	"github.com/hashicorp/waypoint/internal/runner"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/hashicorp/waypoint/pkg/server/singleprocess"
)

func TestStream_single(t *testing.T) {
	log := hclog.L()
	ctx := context.Background()
	require := require.New(t)
	client := singleprocess.TestServer(t)

	// log.SetLevel(hclog.Trace)

	// Create, should get an ID back
	singleprocess.TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)
	resp, err := client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(err)
	require.NotNil(resp)
	require.NotEmpty(resp.JobId)

	// Create a runner
	r := runner.TestRunner(t,
		runner.WithClient(client),
		runner.WithLogger(log),
	)
	defer r.Close()
	require.NoError(r.Start(ctx))
	go r.Accept(ctx)

	// Stream should complete
	result, err := jobstream.Stream(ctx, resp.JobId, jobstream.WithClient(client))
	require.NoError(err)
	require.NotNil(result)
}
