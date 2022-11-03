package singleprocess

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/subtle"
	"io"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/mr-tron/base58"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	empty "google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

const (
	// DefaultUser is the username of the initial user created during bootstrapping. This
	// also is the user that Waypoint server versions prior to 0.5 used with
	// their token, so we use this to detect that scenario as well.
	DefaultUser = serverstate.DefaultUser

	// DefaultUserId is the ID of the initial user created during bootstrapping.
	DefaultUserId = serverstate.DefaultUserId

	// DefaultKeyId is the identifier for the default key to use to generating tokens.
	DefaultKeyId = "k1"

	// tokenMagic is used as a byte sequence prepended to the encoded TokenTransport to identify
	// the token as valid before attempting to decode it. This is mostly a nicity to improve
	// understanding of the token data and error messages.
	tokenMagic = "wp24"
)

var (
	ErrInvalidToken = errors.New("invalid authentication token")

	// unauthenticatedEndpoints are the gRPC service APIs that do not require
	// any authentication. Authenticate doesn't even attempt to parse the
	// token so it can be totally invalid.
	unauthenticatedEndpoints = map[string]struct{}{
		"BootstrapToken":     {},
		"ConvertInviteToken": {},
		"DecodeToken":        {},

		"GetVersionInfo": {},

		"ListOIDCAuthMethods": {},
		"GetOIDCAuthURL":      {},
		"CompleteOIDCAuth":    {},

		"RunnerToken": {},

		"NoAuthRunTrigger": {},
	}
)

type userKey struct{}

// UserWithContext inserts the user value u into the context. This can
// be extracted with UserFromContext.
func UserWithContext(ctx context.Context, u *pb.User) context.Context {
	return context.WithValue(ctx, userKey{}, u)
}

// UserFromContext returns the authenticated user in the request context.
// This will return nil if the user is not authenticated. Note that a user
// may not be authenticated but the request can still be authenticated
// using a non-user token type. The safeste way to check is decodedTokenFromContext.
func (s *Service) UserFromContext(ctx context.Context) *pb.User {
	value, ok := ctx.Value(userKey{}).(*pb.User)
	if !ok && s.superuser {
		value = &pb.User{Id: DefaultUserId, Username: DefaultUser}
	}

	return value
}

type decodedTokenKey struct{}

// DecodedTokenWithContext inserts the decrypted token t into the context.
func DecodedTokenWithContext(ctx context.Context, t *pb.Token) context.Context {
	return context.WithValue(ctx, decodedTokenKey{}, t)
}

// decodedTokenFromContext returns the validated token used with the request.
// The token is guaranteed to be valid, meaning that it successfully
// was signed and decrypted. This will return nil if no token was present
// for the request.
func (s *Service) decodedTokenFromContext(ctx context.Context) *pb.Token {
	value, ok := ctx.Value(decodedTokenKey{}).(*pb.Token)
	if !ok && s.superuser {
		// We are in implicit superuser mode meaning everything is always
		// allowed. Create a login token for the superuser.
		value = &pb.Token{
			Kind: &pb.Token_Login_{
				Login: &pb.Token_Login{
					UserId: DefaultUserId,
				},
			},
		}
	}

	return value
}

// CookieFromRequest returns the server cookie value provided during the request,
// or blank if none (or a blank cookie) is provided.
func CookieFromRequest(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if c, ok := md["wpcookie"]; ok && len(c) > 0 {
			return c[0]
		}
	}

	return ""
}

