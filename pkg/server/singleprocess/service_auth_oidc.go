package singleprocess

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/cap/oidc"
	"github.com/hashicorp/go-bexpr"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	empty "google.golang.org/protobuf/types/known/emptypb"

	wpoidc "github.com/hashicorp/waypoint/pkg/auth/oidc"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

const (
	// oidcAuthExpiry is the duration that an OIDC-based login is valid for.
	// We default this to 30 days for now but that is arbitrary. We can change
	// this default anytime or choose to make it configurable one day on the
	// server.
	oidcAuthExpiry = 30 * 24 * time.Hour

	// oidcReqExpiry is the time that an OIDC auth request is valid for.
	// 5 minutes should be plenty of time to complete auth.
	oidcReqExpiry = 5 * 60 * time.Minute
)

func (s *Service) ListOIDCAuthMethods(
	ctx context.Context,
	req *empty.Empty,
) (*pb.ListOIDCAuthMethodsResponse, error) {
	// We implement this by just requesting all the auth methods. We could
	// index OIDC methods specifically and do this more efficiently but
	// realistically we don't expect there to ever be that many auth methods.
	// Even if there were thousands (why????) this would be okay.
	values, err := s.state(ctx).AuthMethodList()
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to list OIDC auth methods",
		)
	}

	// Go through and extract the auth methods
	var result []*pb.OIDCAuthMethod
	for _, method := range values {
		_, ok := method.Method.(*pb.AuthMethod_Oidc)
		if !ok {
			continue
		}

		result = append(result, &pb.OIDCAuthMethod{
			Name:        method.Name,
			DisplayName: method.DisplayName,
			Kind:        pb.OIDCAuthMethod_UNKNOWN,
		})
	}

	return &pb.ListOIDCAuthMethodsResponse{AuthMethods: result}, nil
}

func (s *Service) GetOIDCAuthURL(
	ctx context.Context,
	req *pb.GetOIDCAuthURLRequest,
) (*pb.GetOIDCAuthURLResponse, error) {
	if err := serverptypes.ValidateGetOIDCAuthURLRequest(req); err != nil {
		return nil, err
	}

	// Get the auth method
	am, err := s.state(ctx).AuthMethodGet(req.AuthMethod)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to get OIDC auth method",
			"auth_method",
			req.AuthMethod.GetName(),
		)
	}

	// The auth method must be OIDC
	amMethod, ok := am.Method.(*pb.AuthMethod_Oidc)
	if !ok {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			status.Errorf(codes.FailedPrecondition,
				"auth method is not OIDC"),
			"auth method is not an OIDC auth method",
			"auth_method",
			am.Name,
		)
	}

	// We need our server config.
	sc, err := s.state(ctx).ServerConfigGet()
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to get server config when looking up OIDC auth method",
		)
	}

	// Get our OIDC provider
	provider, err := s.oidcCache.Get(ctx, am, sc)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to get OIDC auth method",
			"auth_method",
			am.Name,
		)
	}

	// Create a minimal request to get the auth URL
	oidcReqOpts := []oidc.Option{
		oidc.WithNonce(req.Nonce),
	}
	if v := amMethod.Oidc.Scopes; len(v) > 0 {
		oidcReqOpts = append(oidcReqOpts, oidc.WithScopes(v...))
	}
	oidcReq, err := oidc.NewRequest(
		oidcReqExpiry,
		req.RedirectUri,
		oidcReqOpts...,
	)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"OIDC request failed for auth method",
			"auth_method",
			am.Name,
		)
	}

	// Get the auth URL
	url, err := provider.AuthURL(ctx, oidcReq)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to get auth URL for OIDC auth method",
			"auth_method",
			am.Name,
		)
	}

	return &pb.GetOIDCAuthURLResponse{
		Url: url,
	}, nil
}

