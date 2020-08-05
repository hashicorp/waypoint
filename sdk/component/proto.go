package component

import (
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
func ProtoAny(m interface{}) (*any.Any, error) {
	msg, ok := m.(proto.Message)

	// If it isn't a message directly, we accept marshalers
	if !ok {
		pm, ok := m.(ProtoMarshaler)
		if !ok {
			return nil, nil
		}

		msg = pm.Proto()
	}

	// If the message is already an Any, then we're done
	if result, ok := msg.(*any.Any); ok {
		return result, nil
	}

	// Marshal it
	return ptypes.MarshalAny(msg)
}

// ProtoAny returns []*any.Any for the given input slice by encoding
// each result into a proto value.
func ProtoAnySlice(m interface{}) ([]*any.Any, error) {
	val := reflect.ValueOf(m)
	result := make([]*any.Any, val.Len())
	for i := 0; i < val.Len(); i++ {
		var err error
		result[i], err = ProtoAny(val.Index(i).Interface())
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// ProtoAnyUnmarshal attempts to unmarshal a ProtoMarshler implementation
// to another type. This can be used to get more concrete data out of a
// generic component.
func ProtoAnyUnmarshal(m interface{}, out proto.Message) error {
	msg, ok := m.(proto.Message)

	// If it isn't a message directly, we accept marshalers
	if !ok {
		pm, ok := m.(ProtoMarshaler)
		if !ok {
			return status.Errorf(codes.FailedPrecondition,
				"expected value to be a proto message, got %T",
				m)
		}

		msg = pm.Proto()
	}

	result, ok := msg.(*any.Any)
	if !ok {
		return status.Errorf(codes.FailedPrecondition, "expected *any.Any, got %T", msg)
	}

	// Unmarshal
	return ptypes.UnmarshalAny(result, out)
}
