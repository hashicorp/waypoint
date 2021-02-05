package runner

import (
	"context"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/core"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func (r *Runner) executeUpOp(
	ctx context.Context,
	job *pb.Job,
	project *core.Project,
) (*pb.Job_Result, error) {
	app, err := project.App(job.Application.Application)
	if err != nil {
		return nil, err
	}

	// Build it
	app.UI.Output("Building...", terminal.WithHeaderStyle())
	result, err := r.executeBuildOp(ctx, &pb.Job{
		Application: job.Application,
		Operation: &pb.Job_Build{
			Build: &pb.Job_BuildOp{},
		},
	}, project)
	if err != nil {
		return nil, err
	}
	buildResult := result.Build

	// Deploy it
	app.UI.Output("Deploying...", terminal.WithHeaderStyle())
	result, err = r.executeDeployOp(ctx, &pb.Job{
		Application: job.Application,
		Operation: &pb.Job_Deploy{
			Deploy: &pb.Job_DeployOp{
				Artifact: buildResult.Push,
			},
		},
	}, project)
	if err != nil {
		return nil, err
	}
	deployResult := result.Deploy

	// We're releasing, do that too.
	app.UI.Output("Releasing...", terminal.WithHeaderStyle())
	result, err = r.executeDeployOp(ctx, &pb.Job{
		Application: job.Application,
		Operation: &pb.Job_Release{
			Release: &pb.Job_ReleaseOp{
				Deployment: deployResult.Deployment,
				Prune:      true,
			},
		},
	}, project)
	if err != nil {
		return nil, err
	}
	releaseResult := result.Release

	// Try to get the hostname so we can build up the URL.
	var hostname *pb.Hostname
	hostnamesResp, err := r.client.ListHostnames(ctx, &pb.ListHostnamesRequest{
		Target: &pb.Hostname_Target{
			Target: &pb.Hostname_Target_Application{
				Application: &pb.Hostname_TargetApp{
					Application: deployResult.Deployment.Application,
					Workspace:   deployResult.Deployment.Workspace,
				},
			},
		},
	})
	if err == nil && len(hostnamesResp.Hostnames) > 0 {
		hostname = hostnamesResp.Hostnames[0]
	}
	var appUrl, deployUrl string
	if hostname != nil {
		appUrl = "https://" + hostname.Fqdn
	}
	if deployResult.Deployment.Preload.DeployUrl != "" {
		deployUrl = "https://" + deployResult.Deployment.Preload.DeployUrl
	}

	return &pb.Job_Result{
		Build:   buildResult,
		Deploy:  deployResult,
		Release: releaseResult,
		Up: &pb.Job_UpResult{
			ReleaseUrl: releaseResult.Release.Url,
			AppUrl:     appUrl,
			DeployUrl:  deployUrl,
		},
	}, nil
}
