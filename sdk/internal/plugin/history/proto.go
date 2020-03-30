package history

import (
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"

	"github.com/mitchellh/devflow/sdk/component"
)

func unmarshalSlice(in, out interface{}) error {
	outVal := reflect.ValueOf(out)

	inVal := reflect.ValueOf(in)
	for i := 0; i < inVal.Len(); i++ {
		raw := inVal.Index(i).Interface()

		// Get our proto out
		if pm, ok := raw.(component.ProtoMarshaler); ok {
			raw = pm.Proto()
		}

		// We expect an *any.Any directly.
		av, ok := raw.(*any.Any)
		if !ok {
			return fmt.Errorf("non-Any value found in slice: %#v", raw)
		}

		name, err := ptypes.AnyMessageName(av)
		if err != nil {
			return err
		}

		typ := proto.MessageType(name)
		if typ == nil {
			return fmt.Errorf("cannot decode type: %s", name)
		}

		// Allocate the message type. If it is a pointer we want to
		// allocate the actual structure and not the pointer to the structure.
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		v := reflect.New(typ)
		v.Elem().Set(reflect.Zero(typ))

		// Unmarshal directly into our newly allocated structure.
		if err := ptypes.UnmarshalAny(av, v.Interface().(proto.Message)); err != nil {
			return err
		}

		outVal.Elem().Set(reflect.Append(outVal.Elem(), v))
	}

	return nil
}
