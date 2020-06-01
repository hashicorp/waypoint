package plugin

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-argmapper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/component/mocks"
	"github.com/hashicorp/waypoint/sdk/proto"
)

func TestLogPlatformLogs(t *testing.T) {
	require := require.New(t)

	expected := []component.LogEvent{component.LogEvent{
		Partition: "foo",
		Timestamp: time.Now(),
		Message:   "hello",
	}}
	mockLV := &mocks.LogViewer{}
	mockLV.On("NextLogBatch", mock.Anything).Return(expected, nil)

	logsFunc := func(ctx context.Context, args *component.Source) component.LogViewer {
		return mockLV
	}

	mockLP := &mocks.LogPlatform{}
	mockLP.On("LogsFunc").Return(logsFunc)

	plugins := Plugins(WithComponents(mockLP), WithMappers(testDefaultMappers(t)...))
	client, server := plugin.TestPluginGRPCConn(t, plugins[1])
	defer client.Close()
	defer server.Stop()

	raw, err := client.Dispense("log_platform")
	require.NoError(err)
	lp := raw.(component.LogPlatform)
	f := lp.LogsFunc().(*argmapper.Func)
	require.NotNil(f)

	result := f.Call(
		argmapper.Typed(context.Background()),
		argmapper.Typed(&proto.Args_Source{App: "foo"}),
	)
	require.NoError(result.Err())

	raw = result.Out(0)
	require.NotNil(raw)

	value := raw.(component.LogViewer)
	require.NotNil(value)

	entries, err := value.NextLogBatch(context.Background())
	require.NoError(err)
	require.Len(entries, 1)
	require.Equal(expected[0].Message, entries[0].Message)
}
