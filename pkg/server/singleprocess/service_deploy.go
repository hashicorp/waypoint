package singleprocess

import (
	"context"
	"fmt"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

func (s *Service) UpsertDeployment(
	ctx context.Context,
	req *pb.UpsertDeploymentRequest,
) (*pb.UpsertDeploymentResponse, error) {
	if err := serverptypes.ValidateUpsertDeploymentRequest(req); err != nil {
		return nil, err
	}

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

	if err := s.state(ctx).DeploymentPut(!insert, result); err != nil {
		return nil, hcerr.Externalize(hclog.FromContext(ctx), err, "failed to insert deployment for app", "app", req.Deployment.Application, "id", req.Deployment.Id)
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

func (s *Service) ListDeployments(
	ctx context.Context,
	req *pb.ListDeploymentsRequest,
) (*pb.ListDeploymentsResponse, error) {
	result, err := s.state(ctx).DeploymentList(req.Application,
		serverstate.ListWithStatusFilter(req.Status...),
		serverstate.ListWithOrder(req.Order),
		serverstate.ListWithWorkspace(req.Workspace),
		serverstate.ListWithPhysicalState(req.PhysicalState),
	)
	if err != nil {
		return nil, hcerr.Externalize(hclog.FromContext(ctx), err, "failed to list deployments for app", "app", req.Application.Application, "project", req.Application.Project)
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
func (s *Service) GetDeployment(
	ctx context.Context,
	req *pb.GetDeploymentRequest,
) (*pb.Deployment, error) {
	if err := serverptypes.ValidateGetDeploymentRequest(req); err != nil {
		return nil, err
	}

	d, err := s.state(ctx).DeploymentGet(req.Ref)
	if err != nil {
		return nil, hcerr.Externalize(hclog.FromContext(ctx), err, "failed to get deployment", "target", req.Ref.Target)
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

func (s *Service) deploymentPreloadUrl(
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
			(&serverptypes.Deployment{Deployment: d}).URLFragment(),
			strings.TrimPrefix(hostname.Fqdn, hostname.Hostname),
		)
	}

	return nil
}

func (s *Service) deploymentPreloadDetails(
	ctx context.Context,
	req pb.Deployment_LoadDetails,
	d *pb.Deployment,
) error {
	if req <= pb.Deployment_NONE {
		return nil
	}

	pa, err := s.state(ctx).ArtifactGet(&pb.Ref_Operation{
		Target: &pb.Ref_Operation_Id{
			Id: d.ArtifactId,
		},
	})
	if err != nil {
		return err
	}

	d.Preload.Artifact = pa

	if req == pb.Deployment_BUILD {
		build, err := s.state(ctx).BuildGet(&pb.Ref_Operation{
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
