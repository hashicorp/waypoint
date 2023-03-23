// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package inlinekeepalive

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/protobuf/reflect/protoreflect"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

const (
	KeepaliveProtoSignature     = "inline_keepalive"
	GrpcMetaSendKeepalivesKey   = "wp-inline-keepalives"
	GrpcMetaSendKeepalivesValue = "true"
)

// IsInlineKeepalive determines if a given proto message is
// an inline keepalive.
func IsInlineKeepalive(log hclog.Logger, m protoreflect.ProtoMessage) bool {

	// One might think a type assertion, something like `m.(*pb.InlineKeepalive)`,
	// would be sufficient here. In fact, m is _always_ going to be the type that
	// we expect to receive via the proto spec, and never the inlinekeepalive,
	// EVEN IF what we sent was an inlinekeepalive.
	// If this is an inlinekeepalive, all the fields on `m` will be nil, but the
	// protoreflect output will tell us what fields it received that it wasn't expecting.
	// For a inlinekeepalive message, that will include field 100000000.
	unknownFields := m.ProtoReflect().GetUnknown()
	if unknownFields == nil {
		// No unknown fields here, not our keepalive, continue as normal
		return false
	}

	// unknownFields is binary protobuf representation of the unknown fields.
	// For an unknown keepalive, this actually the full protobuf message.
	// See https://developers.google.com/protocol-buffers/docs/encoding#structure,
	// > "The binary version of a message just uses the field's number as the key
	//    – the name and declared type for each field can only be determined on
	//    the decoding end by referencing the message type's definition (i.e. the
	//    .proto file)."
	// The name or type of the value isn't encoded in the protobuf, just the field
	// number. We can take that data, proto-unmarshall it into an inline keepalive message,
	// and if the unknown field number was 10000000, the unmarshal will succeed and
	// we'll get a value.
	keepAlive := pb.InlineKeepalive{}
	err := proto.Unmarshal(unknownFields, &keepAlive)
	if err != nil {
		// couldn't marshal to the keepalive message, continue as normal
		return false
	}

	if keepAlive.Signature != KeepaliveProtoSignature {
		// We had some other protobuf message with an unknown field whose number matched the high-order
		// keepalive message field number. The other message's field 10000000 must have also been of
		// type "string", or we would have failed out earlier.
		// This is a tricky case. We'll return false here, because we're sure that the current message
		// isn't an inline keepalive. HOWEVER, if the other party sends us a real inline keepalive message,
		// we will succeed at marshalling that to the real message's field 10000000. That message will
		// have nil values everywhere _except_ that high-order field, which will have a value of "inline_keepalive".
		// That message will be passed along to the real RPC handler, and there is no predicting what that handler
		// would do with such a message.
		log.Error("Unexpectedly received proto message with an unknown field matching the inlinekeepalive reserved field number. If inline keepalives are sent by the other party, unexpected behavior will occur.")
		return false
	}

	return true
}

// GrpcStream can be either a grpc.ClientStream or a grpc.ServerStream
type GrpcStream interface {
	SendMsg(m interface{}) error
}

// ServeKeepalives sends keepalive messages along the provided grpc stream
// at the rate specified by sendInterval.
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
		// NOTE(izaak): It's possible here that we're attempting to send after CloseSend
		// has been called, but before the context has been canceled. We could avoid that
		// by adding another mutex, but I think it's OK.
		err := stream.SendMsg(&pb.InlineKeepalive{Signature: KeepaliveProtoSignature})
		sendMx.Unlock()

		if err != nil {
			if err == io.EOF {
				log.Trace("topping inline keepalive server - received EOF on SendMsg")
				return
			}
			log.Debug("Failed sending inlinekeepalive", "err", err)
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
