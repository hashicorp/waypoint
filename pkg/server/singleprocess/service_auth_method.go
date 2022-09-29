package singleprocess

import (
	"context"

	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/go-hclog"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func (s *Service) GetAuthMethod(
	ctx context.Context,
	req *pb.GetAuthMethodRequest,
) (*pb.GetAuthMethodResponse, error) {
	if err := serverptypes.ValidateGetAuthMethodRequest(req); err != nil {
		return nil, err
	}

	value, err := s.state(ctx).AuthMethodGet(req.AuthMethod)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to get auth method",
			"auth_method",
			req.AuthMethod.GetName(),
		)
	}

	return &pb.GetAuthMethodResponse{
		AuthMethod: value,
	}, nil
}

func (s *Service) UpsertAuthMethod(
	ctx context.Context,
	req *pb.UpsertAuthMethodRequest,
) (*pb.UpsertAuthMethodResponse, error) {
	if err := serverptypes.ValidateUpsertAuthMethodRequest(req); err != nil {
		return nil, err
	}

	// Display name defaults to the name
	if req.AuthMethod.DisplayName == "" {
		req.AuthMethod.DisplayName = req.AuthMethod.Name
	}

	// Write it
	if err := s.state(ctx).AuthMethodPut(req.AuthMethod); err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to put auth method",
			"auth_method",
			req.AuthMethod.GetName(),
		)
	}

	return &pb.UpsertAuthMethodResponse{
		AuthMethod: req.AuthMethod,
	}, nil
}

func (s *Service) DeleteAuthMethod(
	ctx context.Context,
	req *pb.DeleteAuthMethodRequest,
) (*empty.Empty, error) {
	if err := serverptypes.ValidateDeleteAuthMethodRequest(req); err != nil {
		return nil, err
	}

	// Validate that the auth method exists
	if _, err := s.state(ctx).AuthMethodGet(req.AuthMethod); err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to get auth method when attempting to delete it",
			"auth_method",
			req.AuthMethod.GetName(),
		)
	}

	// There may be a race between deleting and checking above, but that
	// is okay because the delete is idempotent.
	if err := s.state(ctx).AuthMethodDelete(req.AuthMethod); err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to delete auth method",
			"auth_method",
			req.AuthMethod.GetName(),
		)
	}

	// Delete from the cache. If this auth method isn't OIDC that's okay
	// cause this will do nothing.
	s.oidcCache.Delete(ctx, req.AuthMethod.Name)

	return &empty.Empty{}, nil
}

func (s *Service) ListAuthMethods(
	ctx context.Context,
	req *empty.Empty,
) (*pb.ListAuthMethodsResponse, error) {
	values, err := s.state(ctx).AuthMethodList()
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed listing auth methods",
		)
	}

	return &pb.ListAuthMethodsResponse{AuthMethods: values}, nil
}