// Authenticate implements the server.AuthChecker interface.
//
// This checks if the given endpoint should be allowed. This is called during a
// gRPC request. Effects is some information about the endpoint, at present these are
// either ["readonly"] or ["mutable"] to indicate if the endpoint will be only reading
// data or also mutating it.
func (s *Service) Authenticate(
	ctx context.Context, token, endpoint string, effects []string,
) (context.Context, error) {
	_, anonEndpoint := unauthenticatedEndpoints[endpoint]

	// Check the cookie
	if c := CookieFromRequest(ctx); c != "" {
		serverConfig, err := s.state(ctx).ServerConfigGet(ctx)
		if err != nil {
			return nil, err
		}

		if !strings.EqualFold(serverConfig.Cookie, c) {
			return nil, status.Errorf(codes.PermissionDenied, "server cookie does not match")
		}
	}

	// We require a token if this isn't an unauthenticated endpoint
	if !anonEndpoint && token == "" {
		return nil, status.Errorf(codes.Unauthenticated, "Authorization token is not supplied")
	}

	// Try to decode the token
	_, body, err := s.decodeToken(ctx, token)
	if err != nil {
		body = nil

		// We only return an error if this isn't an anonymous endpoint.
		// Otherwise, we ignore token errors because they don't affect
		// guest behavior.
		if !anonEndpoint {
			return nil, err
		}
	}

	// Store the token in the context
	ctx = DecodedTokenWithContext(ctx, body)

	// If we are at an unauthenticated endpoint, no need to verify further.
	if anonEndpoint {
		return ctx, nil
	}

	// "Authentication" depends on the type of token.
	switch k := body.Kind.(type) {
	case *pb.Token_Trigger_:
		// trigger token auth should explicitly not be allowed for gRPC requests
		return nil, status.Errorf(codes.PermissionDenied, "Trigger URL token not "+
			"authorized to make requests on this endpoint.")

	case *pb.Token_Login_:
		return s.authLogin(ctx, body, endpoint)

	case *pb.Token_Runner_:
		return s.authRunner(ctx, k.Runner, endpoint)

	default:
		return nil, ErrInvalidToken
	}
}

// authRunner authenticates runner token types.
func (s *Service) authRunner(
	ctx context.Context, tokenRunner *pb.Token_Runner, endpoint string,
) (context.Context, error) {

	log := hclog.FromContext(ctx)

	runnerId, err := s.decodeId(tokenRunner.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to decode id in runner token")
	}

	// If no ID is set, then the runner is assumed at all times to be adopted.
	// This use case is used to "pre-adopt" runners and avoid the adoption
	// lifecycle completely, such as with infinitely autoscaled runners.
	if runnerId == "" {
		// Authenticated.
		return ctx, nil
	}

	// If this is a runner registration request, we allow it through
	// because those APIs will verify and adopt the runner if they can.
	if endpoint == "RunnerConfig" || endpoint == "RunnerToken" {
		return ctx, nil
	}

	// Get our runner
	r, err := s.state(ctx).RunnerById(ctx, runnerId, nil)
	if status.Code(err) == codes.NotFound {
		err = nil
		r = nil
	}
	if err != nil {
		return nil, err
	}

	// Runner is not adopted if it does not exist or is not in an adopted state
	notAdopted := r == nil ||
		(r.AdoptionState != pb.Runner_ADOPTED &&
			r.AdoptionState != pb.Runner_PREADOPTED)
	if notAdopted {
		if r == nil {
			log.Debug("unknown runner attempted to connect", "id", runnerId)
		} else {
			log.Debug("rejecting runner due to adoption state", "id", runnerId, "state", r.AdoptionState.String())
		}

		// We sleep here to tarpit any runaway runners that are going to thrashing
		// trying to connect over and over and over even though they're not allowed in.
		time.Sleep(5 * time.Second)
		return nil, status.Errorf(codes.PermissionDenied,
			"runner is not adopted")
	}

	// If we have a label hash and it doesn't match the labels on the
	// runner, then error.
	if r != nil && tokenRunner.LabelHash > 0 {
		hash, err := serverptypes.RunnerLabelHash(r.Labels)
		if err != nil {
			return nil, err
		}

		if tokenRunner.LabelHash != hash {
			return nil, status.Errorf(codes.PermissionDenied,
				"runner labels have changed since this token was issued")
		}
	}

	return ctx, nil
}

