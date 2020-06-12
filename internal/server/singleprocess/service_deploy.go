package singleprocess

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/singleprocess/state"
)

func (s *service) UpsertDeployment(
	ctx context.Context,
	req *pb.UpsertDeploymentRequest,
) (*pb.UpsertDeploymentResponse, error) {
	result := req.Deployment

	// If we have no ID, then we're inserting and need to generate an ID.
	insert := result.Id == ""
	if insert {
		// Get the next id
		id, err := server.Id()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "uuid generation failed: %s", err)
		}

		// Specify the id
		result.Id = id
	}

	if err := s.state.DeploymentPut(!insert, result); err != nil {
		return nil, err
	}

	return &pb.UpsertDeploymentResponse{Deployment: result}, nil
}

// TODO: test
func (s *service) ListDeployments(
	ctx context.Context,
	req *pb.ListDeploymentsRequest,
) (*pb.ListDeploymentsResponse, error) {
	result, err := s.state.DeploymentList(req.Application,
		state.ListWithStatusFilter(req.Status...),
		state.ListWithOrder(req.Order),
		state.ListWithWorkspace(req.Workspace),
	)
	if err != nil {
		return nil, err
	}

	return &pb.ListDeploymentsResponse{Deployments: result}, nil
}

// GetDeployment returns a Deployment based on ID
func (s *service) GetDeployment(
	ctx context.Context,
	req *pb.GetDeploymentRequest,
) (*pb.Deployment, error) {
	return s.state.DeploymentGet(req.DeploymentId)
}
