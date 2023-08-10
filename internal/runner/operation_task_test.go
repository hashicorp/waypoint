// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package runner

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/go-hclog"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/singleprocess"
)

func TestRunnerTaskLauncherStart(t *testing.T) {
	if os.Getenv("WAYPOINT_BUILTIN_PLUGIN_EXE") == "" {
		t.Skip("unable to run plugins in tests without setting plugin path")

	}
	require := require.New(t)
	ctx := context.Background()
	client := singleprocess.TestServer(t)

	log := hclog.New(&hclog.LoggerOptions{
		Name:            "test-runner",
		Level:           hclog.Debug,
		IncludeLocation: true,
	})

	// Initialize our runner
	runner, err := New(
		WithClient(client),
		WithLogger(log),
	)
	require.NoError(err)
	defer runner.Close()

	// Start it
	require.NoError(runner.Start(ctx))

	job := &pb.Job{
		Operation: &pb.Job_StartTask{
			StartTask: &pb.Job_StartTaskLaunchOp{
				Params: &pb.Job_TaskPluginParams{
					PluginType: "docker",
					HclConfig:  []byte("force_pull = true\n"),
				},
				Info: &pb.TaskLaunchInfo{
					OciUrl: "ubuntu",
					Arguments: []string{
						"date",
					},
				},
			},
		},
	}

	res, err := runner.executeStartTaskOp(ctx, runner.logger, runner.ui, job)
	require.NoError(err)

	require.NotNil(t, res.StartTask)
	require.NotNil(t, res.StartTask.State)
}
