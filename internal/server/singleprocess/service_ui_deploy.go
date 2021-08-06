package singleprocess

import (
	"context"

	"github.com/hashicorp/waypoint/internal/server/singleprocess/state"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func (s *service) UI_ListDeployments(
	ctx context.Context,
	req *pb.UI_ListDeploymentsRequest,
) (*pb.UI_ListDeploymentsResponse, error) {
	deployList, err := s.state.DeploymentList(req.Application,
		state.ListWithOrder(req.Order),
		state.ListWithWorkspace(req.Workspace),
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

		// Always pre-populate deployment details for bundles
		if err := s.deploymentPreloadDetails(ctx, pb.Deployment_BUILD, deploy); err != nil {
			return nil, err
		}

		bundle := pb.UI_DeploymentBundle{
			Deployment:         deploy,
			LatestStatusReport: matchingStatusReport,
		}
		deployBundles = append(deployBundles, &bundle)
	}
	return &pb.UI_ListDeploymentsResponse{
		Deployments: deployBundles,
	}, nil
}
