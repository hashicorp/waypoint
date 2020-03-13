package funcspec

import (
	"reflect"
	"sort"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"

	"github.com/mitchellh/devflow/sdk/internal-shared/mapper"
	pb "github.com/mitchellh/devflow/sdk/proto"
)

// Args is the type expected for the argument for an argsType. An argument
// of this type to a mapper.Func that has the ArgsMapperType configured will
// automatically populate this argument with matching proto messages.
type Args []*any.Any

// ArgsType is the reflect.Type for ArgsType.
var ArgsType = reflect.TypeOf((Args)(nil))

// ArgsMapperType implements mapper.Type to expect multiple *any.Any
// values with the matching types (as strings). This will match into an
// argument of type []*any.Any allowing the capture of all those values.
type ArgsMapperType struct {
	Expected []string
}

// AppendArgs is used to append proto.Message values to a list of Args.
// This will automatically marshal each message to an Any type.
func AppendArgs(args Args, ms ...proto.Message) (Args, error) {
	for _, m := range ms {
		encoded, err := ptypes.MarshalAny(m)
		if err != nil {
			return nil, err
		}

		args = append(args, encoded)
	}

	return args, nil
}

// makeArgsMapperType can be used with mapper.WithType as a second
// parameter to construct the ArgsMapperType for a func spec.
func makeArgsMapperType(s *pb.FuncSpec) func(int, reflect.Type) mapper.Type {
	return func(int, reflect.Type) mapper.Type {
		return &ArgsMapperType{
			Expected: s.Args,
		}
	}
}

// Match implements mapper.Type by constructing an []*any.Any if there
// exists an *any.Any for all expected types.
func (t *ArgsMapperType) Match(values ...interface{}) interface{} {
	return t.args(values, nil)
}

func (t *ArgsMapperType) Missing(values ...interface{}) []mapper.Type {
	var missing []mapper.Type
	t.args(values, &missing)
	return missing
}

func (t *ArgsMapperType) args(
	values []interface{},
	missing *[]mapper.Type, // if non-nil, will be populated with missing types
) interface{} {
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
		// If we don't need to track missing values, just return
		if missing == nil {
			return nil
		}

		// We are tracking missing values so go through the expect typed
		// and if we know the proto type (it is registered locally) then
		// offer that.
		for name := range expectMap {
			if typ := proto.MessageType(name); typ != nil {
				*missing = append(*missing, &mapper.ReflectType{
					Type: typ,
				})

				delete(expectMap, name)
			}
		}

		// If we have any we still expect, then use the default type
		if len(expectMap) > 0 {
			*missing = nil
		}
	}

	return result
}

func (t *ArgsMapperType) Key() interface{} {
	// Our string value is a unique key that is stable (sorted)
	return t.String()
}

func (t *ArgsMapperType) String() string {
	sort.Strings(t.Expected)
	return "protobuf Any types: " + strings.Join(t.Expected, ", ")
}

// Assertion
var _ mapper.Type = (*ArgsMapperType)(nil)
