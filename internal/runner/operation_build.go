package runner

import (
	"context"

	"github.com/hashicorp/waypoint/internal/core"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func (r *Runner) executeBuildOp(
	ctx context.Context,
	job *pb.Job,
	project *core.Project,
) (*pb.Job_Result, error) {
	log := r.logger
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
	log.Info("DataSourceRef")
	log.Info("%+v\n", job.DataSourceRef)
	build.DataSourceRef = job.DataSourceRef

	return &pb.Job_Result{
		Build: &pb.Job_BuildResult{
			Build: build,
			Push:  push,
			DataSourceRef: &pb.Job_DataSource_Ref{
				Ref: build.DataSourceRef.Ref,
			},
		},
	}, nil
}
