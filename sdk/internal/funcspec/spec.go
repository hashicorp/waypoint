package funcspec

import (
	"context"
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-hclog"

	"github.com/mitchellh/devflow/sdk/internal-shared/mapper"
	pb "github.com/mitchellh/devflow/sdk/proto"
)

// Spec takes a function pointer and generates a FuncSpec from it. The
// function must only take arguments that are proto.Message implementations
// or have a chain of mappers that directly convert to a proto.Message.
func Spec(f interface{}, opts ...Option) (*pb.FuncSpec, error) {
	// Build our configuration
	var cfg config
	for _, opt := range opts {
		opt(&cfg)
	}

	// Build our initial mapper
	mf, err := mapper.NewFunc(f, mapper.WithLogger(cfg.Logger))
	if err != nil {
		return nil, err
	}

	// These check functions are used multiple times to check type impl.
	checkCtx := mapper.CheckReflectType(contextType)
	checkProto := mapper.CheckReflectType(protoMessageType)

	// We need to find a path through that only has protobuf requirements
	// or "context". These are the only given values to the func for plugins.
	types := mf.ChainInputSet(cfg.Mappers, mapper.CheckOr(checkCtx, checkProto))
	if len(types) == 0 {
		return nil, fmt.Errorf(
			"cannot satisfy the function %s. The function takes arguments that "+
				"are not proto.Messages or have no mappers to convert to proto.Messages",
			mf)
	}

	// Build our FuncSpec. The name we use is just the name on this side.
	result := pb.FuncSpec{Name: mf.Name}

	// For each type, get the Any message name for it.
	for _, t := range types {
		// Ignore any non-proto types. We verify above that we only get
		// types we accept so we don't care what the other types are at this
		// point.
		if !checkProto(t) {
			continue
		}

		// If we're here we know its a proto.Message
		result.Args = append(result.Args, typeToMessage(t))
	}

	// Get the result type. If it isn't a proto message, we look for a chain
	// to get us to that proto message.
	out := mf.Out
	if !checkProto(out) {
		chain := mapper.ChainTarget(checkProto, cfg.Mappers)
		if chain == nil {
			return nil, fmt.Errorf(
				"function must output a type that is a proto.Message or has " +
					"a chain of mappers that result in a proto.Message")
		}

		out = chain.Out()
	}
	result.Result = typeToMessage(out)

	return &result, nil
}

type config struct {
	Logger  hclog.Logger
	Mappers []*mapper.Func
}

type Option func(*config)

func WithLogger(v hclog.Logger) Option {
	return func(c *config) { c.Logger = v }
}

func WithMappers(v []*mapper.Func) Option {
	return func(c *config) { c.Mappers = append(c.Mappers, v...) }
}

// typeToMessage converts a mapper.Type to the proto.Message name value.
//
// preconditions:
//   - t is a ReflectType
//   - the typ represented by t is a proto.Message
func typeToMessage(t mapper.Type) string {
	typ := t.(*mapper.ReflectType).Type
	return proto.MessageName(reflect.Zero(typ).Interface().(proto.Message))
}

var (
	contextType      = reflect.TypeOf((*context.Context)(nil)).Elem()
	protoMessageType = reflect.TypeOf((*proto.Message)(nil)).Elem()
)
