package inlinekeepalive

import (
	"context"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/protobuf/reflect/protoreflect"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

const (
	KeepaliveProtoSignature = "inline_keepalive"
	HeaderSendKeepalives    = "wp-inline-keepalives"
	FeatureName             = "inline-keepalives"
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
func ServeKeepalives(
	ctx context.Context,
	log hclog.Logger,
	stream GrpcStream,
) {
	log.Trace("Starting a inlinekeepalive interceptor for request")

	intervalTicker := time.NewTicker(time.Duration(5) * time.Second)
	for {
		err := stream.SendMsg(&pb.InlineKeepalive{Signature: KeepaliveProtoSignature})
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
