package singleprocess

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/go-memdb"

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
		// For more complex deployments, sometimes it takes a bit longer for
		// the instance to respond. This blocks on a response when querying
		// a deployment by id, and lets the user know it is taking a while
		//
		// The default case is no timeout, no blocking. Make a request and return

		timeout := "0"
		if req.ConnectTimeout != "" {
			timeout = req.ConnectTimeout
		}
		connectTimeout, err := time.ParseDuration(timeout)
		if err != nil {
			return nil, err
		}

		if connectTimeout > 0 {
			_, cancel := context.WithTimeout(ctx, connectTimeout)
			defer cancel()
		}

		for {
			ws := memdb.NewWatchSet()
			result, err = s.state.InstancesByDeployment(scope.DeploymentId, ws)

			if err != nil {
				return nil, err
			}

			if len(result) > 0 || connectTimeout == 0 {
				break
			}

			// Wait for any changes
			if err := ws.WatchCtx(ctx); err != nil {
				return nil, err
			}
		}
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
