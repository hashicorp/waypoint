package singleprocess

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/handlers"
)

func (s *service) ListRunners(
	ctx context.Context,
	req *pb.ListRunnersRequest,
) (*pb.ListRunnersResponse, error) {
	return handlers.ListRunners(s, ctx, req)
}

func (s *service) GetRunner(
	ctx context.Context,
	req *pb.GetRunnerRequest,
) (*pb.Runner, error) {
	return handlers.GetRunner(s, ctx, req)
}

func (s *service) RunnerGetDeploymentConfig(
	ctx context.Context,
	req *pb.RunnerGetDeploymentConfigRequest,
) (*pb.RunnerGetDeploymentConfigResponse, error) {
	return handlers.RunnerGetDeploymentConfig(s, ctx, req)
}

func (s *service) AdoptRunner(
	ctx context.Context,
	req *pb.AdoptRunnerRequest,
) (*empty.Empty, error) {
	return handlers.AdoptRunner(s, ctx, req)
}

func (s *service) ForgetRunner(
	ctx context.Context,
	req *pb.ForgetRunnerRequest,
) (*empty.Empty, error) {
	return handlers.ForgetRunner(s, ctx, req)
}

func (s *service) RunnerToken(
	ctx context.Context,
	req *pb.RunnerTokenRequest,
) (*pb.RunnerTokenResponse, error) {
	return handlers.RunnerToken(s, ctx, req)
}

func (s *service) RunnerConfig(
	srv pb.Waypoint_RunnerConfigServer,
) error {
	return handlers.RunnerConfig(s, srv)
}

func (s *service) RunnerJobStream(
	srv pb.Waypoint_RunnerJobStreamServer,
) error {
	return handlers.RunnerJobStream(s, srv)
}
