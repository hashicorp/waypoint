package singleprocess

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

func TestServicePollQueue(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(err)
	client := server.TestServer(t, impl)

	// Create a project
	_, err = client.UpsertProject(ctx, &pb.UpsertProjectRequest{
		Project: serverptypes.TestProject(t, &pb.Project{
			Name: "A",
			DataSource: &pb.Job_DataSource{
				Source: &pb.Job_DataSource_Local{
					Local: &pb.Job_Local{},
				},
			},
			DataSourcePoll: &pb.Project_Poll{
				Enabled:  true,
				Interval: "15ms",
			},
		}),
	})
	require.NoError(err)

	// Wait a bit. The interval is so low that this should trigger
	// multiple loops through the poller. But we want to ensure we
	// have only one poll job queued.
	time.Sleep(50 * time.Millisecond)

	// We should have a single poll job
	var jobs []*pb.Job
	raw, err := testServiceImpl(impl).state.JobList()
	for _, j := range raw {
		if j.State != pb.Job_ERROR {
			jobs = append(jobs, j)
		}
	}
	require.NoError(err)
	require.Len(jobs, 1)

	// Cancel our poller to ensure it stops
	testServiceImpl(impl).Close()

	// Ensure we don't queue more jobs
	time.Sleep(50 * time.Millisecond)
	raw, err = testServiceImpl(impl).state.JobList()
	require.NoError(err)
	time.Sleep(50 * time.Millisecond)
	raw2, err := testServiceImpl(impl).state.JobList()
	require.NoError(err)
	require.Equal(len(raw), len(raw2))
}
