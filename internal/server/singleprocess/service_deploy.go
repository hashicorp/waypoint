package singleprocess

import (
	"context"
	"time"

	"github.com/hashicorp/go-hclog"
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

	// This requires: (1) URL service is enabled (2) auto hostname isn't
	// explicitly set to false in the request and (3) either the server
	// default is true or we explicitly ask for it.
	if s.urlClient != nil &&
		req.AutoHostname != pb.UpsertDeploymentRequest_FALSE &&
		(s.urlConfig.AutomaticAppHostname || req.AutoHostname == pb.UpsertDeploymentRequest_TRUE) {
		// Our hostname target. We need this to automatically create a hostname.
		target := &pb.Hostname_Target{
			Target: &pb.Hostname_Target_Application{
				Application: &pb.Hostname_TargetApp{
					Application: result.Application,
					Workspace:   result.Workspace,
				},
			},
		}

		// Create a context that will timeout relatively quickly
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		// Create the hostname. We ignore errors.
		_, err := s.createHostnameIfNotExist(ctx, target)
		if err != nil {
			log := hclog.FromContext(ctx)
			log.Info("error creating default hostname", "err", err)
		}
	}

	return &pb.UpsertDeploymentResponse{Deployment: result}, nil
}

func (s *service) ListDeployments(
	ctx context.Context,
	req *pb.ListDeploymentsRequest,
) (*pb.ListDeploymentsResponse, error) {
	result, err := s.state.DeploymentList(req.Application,
		state.ListWithStatusFilter(req.Status...),
		state.ListWithOrder(req.Order),
		state.ListWithWorkspace(req.Workspace),
		state.ListWithPhysicalState(req.PhysicalState),
	)
	if err != nil {
		return nil, err
	}

	if req.LoadDetails == pb.Deployment_ARTIFACT || req.LoadDetails == pb.Deployment_BUILD {
		for _, dep := range result {
			pa, err := s.state.ArtifactGet(&pb.Ref_Operation{
				Target: &pb.Ref_Operation_Id{
					Id: dep.ArtifactId,
				},
			})

			if err != nil {
				return nil, err
			}

			dep.Preload.Artifact = pa

			if req.LoadDetails == pb.Deployment_BUILD {
				build, err := s.state.BuildGet(&pb.Ref_Operation{
					Target: &pb.Ref_Operation_Id{
						Id: pa.BuildId,
					},
				})

				if err != nil {
					return nil, err
				}

				dep.Preload.Build = build
			}
		}
	}

	return &pb.ListDeploymentsResponse{Deployments: result}, nil
}

// GetDeployment returns a Deployment based on ID
func (s *service) GetDeployment(
	ctx context.Context,
	req *pb.GetDeploymentRequest,
) (*pb.Deployment, error) {
	return s.state.DeploymentGet(req.Ref)
}
