package runner

import (
	"context"

	"github.com/apex/log"
	"github.com/hashicorp/waypoint/internal/core"
	"github.com/hashicorp/waypoint/internal/telemetry/metrics"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func (r *Runner) executeBuildOp(
	ctx context.Context,
	job *pb.Job,
	project *core.Project,
) (*pb.Job_Result, error) {
	log.Info("==== calling timer for build")
	jt := metrics.StartTimer(metrics.JobBuild)
	defer func() {
		log.Info("==== calling Record for build")
		jt.Record()
	}()

	app, err := project.App(job.Application.Application)
	if err != nil {
		return nil, err
	}

	op, ok := job.Operation.(*pb.Job_Build)
	if !ok {
		// this shouldn't happen since the call to this function is gated
		// on the above type match.
		panic("operation not expected type")
	}

	build, push, err := app.Build(ctx, core.BuildWithPush(!op.Build.DisablePush))
	if err != nil {
		return nil, err
	}

	return &pb.Job_Result{
		Build: &pb.Job_BuildResult{
			Build: build,
			Push:  push,
		},
	}, nil
}
