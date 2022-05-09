package serverclient

import (
	"context"
	"crypto/subtle"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/mr-tron/base58"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

// This is a weird type that only exists to satisify the interface required by
// grpc.WithPerRPCCredentials. That api is designed to incorporate things like OAuth
// but in our case, we really just want to send this static token through, but we still
// need to the dance.
type StaticToken string

func (t StaticToken) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": string(t),
	}, nil
}

func (t StaticToken) RequireTransportSecurity() bool {
	return false
}

// ContextToken implements grpc.WithPerRPCCredentials and extracts the token
// from the context or otherwise falls back to a default string value (which
// might be empty).
type ContextToken string

func (t ContextToken) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	finalT := TokenFromContext(ctx)
	if finalT == "" {
		finalT = string(t)
	}

	// If no token, do nothing.
	if finalT == "" {
		return nil, nil
	}

	return map[string]string{
		"authorization": finalT,
	}, nil
}

func (t ContextToken) RequireTransportSecurity() bool {
	return false
}

type tokenKey struct{}

// TokenWithContext stores a token on the given context. If this context
// is used by a connection created by serverclient, this token will take
// priority for requests.
func TokenWithContext(ctx context.Context, t string) context.Context {
	return context.WithValue(ctx, tokenKey{}, t)
}

// TokenFromContext extracts a token (if any) from the given context.
func TokenFromContext(ctx context.Context) string {
	raw := ctx.Value(tokenKey{})
	if raw == nil {
		return ""
	}

	return raw.(string)
}

const (
	// tokenMagic is used as a byte sequence prepended to the encoded TokenTransport to identify
	// the token as valid before attempting to decode it. This is mostly a nicity to improve
	// understanding of the token data and error messages.
	tokenMagic = "wp24"
)

// TokenDecode provides the ability to decode a waypoint token into
// the embedded protobuf. WARNING: This function is unable to
// verify the token, only the server that generate the server can do that.
func TokenDecode(token string) (*pb.Token, error) {
	data, err := base58.Decode(token)
	if err != nil {
		return nil, err
	}

	if subtle.ConstantTimeCompare(data[:len(tokenMagic)], []byte(tokenMagic)) != 1 {
		return nil, errors.Wrapf(ErrInvalidToken, "bad magic")
	}

	var tt pb.TokenTransport
	err = proto.Unmarshal(data[len(tokenMagic):], &tt)
	if err != nil {
		return nil, err
	}

	// Decode the actual token structure
	var body pb.Token
	err = proto.Unmarshal(tt.Body, &body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}
