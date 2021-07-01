package singleprocess

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/subtle"
	"io"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/mr-tron/base58"
	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/protocolversion"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

const (
	// The user that all tokens are encoded with. The server has a single server model
	// so all actions are mapped to this user.
	DefaultUser = "waypoint"

	// The identifier for the default key to use to generating tokens.
	DefaultKeyId = "k1"

	// Used as a byte sequence prepended to the encoded TokenTransport to identify
	// the token as valid before attempting to decode it. This is mostly a nicity to improve
	// understanding of the token data and error messages.
	tokenMagic = "wp24"

	// The size in bytes that the HMAC keys should be. Each key will contain this number of bytes
	// of data from rand.Reader
	hmacKeySize = 32

	// A prefix added to the key id when looking up the HMAC key from the database
	dbKeyPrefix = "hmacKey:"
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
		"GetVersionInfo":     {},
	}
)

// Authenticate implements the server.AuthChecker interface.
//
// This checks if the given endpoint should be allowed. This is called during a
// gRPC request. Effects is some information about the endpoint, at present these are
// either ["readonly"] or ["mutable"] to indicate if the endpoint will be only reading
// data or also mutating it.
func (s *service) Authenticate(ctx context.Context, token, endpoint string, effects []string) error {
	// Ignore unauthenticated endpoints
	if _, ok := unauthenticatedEndpoints[endpoint]; ok {
		return nil
	}

	if token == "" {
		return status.Errorf(codes.Unauthenticated, "Authorization token is not supplied")
	}

	_, body, err := s.decodeToken(token)
	if err != nil {
		return err
	}

	// Token protocol version must be compatible with this server
	if body.VersionInfo.Api.Current > protocolversion.Current().Api.Current || body.VersionInfo.Api.Current < protocolversion.Current().Api.Minimum {
		return status.Errorf(codes.InvalidArgument,
			"token protocol version %d is incompatible with this server's protocol version %d",
			body.VersionInfo.Api.Current, protocolversion.Current().Api.Current,
		)
	}

	// Token must be a login token to be used for auth
	login, ok := body.Kind.(*pb.Token_Login_)
	if !ok || login == nil {
		return ErrInvalidToken
	}

	// If this is an entrypoint token then we can only access entrypoint APIs.
	if login.Login.Entrypoint != nil && !strings.HasPrefix(endpoint, "Entrypoint") {
		return status.Errorf(codes.Unauthenticated, "Unauthorized endpoint")
	}

	// TODO When we have a user model, this is where you'll check for the user.
	if login.Login.UserId != DefaultUser {
		return ErrInvalidToken
	}

	return nil
}

