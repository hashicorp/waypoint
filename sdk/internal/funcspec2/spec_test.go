package funcspec

import (
	"reflect"
	"testing"

	"github.com/mitchellh/go-argmapper"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/go-hclog"
	pb "github.com/hashicorp/waypoint/sdk/proto"
)

func init() {
	hclog.L().SetLevel(hclog.Trace)
}

func TestSpec(t *testing.T) {
	t.Run("proto to proto", func(t *testing.T) {
		require := require.New(t)

		spec, err := Spec(func(*pb.Empty) *pb.Empty { return nil })
		require.NoError(err)
		require.NotNil(spec)
		require.Len(spec.Args, 1)
		require.Empty(spec.Args[0].Name)
		require.Equal("proto.Empty", spec.Args[0].Type)
		require.Len(spec.Result, 1)
		require.Empty(spec.Result[0].Name)
		require.Equal("proto.Empty", spec.Result[0].Type)
	})

	t.Run("converted args to proto", func(t *testing.T) {
		require := require.New(t)

		type Foo struct{}

		spec, err := Spec(func(*Foo) *pb.Empty { return nil },
			argmapper.Converter(func(*pb.Empty) *Foo { return nil }),
		)
		require.NoError(err)
		require.NotNil(spec)
		require.Len(spec.Args, 1)
		require.Empty(spec.Args[0].Name)
		require.Equal("proto.Empty", spec.Args[0].Type)
		require.Len(spec.Result, 1)
		require.Empty(spec.Result[0].Name)
		require.Equal("proto.Empty", spec.Result[0].Type)
	})

	t.Run("unsatisfied conversion", func(t *testing.T) {
		require := require.New(t)

		type Foo struct{}
		type Bar struct{}

		spec, err := Spec(func(*Foo) *pb.Empty { return nil },
			argmapper.Converter(func(*pb.Empty) *Bar { return nil }),
		)
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
			argmapper.FilterOutput(argmapper.FilterType(reflect.TypeOf(int(0)))),
		)
		require.NoError(err)
		require.NotNil(spec)
		require.Len(spec.Args, 1)
		require.Empty(spec.Args[0].Name)
		require.Equal("proto.Empty", spec.Args[0].Type)
		require.Empty(spec.Result)
	})

	t.Run("WithOutput proto to interface, doesn't implement", func(t *testing.T) {
		require := require.New(t)

		spec, err := Spec(func(*pb.Empty) struct{} { return struct{}{} },
			argmapper.FilterOutput(argmapper.FilterType(reflect.TypeOf((*testSpecInterface)(nil)).Elem())),
		)
		require.Error(err)
		require.Nil(spec)
	})

	t.Run("WithOutput proto to interface", func(t *testing.T) {
		require := require.New(t)

		spec, err := Spec(func(*pb.Empty) *testSpecInterfaceImpl { return nil },
			argmapper.FilterOutput(argmapper.FilterType(reflect.TypeOf((*testSpecInterface)(nil)).Elem())),
		)
		require.NoError(err)
		require.NotNil(spec)
		require.Len(spec.Args, 1)
		require.Empty(spec.Args[0].Name)
		require.Equal("proto.Empty", spec.Args[0].Type)
		require.Empty(spec.Result)
	})

	t.Run("args as extra values", func(t *testing.T) {
		require := require.New(t)

		type Foo struct{}
		type Bar struct{}

		spec, err := Spec(func(*Foo, *Bar) *pb.Empty { return nil },
			argmapper.Converter(func(*pb.Empty) *Foo { return nil }),
			argmapper.Typed(&Bar{}),
		)
		require.NoError(err)
		require.NotNil(spec)
		require.Len(spec.Args, 1)
		require.Empty(spec.Args[0].Name)
		require.Equal("proto.Empty", spec.Args[0].Type)
		require.Len(spec.Result, 1)
		require.Empty(spec.Result[0].Name)
		require.Equal("proto.Empty", spec.Result[0].Type)
	})
}

type testSpecInterface interface {
	hello()
}

type testSpecInterfaceImpl struct{}

func (testSpecInterfaceImpl) hello() {}
