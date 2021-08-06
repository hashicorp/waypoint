package singleprocess

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/ptypes"
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
	if s.urlClient() != nil &&
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
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		// Create the hostname. We ignore errors.
		_, err := s.createHostnameIfNotExist(ctx, target)
		if err != nil {
			log := hclog.FromContext(ctx)
			log.Info("error creating default hostname", "err", err)
		}

		// Populate the URL preload data
		if err := s.deploymentPreloadUrl(ctx, result); err != nil {
			return nil, err
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

	for _, dep := range result {
		setDeploymentUrlIfNeeded(dep)
		if err := s.deploymentPreloadUrl(ctx, dep); err != nil {
			return nil, err
		}
		if err := s.deploymentPreloadDetails(ctx, req.LoadDetails, dep); err != nil {
			return nil, err
		}
	}

	return &pb.ListDeploymentsResponse{Deployments: result}, nil
}

func setDeploymentUrlIfNeeded(d *pb.Deployment) {
	if d.Url != "" {
		if d.Preload == nil {
			d.Preload = &pb.Deployment_Preload{}
		}

		d.Preload.DeployUrl = d.Url
	}
}

// GetDeployment returns a Deployment based on ID
func (s *service) GetDeployment(
	ctx context.Context,
	req *pb.GetDeploymentRequest,
) (*pb.Deployment, error) {
	d, err := s.state.DeploymentGet(req.Ref)
	if err != nil {
		return nil, err
	}

	setDeploymentUrlIfNeeded(d)

	// Populate the URL preload data
	if err := s.deploymentPreloadUrl(ctx, d); err != nil {
		return nil, err
	}
	if err := s.deploymentPreloadDetails(ctx, req.LoadDetails, d); err != nil {
		return nil, err
	}

	return d, nil
}

func (s *service) deploymentPreloadUrl(
	ctx context.Context,
	d *pb.Deployment,
) error {
	// If we had no entrypoint config it is not possible for the preload URL to work.
	if !d.HasEntrypointConfig {
		return nil
	}

	resp, err := s.ListHostnames(ctx, &pb.ListHostnamesRequest{
		Target: &pb.Hostname_Target{
			Target: &pb.Hostname_Target_Application{
				Application: &pb.Hostname_TargetApp{
					Application: d.Application,
					Workspace:   d.Workspace,
				},
			},
		},
	})
	if err == nil && len(resp.Hostnames) > 0 {
		hostname := resp.Hostnames[0]

		d.Preload.DeployUrl = fmt.Sprintf(
			"%s--%s%s",
			hostname.Hostname,
			(&ptypes.Deployment{Deployment: d}).URLFragment(),
			strings.TrimPrefix(hostname.Fqdn, hostname.Hostname),
		)
	}

	return nil
}

func (s *service) deploymentPreloadDetails(
	ctx context.Context,
	req pb.Deployment_LoadDetails,
	d *pb.Deployment,
) error {
	if req <= pb.Deployment_NONE {
		return nil
	}

	pa, err := s.state.ArtifactGet(&pb.Ref_Operation{
		Target: &pb.Ref_Operation_Id{
			Id: d.ArtifactId,
		},
	})
	if err != nil {
		return err
	}

	d.Preload.Artifact = pa

	if req == pb.Deployment_BUILD {
		build, err := s.state.BuildGet(&pb.Ref_Operation{
			Target: &pb.Ref_Operation_Id{
				Id: pa.BuildId,
			},
		})
		if err != nil {
			return err
		}

		d.Preload.Build = build
	}

	return nil
}
