package runner

import (
	"context"
	"github.com/hashicorp/waypoint/internal/core"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func (r *Runner) executeDeployOp(
	ctx context.Context,
	job *pb.Job,
	project *core.Project,
) (*pb.Job_Result, error) {
	app, err := project.App(job.Application.Application)
	if err != nil {
		return nil, err
	}

	op, ok := job.Operation.(*pb.Job_Deploy)
	if !ok {
		// this shouldn't happen since the call to this function is gated
		// on the above type match.
		panic("operation not expected type")
	}

	deploymentResult, err := app.Deploy(ctx, op.Deploy.Artifact)
	if err != nil {
		return nil, err
	}

	// Update to the latest deployment in order to get all the preload data.
	deployment, err := r.client.GetDeployment(ctx, &pb.GetDeploymentRequest{
		Ref: &pb.Ref_Operation{
			Target: &pb.Ref_Operation_Id{Id: deploymentResult.Id},
		},
	})
	if err != nil {
		return nil, err
	}

	return &pb.Job_Result{
		Deploy: &pb.Job_DeployResult{
			Deployment: deployment,
		},
	}, nil
}
