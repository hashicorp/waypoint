package runner

import (
	"context"

	"github.com/hashicorp/waypoint/internal/core"
	servercomponent "github.com/hashicorp/waypoint/internal/server/component"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/sdk/component"
)

func (r *Runner) executeReleaseOp(
	ctx context.Context,
	job *pb.Job,
	project *core.Project,
) (*pb.Job_Result, error) {
	app, err := project.App(job.Application.Application)
	if err != nil {
		return nil, err
	}

	op, ok := job.Operation.(*pb.Job_Release)
	if !ok {
		// this shouldn't happen since the call to this function is gated
		// on the above type match.
		panic("operation not expected type")
	}

	targets := make([]component.ReleaseTarget, len(op.Release.TrafficSplit.Targets))
	for i, split := range op.Release.TrafficSplit.Targets {
		// Get the deployment
		deployment, err := r.client.GetDeployment(ctx, &pb.GetDeploymentRequest{
			DeploymentId: split.DeploymentId,
		})
		if err != nil {
			return nil, err
		}

		targets[i] = component.ReleaseTarget{
			Deployment: servercomponent.Deployment(deployment),
			Percent:    uint(split.Percent),
		}
	}

	release, _, err := app.Release(ctx, targets)
	if err != nil {
		return nil, err
	}

	return &pb.Job_Result{
		Release: &pb.Job_ReleaseResult{
			Release: release,
		},
	}, nil
}
