package plugin

import (
	"context"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/mitchellh/go-argmapper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/component/mocks"
	"github.com/hashicorp/waypoint/sdk/history"
	historymocks "github.com/hashicorp/waypoint/sdk/history/mocks"
	"github.com/hashicorp/waypoint/sdk/internal-shared/protomappers"
	"github.com/hashicorp/waypoint/sdk/internal/testproto"
	pb "github.com/hashicorp/waypoint/sdk/proto"
)

func init() {
	// Set our default log level lower for tests
	hclog.L().SetLevel(hclog.Trace)
}

func TestPlugins(t *testing.T) {
	require := require.New(t)

	mock := &mocks.Builder{}
	plugins := Plugins(WithComponents(mock))
	bp := plugins[1]["builder"].(*BuilderPlugin)
	require.Equal(bp.Impl, mock)
}

func testDefaultMappers(t *testing.T) []*argmapper.Func {
	var mappers []*argmapper.Func
	for _, raw := range protomappers.All {
		f, err := argmapper.NewFunc(raw)
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
		// TODO(mitchellh): uncomment
		//historyClient history.Client,
	) *testproto.Data {
		called = true
		assert.NotNil(ctx)
		assert.Equal("foo", args.App)

		// Test history client
		/*
			assert.NotNil(historyClient)
			_, err := historyClient.Deployments(ctx, nil)
			assert.NoError(err)
		*/

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
	implFunc := getFunc(raw).(*argmapper.Func)

	historyMock := &historymocks.Client{}
	historyMock.On("Deployments", mock.Anything, &history.Lookup{}).Return([]component.Deployment{}, nil)

	// Call our function by building a chain. We use the chain so we
	// have access to the same level of mappers that a default plugin
	// would normally have.
	result := implFunc.Call(
		argmapper.ConverterFunc(mappers...),

		argmapper.Typed(context.Background()),
		argmapper.Typed(hclog.L()),

		argmapper.Typed(&pb.Args_Source{App: "foo"}),
		argmapper.Typed(historyMock),
	)
	require.NoError(result.Err())

	// We only require a result if the function type expects us to return
	// a result. Otherwise, we just expect nil because it is error-only.
	if result.Len() > 0 {
		require.NotNil(result.Out(0))
	}

	require.True(called)
}
