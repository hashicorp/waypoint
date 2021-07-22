package singleprocess

import (
	"context"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/singleprocess/state"
)

func (s *service) UI_ListDeployments(
	ctx context.Context,
	req *pb.UI_ListDeploymentsRequest,
) (*pb.UI_ListDeploymentsResponse, error) {
	deployments, err := s.state.DeploymentList(req.Application,
		state.ListWithWorkspace(req.Workspace),
		state.ListWithOrder(&pb.OperationOrder{
			Order: pb.OperationOrder_START_TIME,
			Desc:  true,
		}),
	)

	if err != nil {
		return nil, err
	}

	statusReports, err := s.state.StatusReportList(req.Application,
		state.ListWithWorkspace(req.Workspace),
		state.ListWithOrder(&pb.OperationOrder{
			Order: pb.OperationOrder_START_TIME,
			Desc:  true,
		}),
	)

	if err != nil {
		return nil, err
	}

	result := make([]*pb.UI_DeploymentBundle, len(deployments))

	for i, deployment := range deployments {
		var statusReport *pb.StatusReport

		for _, s := range statusReports {
			if s.GetDeploymentId() == deployment.Id {
				statusReport = s
				break
			}
		}

		bundle := pb.UI_DeploymentBundle{
			Deployment:         deployment,
			LatestStatusReport: statusReport,
		}

		result[i] = &bundle
	}

	return &pb.UI_ListDeploymentsResponse{Deployments: result}, nil
}

func (s *service) UI_GetDeployment(
	ctx context.Context,
	req *pb.UI_GetDeploymentRequest,
) (*pb.UI_DeploymentBundle, error) {
	return nil, nil
}
