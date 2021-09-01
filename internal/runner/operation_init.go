package runner

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint/internal/core"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func (r *Runner) executeInitOp(
	ctx context.Context,
	log hclog.Logger,
	project *core.Project,
) (*pb.Job_Result, error) {
	client := project.Client()

	for _, name := range project.Apps() {
		_, err := client.UpsertApplication(ctx, &pb.UpsertApplicationRequest{
			Project: project.Ref(),
			Name:    name,
		})

		if err != nil {
			return nil, err
		}
	}

	return &pb.Job_Result{
		Init: &pb.Job_InitResult{},
	}, nil
}
