package inlinekeepalive

import (
	"context"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// KeepaliveClientStream implements grpc.ClientStream
type KeepaliveClientStream struct {
	log     hclog.Logger
	handler grpc.ClientStream
	sendMx  *sync.Mutex
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
	// Concurrent calls to SendMsg are unsafe, and there may be another
	// goroutine sending keepalives. Lock before sending.
	k.sendMx.Lock()
	defer k.sendMx.Unlock()

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

// KeepaliveClientStreamInterceptor returns a stream interceptor
// that sends inline keepalive messages on client streams (if the server
// is compatible), and intercepts inline keepalives from the server.
// This is intended to be invoked once at the beginning of an RPC, may call
// the server's GetVersionInfo RPC, and if the server is compatible and this
// is a ClientStream, will spawn a goroutine that runs for the duration
// of the stream to send inline keepalives.
// Will send a keepalive every sendInterval
func KeepaliveClientStreamInterceptor(sendInterval time.Duration) grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		log := hclog.FromContext(ctx).With("method", method)

		ctx = metadata.AppendToOutgoingContext(ctx, GrpcMetaSendKeepalivesKey, GrpcMetaSendKeepalivesValue)

		handler, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get the client handler when setting up the inlinekeepalive interceptor on method %q", method)
		}

		// Ensures SendMsg is not called concurrently.
		sendMx := &sync.Mutex{}

		// Check compatibility and serve keepalives in a separate goroutine. This way, we don't
		// slow down the initiation of the stream.
		go func() {
			if !desc.ClientStreams {
				return
			}

			client := pb.NewWaypointClient(cc)

			versionInfo, err := client.GetVersionInfo(ctx, &emptypb.Empty{})
			if err != nil {
				if status.Code(err) == codes.Canceled || status.Code(err) == codes.Unavailable {
					log.Trace("context canceled while determining if server is compatible with inline keepalives - will not send.", "err", err)
					return
				}
				log.Warn("failed getting version info to determine if server is inline-keepalive compatible - will not send them", "err", err)
				return
			}

			isCompatible := false
			if versionInfo.ServerFeatures != nil {
				for _, feature := range versionInfo.ServerFeatures.Features {
					if feature == pb.ServerFeatures_FEATURE_INLINE_KEEPALIVES {
						isCompatible = true
						break
					}
				}
			}

			if !isCompatible {
				log.Trace("server not compatible with inline keepalives - will not send them")
				return
			}

			// Send keepalives for as long as the handler has the connection open.
			ServeKeepalives(ctx, log, handler, sendInterval, sendMx)
			log.Trace("stopped sending inline keepalives")
		}()

		return &KeepaliveClientStream{
			handler: handler,
			log:     log.With("interceptor", "inlinekeepalive"),
			sendMx:  sendMx,
		}, nil
	}
}
