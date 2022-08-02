package inlinekeepalive

import (
	"context"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/protobuf/reflect/protoreflect"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

const (
	KeepaliveProtoSignature   = "inline_keepalive"
	HeaderSendKeepalivesKey   = "wp-inline-keepalives"
	HeaderSendKeepalivesValue = "true"
	FeatureName               = "inline-keepalives"
)

// IsInlineKeepalive determines if a given proto message is
// an inline keepalive.
func IsInlineKeepalive(log hclog.Logger, m protoreflect.ProtoMessage) bool {
	unknownFields := m.ProtoReflect().GetUnknown()
	if unknownFields == nil {
		// No unknown fields here, not our keepalive, continue as normal
		return false
	}

	keepAlive := pb.InlineKeepalive{}
	err := proto.Unmarshal(unknownFields, &keepAlive)
	if err != nil {
		// couldn't marshal to the keepalive message, continue as normal
		return false
	}

	if keepAlive.Signature != KeepaliveProtoSignature {
		// We had some other protobuf message with an unknown field whose number matched the high-order
		// keepalive message field number. This seems very unlikely.
		log.Warn("Unexpectedly received proto message with an unknown field matching the inlinekeepalive reserved field number")
		return false
	}

	return true
}

// GrpcStream can be either a grpc.ClientStream or a grpc.ServerStream
type GrpcStream interface {
	SendMsg(m interface{}) error
}

// ServeKeepalives sends keepalive messages along the provided grpc stream
// at a rate of one every five seconds.
// It returns when the context is cancelled.
// NOTE: this will call SendMsg, and concurrent calls to SendMsg are unsafe.
// This will not call SendMsg unless it holds the sendMx lock.
func ServeKeepalives(
	ctx context.Context,
	log hclog.Logger,
	stream GrpcStream,
	sendInterval time.Duration,
	sendMx *sync.Mutex,
) {
	log.Trace("Starting a inlinekeepalive interceptor for request")

	intervalTicker := time.NewTicker(sendInterval)
	for {
		sendMx.Lock()
		err := stream.SendMsg(&pb.InlineKeepalive{Signature: KeepaliveProtoSignature})
		sendMx.Unlock()
		if err != nil {
			log.Warn("Failed sending inlinekeepalive", "err", err)
		}

		select {
		case <-intervalTicker.C:
			continue
		case <-ctx.Done():
			log.Trace("Request complete - stopping inlinekeepalive interceptor")
			return
		}
	}
}
