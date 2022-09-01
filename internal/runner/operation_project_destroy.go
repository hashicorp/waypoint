package runner

import (
	"context"
	"github.com/hashicorp/go-hclog"
	projConfig "github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/core"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

func (r *Runner) executeDestroyProjectOp(
	ctx context.Context,
	log hclog.Logger,
	job *pb.Job,
	project *core.Project,
	cfg *projConfig.Config,
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
				Name:  destroyProjectOp.DestroyProject.Project.Project,
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
						Project:     destroyProjectOp.DestroyProject.Project.Project,
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
				// We get a copy of the project with the current workspace set to destroy the
				// deployments in each workspace
				newProject, err := project.Copy(ctx, workspace.Name, cfg)
				if err != nil {
					return nil, err
				}
				if _, err := r.executeDestroyOp(ctx,
					&pb.Job{
						Application: &pb.Ref_Application{
							Application: app,
							Project:     destroyProjectOp.DestroyProject.Project.Project,
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
					}, newProject,
				); err != nil {
					return nil, err
				}

				// Get the hostnames for our app/workspace and delete them
				if hostnames, err := client.ListHostnames(ctx, &pb.ListHostnamesRequest{
					Target: &pb.Hostname_Target{
						Target: &pb.Hostname_Target_Application{
							Application: &pb.Hostname_TargetApp{
								Application: &pb.Ref_Application{
									Application: app,
									Project:     destroyProjectOp.DestroyProject.Project.Project,
								},
								Workspace: &pb.Ref_Workspace{
									Workspace: workspace.Name,
								},
							},
						},
					},
				}); err != nil {
					if status.Code(err) == codes.FailedPrecondition {
						urlDisabledErr := "rpc error: code = FailedPrecondition desc = server doesn't have the URL service enabled"
						if err.Error() == urlDisabledErr {
							//means that the server doesn't have the URL service enabled
							//so no hostnames to list/delete
							break
						}
					}
					return nil, err
				} else {
					for _, hostname := range hostnames.Hostnames {
						if _, err = client.DeleteHostname(ctx, &pb.DeleteHostnameRequest{Hostname: hostname.Hostname}); err != nil {
							return nil, err
						}
					}
				}
			}

		}
	}

	log.Debug("Deleting DB records for project")
	// Delete the project from the database
	_, err = client.DestroyProject(ctx, &pb.DestroyProjectRequest{
		Project: destroyProjectOp.DestroyProject.Project,
	})
	// If there is an error destroying a project, reset its state to ACTIVE
	if err != nil {
		_, err := client.UpsertProject(ctx,
			&pb.UpsertProjectRequest{
				Project: &pb.Project{
					Name:  destroyProjectOp.DestroyProject.Project.Project,
					State: pb.Project_ACTIVE,
				},
			},
		)
		return nil, err
	}

	return &pb.Job_Result{
		ProjectDestroy: &pb.Job_ProjectDestroyResult{
			JobId: job.Id,
		},
	}, nil
}
