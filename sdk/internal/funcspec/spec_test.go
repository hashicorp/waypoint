package funcspec

import (
	"reflect"
	"testing"

	"github.com/hashicorp/waypoint/sdk/internal-shared/mapper"
	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/sdk/proto"
)

func TestSpec(t *testing.T) {
	t.Run("proto to proto", func(t *testing.T) {
		require := require.New(t)

		spec, err := Spec(func(*pb.Empty) *pb.Empty { return nil })
		require.NoError(err)
		require.NotNil(spec)
		require.Equal([]string{"proto.Empty"}, spec.Args)
		require.Equal("proto.Empty", spec.Result)
	})

	t.Run("converted args to proto", func(t *testing.T) {
		require := require.New(t)

		type Foo struct{}

		spec, err := Spec(func(*Foo) *pb.Empty { return nil }, WithMappers([]*mapper.Func{
			mustFunc(t, func(*pb.Empty) *Foo { return nil }),
		}))
		require.NoError(err)
		require.NotNil(spec)
		require.Equal([]string{"proto.Empty"}, spec.Args)
		require.Equal("proto.Empty", spec.Result)
	})

	t.Run("unsatisfied conversion", func(t *testing.T) {
		require := require.New(t)

		type Foo struct{}
		type Bar struct{}

		spec, err := Spec(func(*Foo) *pb.Empty { return nil }, WithMappers([]*mapper.Func{
			mustFunc(t, func(*pb.Empty) *Bar { return nil }),
		}))
		require.Error(err)
		require.Nil(spec)
	})

	t.Run("proto to int", func(t *testing.T) {
		require := require.New(t)

		spec, err := Spec(func(*pb.Empty) int { return 0 })
		require.Error(err)
		require.Nil(spec)
	})

	t.Run("WithOutput proto to int", func(t *testing.T) {
		require := require.New(t)

		spec, err := Spec(func(*pb.Empty) int { return 0 },
			WithOutput(reflect.TypeOf(int(0))))
		require.NoError(err)
		require.NotNil(spec)
		require.Equal([]string{"proto.Empty"}, spec.Args)
		require.Empty(spec.Result)
	})

	t.Run("WithOutput proto to interface, doesn't implement", func(t *testing.T) {
		require := require.New(t)

		spec, err := Spec(func(*pb.Empty) struct{} { return struct{}{} },
			WithOutput(reflect.TypeOf((*testSpecInterface)(nil)).Elem()))
		require.Error(err)
		require.Nil(spec)
	})

	t.Run("WithOutput proto to interface", func(t *testing.T) {
		require := require.New(t)

		spec, err := Spec(func(*pb.Empty) *testSpecInterfaceImpl { return nil },
			WithOutput(reflect.TypeOf((*testSpecInterface)(nil)).Elem()))
		require.NoError(err)
		require.NotNil(spec)
		require.Equal([]string{"proto.Empty"}, spec.Args)
		require.Empty(spec.Result)
	})

	t.Run("args as extra values", func(t *testing.T) {
		require := require.New(t)

		type Foo struct{}
		type Bar struct{}

		spec, err := Spec(func(*Foo, *Bar) *pb.Empty { return nil }, WithMappers([]*mapper.Func{
			mustFunc(t, func(*pb.Empty) *Foo { return nil }),
		}), WithValues(&Bar{}))
		require.NoError(err)
		require.NotNil(spec)
		require.Equal([]string{"proto.Empty"}, spec.Args)
		require.Equal("proto.Empty", spec.Result)
	})
}

func mustFunc(t *testing.T, f interface{}) *mapper.Func {
	t.Helper()

	result, err := mapper.NewFunc(f)
	require.NoError(t, err)
	return result
}

type testSpecInterface interface {
	hello()
}

type testSpecInterfaceImpl struct{}

func (testSpecInterfaceImpl) hello() {}
