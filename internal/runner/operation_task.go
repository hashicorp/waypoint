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

	params := op.StartTask.Params

	pi, c, err := plugin.Open(ctx, log, &plugin.PluginRequest{
		Config: plugin.Config{
			Name: params.PluginType,
		},
		Dir:        "/tmp",
		ConfigData: params.HclConfig,
		JsonConfig: params.HclFormat == pb.Project_JSON,
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

	val, err := pi.Invoke(ctx, log, fn, tli)
	if err != nil {
		return nil, err
	}

	result, err := component.ProtoAny(val)
	if err != nil {
		return nil, err
	}

	return &pb.Job_Result{
		StartTask: &pb.Job_StartTaskResult{
			State: result,
		},
	}, nil
}

func (r *Runner) executeStopTaskOp(
	ctx context.Context,
	log hclog.Logger,
	ui terminal.UI,
	job *pb.Job,
) (*pb.Job_Result, error) {
	op, ok := job.Operation.(*pb.Job_StopTask)
	if !ok {
		// this shouldn't happen since the call to this function is gated
		// on the above type match.
		panic("operation not expected type")
	}

	params := op.StopTask.Params

	pi, c, err := plugin.Open(ctx, log, &plugin.PluginRequest{
		Config: plugin.Config{
			Name: params.PluginType,
		},
		Dir:        "/tmp",
		ConfigData: params.HclConfig,
		JsonConfig: params.HclFormat == pb.Project_JSON,
		Type:       component.TaskLauncherType,
	})
	if err != nil {
		return nil, err
	}

	defer pi.Close()

	stop := c.(component.TaskLauncher).StopTaskFunc()

	_, err = pi.Invoke(ctx, log, stop, plugin.ArgNamedAny("state", op.StopTask.State))
	if err != nil {
		return nil, err
	}

	return &pb.Job_Result{}, nil
}
