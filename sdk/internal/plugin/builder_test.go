package plugin

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mitchellh/devflow/sdk/component"
	"github.com/mitchellh/devflow/sdk/component/mocks"
	"github.com/mitchellh/devflow/sdk/pkg/mapper"
	"github.com/mitchellh/devflow/sdk/proto"
)

func TestBuilderBuild(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	called := false
	buildFunc := func(ctx context.Context, args *proto.Args_Source) *proto.Empty {
		called = true
		assert.NotNil(ctx)
		assert.Equal("foo", args.App)
		return &proto.Empty{}
	}

	mockB := &mocks.Builder{}
	mockB.On("BuildFunc").Return(buildFunc)

	plugins := Plugins(mockB)
	client, server := plugin.TestPluginGRPCConn(t, plugins[1])
	defer client.Close()
	defer server.Stop()

	raw, err := client.Dispense("builder")
	require.NoError(err)
	builder := raw.(component.Builder)
	f := builder.BuildFunc().(*mapper.Func)
	require.NotNil(f)

	raw, err = f.Call(context.Background(), &proto.Args_Source{App: "foo"})
	require.NoError(err)
	require.NotNil(raw)

	result := raw.(*any.Any)
	name, err := ptypes.AnyMessageName(result)
	require.NoError(err)
	require.Equal("proto.Empty", name)

	require.True(called)
}
