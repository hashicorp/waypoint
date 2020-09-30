package protocol

import (
	"context"
)

type contextKeyType string

// WithContext stores the protocol version in the context.
func WithContext(ctx context.Context, typ Type, vsn uint32) context.Context {
	return context.WithValue(ctx, contextKeyType(typ.String()), vsn)
}

// FromContext retrieves the protocol version from the context, or returns
// zero if no version was present.
func FromContext(ctx context.Context, typ Type) uint32 {
	v, ok := ctx.Value(contextKeyType(typ.String())).(uint32)
	if !ok {
		v = 0
	}

	return v
}
