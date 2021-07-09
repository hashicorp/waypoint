package singleprocess

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

func (s *service) GetAuthMethod(
	ctx context.Context,
	req *pb.GetAuthMethodRequest,
) (*pb.GetAuthMethodResponse, error) {
	value, err := s.state.AuthMethodGet(req.AuthMethod)
	if err != nil {
		return nil, err
	}

	return &pb.GetAuthMethodResponse{
		AuthMethod: value,
	}, nil
}

func (s *service) UpsertAuthMethod(
	ctx context.Context,
	req *pb.UpsertAuthMethodRequest,
) (*pb.UpsertAuthMethodResponse, error) {
	if err := serverptypes.ValidateUpsertAuthMethodRequest(req); err != nil {
		return nil, err
	}

	// Write it
	if err := s.state.AuthMethodPut(req.AuthMethod); err != nil {
		return nil, err
	}

	return &pb.UpsertAuthMethodResponse{
		AuthMethod: req.AuthMethod,
	}, nil
}

func (s *service) DeleteAuthMethod(
	ctx context.Context,
	req *pb.DeleteAuthMethodRequest,
) (*empty.Empty, error) {
	if err := serverptypes.ValidateDeleteAuthMethodRequest(req); err != nil {
		return nil, err
	}

	if err := s.state.AuthMethodDelete(req.AuthMethod); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *service) ListAuthMethods(
	ctx context.Context,
	req *empty.Empty,
) (*pb.ListAuthMethodsResponse, error) {
	values, err := s.state.AuthMethodList()
	if err != nil {
		return nil, err
	}

	return &pb.ListAuthMethodsResponse{AuthMethods: values}, nil
}

func (s *service) ListOIDCAuthMethods(
	ctx context.Context,
	req *empty.Empty,
) (*pb.ListOIDCAuthMethodsResponse, error) {
	// We implement this by just requesting all the auth methods. We could
	// index OIDC methods specifically and do this more efficiently but
	// realistically we don't expect there to ever be that many auth methods.
	// Even if there were thousands (why????) this would be okay.
	values, err := s.state.AuthMethodList()
	if err != nil {
		return nil, err
	}

	// Go through and extract the auth methods
	var result []*pb.OIDCAuthMethod
	for _, method := range values {
		_, ok := method.Method.(*pb.AuthMethod_Oidc)
		if !ok {
			continue
		}

		// Attempt to sniff the kind from the OIDC discovery URL

		result = append(result, &pb.OIDCAuthMethod{
			Name:        method.Name,
			DisplayName: method.DisplayName,
			Kind:        pb.OIDCAuthMethod_UNKNOWN,
		})
	}

	return &pb.ListOIDCAuthMethodsResponse{AuthMethods: result}, nil
}