// authLogin authenticates login token types.
func (s *Service) authLogin(
	ctx context.Context, body *pb.Token, endpoint string,
) (context.Context, error) {
	log := hclog.FromContext(ctx)
	// Token must be a login token to be used for auth
	login, ok := body.Kind.(*pb.Token_Login_)
	if !ok || login == nil {
		return nil, ErrInvalidToken
	}

	// If this is an entrypoint token then we can only access entrypoint APIs.
	if login.Login.Entrypoint != nil && !strings.HasPrefix(endpoint, "Entrypoint") {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthorized endpoint")
	}

	userId, err := s.decodeId(login.Login.UserId)
	if err != nil {
		msg := "failed to decode id when authenticating login token"
		log.Error(msg, "id", login.Login.UserId, "err", err)
		return nil, status.Errorf(codes.Internal, msg)
	}

	// Look up the user that this token is for.
	user, err := s.state(ctx).UserGet(ctx, &pb.Ref_User{
		Ref: &pb.Ref_User_Id{
			Id: &pb.Ref_UserId{Id: userId},
		},
	})
	if status.Code(err) == codes.NotFound {
		// If we are a legacy token logging into a WP 0.5+ server for the
		// first time, we need to create the initial Waypoint user. This is
		// purely a backwards compatibility case that we should drop at some
		// point.
		if !body.UnusedLogin || userId != DefaultUserId {
			return nil, status.Errorf(codes.Unauthenticated,
				"Pre-Waypoint 0.5 token must be for the default user. This should always "+
					"be the case so the token is likely corrupt.")
		}

		// Bootstrap our user
		if _, err := s.bootstrapUser(ctx); err != nil {
			// If we're already bootstrapped, give a slightly better error.
			if status.Code(err) == codes.PermissionDenied {
				return nil, status.Errorf(codes.Unauthenticated,
					"Pre-Waypoint 0.5 token no longer accepted once the bootstrap user is deleted")
			}

			return nil, err
		}

		// Look up the user again
		user, err = s.state(ctx).UserGet(ctx, &pb.Ref_User{
			Ref: &pb.Ref_User_Id{
				Id: &pb.Ref_UserId{Id: userId},
			},
		})
	}
	if err != nil {
		return nil, err
	}

	return UserWithContext(ctx, user), nil
}

// decodeToken parses the string and validates it as a valid token. If the token
// has a validity period attached to it, the period is checked here.
//
// This will accept older (pre-Waypoint 0.5) tokens and automatically
// upgrade them to the 0.5 format if it is able. The "unused" fields are
// left untouched in this case.
func (s *Service) decodeToken(ctx context.Context, token string) (*pb.TokenTransport, *pb.Token, error) {
	data, err := base58.Decode(token)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to base58 decode token")
	}

	if subtle.ConstantTimeCompare(data[:len(tokenMagic)], []byte(tokenMagic)) != 1 {
		return nil, nil, errors.Errorf("bad magic")
	}

	var tt pb.TokenTransport
	err = proto.Unmarshal(data[len(tokenMagic):], &tt)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to proto unmarshal token")
	}

	isValid, err := s.state(ctx).TokenSignatureVerify(ctx, tt.Body, tt.Signature, tt.KeyId)
	if err != nil {
		return nil, nil, errors.Wrapf(ErrInvalidToken, err.Error())
	}
	if !isValid {
		return nil, nil, errors.Wrapf(ErrInvalidToken, "bad token signature")
	}

	// Decode the actual token structure
	var body pb.Token
	err = proto.Unmarshal(tt.Body, &body)
	if err != nil {
		return nil, nil, err
	}

	if body.ValidUntil != nil {
		vt := time.Unix(body.ValidUntil.Seconds, int64(body.ValidUntil.Nanos))

		now := time.Now()
		if vt.Before(now.UTC()) {
			return nil, nil, errors.Wrapf(ErrInvalidToken,
				"Token has expired. The token was valid until: %s. Please "+
					"reauthenticate to continue accessing Waypoint.", vt)
		}
	}

	// Determine if we have an old token and upgrade it. We ignore unused
	// fields if we have a Kind set, since this was introduced in WP 0.5.
	if body.Kind == nil {
		login := &pb.Token_Login{
			UserId:     DefaultUserId,
			Entrypoint: body.UnusedEntrypoint,
		}

		switch {
		case body.UnusedLogin:
			body.Kind = &pb.Token_Login_{
				Login: login,
			}

		case body.UnusedInvite:
			body.Kind = &pb.Token_Invite_{
				Invite: &pb.Token_Invite{
					Login: login,
				},
			}

		default:
			return nil, nil, errors.Wrapf(ErrInvalidToken, "token kind is not set")
		}
	}

	return &tt, &body, nil
}

