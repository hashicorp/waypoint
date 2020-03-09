package plugin

import (
	"context"
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hashicorp/go-hclog"

	"github.com/mitchellh/devflow/sdk/internal-shared/mapper"
)

// callDynamicFunc calls a dynamic (mapper-based) function with the
// given input arguments. This is a helper that is expected to be used
// by most component gRPC servers to implement their function calls.
func callDynamicFunc(
	ctx context.Context,
	log hclog.Logger,
	args []*any.Any,
	f interface{},
	mappers []*mapper.Func,
) (interface{}, error) {
	// Decode all our arguments. We are on the plugin side now so we expect
	// to be able to decode all types sent to us.
	decoded := make([]interface{}, len(args)+1)
	decoded[0] = ctx
	for idx, arg := range args {
		name, err := ptypes.AnyMessageName(arg)
		if err != nil {
			return nil, err
		}

		typ := proto.MessageType(name)
		if typ == nil {
			return nil, fmt.Errorf("cannot decode type: %s", name)
		}

		// Allocate the message type. If it is a pointer we want to
		// allocate the actual structure and not the pointer to the structure.
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		v := reflect.New(typ)
		v.Elem().Set(reflect.Zero(typ))

		// Unmarshal directly into our newly allocated structure.
		if err := ptypes.UnmarshalAny(arg, v.Interface().(proto.Message)); err != nil {
			return nil, err
		}

		decoded[idx+1] = v.Interface()
	}

	// Build our mapper function and find the chain to get us to the required
	// arguments if possible. This chain will do things like convert from
	// our raw proto types to richer structures if the plugin expects that.
	mf, err := mapper.NewFunc(f, mapper.WithLogger(log))
	if err != nil {
		return nil, err
	}

	chain, err := mf.Chain(mappers, decoded...)
	if err != nil {
		return nil, err
	}

	return chain.Call()
}

// callDynamicFuncAny is callDynamicFunc that automatically encodes the
// result to an *any.Any.
func callDynamicFuncAny(
	ctx context.Context,
	log hclog.Logger,
	args []*any.Any,
	f interface{},
	mappers []*mapper.Func,
) (*any.Any, error) {
	result, err := callDynamicFunc(ctx, log, args, f, mappers)
	if err != nil {
		return nil, err
	}

	// We expect the final result to always be a proto message so we can
	// send it back over the wire.
	//
	// NOTE(mitchellh): If we wanted to in the future, we can probably change
	// this to be any type that has a mapper that can take it to be a
	// proto.Message.
	msg, ok := result.(proto.Message)
	if !ok {
		return nil, fmt.Errorf(
			"result of plugin-based function must be a proto.Message, got %T", msg)
	}

	return ptypes.MarshalAny(msg)
}
