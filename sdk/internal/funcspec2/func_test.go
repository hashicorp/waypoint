package funcspec

import (
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/mitchellh/go-argmapper"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/go-hclog"
	pb "github.com/hashicorp/waypoint/sdk/proto"
)

func init() {
	hclog.L().SetLevel(hclog.Trace)
}

func TestFunc(t *testing.T) {
	t.Run("single any result", func(t *testing.T) {
		require := require.New(t)

		spec, err := Spec(func(*pb.Empty) *pb.Empty { return &pb.Empty{} })
		require.NoError(err)
		require.NotNil(spec)

		f, err := Func(spec, func(args Args) (*any.Any, error) {
			require.Len(args, 1)
			require.NotNil(args[0])

			// At this point we'd normally RPC out.
			return ptypes.MarshalAny(&pb.Empty{})
		})
		require.NoError(err)

		msg, err := ptypes.MarshalAny(&pb.Empty{})
		require.NoError(err)

		result := f.Call(argmapper.TypedSubtype(msg, proto.MessageName(&pb.Empty{})))
		require.NoError(result.Err())
		require.Equal(reflect.Struct, reflect.ValueOf(result.Out(0)).Kind())
	})

	t.Run("single missing requirement", func(t *testing.T) {
		require := require.New(t)

		spec, err := Spec(func(*pb.Empty) *pb.Empty { return &pb.Empty{} })
		require.NoError(err)
		require.NotNil(spec)

		f, err := Func(spec, func(args Args) (*any.Any, error) {
			require.Len(args, 1)
			require.NotNil(args[0])

			// At this point we'd normally RPC out.
			return ptypes.MarshalAny(&pb.Empty{})
		})
		require.NoError(err)

		// Create an argument with the wrong type
		msg, err := ptypes.MarshalAny(&pb.FuncSpec{})
		require.NoError(err)
		result := f.Call(argmapper.TypedSubtype(msg, proto.MessageName(&pb.FuncSpec{})))

		// We should have an error
		require.Error(result.Err())
		require.Contains(result.Err().Error(), "argument cannot")
	})

	t.Run("match callback output if no results", func(t *testing.T) {
		require := require.New(t)

		spec, err := Spec(func(*pb.Empty) *pb.Empty { return &pb.Empty{} })
		require.NoError(err)
		require.NotNil(spec)

		// No results
		spec.Result = nil

		// Build our func to return a primitive
		f, err := Func(spec, func(args Args) int {
			require.Len(args, 1)
			require.NotNil(args[0])
			return 42
		})
		require.NoError(err)

		// Call the function with the proto type we expect
		msg, err := ptypes.MarshalAny(&pb.Empty{})
		require.NoError(err)
		result := f.Call(argmapper.TypedSubtype(msg, proto.MessageName(&pb.Empty{})))

		// Should succeed and give us our primitive
		require.NoError(result.Err())
		require.Equal(42, result.Out(0))
	})
}
