package funcspec

import (
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hashicorp/go-argmapper"
)

// anyConvGen is an argmapper.ConverterGenFunc that dynamically creates
// converters to *any.Any for types that implement proto.Message. This
// allows automatic conversion to *any.Any.
//
// This is automatically injected for all funcspec.Func calls.
func anyConvGen(v argmapper.Value) (*argmapper.Func, error) {
	anyType := reflect.TypeOf((*any.Any)(nil))
	protoMessageType := reflect.TypeOf((*proto.Message)(nil)).Elem()
	if !v.Type.Implements(protoMessageType) {
		return nil, nil
	}

	// We take this value as our input.
	inputSet, err := argmapper.NewValueSet([]argmapper.Value{v})
	if err != nil {
		return nil, err
	}

	// Generate an int with the subtype of the string value
	outputSet, err := argmapper.NewValueSet([]argmapper.Value{argmapper.Value{
		Name:    v.Name,
		Type:    anyType,
		Subtype: proto.MessageName(reflect.Zero(v.Type).Interface().(proto.Message)),
	}})
	if err != nil {
		return nil, err
	}

	return argmapper.BuildFunc(inputSet, outputSet, func(in, out *argmapper.ValueSet) error {
		anyVal, err := ptypes.MarshalAny(inputSet.Typed(v.Type).Value.Interface().(proto.Message))
		if err != nil {
			return err
		}

		outputSet.Typed(anyType).Value = reflect.ValueOf(anyVal)
		return nil
	})
}
