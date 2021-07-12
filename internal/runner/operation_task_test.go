package runner

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/go-hclog"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/singleprocess"
)

func TestRunnerTaskLauncherStart(t *testing.T) {
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
	require.NoError(runner.Start())

	job := &pb.Job{
		Operation: &pb.Job_StartTask{
			StartTask: &pb.Job_StartTaskLaunchOp{
				PluginType: "docker",
				HclConfig:  []byte("force_pull = true\n"),
				Info: &pb.TaskLaunchInfo{
					OciUrl: "ubuntu",
					Arguments: []string{
						"date",
					},
				},
			},
		},
	}

	_, err = runner.executeStartTaskOp(ctx, runner.logger, runner.ui, job)
	require.NoError(err)

}
