package singleprocess

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/singleprocess/state"
)

func (s *service) ListInstances(
	ctx context.Context,
	req *pb.ListInstancesRequest,
) (*pb.ListInstancesResponse, error) {
	var result []*state.Instance
	var err error
	switch scope := req.Scope.(type) {
	case *pb.ListInstancesRequest_DeploymentId:
		result, err = s.state.InstancesByDeployment(scope.DeploymentId, nil)

	case *pb.ListInstancesRequest_Application_:
		result, err = s.state.InstancesByApp(
			scope.Application.Application,
			scope.Application.Workspace,
			nil,
		)

	default:
		return nil, status.Errorf(codes.FailedPrecondition,
			"scope is invalid")
	}
	if err != nil {
		return nil, err
	}

	final := make([]*pb.Instance, len(result))
	for i, r := range result {
		final[i] = r.Proto()
	}

	return &pb.ListInstancesResponse{Instances: final}, nil
}
