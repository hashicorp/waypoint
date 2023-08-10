// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package server

import (
	"context"
	"strings"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/hashicorp/opaqueany"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/hashicorp/waypoint/pkg/nullify"
)

// GWNullAnyUnaryInterceptor returns a gRPC unary interceptor that replaces all
// *any.Any fields in structs with null ONLY FOR grpc-gateway requests.
//
// grpc-gateway requests are detected by the presence of any metadata that
// starts with the grpc-gateway prefix (HTTP headers).
func GWNullAnyUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		// Invoke the handler.
		resp, err := handler(ctx, req)

		// If we have no response, do nothing.
		if resp == nil {
			return resp, err
		}

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			// Detect if this is a gRPC gateway request. We check if any
			// incoming metadata has the metadata prefix. This should be set
			// for the official HTTP headers that should be set on all gateway
			// requests.
			gw := false
			for k := range md {
				if strings.HasPrefix(k, gwruntime.MetadataPrefix) {
					gw = true
					break
				}
			}

			if gw {
				if nerr := nullify.Nullify(resp, (*opaqueany.Any)(nil)); nerr != nil {
					return nil, nerr
				}
			}
		}

		return resp, err
	}
}
