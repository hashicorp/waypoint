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

	_, _, err = app.Build(ctx, core.BuildWithPush(true))
	if err != nil {
		return err
	}

	return nil
}