// decodeToken parses the string and validates it as a valid token. If the token
// has a validity period attached to it, the period is checked here.
//
// This will accept older (pre-Waypoint 0.5) tokens and automatically
// upgrade them to the 0.5 format if it is able. The "unused" fields are
// left untouched in this case.
func (s *service) decodeToken(token string) (*pb.TokenTransport, *pb.Token, error) {
	data, err := base58.Decode(token)
	if err != nil {
		return nil, nil, err
	}

	if subtle.ConstantTimeCompare(data[:len(tokenMagic)], []byte(tokenMagic)) != 1 {
		return nil, nil, errors.Wrapf(ErrInvalidToken, "bad magic")
	}

	var tt pb.TokenTransport
	err = proto.Unmarshal(data[len(tokenMagic):], &tt)
	if err != nil {
		return nil, nil, err
	}

	key, err := s.state.HMACKeyGet(tt.KeyId)
	if err != nil || key == nil {
		return nil, nil, errors.Wrapf(ErrInvalidToken, "unknown key")
	}

	// Hash the token body using the HMAC key so that we can compare
	// with our signature to ensure this hasn't been tampered with.
	h, err := blake2b.New256(key.Key)
	if err != nil {
		return nil, nil, err
	}

	h.Write(tt.Body)
	sum := h.Sum(nil)

	if subtle.ConstantTimeCompare(sum, tt.Signature) != 1 {
		return nil, nil, errors.Wrapf(ErrInvalidToken, "bad signature")
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
			return nil, nil, errors.Wrapf(ErrInvalidToken, "token has expired. %s < %s", vt, now)
		}
	}

	// Determine if we have an old token and upgrade it. We ignore unused
	// fields if we have a Kind set, since this was introduced in WP 0.5.
	if body.Kind == nil {
		login := &pb.Token_Login{
			UserId:     DefaultUser,
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

// Encode the given token with the given key and metadata.
// keyId controls which key is used to sign the key (key values are generated lazily).
// metadata is attached to the token transport as configuration style information
func (s *service) encodeToken(keyId string, metadata map[string]string, body *pb.Token) (string, error) {
	// Get the key material
	key, err := s.state.HMACKeyCreateIfNotExist(keyId, hmacKeySize)
	if err != nil {
		return "", err
	}

	// Proto encode the body, this is what we sign.
	bodyData, err := proto.Marshal(body)
	if err != nil {
		return "", err
	}

	// Sign it
	h, err := blake2b.New256(key.Key)
	if err != nil {
		return "", err
	}
	h.Write(bodyData)

	// Build our wrapper which is not signed or encrypted.
	var tt pb.TokenTransport
	tt.Body = bodyData
	tt.KeyId = keyId
	tt.Metadata = metadata
	tt.Signature = h.Sum(nil)

	// Marshal the wrapper and base58 encode it.
	ttData, err := proto.Marshal(&tt)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	buf.WriteString(tokenMagic)
	buf.Write(ttData)

	return base58.Encode(buf.Bytes()), nil
}

// Create a new login token. This is just a gRPC wrapper around newToken.
func (s *service) GenerateLoginToken(
	ctx context.Context, req *pb.LoginTokenRequest,
) (*pb.NewTokenResponse, error) {
	// We don't currently have any other users.
	if req.User != nil {
		return nil, status.Errorf(codes.PermissionDenied, "cannot create token for other users")
	}

	// If we have a duration set, set the expiry
	var dur time.Duration
	if d := req.Duration; d != "" {
		var err error
		dur, err = time.ParseDuration(d)
		if err != nil {
			return nil, err
		}
	}

	login := &pb.Token_Login{
		UserId: DefaultUser,
	}

	token, err := s.newToken(dur, DefaultKeyId, nil, &pb.Token{
		Kind: &pb.Token_Login_{Login: login},
	})
	if err != nil {
		return nil, err
	}

	return &pb.NewTokenResponse{Token: token}, nil
}

// newToken is the generic internal function to create and encode a new
// token. The final parameter "body" should be set to the initial value
// of the token body, most importantly the "Kind" field should be set.
func (s *service) newToken(
	duration time.Duration,
	keyId string,
	metadata map[string]string,
	body *pb.Token,
) (string, error) {
	body.VersionInfo = protocolversion.Current()

	body.IssuedTime = ptypes.TimestampNow()

	// If this token expires at some point, set an expiry
	if duration > 0 {
		now := time.Now().UTC().Add(duration)
		body.ValidUntil = &timestamp.Timestamp{
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

	return s.encodeToken(keyId, metadata, body)
}

// Create a new invite token.
func (s *service) GenerateInviteToken(
	ctx context.Context, req *pb.InviteTokenRequest,
) (*pb.NewTokenResponse, error) {
	// Old behavior, if we have the entrypoint set, we convert that to
	// a request in the new (WP 0.5+) style. We do this right away so the rest of the
	// request can assume the new style.
	if ep := req.UnusedEntrypoint; ep != nil {
		req.Signup = nil // not a signup token
		req.Login = &pb.Token_Login{
			UserId:     DefaultUser,
			Entrypoint: ep,
		}
	}

	if req.Login == nil {
		req.Login = &pb.Token_Login{
			// TODO(mitchellh): when we have a login system, set this
			// to the currently logged in user.
			UserId: DefaultUser,
		}
	}

	// We don't currently have any other users.
	if req.Login.UserId != DefaultUser {
		return nil, status.Errorf(codes.PermissionDenied, "cannot create invite for other users")
	}

	dur, err := time.ParseDuration(req.Duration)
	if err != nil {
		return nil, err
	}

	invite := &pb.Token_Invite{
		// TODO(mitchellh): when we have a user system, set this.
		FromUserId: "",
		Login:      req.Login,
		Signup:     req.Signup,
	}

	token, err := s.newToken(dur, DefaultKeyId, nil, &pb.Token{
		Kind: &pb.Token_Invite_{Invite: invite},
	})
	if err != nil {
		return nil, err
	}

	return &pb.NewTokenResponse{Token: token}, nil
}

// Given an invite token, validate it and return a login token. This is a gRPC wrapper around ExchangeInvite.
func (s *service) ConvertInviteToken(ctx context.Context, req *pb.ConvertInviteTokenRequest) (*pb.NewTokenResponse, error) {
	_, body, err := s.decodeToken(req.Token)
	if err != nil {
		return nil, err
	}

	kind, ok := body.Kind.(*pb.Token_Invite_)
	if !ok || kind == nil {
		return nil, errors.Wrapf(ErrInvalidToken, "not an invite token")
	}
	invite := kind.Invite

	// Our login token is just the login token on the invite.
	login := invite.Login

	// If we have a signup invite, then error for now until we have an account system.
	if invite.Signup != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "signup tokens not allowed")
	}

	token, err := s.newToken(0, DefaultKeyId, nil, &pb.Token{
		Kind: &pb.Token_Login_{Login: login},
	})
	if err != nil {
		return nil, err
	}

	return &pb.NewTokenResponse{Token: token}, nil
}

// BootstrapToken RPC call.
func (s *service) BootstrapToken(ctx context.Context, req *empty.Empty) (*pb.NewTokenResponse, error) {
	if !s.state.HMACKeyEmpty() {
		return nil, status.Errorf(codes.PermissionDenied, "server is already bootstrapped")
	}

	token, err := s.newToken(0, DefaultKeyId, nil, &pb.Token{
		Kind: &pb.Token_Login_{
			Login: &pb.Token_Login{
				UserId: DefaultUser,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &pb.NewTokenResponse{Token: token}, nil
}

// Bootstrapped returns true if the server is already bootstrapped. If
// this returns true then BootstrapToken can no longer be called.
func (s *service) Bootstrapped() bool {
	return !s.state.HMACKeyEmpty()
}

// DecodeToken RPC call.
func (s *service) DecodeToken(
	ctx context.Context, req *pb.DecodeTokenRequest,
) (*pb.DecodeTokenResponse, error) {
	tt, body, err := s.decodeToken(req.Token)
	if err != nil {
		return nil, err
	}

	return &pb.DecodeTokenResponse{
		Token:     body,
		Transport: tt,
	}, nil
}
