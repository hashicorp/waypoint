package serverclient

import (
	"context"
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
