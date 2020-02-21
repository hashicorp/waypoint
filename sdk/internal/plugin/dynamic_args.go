package plugin

import (
	"reflect"
	"sort"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"

	"github.com/mitchellh/devflow/sdk/pkg/mapper"
	pb "github.com/mitchellh/devflow/sdk/proto"
)

// dynamicArgs is the type expected for the argument for a dynamicArgsType.
// This will automatically map all the expected dynamic arguments into the
// function.
type dynamicArgs []*any.Any

// dynamicArgsType is the reflect.Type for dynamicArgs.
var dynamicArgsType = reflect.TypeOf((dynamicArgs)(nil))

// dynamicArgsMapperType implements mapper.Type to expect multiple *any.Any
// values with the matching types (as strings). This will match into an
// argument of type []*any.Any allowing the capture of all those values.
type dynamicArgsMapperType struct {
	Expected []string
}

// makeDynamicArgsMapperType can be used with mapper.WithType as a second
// parameter to construct the dynamicArgsMapperType for a func spec.
func makeDynamicArgsMapperType(s *pb.FuncSpec) func(int, reflect.Type) mapper.Type {
	return func(int, reflect.Type) mapper.Type {
		return &dynamicArgsMapperType{
			Expected: s.Args,
		}
	}
}

// Match implements mapper.Type by constructing an []*any.Any if there
// exists an *any.Any for all expected types.
func (t *dynamicArgsMapperType) Match(values ...interface{}) interface{} {
	expectMap := make(map[string]struct{})
	for _, v := range t.Expected {
		expectMap[v] = struct{}{}
	}

	var result []*any.Any
	for _, raw := range values {
		av, ok := raw.(*any.Any)

		// If this isn't an *any.Any, we can still take a proto.Message
		// and manually encode it. This path is really only used for our
		// built-in types since any custom types are never going to be
		// decoded in core (since we don't link against plugins directly).
		if !ok {
			pv, ok := raw.(proto.Message)
			if !ok {
				continue
			}

			// If we don't expect this value, then ignore
			if _, ok := expectMap[proto.MessageName(pv)]; !ok {
				continue
			}

			var err error
			av, err = ptypes.MarshalAny(pv)
			if err != nil {
				continue
			}
		}

		// If this value isn't an Any then ignore it
		if av == nil {
			continue
		}

		// If this value isn't in the map of expected types, ignore it
		key, err := ptypes.AnyMessageName(av)
		if err != nil {
			continue
		}
		if _, ok := expectMap[key]; !ok {
			continue
		}

		// A match, record it.
		result = append(result, av)

		// Delete the value from the map so we don't match it again.
		// We only take the first matching any type since there should be
		// exactly one. This is how mapper works: type matching.
		delete(expectMap, key)
	}

	// If we're missing any expected values, then we can't match.
	if len(expectMap) > 0 {
		return nil
	}

	return result
}

func (t *dynamicArgsMapperType) Key() interface{} {
	// Our string value is a unique key that is stable (sorted)
	return t.String()
}

func (t *dynamicArgsMapperType) String() string {
	sort.Strings(t.Expected)
	return "protobuf Any types: " + strings.Join(t.Expected, ", ")
}

var _ mapper.Type = (*dynamicArgsMapperType)(nil)
