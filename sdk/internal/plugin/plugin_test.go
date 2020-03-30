package plugin

import (
	"context"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mitchellh/devflow/sdk/component"
	"github.com/mitchellh/devflow/sdk/component/mocks"
	"github.com/mitchellh/devflow/sdk/history"
	historymocks "github.com/mitchellh/devflow/sdk/history/mocks"
	"github.com/mitchellh/devflow/sdk/internal-shared/mapper"
	"github.com/mitchellh/devflow/sdk/internal-shared/protomappers"
	"github.com/mitchellh/devflow/sdk/internal/plugincomponent"
	"github.com/mitchellh/devflow/sdk/internal/testproto"
	pb "github.com/mitchellh/devflow/sdk/proto"
)

func init() {
	// Set our default log level lower for tests
	hclog.L().SetLevel(hclog.Debug)
}

func TestPlugins(t *testing.T) {
	require := require.New(t)

	mock := &mocks.Builder{}
	plugins := Plugins(WithComponents(mock))
	bp := plugins[1]["builder"].(*BuilderPlugin)
	require.Equal(bp.Impl, mock)
}

func testDefaultMappers(t *testing.T) []*mapper.Func {
	var mappers []*mapper.Func
	for _, raw := range protomappers.All {
		f, err := mapper.NewFunc(raw)
		require.NoError(t, err)
		mappers = append(mappers, f)
	}

	return mappers
}

// testDynamicFunc ensures that the dynamic function capabilities work
// properly. This should be called for each individual dynamic function
// the component exposes.
func testDynamicFunc(
	t *testing.T,
	typ string,
	value interface{},
	setFunc func(interface{}, interface{}), // set the function on your mock
	getFunc func(interface{}) interface{}, // get the function
) {
	require := require.New(t)
	assert := assert.New(t)

	// Our callback that we verify. We specify a LOT of args here because
	// we want to verify that each one will work properly. This is the core
	// of this test.
	called := false
	setFunc(value, func(
		ctx context.Context,
		args *component.Source,
		historyClient history.Client,
	) *testproto.Data {
		called = true
		assert.NotNil(ctx)
		assert.Equal("foo", args.App)

		// Test history client
		assert.NotNil(historyClient)
		resp, err := historyClient.Deployments(ctx, &history.Lookup{
			Type: (*[]*plugincomponent.Artifact)(nil),
		})
		if assert.NoError(err) {
			assert.NotNil(resp)
		}

		return &testproto.Data{Value: "hello"}
	})

	// Get the mappers
	mappers := testDefaultMappers(t)

	// Init the plugin server
	plugins := Plugins(WithComponents(value), WithMappers(mappers...))
	client, server := plugin.TestPluginGRPCConn(t, plugins[1])
	defer client.Close()
	defer server.Stop()

	// Dispense the plugin
	raw, err := client.Dispense(typ)
	require.NoError(err)
	implFunc := getFunc(raw).(*mapper.Func)

	historyMock := &historymocks.Client{}
	historyMock.On("Deployments", mock.Anything, (*history.Lookup)(nil)).Return([]component.Deployment{}, nil)

	// Call our function by building a chain. We use the chain so we
	// have access to the same level of mappers that a default plugin
	// would normally have.
	chain, err := implFunc.Chain(mappers,
		context.Background(),
		hclog.L(),

		&pb.Args_Source{App: "foo"},
		historyMock,
	)
	require.NoError(err)
	require.NotNil(chain)

	// Call our function chain
	raw, err = chain.Call()
	require.NoError(err)
	require.NotNil(raw)

	require.True(called)
}
