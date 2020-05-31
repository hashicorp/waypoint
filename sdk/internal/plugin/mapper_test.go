package plugin

import (
	"context"
	"testing"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-argmapper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/sdk/internal/funcspec"
	"github.com/hashicorp/waypoint/sdk/internal/testproto"
	pb "github.com/hashicorp/waypoint/sdk/proto"
)

func TestMapperClient(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	mA, err := argmapper.NewFunc(func(a *testproto.A) *testproto.B {
		return &testproto.B{Value: a.Value + 1}
	})
	require.NoError(err)

	plugins := Plugins(WithMappers(append(testDefaultMappers(t), mA)...))
	client, server := plugin.TestPluginGRPCConn(t, plugins[1])
	defer client.Close()
	defer server.Stop()

	raw, err := client.Dispense("mapper")
	require.NoError(err)
	mapper := raw.(*MapperClient)

	mappers, err := mapper.Mappers()
	require.NoError(err)
	require.NotEmpty(mappers)

	targetSpec := &pb.FuncSpec{
		Args: []*pb.FuncSpec_Value{
			&pb.FuncSpec_Value{
				Type: "testproto.B",
			},
		},

		Result: []*pb.FuncSpec_Value{
			&pb.FuncSpec_Value{
				Type: "testproto.Data",
			},
		},
	}

	called := false
	target := funcspec.Func(targetSpec, func(args funcspec.Args) (interface{}, error) {
		cb := func(v *testproto.B) *testproto.Data {
			called = true
			assert.Equal(int32(2), v.Value)
			return &testproto.Data{}
		}

		return callDynamicFunc2(cb, args,
			argmapper.Typed(context.Background()),
			argmapper.ConverterFunc(mappers...),
		)
	})

	result := target.Call(
		argmapper.Typed(context.Background()),
		argmapper.Typed(&testproto.A{Value: 1}),
		argmapper.ConverterFunc(mappers...),
	)
	require.NoError(result.Err())
	require.True(called)
}
