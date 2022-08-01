package runner

import (
	"context"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint/internal/core"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func (r *Runner) executeDestroyProjectOp(
	ctx context.Context,
	log hclog.Logger,
	job *pb.Job,
	project *core.Project,
) (*pb.Job_Result, error) {
	client := project.Client()
	destroyProjectOp, ok := job.Operation.(*pb.Job_DestroyProject)
	if !ok {
		// this shouldn't happen since the call to this function is gated
		// on the above type match.
		panic("operation not expected type")
	}

	// Update the project state to indicate that it's being destroyed
	// TODO: Do we need to do something with the response here?
	_, err := client.UpsertProject(ctx,
		&pb.UpsertProjectRequest{
			Project: &pb.Project{
				Name:  destroyProjectOp.DestroyProject.Project.Name,
				State: pb.Project_DESTROYING,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	if !destroyProjectOp.DestroyProject.SkipDestroyResources {
		for _, app := range project.Apps() {
			// TODO: Do we need to do something with the response here?
			_, err := r.executeDestroyOp(ctx,
				&pb.Job{
					Application: &pb.Ref_Application{
						Application: app,
					},
					Operation: &pb.Job_Destroy{
						Destroy: &pb.Job_DestroyOp{
							// TODO: Do we need to specify a workspace? We want resources
							// in all workspaces to be deleted
							Target: &pb.Job_DestroyOp_Workspace{},
						},
					},
				}, project,
			)
			if err != nil {
				return nil, err
			}
		}
	}

	// Delete the project from the database
	// TODO: Do we need to do something with the response here?
	_, err = client.DestroyProject(ctx, &pb.DestroyProjectRequest{
		Project: &pb.Ref_Project{
			Project: destroyProjectOp.DestroyProject.Project.Name,
		},
	})
	if err != nil {
		return nil, err
	}

	// TODO: Need specific job result type?
	return &pb.Job_Result{}, nil
}
