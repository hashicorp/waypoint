package singleprocess

import (
	"context"
	"strings"
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

	// Create our OIDC provider.
	// NOTE(mitchellh): we should cache this because this has to make an
	// external HTTP request and do a bunch of other stateful things.
	oidcCfg, err := oidc.NewConfig(
		amMethod.Oidc.DiscoveryUrl,
		amMethod.Oidc.ClientId,
		oidc.ClientSecret(amMethod.Oidc.ClientSecret),
		[]oidc.Alg{ // TODO(mitchellh): make this configurable
			oidc.EdDSA,
		},
		amMethod.Oidc.AllowedRedirectUris,
		oidc.WithAudiences(amMethod.Oidc.Auds...),
		oidc.WithProviderCA(strings.Join(amMethod.Oidc.DiscoveryCaPem, "\n")),
	)
	if err != nil {
		return nil, err
	}

	provider, err := oidc.NewProvider(oidcCfg)
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"error initializing OIDC provider: %s", err)
	}
	defer provider.Done()

	// Create a minimal request to get the auth URL
	oidcReqOpts := []oidc.Option{
		oidc.WithNonce(req.ClientNonce),
	}
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
