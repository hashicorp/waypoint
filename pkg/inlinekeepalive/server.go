package inlinekeepalive

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/reflect/protoreflect"
)

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

func KeepaliveServerStreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {

		// TODO(izaak): check client version

		ctx := ss.Context()
		log := hclog.FromContext(ctx).With("method", info.FullMethod)

		// Only send keepalives if this is a server stream - not allowed otherwise
		if info.IsServerStream {
			go ServeKeepalives(ctx, log, ss)
		}

		return handler(srv, &KeepaliveServerStream{
			ss:  ss,
			log: hclog.FromContext(ss.Context()).With("interceptor", "inlinekeepalive"),
		})
	}
}
