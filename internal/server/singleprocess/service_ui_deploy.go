package singleprocess

import (
	"context"
	"fmt"
	"strings"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/ptypes"
	"github.com/hashicorp/waypoint/internal/server/singleprocess/state"
)

func (s *service) UI_ListDeployments(
	ctx context.Context,
	req *pb.UI_ListDeploymentsRequest,
) (*pb.UI_ListDeploymentsResponse, error) {
	deployList, err := s.state.DeploymentList(req.Application,
		state.ListWithOrder(req.Order),
		state.ListWithWorkspace(req.Workspace),
		state.ListWithStatusFilter(req.Status...),
		state.ListWithPhysicalState(req.PhysicalState),
	)
	if err != nil {
		return nil, err
	}

	// TODO: make this more efficient. We should be able to just grab the relevant status report in one go, not have to
	// iterate over all of them.
	// NOTE(brian): We need to redo how GetLatestStatusReport is implemented. Right now it just calls its inherited func
	// from app operation to get the latest item in the database. For us to target a deployment or release status report
	// we'll have to not use that abstraction and instead write our own query for grabbing a status report if a target
	// is requested.
	statusReports, err := s.state.StatusReportList(
		req.Application,
		// NOTE(izaak): the only implemented order for list is pb.OperationOrder_COMPLETE_TIME, which doesn't apply here.
		state.ListWithWorkspace(req.Workspace),
	)
	if err != nil {
		return nil, err
	}

	var deployBundles []*pb.UI_DeploymentBundle

	for _, deploy := range deployList {

		bundle := pb.UI_DeploymentBundle{
			Deployment: deploy,
		}

		// Find status report
		var matchingStatusReport *pb.StatusReport
		for _, statusReport := range statusReports {
			switch target := statusReport.TargetId.(type) {
			case *pb.StatusReport_DeploymentId:
				if target.DeploymentId == deploy.Id {
					// We need to find the _latest_ status report that matches. Another opportunity for efficiency by improving the statue query
					if matchingStatusReport == nil || statusReport.GeneratedTime.GetSeconds() > matchingStatusReport.GeneratedTime.Seconds {
						matchingStatusReport = statusReport
					}
				}
			}
		}
		bundle.LatestStatusReport = matchingStatusReport

		// Find artifact
		bundle.Artifact, err = s.state.ArtifactGet(&pb.Ref_Operation{
			Target: &pb.Ref_Operation_Id{
				Id: deploy.ArtifactId,
			},
		})
		if err != nil {
			return nil, err
		}

		// Find build
		if bundle.Artifact != nil {
			bundle.Build, err = s.state.BuildGet(&pb.Ref_Operation{
				Target: &pb.Ref_Operation_Id{
					Id: bundle.Artifact.BuildId,
				},
			})
			if err != nil {
				return nil, err
			}
		}

		// Find deployment url
		// If we had no entrypoint config it is not possible for the preload URL to work.
		if deploy.HasEntrypointConfig {
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
			if err == nil && len(resp.Hostnames) > 0 {
				hostname := resp.Hostnames[0]

				bundle.DeployUrl = fmt.Sprintf(
					"%s--%s%s",
					hostname.Hostname,
					(&ptypes.Deployment{Deployment: deploy}).URLFragment(),
					strings.TrimPrefix(hostname.Fqdn, hostname.Hostname),
				)
			}
		}

		deployBundles = append(deployBundles, &bundle)
	}
	return &pb.UI_ListDeploymentsResponse{
		Deployments: deployBundles,
	}, nil
}
