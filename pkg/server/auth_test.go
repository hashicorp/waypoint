// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type trivialAuth struct {
	method        string
	token         string
	effects       []string
	contextReturn context.Context
}

// Called before each RPC to authenticate it.
func (t *trivialAuth) Authenticate(ctx context.Context, token string, endpoint string, effects []string) (context.Context, error) {
	t.method = endpoint
	t.token = token
	t.effects = effects
	return t.contextReturn, nil
}

func TestAuthUnaryInterceptor(t *testing.T) {
	require := require.New(t)

	var chk trivialAuth

	f := AuthUnaryInterceptor(&chk)

	ctx := context.Background()

	tokenVal := "this-is-a-token"

	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"authorization": []string{tokenVal},
	})

	// Empty context
	called := false
	resp, err := f(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/foo/bar"},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			called = true
			return "hello", nil
		},
	)

	require.True(called)
	require.Equal("hello", resp)
	require.NoError(err)

	require.Equal(tokenVal, chk.token)
	require.Equal("bar", chk.method)
	require.Equal(DefaultEffects, chk.effects)
}

func TestAuthUnaryInterceptor_replaceContext(t *testing.T) {
	require := require.New(t)

	ctx2 := context.Background()

	var chk trivialAuth
	chk.contextReturn = ctx2

	f := AuthUnaryInterceptor(&chk)

	ctx := context.Background()

	tokenVal := "this-is-a-token"

	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"authorization": []string{tokenVal},
	})

	// Empty context
	var got context.Context
	_, err := f(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/foo/bar"},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			got = ctx
			return "hello", nil
		},
	)

	require.NoError(err)
	require.Equal(got, ctx2)
}
