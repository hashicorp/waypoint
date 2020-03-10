package funcspec

import (
	"testing"

	"github.com/mitchellh/devflow/sdk/internal-shared/mapper"
	"github.com/stretchr/testify/require"

	pb "github.com/mitchellh/devflow/sdk/proto"
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
}

func mustFunc(t *testing.T, f interface{}) *mapper.Func {
	t.Helper()

	result, err := mapper.NewFunc(f)
	require.NoError(t, err)
	return result
}
