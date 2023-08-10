// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package inlinekeepalive

import (
	"context"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// KeepaliveClientStream implements grpc.ServerStream
type KeepaliveServerStream struct {
	log    hclog.Logger
	ss     grpc.ServerStream
	sendMx *sync.Mutex
}

func (k *KeepaliveServerStream) SetHeader(md metadata.MD) error {
	return k.ss.SetHeader(md)
}

func (k *KeepaliveServerStream) SendHeader(md metadata.MD) error {
	return k.ss.SendHeader(md)
}

func (k *KeepaliveServerStream) SetTrailer(md metadata.MD) {
	k.ss.SetTrailer(md)
}

func (k *KeepaliveServerStream) Context() context.Context {
	return k.ss.Context()
}

func (k *KeepaliveServerStream) SendMsg(m interface{}) error {
	// Concurrent calls to SendMsg are unsafe, and there may be another
	// goroutine sending keepalives. Lock before sending.
	k.sendMx.Lock()
	defer k.sendMx.Unlock()

	return k.ss.SendMsg(m)
}

// RecvMsg intercepts keepalive messages and does not pass them
// along to the handler.
func (k *KeepaliveServerStream) RecvMsg(m interface{}) error {
	for {
		err := k.ss.RecvMsg(m)
		if err != nil {
			// NOTE(izaak): If we received an inline keepalive, but m is a proto message that also has a field
			// number 10000000, but a different type, we will get (as of this writing)
			// a GRPC status error (code 13) with this message:
			// > "grpc: failed to unmarshal the received message proto: cannot parse invalid wire-format data"
			// At this point, I don't see a way to distinguish between a waypoint developer accidentally adding
			// a field 10000000 to a stream message, vs any other cause for a parse error. It may be
			// possible to pull the expected field numbers out of m and notice that it's expecting a
			// field 10000000, but I don't see a way to get that via protoreflect.
			return err
		}

		pm, ok := m.(protoreflect.ProtoMessage)
		if !ok {
			// Weird, not a protobuf message, but not our keepalive, so continue as normal
			return nil
		}

		if !IsInlineKeepalive(k.log, pm) {
			return nil
		}

		// It's a keepalive message! Ignore it and recv again
		continue
	}
}

// isClientCompatible determines if a client is able to receive inline keepalives
// by examining context metadata.
func isClientCompatible(ctx context.Context) bool {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return false
	}

	vals := md.Get(GrpcMetaSendKeepalivesKey)
	for _, val := range vals {
		if val == GrpcMetaSendKeepalivesValue {
			return true
		}
	}
	return false
}

// KeepaliveServerStreamInterceptor returns a stream interceptor
// that sends inline keepalive messages on server streams (if the client
// is compatible), and intercepts inline keepalives from the client.
// This is intended to be invoked once at the beginning of an RPC. If
// the client is compatible and this is a ServerStream, will spawn a
// goroutine that runs for the duration of the stream to send inline keepalives.
// Will send a keepalive every sendInterval.
func KeepaliveServerStreamInterceptor(sendInterval time.Duration) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		ctx := ss.Context()
		log := hclog.FromContext(ctx).With("method", info.FullMethod)

		// Ensures SendMsg is not called concurrently.
		sendMx := &sync.Mutex{}

		// Only send keepalives if this is a server stream - not allowed otherwise
		if info.IsServerStream && isClientCompatible(ctx) {
			go func() {
				ServeKeepalives(ctx, log, ss, sendInterval, sendMx)
				log.Trace("stopped sending inline keepalives")
			}()
		}

		return handler(srv, &KeepaliveServerStream{
			ss:     ss,
			log:    log.With("interceptor", "inlinekeepalive"),
			sendMx: sendMx,
		})
	}
}