// encodeToken Encodes the given token with the given key and metadata.
// keyId controls which key is used to sign the key (key values are generated lazily).
// metadata is attached to the token transport as configuration style information
func (s *Service) encodeToken(ctx context.Context, tt *pb.TokenTransport, body *pb.Token) (string, error) {
	// Proto encode the token, this is what we encrypt (HCP) or sign/hash (OSS).
	tokenBodyData, err := proto.Marshal(body)
	if err != nil {
		return "", errors.Wrapf(err, "failed to proto marshal the token")
	}

	tokenSignature, err := s.state(ctx).TokenSignature(ctx, tokenBodyData, tt.KeyId)
	if err != nil {
		return "", errors.Wrapf(err, "failed getting token signature")
	}

	// Build our wrapper which is not signed or encrypted.
	tt.Body = tokenBodyData
	tt.Signature = tokenSignature

	// Marshal the wrapper and base58 encode it.
	ttData, err := proto.Marshal(tt)
	if err != nil {
		return "", errors.Wrapf(err, "failed proto marshalling token transport")
	}

	var buf bytes.Buffer
	buf.WriteString(tokenMagic)
	buf.Write(ttData)

	// Encode token
	return base58.Encode(buf.Bytes()), nil
}

// Create a new login token. This is just a gRPC wrapper around newToken.
func (s *Service) GenerateLoginToken(
	ctx context.Context, req *pb.LoginTokenRequest,
) (*pb.NewTokenResponse, error) {
	log := hclog.FromContext(ctx)

	// Get our user, that's what we log in as
	currentUser := s.UserFromContext(ctx)

	// If we have a duration set, set the expiry
	var dur time.Duration
	if d := req.Duration; d != "" {
		var err error
		dur, err = time.ParseDuration(d)
		if err != nil {
			return nil, err
		}
	}

	encodedId, err := s.encodeId(ctx, currentUser.Id)
	if err != nil {
		msg := "failed to encode id when generating a login token"
		log.Error(msg, "currentUser.Id", currentUser.Id, "err", err)
		return nil, status.Error(codes.Internal, msg)
	}

	login := &pb.Token_Login{
		UserId: encodedId,
	}

	// If we're authing as another user, we have to get that user
	if req.User != nil {
		user, err := s.state(ctx).UserGet(ctx, req.User)
		if err != nil {
			return nil, err
		}

		login.UserId, err = s.encodeId(ctx, user.Id)
		if err != nil {
			msg := "failed to encode id when generating a login token"
			log.Error(msg, "user.Id", currentUser.Id, "err", err)
			return nil, status.Error(codes.Internal, msg)
		}
	}

	createToken := &pb.Token{}

	if req.Trigger {
		// NOTE(briancain): when Waypoint has proper access control, this should
		// only be allowed to be generated by super users. The auth call might
		// be checked before the client gets to make this request though, but for
		// now the note about it is here because every user has the same permissions.
		createToken.Kind = &pb.Token_Trigger_{
			Trigger: &pb.Token_Trigger{
				FromUserId: encodedId,
			},
		}
	} else {
		// Default token type
		createToken.Kind = &pb.Token_Login_{Login: login}
	}

	token, err := s.newToken(ctx, dur, s.activeAuthKeyId, nil, createToken)
	if err != nil {
		return nil, err
	}

	return &pb.NewTokenResponse{Token: token}, nil
}