func (s *Service) CompleteOIDCAuth(
	ctx context.Context,
	req *pb.CompleteOIDCAuthRequest,
) (*pb.CompleteOIDCAuthResponse, error) {
	log := hclog.FromContext(ctx)

	if err := serverptypes.ValidateCompleteOIDCAuthRequest(req); err != nil {
		return nil, err
	}

	// Get the auth method
	am, err := s.state(ctx).AuthMethodGet(req.AuthMethod)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to get OIDC auth method",
			"auth_method",
			req.AuthMethod.GetName(),
		)
	}

	// The auth method must be OIDC
	amMethod, ok := am.Method.(*pb.AuthMethod_Oidc)
	if !ok {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			status.Errorf(codes.FailedPrecondition,
				"auth method is not OIDC"),
			"auth method is not an OIDC auth method",
			"auth_method",
			am.Name,
		)
	}

	// We need our server config.
	sc, err := s.state(ctx).ServerConfigGet()
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to get server config when looking up OIDC auth method",
		)
	}

	// Get our OIDC provider
	provider, err := s.oidcCache.Get(ctx, am, sc)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to get OIDC auth method",
			"auth_method",
			am.Name,
		)
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
		oidcReqExpiry,
		req.RedirectUri,
		oidcReqOpts...,
	)
	if err != nil {
		return nil, err
	}

	// Exchange our code for our token
	oidcToken, err := provider.Exchange(ctx, oidcReq, req.State, req.Code)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"OIDC request failed to exchange for token for auth method",
			"auth_method",
			am.Name,
		)
	}

	// Extract the claims as a raw JSON message.
	var jsonClaims json.RawMessage
	if err := oidcToken.IDToken().Claims(&jsonClaims); err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to retrieve OIDC claims from auth method token",
			"auth_method",
			am.Name,
		)
	}

	// Structurally extract only the claim fields we care about.
	var idClaimVals idClaims
	if err := json.Unmarshal([]byte(jsonClaims), &idClaimVals); err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to extract OIDC claims from auth method token",
			"auth_method",
			am.Name,
		)
	}

	// Valid OIDC providers should never behave this way.
	if idClaimVals.Iss == "" || idClaimVals.Sub == "" {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			errors.New("OIDC provider returned empty issuer or subscriber ID"),
			"OIDC provider returned empty issuer or subscriber ID",
			"auth_method",
			am.Name,
		)
	}

	// From this point forward, log all data
	log = log.With(
		"iss", idClaimVals.Iss,
		"sub", idClaimVals.Sub,
	)

	// Get the user info if we have a user account, and merge those claims in.
	// User claims override all the claims in the ID token.
	var userClaims json.RawMessage
	if userTokenSource := oidcToken.StaticTokenSource(); userTokenSource != nil {
		if err := provider.UserInfo(ctx, userTokenSource, idClaimVals.Sub, &userClaims); err != nil {
			return nil, hcerr.Externalize(
				hclog.FromContext(ctx),
				err,
				"failed to extract userinfo from claims",
				"auth_method",
				am.Name,
			)
		}
	}

	// Verify this user is allowed to auth at all.
	if am.AccessSelector != "" {
		// Get our data
		selectorData, err := wpoidc.SelectorData(amMethod.Oidc, jsonClaims, userClaims)
		if err != nil {
			return nil, hcerr.Externalize(
				hclog.FromContext(ctx),
				err,
				"error processiong OIDC claims",
				"auth_method",
				am.Name,
			)
		}

		eval, err := bexpr.CreateEvaluator(am.AccessSelector)
		if err != nil {
			// This shouldn't happen since we validate on auth method create.
			return nil, hcerr.Externalize(
				hclog.FromContext(ctx),
				err,
				"error creating OIDC evaluator",
				"auth_method",
				am.Name,
			)
		}

		allowed, err := eval.Evaluate(selectorData)
		if err != nil {
			return nil, hcerr.Externalize(
				hclog.FromContext(ctx),
				err,
				"error evaluating OIDC selector claims",
				"auth_method",
				am.Name,
			)
		}

		if !allowed {
			// Warn so an operator can detect
			log.Warn("rejected OIDC login based on access selector",
				"claims_json", string(jsonClaims),
				"selector", am.AccessSelector,
			)

			return nil, hcerr.Externalize(
				hclog.FromContext(ctx),
				status.Errorf(codes.PermissionDenied, "rejected OIDC login based on access selector"),
				"Your account was denied access. Please contact your Waypoint "+
					"server administrator for more information.",
				"auth_method",
				am.Name,
			)
		}
	}

	// Look up a user by sub.
	user, err := s.oidcInitUser(ctx, log, &idClaimVals)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"error finding or creating OIDC user",
			"auth_method",
			am.Name,
			// not sure if we should log out idClaimVals.Email here as well
		)
	}

	// Generate a token for this user
	token, err := s.newToken(ctx, oidcAuthExpiry, s.activeAuthKeyId, nil, &pb.Token{
		Kind: &pb.Token_Login_{
			Login: &pb.Token_Login{UserId: user.Id},
		},
	})
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"error generating new token for OIDC user",
			"auth_method",
			am.Name,
		)
	}

	return &pb.CompleteOIDCAuthResponse{
		Token:          token,
		User:           user,
		IdClaimsJson:   string(jsonClaims),
		UserClaimsJson: string(userClaims),
	}, nil
}

// oidcInitUser finds or creates the user for the given OIDC information.
func (s *Service) oidcInitUser(ctx context.Context, log hclog.Logger, claims *idClaims) (*pb.User, error) {
	// This method attempts to find, link, or create a new user to the
	// given OIDC result in the following order:
	//
	//   (1) find user with exact account link (iss, sub match)
	//   (2) find user with matching email and then link it
	//   (3) create new user and link it
	//

	// The email for the user. We only set this if the email is known and
	// verified. This prevents impersonation.
	var email string
	if claims.Email != "" && claims.EmailVerified {
		email = claims.Email
	}

	// First look up by exact account link.
	user, err := s.state(ctx).UserGetOIDC(claims.Iss, claims.Sub)
	if err != nil {
		if status.Code(err) != codes.NotFound {
			return nil, err
		}

		// Just ensure user is nil cause that's the check we'll keep using.
		user = nil
	}
	if user != nil {
		return user, nil
	}

	// Look up the user by email if we don't have a user by sub.
	if email != "" {
		user, err = s.state(ctx).UserGetEmail(email)
		if err != nil {
			if status.Code(err) != codes.NotFound {
				return nil, err
			}

			// Just ensure user is nil cause that's the check we'll keep using.
			user = nil
		}
	}

	// If the user still doesn't exist, we create a new user.
	if user == nil {
		// Random username to start.
		// NOTE(mitchellh): we can improve this in a ton of ways in
		// the future by using their preferred username claim, first name,
		// etc.
		username := fmt.Sprintf("user_%d", time.Now().Unix())

		user = &pb.User{
			Username: username,
			Email:    email,
		}
	}

	// Setup their link
	user.Links = append(user.Links, &pb.User_Link{
		Method: &pb.User_Link_Oidc{
			Oidc: &pb.User_Link_OIDC{
				Iss: claims.Iss,
				Sub: claims.Sub,
			},
		},
	})

	if err := s.state(ctx).UserPut(user); err != nil {
		return nil, err
	}

	log.Info("new OIDC user linked",
		"user_id", user.Id,
		"username", user.Username,
	)

	return user, nil
}

// idClaims are the claims for the ID token that we care about. There
// are many more claims[1] but we only add what we need.
//
// [1]: https://openid.net/specs/openid-connect-core-1_0.html#StandardClaims
type idClaims struct {
	Iss           string `json:"iss"`
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}
