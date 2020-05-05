package funcspec

import (
	"context"
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint/sdk/internal-shared/mapper"
	pb "github.com/hashicorp/waypoint/sdk/proto"
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

	// Defaults
	if cfg.Logger == nil {
		cfg.Logger = hclog.L()
	}
	if cfg.Output == nil {
		cfg.Output = protoMessageType
	}

	// Build our initial mapper
	mf, err := mapper.NewFunc(f, mapper.WithLogger(cfg.Logger), mapper.WithValues(cfg.Values...))
	if err != nil {
		return nil, err
	}

	// These check functions are used multiple times to check type impl.
	checkCtx := mapper.CheckReflectType(contextType)
	checkProto := mapper.CheckReflectType(protoMessageType)

	// Our input set can be ctx, proto, or any of the extra values given
	inputCheck := mapper.CheckOr(checkCtx, checkProto)
	for _, v := range cfg.Values {
		inputCheck = mapper.CheckOr(inputCheck, mapper.CheckReflectType(reflect.TypeOf(v)))
	}

	// We need to find a path through that only has protobuf requirements
	// or "context". These are the only given values to the func for plugins.
	types := mf.ChainInputSet(cfg.Mappers, inputCheck)
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
	checkOutput := mapper.CheckReflectType(cfg.Output)
	out := mf.Out
	if !cfg.NoOutput && !checkOutput(out) {
		chain := mapper.ChainTarget(checkOutput, cfg.Mappers)
		if chain == nil {
			return nil, fmt.Errorf(
				"function must output a type that is a %[1]s or has "+
					"a chain of mappers that result in a %[1]s", cfg.Output.String())
		}

		out = chain.Out()
	}
	if checkProto(out) {
		result.Result = typeToMessage(out)
	}

	return &result, nil
}

type config struct {
	Logger   hclog.Logger
	Mappers  []*mapper.Func
	Output   reflect.Type
	NoOutput bool
	Values   []interface{}
}

type Option func(*config)

// WithLogger specifies a logger. If this isn't specified then the default
// logger will be used.
func WithLogger(v hclog.Logger) Option {
	return func(c *config) { c.Logger = v }
}

// WithMappers specifies mappers and appends the mappers to the configuration.
// For Spec, these mappers will be used to verify that the arguments and
// result can be converted to proto.Message values. For Func, this has no
// effect.
func WithMappers(v []*mapper.Func) Option {
	return func(c *config) { c.Mappers = append(c.Mappers, v...) }
}

// WithOutput specifies the expected output type of the function. This
// defaults to a proto.Message by default. If this type is NOT a proto.Message
// then the protobuf FuncSpec "out" field will be set to blank indicating
// that it is some other type.
func WithOutput(t reflect.Type) Option {
	return func(c *config) { c.Output = t }
}

// WithNoOutput specifies that there is no output type expected, it is
// just a function that returns an error type.
func WithNoOutput() Option {
	return func(c *config) { c.NoOutput = true }
}

// WithValues specifies extra values, this just passes through to the
// mapper.WithValues.
func WithValues(vs ...interface{}) Option {
	return func(c *config) { c.Values = vs }
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
