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

// TODO(briancain): This is where we can at least update Task's start and
// stop jobs. Need to check if (Start|Stop)TaskFunc() blocks on the job
// finishing or if it fires and moves on so we can upate the Task state

// Not really sure where or how we're gonna update the RunJob though since
// that would just be running a regular operation on a runner.

// What if we mark a RunJob as a task?? Then when the job state updates, it checks
// if it's a task, looks up task by job id (self id), so it can update the Task triple state

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

	var task *pb.Task
	if job.Task != nil {
		log.Debug("updating task to starting state")
		taskResp, err := r.client.GetTask(ctx, &pb.GetTaskRequest{
			Ref: job.Task,
		})
		if err != nil {
			return nil, err
		} else {
			task = taskResp.Task
		}

		task.JobState = pb.Task_STARTING

		// Update Task state to "starting"!
		_, err = r.client.UpsertTask(ctx, &pb.UpsertTaskRequest{
			Task: task,
		})
		if err != nil {
			return nil, err
		}
	} else {
		log.Warn("no task set on job while executing start task! This is likely a bug")
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

	if task != nil {
		task.JobState = pb.Task_STARTED

		// TODO(briancain): We need to update the return value of the Task plugin
		// system to include the resource ID field that it spawned the task with.
		// Currently it only returns an Any proto, with that exact id. Instead,
		// we should make sure the Task plugin system always returns an id we can
		// extract here and place on the Task. This could work similar to how
		// the SDK system handles optional deployment URLs.
		// task.ResourceName = resourceName

		// Update Task state to "started"!
		_, err = r.client.UpsertTask(ctx, &pb.UpsertTaskRequest{
			Task: task,
		})
		if err != nil {
			return nil, err
		}
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

	var task *pb.Task
	if job.Task != nil {
		log.Debug("updating task to stopping state")
		taskResp, err := r.client.GetTask(ctx, &pb.GetTaskRequest{
			Ref: job.Task,
		})
		if err != nil {
			return nil, err
		} else {
			task = taskResp.Task
		}

		task.JobState = pb.Task_STOPPING

		// Update Task state to "stopping"!
		_, err = r.client.UpsertTask(ctx, &pb.UpsertTaskRequest{
			Task: task,
		})
		if err != nil {
			return nil, err
		}
	} else {
		log.Warn("no task set on job while executing stop task! This is likely a bug")
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

	if task != nil {
		task.JobState = pb.Task_STOPPED

		// Update Task state to "stopped"!
		_, err = r.client.UpsertTask(ctx, &pb.UpsertTaskRequest{
			Task: task,
		})
		if err != nil {
			return nil, err
		}
	}

	return &pb.Job_Result{}, nil
}
