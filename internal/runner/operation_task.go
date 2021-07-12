package runner

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/plugin"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func (r *Runner) executeStartTaskOp(
	ctx context.Context,
	log hclog.Logger,
	ui terminal.UI,
	job *pb.Job,
) (*pb.Job_Result, error) {
	op, ok := job.Operation.(*pb.Job_StartTask)
	if !ok {
		// this shouldn't happen since the call to this function is gated
		// on the above type match.
		panic("operation not expected type")
	}

	pi, c, err := plugin.OpenPlugin(ctx, log, &plugin.PluginRequest{
		Config: plugin.Config{
			Name: op.StartTask.PluginType,
		},
		Dir:        "/tmp",
		ConfigData: op.StartTask.HclConfig,
		JsonConfig: op.StartTask.HclFormat == pb.Project_JSON,
		Type:       component.TaskLauncherType,
	})
	if err != nil {
		return nil, err
	}

	defer pi.Close()

	tli := &component.TaskLaunchInfo{}

	sti := op.StartTask.Info

	if sti != nil {
		tli.OciUrl = sti.OciUrl
		tli.EnvironmentVariables = sti.EnvironmentVariables
		tli.Arguments = sti.Arguments
	}

	fn := c.(component.TaskLauncher).StartTaskFunc()

	_, err = pi.Invoke(ctx, log, fn, tli)
	if err != nil {
		return nil, err
	}

	return &pb.Job_Result{}, nil
}
