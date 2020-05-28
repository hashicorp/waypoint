package funcspec

import (
	"reflect"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/mitchellh/go-argmapper"

	pb "github.com/hashicorp/waypoint/sdk/proto"
)

// Args is a type that will be populated with all the expected args of
// the FuncSpec. This can be used in the callback (cb) to Func.
type Args []*any.Any

// Func takes a FuncSpec and returns a *mapper.Func that can be called
// to invoke this function. The callback can have an argument type of Args
// in order to get access to the required dynamic proto.Any types of the
// FuncSpec.
func Func(s *pb.FuncSpec2, cb interface{}) (*argmapper.Func, error) {
	// Build a Func around our callback so that we can inspect the
	// input/output sets since we want to merge with that.
	cbFunc, err := argmapper.NewFunc(cb)
	if err != nil {
		return nil, err
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
	// that later.
	for i, v := range inputValues {
		if v.Type == argsType {
			inputValues[i] = inputValues[len(inputValues)-1]
			inputValues = inputValues[:len(inputValues)-1]
			break
		}
	}

	inputSet, err := argmapper.NewValueSet(inputValues)
	if err != nil {
		return nil, err
	}

	// Build our output values based on the advertised result value.
	outputSet := cbFunc.Output()
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
			return nil, err
		}
	}

	return argmapper.BuildFunc(inputSet, outputSet, func(in, out *argmapper.ValueSet) error {
		// Build up our callArgs which we'll pass to our callback.
		var args Args
		var callArgs []argmapper.Arg
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

		// Call into our callback. We should have the exact arguments
		// required in our args list since we merged them earlier. We
		// extract this into our callback output results.
		cbOut := cbFunc.Output()
		if err := cbOut.FromResult(cbFunc.Call(callArgs...)); err != nil {
			return err
		}

		// Go through our output types from the callback
		// TODO docs
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
	})
}

var (
	anyType  = reflect.TypeOf((*any.Any)(nil))
	argsType = reflect.TypeOf(Args(nil))
)
