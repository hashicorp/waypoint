// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package protocolversion

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func TestUnaryClientInterceptor(t *testing.T) {
	require := require.New(t)

	f := UnaryClientInterceptor(&pb.VersionInfo{
		Api: &pb.VersionInfo_ProtocolVersion{
			Current: 10,
			Minimum: 1,
		},

		Entrypoint: &pb.VersionInfo_ProtocolVersion{
			Current: 15,
			Minimum: 1,
		},

		Version: "1.2.3",
	})

	// Call
	var actual context.Context
	called := false
	err := f(context.Background(), "", nil, nil, nil,
		func(
			ctx context.Context,
			method string,
			req, reply interface{},
			cc *grpc.ClientConn,
			opts ...grpc.CallOption) error {
			called = true
			actual = ctx
			return nil
		})
	require.True(called)
	require.NoError(err)

	// Verify
	md, ok := metadata.FromOutgoingContext(actual)
	require.True(ok)

	{
		vs := md.Get(HeaderClientApiProtocol)
		require.Len(vs, 1)
		require.Equal(vs[0], "1,10")
	}
	{
		vs := md.Get(HeaderClientEntrypointProtocol)
		require.Len(vs, 1)
		require.Equal(vs[0], "1,15")
	}
	{
		vs := md.Get(HeaderClientVersion)
		require.Len(vs, 1)
		require.Equal(vs[0], "1.2.3")
	}
}

func TestStreamClientInterceptor(t *testing.T) {
	require := require.New(t)

	f := StreamClientInterceptor(&pb.VersionInfo{
		Api: &pb.VersionInfo_ProtocolVersion{
			Current: 10,
			Minimum: 1,
		},

		Entrypoint: &pb.VersionInfo_ProtocolVersion{
			Current: 15,
			Minimum: 1,
		},

		Version: "1.2.3",
	})

	// Call
	var actual context.Context
	called := false
	_, err := f(context.Background(), nil, nil, "",
		func(
			ctx context.Context,
			desc *grpc.StreamDesc,
			cc *grpc.ClientConn,
			method string,
			opts ...grpc.CallOption) (grpc.ClientStream, error) {
			called = true
			actual = ctx
			return nil, nil
		})
	require.True(called)
	require.NoError(err)

	// Verify
	md, ok := metadata.FromOutgoingContext(actual)
	require.True(ok)

	{
		vs := md.Get(HeaderClientApiProtocol)
		require.Len(vs, 1)
		require.Equal(vs[0], "1,10")
	}
	{
		vs := md.Get(HeaderClientEntrypointProtocol)
		require.Len(vs, 1)
		require.Equal(vs[0], "1,15")
	}
	{
		vs := md.Get(HeaderClientVersion)
		require.Len(vs, 1)
		require.Equal(vs[0], "1.2.3")
	}
}
