package runner

import (
	"context"

	"github.com/hashicorp/waypoint/internal/core"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func (r *Runner) executeBuildOp(ctx context.Context, job *pb.Job, project *core.Project) error {
	app, err := project.App(job.Application.Application)
	if err != nil {
		return err
	}

	op, ok := job.Operation.(*pb.Job_Build)
	if !ok {
		// this shouldn't happen since the call to this function is gated
		// on the above type match.
		panic("operation not expected type")
	}

	_, _, err = app.Build(ctx, core.BuildWithPush(!op.Build.DisablePush))
	if err != nil {
		return err
	}

	return nil
}
