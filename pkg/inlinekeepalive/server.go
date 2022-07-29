package inlinekeepalive

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// KeepaliveClientStream implements grpc.ServerStream
type KeepaliveServerStream struct {
	log hclog.Logger
	ss  grpc.ServerStream
}

func (k *KeepaliveServerStream) SetHeader(md metadata.MD) error {
	return k.ss.SetHeader(md)
}

func (k *KeepaliveServerStream) SendHeader(md metadata.MD) error {
	return k.ss.SendHeader(md)
}

func (k *KeepaliveServerStream) SetTrailer(md metadata.MD) {
	k.ss.SetTrailer(md)
	return
}

func (k *KeepaliveServerStream) Context() context.Context {
	return k.ss.Context()
}

func (k *KeepaliveServerStream) SendMsg(m interface{}) error {
	return k.ss.SendMsg(m)
}

// RecvMsg intercepts keepalive messages and does not pass them
// along to the handler.
func (k *KeepaliveServerStream) RecvMsg(m interface{}) error {
	for {
		err := k.ss.RecvMsg(m)
		if err != nil {
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

	vals := md.Get(HeaderSendKeepalives)
	if len(vals) == 0 {
		return false
	}
	if vals[0] != "true" {
		return false
	}

	return true
}

// KeepaliveServerStreamInterceptor returns a stream interceptor
// that sends inline keepalive messages on server streams (if the client
// is compatible), and intercepts inline keepalives from the client.
func KeepaliveServerStreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if err := ss.SetHeader(metadata.MD{"send_inline_keepalives": []string{"true"}}); err != nil {
			return errors.Wrap(err, "failed setting inline keepalive header")
		}

		//ctx := metadata.AppendToOutgoingContext(ss.Context(), HeaderSendKeepalives, "true")
		ctx := ss.Context()
		log := hclog.FromContext(ctx).With("method", info.FullMethod)

		// Only send keepalives if this is a server stream - not allowed otherwise
		if info.IsServerStream && isClientCompatible(ctx) {
			go ServeKeepalives(ctx, log, ss)
		}

		return handler(srv, &KeepaliveServerStream{
			ss:  ss,
			log: hclog.FromContext(ctx).With("interceptor", "inlinekeepalive"),
		})
	}
}
