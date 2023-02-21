package runner

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/core"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func (r *Runner) executeReleaseOp(
	ctx context.Context,
	log hclog.Logger,
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

	// Our target deployment
	var target *pb.Deployment
	if op.Release.Deployment == nil {
		log.Debug("no deployment specified, using latest deployment")

		// TODO(briancain): we need a GetLatestDeployment endpoint. That would be way better.
		resp, err := r.client.ListDeployments(ctx, &pb.ListDeploymentsRequest{
			Application: job.Application,
			Workspace:   job.Workspace,
			Order: &pb.OperationOrder{
				Limit: 1, // we just care about the latest deployment
			},

			// NOTE(briancain): we _MUST_ preload the deployment details here. This is
			// because when the runner goes to invoke the app_deploy operation, it
			// assumes that the Deployment has pre-loaded the Artifact details.
			LoadDetails: pb.Deployment_ARTIFACT,
		})
		if err != nil {
			return nil, err
		}

		if len(resp.Deployments) == 0 {
			return nil, status.Error(codes.FailedPrecondition, "there are no deployments to release")
		}

		target = resp.Deployments[0]
	} else {
		target = op.Release.Deployment
	}

	// Get our last release. If it's the same generation, then release is
	// a no-op and return this value. We only do this if we have a generation.
	// We SHOULD but if we have an old client, it's possible we don't.
	var release *pb.Release
	if target.Generation != nil {
		resp, err := r.client.GetLatestRelease(ctx, &pb.GetLatestReleaseRequest{
			Application: app.Ref(),
			Workspace:   project.WorkspaceRef(),
			LoadDetails: pb.Release_DEPLOYMENT,
		})
		if status.Code(err) == codes.NotFound {
			err = nil
			resp = nil
		}
		if err != nil {
			return nil, err
		}
		if resp != nil {
			if resp.Preload != nil && resp.Preload.Deployment != nil {
				d := resp.Preload.Deployment
				if d.Generation != nil && d.Generation.Id == target.Generation.Id {
					release = resp
				}
			}
		}
	}

	// If we're pruning, then let's query the deployments we want to prune
	// ahead of time so that fails fast.
	var pruneDeploys []*pb.Deployment
	// If we are pruning deployments, we also prune the releases which
	// released the deployments we're pruning
	var pruneReleases []*pb.Release
	if op.Release.Prune {
		// Determine the number of deployments to keep around.
		retain := 2
		if op.Release.PruneRetainOverride {
			retain = int(op.Release.PruneRetain) + 1 // add 1 to make this the total number
		}

		log.Debug("pruning requested, gathering deployments to prune",
			"retain", retain)
		resp, err := r.client.ListDeployments(ctx, &pb.ListDeploymentsRequest{
			Application:   app.Ref(),
			Workspace:     project.WorkspaceRef(),
			PhysicalState: pb.Operation_CREATED,
			Order: &pb.OperationOrder{
				Order: pb.OperationOrder_COMPLETE_TIME,
				Desc:  true,
			},
		})
		if err != nil {
			return nil, err
		}

		// If we have less than the prune amount, then we do nothing. Otherwise
		// we prune away the ones we're definitely keeping.
		if len(resp.Deployments) <= retain {
			log.Debug("less than the limit deployments exists, no pruning")
			resp.Deployments = nil
		} else {
			resp.Deployments = resp.Deployments[retain:]
		}

		// Assign to short character var since we'll manipulate it a lot
		ds := make([]*pb.Deployment, 0, len(resp.Deployments))
		for _, d := range resp.Deployments {
			// If this is the deployment we're releasing, then do NOT delete it.
			if d.Id == target.Id {
				continue
			}

			// Ignore deployments with the same generation, because this
			// means that they share underlying resources.
			if target.Generation != nil && d.Generation != nil {
				if target.Generation.Id == d.Generation.Id {
					continue
				}
			}

			// TODO this should instead check against the app's platform component
			// and ignore any deployments that are NOT the app's current platform
			// component (ya dig?)
			if d.Component.Name != target.Component.Name {
				continue
			}

			// Mark for deletion
			ds = append(ds, d)
		}
		rl, err := r.client.ListReleases(ctx, &pb.ListReleasesRequest{
			Application: app.Ref(),
			Workspace:   project.WorkspaceRef(),
		})
		if err != nil {
			return nil, err
		}

		var rs []*pb.Release
		for _, release := range rl.Releases {
			for _, d := range ds {
				if release.DeploymentId == d.Id {
					rs = append(rs, release)
				}
			}

		}

		log.Info("will prune deploys", "len", len(ds))
		pruneDeploys = ds
		log.Info("will prune releases", "len", len(rs))
		pruneReleases = rs
	}

	// Do the release
	if release == nil {
		release, _, err = app.Release(ctx, target)
		if err != nil {
			return nil, err
		}
	} else {
		log.Info("not releasing since last released deploy has a matching generation",
			"gen", target.Generation.Id)
	}

	result := &pb.Job_Result{
		Release: &pb.Job_ReleaseResult{
			Release: release,
		},
	}

	// Do the pruning
	if len(pruneDeploys) > 0 {
		log.Info("pruning deploys", "len", len(pruneDeploys))
		app.UI.Output("Pruning old deployments...", terminal.WithHeaderStyle())
		for _, d := range pruneDeploys {
			app.UI.Output("Deployment: %s (v%d)", d.Id, d.Sequence, terminal.WithInfoStyle())
			if err := app.DestroyDeploy(ctx, d); err != nil {
				return result, err
			}
		}
	}
	if len(pruneReleases) > 0 {
		log.Info("pruning releases", "len", len(pruneReleases))
		app.UI.Output("Pruning old releases...", terminal.WithHeaderStyle())
		for _, release := range pruneReleases {
			app.UI.Output("Release: %s (v%d)", release.Id, release.Sequence, terminal.WithInfoStyle())
			if err := app.DestroyRelease(ctx, release); err != nil {
				return result, err
			}
		}
	}

	// Run a status report operation on the recent release
	_, err = app.ReleaseStatusReport(ctx, release)
	if err != nil {
		return nil, err
	}

	return result, nil
}
