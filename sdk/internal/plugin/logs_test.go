package plugin

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mitchellh/devflow/sdk/component"
	"github.com/mitchellh/devflow/sdk/component/mocks"
	"github.com/mitchellh/devflow/sdk/internal-shared/mapper"
	"github.com/mitchellh/devflow/sdk/proto"
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
	f := lp.LogsFunc().(*mapper.Func)
	require.NotNil(f)

	raw, err = f.Call(context.Background(), &proto.Args_Source{App: "foo"})
	require.NoError(err)
	require.NotNil(raw)

	result := raw.(component.LogViewer)
	require.NotNil(result)

	entries, err := result.NextLogBatch(context.Background())
	require.NoError(err)
	require.Len(entries, 1)
	require.Equal(expected[0].Message, entries[0].Message)
}
