package inlinekeepalive

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// KeepaliveClientStream implements grpc.ClientStream
type KeepaliveClientStream struct {
	log     hclog.Logger
	handler grpc.ClientStream
}

func (k *KeepaliveClientStream) Header() (metadata.MD, error) {
	return k.handler.Header()
}

func (k *KeepaliveClientStream) Trailer() metadata.MD {
	return k.handler.Trailer()
}

func (k *KeepaliveClientStream) CloseSend() error {
	return k.handler.CloseSend()
}

func (k *KeepaliveClientStream) Context() context.Context {
	return k.handler.Context()
}

func (k *KeepaliveClientStream) SendMsg(m interface{}) error {
	return k.handler.SendMsg(m)
}

func (k *KeepaliveClientStream) RecvMsg(m interface{}) error {
	for {
		err := k.handler.RecvMsg(m)
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

func KeepaliveClientStreamInterceptor() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		log := hclog.FromContext(ctx).With("method", method)

		handler, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get the client handler when setting up the inlinekeepalive interceptor on method %q", method)
		}

		// Only send keepalives if this is a client stream - not allowed otherwise.
		if desc.ClientStreams {
			// Send keepalives for as long as the handler has the connection open
			go ServeKeepalives(ctx, log, handler)
		}

		return &KeepaliveClientStream{handler: handler}, nil
	}
}
