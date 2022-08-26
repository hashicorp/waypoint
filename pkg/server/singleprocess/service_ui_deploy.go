package singleprocess

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

func (s *Service) UI_GetDeployment(
	ctx context.Context,
	req *pb.UI_GetDeploymentRequest,
) (*pb.UI_GetDeploymentResponse, error) {
	log := hclog.FromContext(ctx)

	getDeployReq := &pb.GetDeploymentRequest{Ref: req.Ref, LoadDetails: req.LoadDetails}

	deploy, err := s.GetDeployment(ctx, getDeployReq)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting deployment %q", deploy.Id)
	}

	bundle, err := s.getDeploymentBundle(ctx, deploy.Application, deploy.Workspace, deploy)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting bundle for deployment %q", deploy.Id)
	}

	if bundle.Deployment.HasEntrypointConfig {
		bundle.DeployUrl, err = s.getDeployUrl(ctx, deploy)
		if err != nil {
			log.Warn(
				"failed getting horizon hostname for deployment, but will attempt again for next deployment",
				"deployment id", deploy.Id,
				"error", err,
			)
		}
	}

	return &pb.UI_GetDeploymentResponse{
		Deployment: bundle,
	}, nil
}

func (s *Service) UI_ListDeployments(
	ctx context.Context,
	req *pb.UI_ListDeploymentsRequest,
) (*pb.UI_ListDeploymentsResponse, error) {
	log := hclog.FromContext(ctx)

	deployList, err := s.state(ctx).DeploymentList(req.Application,
		serverstate.ListWithOrder(req.Order),
		serverstate.ListWithWorkspace(req.Workspace),
		serverstate.ListWithStatusFilter(req.Status...),
		serverstate.ListWithPhysicalState(req.PhysicalState),
	)
	if err != nil {
		return nil, err
	}

	var (
		deployBundles []*pb.UI_DeploymentBundle
		deploymentUrl string
	)

	for _, deploy := range deployList {
		bundle, err := s.getDeploymentBundle(ctx, req.Application, req.Workspace, deploy)

		/*
			NOTE(briancain): Look up horizon URL ONCE, then for each deployment, append the deploy sequence
			This assumes every app deployment has the same URL, which it generally always does. We don't
			need to look up the hostname for every deployment for an application.
			We do this because `ListHostnames` service makes an HTTP request to Horizon. Doing
			this here means _for each_ deployment, we'd query Horizon for the same data, the horizon URL
			and vanity handle. We can grab it once off the first deployment, and use it for the rest of
			the bundle instead.
		*/
		if bundle.Deployment.HasEntrypointConfig {
			if deploymentUrl != "" {
				deploymentUrl, err = s.getDeployUrl(ctx, deploy)
				if err != nil {
					log.Warn(
						"failed getting horizon hostname for deployment, but will attempt again for next deployment",
						"deployment id",
						deploy.Id,
						"error", err,
					)
				}
			}
			bundle.DeployUrl = deploymentUrl
		}
		deployBundles = append(deployBundles, bundle)
	}
	return &pb.UI_ListDeploymentsResponse{
		Deployments: deployBundles,
	}, nil
}

func (s *Service) getDeploymentBundle(
	ctx context.Context,
	application *pb.Ref_Application,
	workspace *pb.Ref_Workspace,
	deploy *pb.Deployment,
) (*pb.UI_DeploymentBundle, error) {
	/*
		TODO(izaak): make this more efficient. We should be able to just grab the relevant status
		report in one go, not have to iterate over all of them.
		NOTE(brian): We need to redo how GetLatestStatusReport is implemented. Right now it just calls its
		inherited func from app operation to get the latest item in the database. For us to target a deployment or
		release status report we'll have to not use that abstraction and instead write our own query for grabbing a
		status report if a target is requested.
	*/
	statusReports, err := s.state(ctx).StatusReportList(
		application,
		// NOTE(izaak): the only implemented order for list is pb.OperationOrder_COMPLETE_TIME, which doesn't apply here.
		serverstate.ListWithWorkspace(workspace),
	)
	if err != nil {
		return nil, err
	}

	bundle := pb.UI_DeploymentBundle{
		Deployment: deploy,
	}

	var matchingStatusReport *pb.StatusReport
	for _, statusReport := range statusReports {
		switch target := statusReport.TargetId.(type) {
		case *pb.StatusReport_DeploymentId:
			if target.DeploymentId == deploy.Id {
				// (izaak) We need to find the _latest_ status report that matches. Another opportunity for efficiency by improving the statue query
				if matchingStatusReport == nil || statusReport.GeneratedTime.GetSeconds() > matchingStatusReport.GeneratedTime.Seconds {
					matchingStatusReport = statusReport
				}
			}
		}
	}
	bundle.LatestStatusReport = matchingStatusReport

	// Find artifact
	bundle.Artifact, err = s.state(ctx).ArtifactGet(&pb.Ref_Operation{
		Target: &pb.Ref_Operation_Id{
			Id: deploy.ArtifactId,
		},
	})
	if err != nil {
		return nil, err
	}

	// Find build
	if bundle.Artifact != nil {
		bundle.Build, err = s.state(ctx).BuildGet(&pb.Ref_Operation{
			Target: &pb.Ref_Operation_Id{
				Id: bundle.Artifact.BuildId,
			},
		})
		if err != nil {
			return nil, err
		}
	}
	return &bundle, nil
}

func (s *Service) getDeployUrl(ctx context.Context, deploy *pb.Deployment) (string, error) {
	// If we had no entrypoint config it is not possible for the preload URL to work.
	if !deploy.HasEntrypointConfig {
		return "", fmt.Errorf("cannot get deploymentURL for a deployment with no entrypoint config")
	}
	resp, err := s.ListHostnames(ctx, &pb.ListHostnamesRequest{
		Target: &pb.Hostname_Target{
			Target: &pb.Hostname_Target_Application{
				Application: &pb.Hostname_TargetApp{
					Application: deploy.Application,
					Workspace:   deploy.Workspace,
				},
			},
		},
	})
	if err != nil {
		return "", errors.Wrapf(err, "error listing hostnames for deployment with app %q and workspace %q", deploy.Application, deploy.Workspace)
	}
	if len(resp.Hostnames) == 0 {
		return "", fmt.Errorf("no hostnames found for app %q and workspace %q", deploy.Application, deploy.Workspace)
	}
	hostname := resp.Hostnames[0]

	return fmt.Sprintf(
		"%s--%s%s",
		hostname.Hostname,
		(&ptypes.Deployment{Deployment: deploy}).URLFragment(), // Deployment Sequence Number
		strings.TrimPrefix(hostname.Fqdn, hostname.Hostname),
	), nil
}
