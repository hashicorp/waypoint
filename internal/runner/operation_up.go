// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package runner

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/core"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func (r *Runner) executeUpOp(
	ctx context.Context,
	log hclog.Logger,
	job *pb.Job,
	project *core.Project,
) (*pb.Job_Result, error) {
	app, err := project.App(job.Application.Application)
	if err != nil {
		return nil, err
	}

	opRaw, ok := job.Operation.(*pb.Job_Up)
	if !ok {
		// this shouldn't happen since the call to this function is gated
		// on the above type match.
		panic("operation not expected type")
	}
	op := opRaw.Up

	// Setup our default options
	if op.Release == nil {
		op.Release = &pb.Job_ReleaseOp{Prune: true}
	}

	// TODO: output the context in use (maybe only if non-default)
	appName := app.Ref().Application

	// Build it
	app.UI.Output("Building %s...", appName, terminal.WithHeaderStyle())
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
	app.UI.Output("Deploying %s...", appName, terminal.WithHeaderStyle())
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

	// Status Report for Deployments
	app.UI.Output("")
	result, err = r.executeStatusReportOp(ctx, log, &pb.Job{
		Application: job.Application,
		Operation: &pb.Job_StatusReport{
			StatusReport: &pb.Job_StatusReportOp{
				Target: &pb.Job_StatusReportOp_Deployment{
					Deployment: deployResult.Deployment,
				},
			},
		},
	}, project)
	if err != nil {
		return nil, err
	}
	statusReportResult := result.StatusReport

	// We're releasing, do that too.
	app.UI.Output("Releasing %s...", appName, terminal.WithHeaderStyle())
	op.Release.Deployment = deployResult.Deployment
	result, err = r.executeReleaseOp(ctx, log, &pb.Job{
		Application: job.Application,
		Operation: &pb.Job_Release{
			Release: op.Release,
		},
	}, project)
	if err != nil {
		return nil, err
	}
	releaseResult := result.Release

	if releaseResult.Release.Unimplemented {
		app.UI.Output("No release phase specified, skipping...")
	}

	// NOTE(briancain): Because executeReleaseOp returns an initialized struct
	// of release results, we need this deep check here to really ensure that a
	// release actually happened, otherwise we'd attempt to run a status report
	// on a nil release
	if releaseResult != nil && releaseResult.Release != nil &&
		releaseResult.Release.Release != nil {
		// Status Report for Releases
		app.UI.Output("")
		result, err = r.executeStatusReportOp(ctx, log, &pb.Job{
			Application: job.Application,
			Operation: &pb.Job_StatusReport{
				StatusReport: &pb.Job_StatusReportOp{
					Target: &pb.Job_StatusReportOp_Release{
						Release: releaseResult.Release,
					},
				},
			},
		}, project)
		if err != nil {
			return nil, err
		}
		statusReportResult = result.StatusReport
	}

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
		StatusReport: statusReportResult,
	}, nil
}
