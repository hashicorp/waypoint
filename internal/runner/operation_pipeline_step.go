package runner

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/core"
	"github.com/hashicorp/waypoint/internal/jobstream"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func (r *Runner) executePipelineStepOp(
	ctx context.Context,
	log hclog.Logger,
	job *pb.Job,
	project *core.Project,
) (*pb.Job_Result, error) {
	op, ok := job.Operation.(*pb.Job_PipelineStep)
	if !ok {
		// this shouldn't happen since the call to this function is gated
		// on the above type match.
		panic("operation not expected type")
	}

	switch kind := op.PipelineStep.Step.Kind.(type) {
	case *pb.Pipeline_Step_Exec_:
		return r.executePipelineStepExec(ctx, log, job, project, kind.Exec)

	default:
		return nil, status.Errorf(codes.FailedPrecondition,
			"invalid step type: %T", op.PipelineStep.Step.Kind)
	}
}

// After crafting a Pipeline Step job, the runner can simply call this func
// to queue the job and watch the job stream output to be returned.
func (r *Runner) queueAndHandleJob(
	ctx context.Context,
	log hclog.Logger,
	project *core.Project,
	job *pb.Job,
) (*pb.Job_Result, error) {
	// Queue our job
	queueResp, err := r.client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: job,
	})
	if err != nil {
		log.Warn("error queueing job", "err", err)
		return nil, err
	}

	// Get the stream
	log.Debug("opening job stream")
	stream, err := r.client.GetJobStream(ctx, &pb.GetJobStreamRequest{
		JobId: queueResp.JobId,
	})
	if err != nil {
		return nil, err
	}

	// Watch job
	streamUI := &jobstream.UI{
		UI:  project.UI,
		Log: log.Named("ui"),
	}
	for {
		resp, err := stream.Recv()
		if err != nil {
			return nil, err
		}
		if resp == nil {
			// This shouldn't happen, but if it does, just ignore it.
			log.Warn("nil response received, ignoring")
			continue
		}

		switch event := resp.Event.(type) {
		case *pb.GetJobStreamResponse_Complete_:
			if event.Complete.Error == nil {
				result := event.Complete.Result
				result.PipelineStep = &pb.Job_PipelineStepResult{
					Result: status.New(codes.OK, "").Proto(),
				}

				return result, nil
			}

			st := status.FromProto(event.Complete.Error)
			log.Warn("job failed", "code", st.Code(), "message", st.Message())
			return &pb.Job_Result{
				PipelineStep: &pb.Job_PipelineStepResult{
					Result: event.Complete.Error,
				},
			}, nil

		case *pb.GetJobStreamResponse_Error_:
			st := status.FromProto(event.Error.Error)
			log.Warn("job failed", "code", st.Code(), "message", st.Message())
			return &pb.Job_Result{
				PipelineStep: &pb.Job_PipelineStepResult{
					Result: event.Error.Error,
				},
			}, nil

		case *pb.GetJobStreamResponse_Terminal_:
			if err := streamUI.Write(event.Terminal.Events); err != nil {
				log.Warn("job stream UI failure", "err", err)
			}

		case *pb.GetJobStreamResponse_State_:
			// Ignore state changes
			log.Debug("child job state change", "state", event)

		default:
			log.Warn("unknown stream event", "event", resp.Event)
		}
	}
}

func (r *Runner) executePipelineStepExec(
	ctx context.Context,
	log hclog.Logger,
	job *pb.Job,
	project *core.Project,
	exec *pb.Pipeline_Step_Exec,
) (*pb.Job_Result, error) {
	var entrypoint []string
	if exec.Command != "" {
		entrypoint = []string{exec.Command}
	}

	// Create a new job that launches our task to run. This is heavily based
	// on the incoming job so we can inherit a lot of the properties. The key
	// change is that we specify a noop operation and task override so that
	// we run ODR with a custom task (which does something) and the actual
	// operation does nothing.
	newJob := &pb.Job{
		Application:         job.Application,
		Workspace:           job.Workspace,
		OndemandRunner:      job.OndemandRunner,
		Labels:              job.Labels,
		DataSource:          job.DataSource,
		DataSourceOverrides: job.DataSourceOverrides,
		WaypointHcl:         job.WaypointHcl,
		Variables:           job.Variables,

		// Must target "any" runner so we get ODR to work.
		TargetRunner: &pb.Ref_Runner{
			Target: &pb.Ref_Runner_Any{
				Any: &pb.Ref_RunnerAny{},
			},
		},

		// Noop so we can skip it
		Operation: &pb.Job_Noop_{
			Noop: &pb.Job_Noop{},
		},

		// Our custom overrides so we can run our task
		OndemandRunnerTask: &pb.Job_TaskOverride{
			// Skip because its just a noop.
			SkipOperation: true,

			LaunchInfo: &pb.TaskLaunchInfo{
				OciUrl:     exec.Image,
				Entrypoint: entrypoint,
				Arguments:  exec.Args,
			},
		},
	}

	return r.queueAndHandleJob(ctx, log, project, newJob)
}
