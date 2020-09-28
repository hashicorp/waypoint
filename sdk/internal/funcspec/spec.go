package funcspec

import (
	"context"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-argmapper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/sdk/proto"
)

// Spec takes a function pointer and generates a FuncSpec from it. The
// function must only take arguments that are proto.Message implementations
// or have a chain of converters that directly convert to a proto.Message.
func Spec(fn interface{}, args ...argmapper.Arg) (*pb.FuncSpec, error) {
	if fn == nil {
		return nil, status.Errorf(codes.Unimplemented, "required plugin type not implemented")
	}

	filterProto := argmapper.FilterType(protoMessageType)

	// Copy our args cause we're going to use append() and we don't
	// want to modify our caller.
	args = append([]argmapper.Arg{
		argmapper.FilterOutput(filterProto),
	}, args...)

	// Build our function
	f, err := argmapper.NewFunc(fn)
	if err != nil {
		return nil, err
	}

	filter := argmapper.FilterOr(
		argmapper.FilterType(contextType),
		filterProto,
	)

	// Redefine the function in terms of protobuf messages. "Redefine" changes
	// the inputs of a function to only require values that match our filter
	// function. In our case, that is protobuf messages.
	f, err = f.Redefine(append(args,
		argmapper.FilterInput(filter),
	)...)
	if err != nil {
		return nil, err
	}

	// Grab the input set of the function and build up our funcspec
	result := pb.FuncSpec{Name: f.Name()}
	for _, v := range f.Input().Values() {
		if !filterProto(v) {
			continue
		}

		result.Args = append(result.Args, &pb.FuncSpec_Value{
			Name: v.Name,
			Type: typeToMessage(v.Type),
		})
	}

	// Grab the output set and store that
	for _, v := range f.Output().Values() {
		// We only advertise proto types in output since those are the only
		// types we can send across the plugin boundary.
		if !filterProto(v) {
			continue
		}

		result.Result = append(result.Result, &pb.FuncSpec_Value{
			Name: v.Name,
			Type: typeToMessage(v.Type),
		})
	}

	return &result, nil
}

func typeToMessage(typ reflect.Type) string {
	return proto.MessageName(reflect.Zero(typ).Interface().(proto.Message))
}

var (
	contextType      = reflect.TypeOf((*context.Context)(nil)).Elem()
	protoMessageType = reflect.TypeOf((*proto.Message)(nil)).Elem()
)
