package singleprocess

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/cap/oidc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

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

		// TODO(mitchellh): sniff the kind from the discovery URL.

		result = append(result, &pb.OIDCAuthMethod{
			Name:        method.Name,
			DisplayName: method.DisplayName,
			Kind:        pb.OIDCAuthMethod_UNKNOWN,
		})
	}

	return &pb.ListOIDCAuthMethodsResponse{AuthMethods: result}, nil
}

func (s *service) GetOIDCAuthURL(
	ctx context.Context,
	req *pb.GetOIDCAuthURLRequest,
) (*pb.GetOIDCAuthURLResponse, error) {
	if err := serverptypes.ValidateGetOIDCAuthURLRequest(req); err != nil {
		return nil, err
	}

	// Get the auth method
	am, err := s.state.AuthMethodGet(req.AuthMethod)
	if err != nil {
		return nil, err
	}

	// The auth method must be OIDC
	amMethod, ok := am.Method.(*pb.AuthMethod_Oidc)
	if !ok {
		return nil, status.Errorf(codes.FailedPrecondition,
			"auth method is not OIDC")
	}

	// Get our OIDC provider
	provider, err := s.oidcCache.Get(ctx, am)
	if err != nil {
		return nil, err
	}

	// Create a minimal request to get the auth URL
	oidcReqOpts := []oidc.Option{}
	if v := amMethod.Oidc.Scopes; len(v) > 0 {
		oidcReqOpts = append(oidcReqOpts, oidc.WithScopes(v...))
	}
	oidcReq, err := oidc.NewRequest(
		5*60*time.Second,
		req.RedirectUri,
		oidcReqOpts...,
	)
	if err != nil {
		return nil, err
	}

	// Get the auth URL
	url, err := provider.AuthURL(ctx, oidcReq)
	if err != nil {
		return nil, err
	}

	return &pb.GetOIDCAuthURLResponse{
		Url: url,
	}, nil
}

func (s *service) CompleteOIDCAuth(
	ctx context.Context,
	req *pb.CompleteOIDCAuthRequest,
) (*pb.CompleteOIDCAuthResponse, error) {
	if err := serverptypes.ValidateCompleteOIDCAuthRequest(req); err != nil {
		return nil, err
	}

	// Get the auth method
	am, err := s.state.AuthMethodGet(req.AuthMethod)
	if err != nil {
		return nil, err
	}

	// The auth method must be OIDC
	amMethod, ok := am.Method.(*pb.AuthMethod_Oidc)
	if !ok {
		return nil, status.Errorf(codes.FailedPrecondition,
			"auth method is not OIDC")
	}

	// Get our OIDC provider
	provider, err := s.oidcCache.Get(ctx, am)
	if err != nil {
		return nil, err
	}

	// Create a minimal request to get the auth URL
	oidcReqOpts := []oidc.Option{
		oidc.WithNonce(req.Nonce),
		oidc.WithState(req.State),
	}
	if v := amMethod.Oidc.Scopes; len(v) > 0 {
		oidcReqOpts = append(oidcReqOpts, oidc.WithScopes(v...))
	}
	if v := amMethod.Oidc.Auds; len(v) > 0 {
		oidcReqOpts = append(oidcReqOpts, oidc.WithAudiences(v...))
	}
	oidcReq, err := oidc.NewRequest(
		5*60*time.Second,
		req.RedirectUri,
		oidcReqOpts...,
	)
	if err != nil {
		return nil, err
	}

	// Exchange our code for our token
	oidcToken, err := provider.Exchange(ctx, oidcReq, req.State, req.Code)
	if err != nil {
		return nil, err
	}

	// Extract the claims for this token
	var idClaims map[string]interface{}
	if err := oidcToken.IDToken().Claims(&idClaims); err != nil {
		return nil, err
	}

	// TODO: look up by sub, look up by email, create new user

	return &pb.CompleteOIDCAuthResponse{}, nil
}
