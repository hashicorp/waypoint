package plugin

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/require"

	"github.com/mitchellh/devflow/sdk/internal-shared/mapper"
	pb "github.com/mitchellh/devflow/sdk/proto"
)

func TestDynamicArgsMapperType_Match(t *testing.T) {
	t.Run("simple match", func(t *testing.T) {
		require := require.New(t)
		f, args := testDynamicArgsFunc(t, []string{"foo"})
		result, err := f.Call(&any.Any{TypeUrl: "example.com/foo"}, &any.Any{TypeUrl: "example.com/bar"})
		require.NoError(err)
		require.NotNil(result)
		require.Len(*args, 1)
		require.Contains(*args, &any.Any{TypeUrl: "example.com/foo"})
	})

	t.Run("no matches", func(t *testing.T) {
		require := require.New(t)
		f, _ := testDynamicArgsFunc(t, []string{"foo"})
		require.Panics(func() {
			f.Call(&any.Any{TypeUrl: "example.com/baz"}, &any.Any{TypeUrl: "example.com/bar"})
		})
	})

	t.Run("match proto.Message", func(t *testing.T) {
		data := &pb.Args_Source{}

		require := require.New(t)
		f, args := testDynamicArgsFunc(t, []string{proto.MessageName(data)})
		result, err := f.Call(data)
		require.NoError(err)
		require.NotNil(result)
		require.Len(*args, 1)

		var value pb.Args_Source
		require.NoError(ptypes.UnmarshalAny((*args)[0], &value))
		require.Equal(data, &value)
	})

	t.Run("multiple requirements: match", func(t *testing.T) {
		require := require.New(t)
		f, args := testDynamicArgsFunc(t, []string{"foo", "bar"})
		result, err := f.Call(
			&any.Any{TypeUrl: "example.com/baz"},
			&any.Any{TypeUrl: "example.com/bar"},
			&any.Any{TypeUrl: "example.com/foo"},
		)

		require.NoError(err)
		require.NotNil(result)
		require.Len(*args, 2)
		require.Contains(*args, &any.Any{TypeUrl: "example.com/foo"})
		require.Contains(*args, &any.Any{TypeUrl: "example.com/bar"})
	})

	t.Run("multiple requirements: missing one", func(t *testing.T) {
		require := require.New(t)
		f, _ := testDynamicArgsFunc(t, []string{"foo", "bar"})
		require.Panics(func() {
			f.Call(
				&any.Any{TypeUrl: "example.com/baz"},
				&any.Any{TypeUrl: "example.com/foo"},
			)
		})
	})
}

func TestDynamicArgsMapperType_Missing(t *testing.T) {
	t.Run("missing only known types", func(t *testing.T) {
		require := require.New(t)

		typ := &dynamicArgsMapperType{Expected: []string{
			proto.MessageName(&pb.Args_Source{}),
			proto.MessageName(&pb.Args_DataDir_Project{}),
		}}

		types := typ.Missing()
		require.NotNil(types)
		require.Len(types, 2)
	})

	t.Run("missing unregistered type", func(t *testing.T) {
		require := require.New(t)

		typ := &dynamicArgsMapperType{Expected: []string{
			proto.MessageName(&pb.Args_Source{}),
			"bar",
		}}

		types := typ.Missing()
		require.Nil(types)
	})
}

// testDynamicArgsFunc returns a new mapper func that when called sets the
// returned dynamicArgs pointer value to the value called.
//
// The types argument is the types that are expected.
func testDynamicArgsFunc(t *testing.T, types []string) (*mapper.Func, *dynamicArgs) {
	var result dynamicArgs
	f, err := mapper.NewFunc(func(args dynamicArgs) int {
		result = args
		return 0
	}, mapper.WithType(dynamicArgsType, makeDynamicArgsMapperType(&pb.FuncSpec{
		Args: types,
	})))
	require.NoError(t, err)

	return f, &result
}