// Create a new runner token.
func (s *Service) GenerateRunnerToken(
	ctx context.Context, req *pb.GenerateRunnerTokenRequest,
) (*pb.NewTokenResponse, error) {
	log := hclog.FromContext(ctx)
	// If we have a duration set, set the expiry
	var dur time.Duration
	if d := req.Duration; d != "" {
		var err error
		dur, err = time.ParseDuration(d)
		if err != nil {
			return nil, err
		}
	}

	encodedId, err := s.encodeId(ctx, req.Id)
	if err != nil {
		msg := "failed to encode id when generating a runner token"
		log.Error(msg, "req.Id", req.Id, "err", err)
		return nil, status.Error(codes.InvalidArgument, msg)
	}

	var hash uint64 = 0
	if len(req.Labels) > 0 {
		var err error
		hash, err = serverptypes.RunnerLabelHash(req.Labels)
		if err != nil {
			return nil, err
		}
	}

	createToken := &pb.Token{
		Kind: &pb.Token_Runner_{
			Runner: &pb.Token_Runner{
				Id:        encodedId,
				LabelHash: hash,
			},
		},
	}

	token, err := s.newToken(ctx, dur, s.activeAuthKeyId, nil, createToken)
	if err != nil {
		return nil, err
	}

	return &pb.NewTokenResponse{Token: token}, nil
}

// newToken is the generic internal function to create and encode a new
// token. The final parameter "body" should be set to the initial value
// of the token body, most importantly the "Kind" field should be set.
func (s *Service) newToken(
	ctx context.Context,
	duration time.Duration,
	keyId string,
	metadata map[string]string,
	body *pb.Token,
) (string, error) {
	body.IssuedTime = timestamppb.Now()

	// If this token expires at some point, set an expiry
	if duration > 0 {
		now := time.Now().UTC().Add(duration)
		body.ValidUntil = &timestamppb.Timestamp{
			Seconds: now.Unix(),
			Nanos:   int32(now.Nanosecond()),
		}
	}

	// Set the accessor ID
	body.AccessorId = make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, body.AccessorId)
	if err != nil {
		return "", err
	}

	// Build our wrapper which is not signed or encrypted.
	var tt pb.TokenTransport
	tt.KeyId = keyId
	tt.Metadata = metadata

	if s.processToken != nil {
		body, err = s.processToken(ctx, &tt, body)
		if err != nil {
			return "", err
		}
	}

	return s.encodeToken(ctx, &tt, body)
}

// Create a new invite token.
func (s *Service) GenerateInviteToken(
	ctx context.Context, req *pb.InviteTokenRequest,
) (*pb.NewTokenResponse, error) {
	log := hclog.FromContext(ctx)

	currentUser := s.UserFromContext(ctx)
	if currentUser == nil {
		return nil, status.Errorf(codes.Unauthenticated, "current user is not authenticated")
	}

	// Old behavior, if we have the entrypoint set, we convert that to
	// a request in the new (WP 0.5+) style. We do this right away so the rest of the
	// request can assume the new style.
	if ep := req.UnusedEntrypoint; ep != nil {
		// NOTE(mitchellh): in the future, we will need to do some policy
		// checks. For now, we allow all users to do this.

		req.Signup = nil // not a signup token
		req.Login = &pb.Token_Login{
			UserId:     DefaultUserId,
			Entrypoint: ep,
		}
	}

	if req.Login == nil {
		req.Login = &pb.Token_Login{
			UserId: currentUser.Id,
		}
	}

	// If we're creating a login token for another user and this is not
	// a signup token, then we need to verify that user exists.
	// req.Login.UserId is authored by the caller and won't be an encoded id
	// so we don't decode it, but we will be sure the resulting token is encoded.
	if req.Login.UserId != currentUser.Id && req.Signup == nil {
		_, err := s.state(ctx).UserGet(ctx, &pb.Ref_User{
			Ref: &pb.Ref_User_Id{
				Id: &pb.Ref_UserId{Id: req.Login.UserId},
			},
		})
		if err != nil {
			return nil, err
		}
	}

	dur, err := time.ParseDuration(req.Duration)
	if err != nil {
		return nil, err
	}

	// TODO(mitchellh): when we have a policy system, we need to ensure only
	// management tokens can signup other users.

	loginUserId, err := s.encodeId(ctx, req.Login.UserId)
	if err != nil {
		msg := "failed to encode id when generating an invite token"
		log.Error(msg, "req.Login.UserId", req.Login.UserId, "err", err)
		return nil, status.Error(codes.InvalidArgument, msg)
	}
	req.Login.UserId = loginUserId

	fromUserId, err := s.encodeId(ctx, currentUser.Id)
	if err != nil {
		msg := "failed to encode the 'from' user's id when generating an invite token"
		log.Error(msg, "currentUser.Id", currentUser.Id, "err", err)
		return nil, status.Error(codes.InvalidArgument, msg)
	}

	invite := &pb.Token_Invite{
		FromUserId: fromUserId,
		Login:      req.Login,
		Signup:     req.Signup,
	}

	token, err := s.newToken(ctx, dur, s.activeAuthKeyId, nil, &pb.Token{
		Kind: &pb.Token_Invite_{Invite: invite},
	})
	if err != nil {
		return nil, err
	}

	return &pb.NewTokenResponse{Token: token}, nil
}

