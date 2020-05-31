package funcspec

import (
	"reflect"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hashicorp/go-argmapper"

	pb "github.com/hashicorp/waypoint/sdk/proto"
)

// Args is a type that will be populated with all the expected args of
// the FuncSpec. This can be used in the callback (cb) to Func.
type Args []*any.Any

// Func takes a FuncSpec and returns a *mapper.Func that can be called
// to invoke this function. The callback can have an argument type of Args
// in order to get access to the required dynamic proto.Any types of the
// FuncSpec.
func Func(s *pb.FuncSpec, cb interface{}, args ...argmapper.Arg) *argmapper.Func {
	// Build a Func around our callback so that we can inspect the
	// input/output sets since we want to merge with that.
	cbFunc, err := argmapper.NewFunc(cb)
	if err != nil {
		panic(err)
	}

	// Create the argmapper input values. All our args are expected to be
	// protobuf Any types that have a subtype matching our string name.
	// We append them directly to our expected values for the callback.
	// This lets us get our callback types in addition to our funcspec types.
	inputValues := cbFunc.Input().Values()
	for _, arg := range s.Args {
		inputValues = append(inputValues, argmapper.Value{
			Name:    arg.Name,
			Type:    anyType,
			Subtype: arg.Type,
		})
	}

	// Remove the Args value if there is one, since we're going to populate
	// that later and we don't need it for the initial call.
	for i, v := range inputValues {
		if v.Type == argsType {
			inputValues[i] = inputValues[len(inputValues)-1]
			inputValues = inputValues[:len(inputValues)-1]
			break
		}
	}

	inputSet, err := argmapper.NewValueSet(inputValues)
	if err != nil {
		panic(err)
	}

	// Build our output set. By default this just matches our output function.
	outputSet := cbFunc.Output()

	// If we have results specified on the Spec, then we expect this to represent
	// a mapper. Mapper callbacks MUST return *any.Any or []*any.Any. When we
	// have a mapper, we change the output type to be all the values we're
	// mapping to.
	if len(s.Result) > 0 {
		var outputValues []argmapper.Value
		for _, result := range s.Result {
			outputValues = append(outputValues, argmapper.Value{
				Name:    result.Name,
				Type:    anyType,
				Subtype: result.Type,
			})
		}

		outputSet, err = argmapper.NewValueSet(outputValues)
		if err != nil {
			panic(err)
		}
	}

	result, err := argmapper.BuildFunc(inputSet, outputSet, func(in, out *argmapper.ValueSet) error {
		callArgs := make([]argmapper.Arg, 0, len(args)+len(in.Values()))

		// Build up our callArgs which we'll pass to our callback. We pass
		// through all args except for *any.Any values. For *any values, we
		// add them to our Args list.
		var args Args
		for _, v := range in.Values() {
			// If we have any *any.Any then we append it to args
			if v.Type == anyType {
				args = append(args, v.Value.Interface().(*any.Any))
				continue
			}

			// If we have any other type, then we set it directly.
			callArgs = append(callArgs, v.Arg())
		}

		// Add our grouped Args type.
		callArgs = append(callArgs, argmapper.Typed(args))

		// Call into our callback. This populates our callback function output.
		cbOut := cbFunc.Output()
		if err := cbOut.FromResult(cbFunc.Call(callArgs...)); err != nil {
			return err
		}

		// If we aren't a mapper, we return now since we've populated our callback.
		if len(s.Result) == 0 {
			return nil
		}

		// We're a mapper, so we have to go through our values and look
		// for the *any.Any value or []*any.Any and populate our expected
		// outputs.
		for _, v := range cbOut.Values() {
			switch v.Type {
			case anyType:
				// We're seeing an *any.Any. So we encode this and try
				// to match it to any value that we have.
				anyVal := v.Value.Interface().(*any.Any)
				st, err := ptypes.AnyMessageName(anyVal)
				if err != nil {
					return err
				}

				expected := out.TypedSubtype(v.Type, st)
				if expected == nil {
					continue
				}

				expected.Value = v.Value
			}
		}

		// Go through our callback output looking
		return nil
	}, append([]argmapper.Arg{
		argmapper.ConverterGen(anyConvGen),
	}, args...)...)
	if err != nil {
		panic(err)
	}

	return result
}

var (
	anyType  = reflect.TypeOf((*any.Any)(nil))
	argsType = reflect.TypeOf(Args(nil))
)
