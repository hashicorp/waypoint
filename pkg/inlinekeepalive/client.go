package inlinekeepalive

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
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

// RecvMsg intercepts keepalive messages and does not pass them
// along to the handler.
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

// isServerCompatible determines if a server is able to receive inline keepalives
// by examining the features it advertises on GetVersionInfo. Triggers an RPC call.
func isServerCompatible(ctx context.Context, cc *grpc.ClientConn) (bool, error) {
	client := pb.NewWaypointClient(cc)

	versionInfo, err := client.GetVersionInfo(ctx, &emptypb.Empty{})
	if err != nil {
		return false, errors.Wrapf(err, "failed getting version info to determine if server is inline-keepalive compatible")
	}

	for _, feature := range versionInfo.Features {
		if feature == FeatureName {
			return true, nil
		}
	}

	return false, nil
}

// KeepaliveClientStreamInterceptor returns a stream interceptor
// that sends inline keepalive messages on client streams (if the server
// is compatible), and intercepts inline keepalives from the server.
// This is intended to be invoked once at the beginning of an RPC, may call
// the server's GetVersionInfo RPC, and if the server is compatible and this
// is a ClientStream, will spawn a goroutine that runs for the duration
// of the stream to send inline keepalives.
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

		ctx = metadata.AppendToOutgoingContext(ctx, HeaderSendKeepalives, "true")

		handler, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get the client handler when setting up the inlinekeepalive interceptor on method %q", method)
		}

		// Check compatibility and serve keepalives in a separate goroutine. This way, we don't
		// slow down the initiation of the stream.
		go func() {
			if !desc.ClientStreams {
				return
			}
			serverCompatible, err := isServerCompatible(ctx, cc)
			if err != nil {
				log.Warn("Failed to determine if server is capatible with inline keepalives - will not send them.", "err", err)
			}

			// Only send keepalives if this is a client stream - not allowed otherwise.
			if serverCompatible {
				// Send keepalives for as long as the handler has the connection open.
				ServeKeepalives(ctx, log, handler)
			}
		}()

		return &KeepaliveClientStream{handler: handler}, nil
	}
}