// Given an invite token, validate it and return a login token. This is a gRPC wrapper around ExchangeInvite.
func (s *Service) ConvertInviteToken(ctx context.Context, req *pb.ConvertInviteTokenRequest) (*pb.NewTokenResponse, error) {
	log := hclog.FromContext(ctx)
	_, body, err := s.decodeToken(ctx, req.Token)
	if err != nil {
		return nil, err
	}

	kind, ok := body.Kind.(*pb.Token_Invite_)
	if !ok || kind == nil {
		return nil, errors.Wrapf(ErrInvalidToken, "not an invite token")
	}
	invite := kind.Invite

	// If we have a signup invite, then create a new user.
	if signup := invite.Signup; signup != nil {
		user := &pb.User{Username: signup.InitialUsername}
		if err := s.state(ctx).UserPut(ctx, user); err != nil {
			return nil, err
		}

		// Setup the login information for the new user
		invite.Login.UserId, err = s.encodeId(ctx, user.Id)
		if err != nil {
			msg := "failed to encode the current user's id when converting an invite token"
			log.Error(msg, "user.Id", user.Id, "err", err)
			return nil, status.Error(codes.InvalidArgument, msg)
		}
	}

	// Our login token is just the login token on the invite.
	login := invite.Login

	token, err := s.newToken(ctx, 0, s.activeAuthKeyId, nil, &pb.Token{
		Kind: &pb.Token_Login_{Login: login},
	})
	if err != nil {
		return nil, err
	}

	return &pb.NewTokenResponse{Token: token}, nil
}

// BootstrapToken RPC call.
func (s *Service) BootstrapToken(ctx context.Context, req *empty.Empty) (*pb.NewTokenResponse, error) {
	if !s.state(ctx).HMACKeyEmpty(ctx) {
		return nil, status.Errorf(codes.PermissionDenied, "server is already bootstrapped")
	}

	// Create a default user
	user, err := s.bootstrapUser(ctx)
	if err != nil {
		return nil, err
	}

	// Create a new token pointed to our existing user
	token, err := s.newToken(ctx, 0, s.activeAuthKeyId, nil, &pb.Token{
		Kind: &pb.Token_Login_{
			Login: &pb.Token_Login{
				UserId: user.Id,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &pb.NewTokenResponse{Token: token}, nil
}

// bootstrapUser creates the initial default user. This will always attempt
// to create the user so gating logic to prevent that is up to the caller.
func (s *Service) bootstrapUser(ctx context.Context) (*pb.User, error) {
	empty, err := s.state(ctx).UserEmpty(ctx)
	if err != nil {
		return nil, err
	}
	if !empty {
		return nil, status.Errorf(codes.PermissionDenied, "server is already bootstrapped")
	}

	// Create a default user
	user := &pb.User{
		Id:       DefaultUserId,
		Username: DefaultUser,
	}

	return user, s.state(ctx).UserPut(ctx, user)
}

// Bootstrapped returns true if the server is already bootstrapped. If
// this returns true then BootstrapToken can no longer be called.
func (s *Service) Bootstrapped(ctx context.Context) bool {
	return !s.state(ctx).HMACKeyEmpty(ctx)
}

// DecodeToken RPC call.
func (s *Service) DecodeToken(
	ctx context.Context, req *pb.DecodeTokenRequest,
) (*pb.DecodeTokenResponse, error) {
	tt, body, err := s.decodeToken(ctx, req.Token)
	if err != nil {
		return nil, err
	}

	return &pb.DecodeTokenResponse{
		Token:     body,
		Transport: tt,
	}, nil
}
