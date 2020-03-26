package component

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
)

// ProtoMarshaler is the interface required by objects that must support
// protobuf marshaling. This expects the object to go to a proto.Message
// which is converted to a proto Any value[1]. The plugin is expected to
// register a proto type that can decode this Any value.
//
// This enables the project to encode intermediate objects (such as artifacts)
// and store them in a database.
//
// [1]: https://developers.google.com/protocol-buffers/docs/proto3#any
type ProtoMarshaler interface {
	// Proto returns a proto.Message of this structure. This may also return
	// a proto Any value and it will not be re-wrapped with Any.
	Proto() proto.Message
}

// ProtoAny returns an *any.Any for the given ProtoMarshaler object.
func ProtoAny(m ProtoMarshaler) (*any.Any, error) {
	msg := m.Proto()

	// If the message is already an Any, then we're done
	if result, ok := msg.(*any.Any); ok {
		return result, nil
	}

	// Marshal it
	return ptypes.MarshalAny(msg)
}
