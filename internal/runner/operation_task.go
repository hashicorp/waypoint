// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package runner

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/opaqueany"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/plugin"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
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
		JsonConfig: params.HclFormat == pb.Hcl_JSON,
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
		tli.Entrypoint = sti.Entrypoint
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

	// Get our state first
	// TODO(briancain): Update this to use the Task job id ref instead of looking
	// at the task state proto
	var state *opaqueany.Any
	switch v := op.StopTask.State.(type) {
	case *pb.Job_StopTaskLaunchOp_Direct:
		log.Debug("using directly provided state")
		state = v.Direct

	case *pb.Job_StopTaskLaunchOp_StartJobId:
		// Look up the state from a start job.
		log.Debug("looking up start job to get state", "start-id", v.StartJobId)
		job, err := r.client.GetJob(ctx, &pb.GetJobRequest{
			JobId: v.StartJobId,
		})

		log.Trace(fmt.Sprintf("Got this job back: %+v", job))

		// If the job is not found, this is not an error. This means the
		// start job never ran for whatever reason and we should not stop
		// anything.
		if status.Code(err) == codes.NotFound {
			log.Warn("start job not found, not stopping anything",
				"start-id", v.StartJobId)
			return nil, nil
		} else if err != nil {
			return nil, errors.Wrapf(err, "failed to look up job with id %s", v.StartJobId)
		}

		// If the job is not in a terminal state, then its an error.
		if job.State != pb.Job_SUCCESS && job.State != pb.Job_ERROR {
			return nil, status.Errorf(codes.FailedPrecondition,
				"cannot stop task when the start job is not terminal: %q",
				job.State)
		}

		// If the job is not a start task launch operation, then error.
		_, ok := job.Operation.(*pb.Job_StartTask)
		if !ok {
			return nil, status.Errorf(codes.FailedPrecondition,
				"start job ID must reference a job with a StartTask op, got %T",
				job.Operation)
		}

		// If we have no result, do nothing.
		if job.Result == nil {
			log.Warn("start job has no result, ignoring")
			return nil, nil
		}

		result := job.Result.StartTask
		if result == nil || result.State == nil {
			log.Warn("start job has no state, ignoring")
			return nil, nil
		}

		// The state we use is the resulting state
		state = result.State
	}

	// At this point, state should not be nil. There are cases earlier
	// where we may exit early with a nil state, but we do not allow a
	// nil state here.
	if state == nil {
		return nil, status.Errorf(codes.FailedPrecondition,
			"nil start task state provided")
	}

	// Launch our plugin
	params := op.StopTask.Params
	pi, c, err := plugin.Open(ctx, log, &plugin.PluginRequest{
		Config: plugin.Config{
			Name: params.PluginType,
		},
		Dir:        "/tmp",
		ConfigData: params.HclConfig,
		JsonConfig: params.HclFormat == pb.Hcl_JSON,
		Type:       component.TaskLauncherType,
	})
	if err != nil {
		return nil, err
	}
	defer pi.Close()

	stop := c.(component.TaskLauncher).StopTaskFunc()
	_, err = pi.Invoke(ctx, log, stop, plugin.ArgNamedAny("state", state))
	if err != nil {
		return nil, err
	}

	return &pb.Job_Result{}, nil
}

func (r *Runner) executeWatchTaskOp(
	ctx context.Context,
	log hclog.Logger,
	ui terminal.UI,
	job *pb.Job,
) (*pb.Job_Result, error) {
	op, ok := job.Operation.(*pb.Job_WatchTask)
	if !ok {
		// this shouldn't happen since the call to this function is gated
		// on the above type match.
		panic("operation not expected type")
	}

	// Look up the state from a start job.
	startId := op.WatchTask.StartJob.Id
	log = log.With("start-job-id", startId)
	log.Debug("looking up start job to get state")
	job, err := r.client.GetJob(ctx, &pb.GetJobRequest{
		JobId: startId,
	})

	// If the job is not found, this is not an error. This means the
	// start job never ran for whatever reason and we should not watch
	// anything.
	if status.Code(err) == codes.NotFound {
		log.Warn("start job not found, skipping watch")
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, "failed to look up job with id %s", startId)
	}

	// If the job is not in a terminal state, then its an error.
	if job.State != pb.Job_SUCCESS && job.State != pb.Job_ERROR {
		return nil, status.Errorf(codes.FailedPrecondition,
			"cannot watch task when the start job is not terminal: %q",
			job.State)
	}

	// If the job is not a start task launch operation, then error.
	startOp, ok := job.Operation.(*pb.Job_StartTask)
	if !ok {
		return nil, status.Errorf(codes.FailedPrecondition,
			"start job ID must reference a job with a StartTask op, got %T",
			job.Operation)
	}

	// If we have no result, do nothing, assume start failed.
	if job.Result == nil {
		log.Warn("start job has no result, ignoring")
		return nil, nil
	}

	result := job.Result.StartTask
	if result == nil || result.State == nil {
		log.Warn("start job has no state, ignoring")
		return nil, nil
	}

	// The state we use is the resulting state
	state := result.State

	// At this point, state should not be nil. There are cases earlier
	// where we may exit early with a nil state, but we do not allow a
	// nil state here.
	if state == nil {
		return nil, status.Errorf(codes.FailedPrecondition,
			"nil start task state provided")
	}

	// We copy the launch params from the start task because we should
	// be using the same task launcher plugin.
	params := startOp.StartTask.Params

	// Launch our plugin
	pi, c, err := plugin.Open(ctx, log, &plugin.PluginRequest{
		Config: plugin.Config{
			Name: params.PluginType,
		},
		Dir:        "/tmp",
		ConfigData: params.HclConfig,
		JsonConfig: params.HclFormat == pb.Hcl_JSON,
		Type:       component.TaskLauncherType,
	})
	if err != nil {
		return nil, err
	}
	defer pi.Close()

	watch := c.(component.TaskLauncher).WatchTaskFunc()
	output, err := pi.Invoke(ctx, log, watch,
		plugin.ArgNamedAny("state", state),
		ui,
	)
	if err != nil {
		return nil, err
	}

	taskResult, ok := output.(*component.TaskResult)
	if !ok {
		return nil, status.Errorf(codes.FailedPrecondition,
			"plugin should've returned TaskResult, got %T", output)
	}

	return &pb.Job_Result{
		WatchTask: &pb.Job_WatchTaskResult{
			ExitCode: int32(taskResult.ExitCode),
		},
	}, nil
}
