package plugin

import (
	"testing"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/require"

	"github.com/mitchellh/devflow/sdk/pkg/mapper"
	"github.com/mitchellh/devflow/sdk/proto"
)

func TestDynamicArgsMapperType(t *testing.T) {
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

// testDynamicArgsFunc returns a new mapper func that when called sets the
// returned dynamicArgs pointer value to the value called.
//
// The types argument is the types that are expected.
func testDynamicArgsFunc(t *testing.T, types []string) (*mapper.Func, *dynamicArgs) {
	var result dynamicArgs
	f, err := mapper.NewFunc(func(args dynamicArgs) int {
		result = args
		return 0
	}, mapper.WithType(dynamicArgsType, makeDynamicArgsMapperType(&proto.FuncSpec{
		Args: types,
	})))
	require.NoError(t, err)

	return f, &result
}
