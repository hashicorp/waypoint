package client

import (
	"context"
	"testing"

	"github.com/hashicorp/waypoint/internal/runner"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/server/singleprocess"
)

func init() {
	hclog.L().SetLevel(hclog.Trace)
}

func TestProjectNoop(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)
	client := singleprocess.TestServer(t)

	// Start a local runner
	testRunner, err := runner.New(runner.WithClient(client))
	require.Nil(err)

	require.NoError(testRunner.Start())

	go func() {
		require.NoError(testRunner.Accept(ctx))
	}()

	// Build our client
	c := TestProject(t, client, WithExecuteJobsLocally(testRunner.Id()))
	app := c.App(TestApp(t, c))

	// TODO(mitchellh): once we have an API to list jobs, verify we have
	// no jobs, and then verify we execute a job after.

	// Noop
	require.NoError(app.Noop(ctx))
}
