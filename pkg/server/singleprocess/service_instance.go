// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package singleprocess

import (
	"context"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-memdb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

func (s *Service) ListInstances(
	ctx context.Context,
	req *pb.ListInstancesRequest,
) (*pb.ListInstancesResponse, error) {
	var result []*serverstate.Instance
	var err error
	switch scope := req.Scope.(type) {
	case *pb.ListInstancesRequest_DeploymentId:
		// For more complex deployments, sometimes it takes a bit longer for
		// the instance to respond. This blocks on a response when querying
		// a deployment by id, and lets the user know it is taking a while
		//
		// The default case is no timeout, no blocking. Make a request and return

		if req.WaitTimeout != "" {
			connectTimeout, err := time.ParseDuration(req.WaitTimeout)
			if err != nil {
				return nil, status.Errorf(codes.FailedPrecondition, "Error parsing wait_timeout: %s", err)
			}

			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, connectTimeout)
			defer cancel()
		}

		for {
			ws := memdb.NewWatchSet()
			result, err = s.state(ctx).InstancesByDeployment(ctx, scope.DeploymentId, ws)

			if err != nil {
				return nil, hcerr.Externalize(
					hclog.FromContext(ctx),
					err,
					"failed to list instances by deployment",
					"deployment_id",
					scope.DeploymentId,
				)
			}

			if len(result) > 0 || req.WaitTimeout == "" {
				break
			}

			// Wait for any changes
			if err := ws.WatchCtx(ctx); err != nil {
				return nil, hcerr.Externalize(
					hclog.FromContext(ctx),
					err,
					"error watching for instance changes",
				)
			}
		}
	case *pb.ListInstancesRequest_Application_:
		result, err = s.state(ctx).InstancesByApp(
			ctx,
			scope.Application.Application,
			scope.Application.Workspace,
			nil,
		)

	default:
		return nil, status.Errorf(codes.FailedPrecondition,
			"scope is invalid")
	}
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"error listing instances",
		)
	}

	final := make([]*pb.Instance, len(result))
	for i, r := range result {
		final[i] = r.Proto()
	}

	return &pb.ListInstancesResponse{Instances: final}, nil
}
