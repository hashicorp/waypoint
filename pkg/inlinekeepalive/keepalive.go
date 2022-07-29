package inlinekeepalive

import (
	"context"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/protobuf/reflect/protoreflect"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

const KEEPALIVE_PROTO_SIGNATURE = "inline_keepalive"

func IsInlineKeepalive(log hclog.Logger, m protoreflect.ProtoMessage) bool {
	// performance enhancement
	unknownFields := m.ProtoReflect().GetUnknown()
	if unknownFields == nil {
		// No unknown fields here, not our keepalive, continue as normal
		return false
	}

	keepAlive := pb.KeepAlive{}
	err := proto.Unmarshal(unknownFields, &keepAlive)
	if err != nil {
		// couldn't marshal to the keepalive message, continue as normal
		return false
	}

	if keepAlive.Signature != KEEPALIVE_PROTO_SIGNATURE {
		// We had some other protobuf message with an unknown field whose number matched the high-order
		// keepalive message field number. This seems very unlikely.
		log.Warn("Unexpectedly received proto message with an unknown field matching the inlinekeepalive reserved field number")
		return false
	}

	return true
}

type GrpcStream interface {
	SendMsg(m interface{}) error
}

func ServeKeepalives(ctx context.Context, log hclog.Logger, stream GrpcStream) {
	log.Trace("Starting a inlinekeepalive interceptor for request")

	ticker := time.NewTicker(time.Duration(1) * time.Second)
	for {
		err := stream.SendMsg(&pb.KeepAlive{Signature: KEEPALIVE_PROTO_SIGNATURE})
		if err != nil {
			log.Trace("Failed sending inlinekeepalive", "err", err)
		}

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			log.Trace("Request complete - stopping inlinekeepalive interceptor")
			return
		}
	}
}
