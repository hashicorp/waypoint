package runner

import (
	"context"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint/internal/core"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	empty "google.golang.org/protobuf/types/known/emptypb"
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
			// Get the workspaces for the app we're currently destroying
			workspaces, err := client.ListWorkspaces(ctx, &pb.ListWorkspacesRequest{
				Scope: &pb.ListWorkspacesRequest_Application{
					Application: &pb.Ref_Application{
						Project:     destroyProjectOp.DestroyProject.Project.Name,
						Application: app,
					},
				},
			})
			if err != nil {
				return nil, err
			}
			// Destroy the resources in the workspaces for the app
			for _, workspace := range workspaces.Workspaces {
				log.Debug("Destroying resources in workspace %s", workspace.Name)
				// TODO: The destroy operation doesn't currently respect the Workspace
				// being set on the Job - this should be updated
				_, err := r.executeDestroyOp(ctx,
					&pb.Job{
						Application: &pb.Ref_Application{
							Application: app,
						},
						Operation: &pb.Job_Destroy{
							Destroy: &pb.Job_DestroyOp{
								Target: &pb.Job_DestroyOp_Workspace{
									Workspace: &empty.Empty{},
								},
							},
						},
						Workspace: &pb.Ref_Workspace{
							Workspace: workspace.Name,
						},
					}, project,
				)
				if err != nil {
					return nil, err
				}
			}

		}
	}

	// Delete the project from the database
	_, err = client.DestroyProject(ctx, &pb.DestroyProjectRequest{
		Project: &pb.Ref_Project{
			Project: destroyProjectOp.DestroyProject.Project.Name,
		},
	})
	if err != nil {
		return nil, err
	}

	return &pb.Job_Result{}, nil
}
