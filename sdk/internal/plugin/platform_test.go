package plugin

import (
	"testing"

	"github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/require"

	"github.com/mitchellh/devflow/sdk/component"
	"github.com/mitchellh/devflow/sdk/component/mocks"
)

func TestPlatform_logsPlatform(t *testing.T) {
	t.Run("implements", func(t *testing.T) {
		require := require.New(t)

		mockV := &mockPlatformLog{}

		plugins := Plugins(WithComponents(mockV), WithMappers(testDefaultMappers(t)...))
		client, server := plugin.TestPluginGRPCConn(t, plugins[1])
		defer client.Close()
		defer server.Stop()

		raw, err := client.Dispense("platform")
		require.NoError(err)
		require.Implements((*component.Platform)(nil), raw)
		require.Implements((*component.LogPlatform)(nil), raw)
	})

	t.Run("doesn't implement", func(t *testing.T) {
		require := require.New(t)

		mockV := &mocks.Platform{}

		plugins := Plugins(WithComponents(mockV), WithMappers(testDefaultMappers(t)...))
		client, server := plugin.TestPluginGRPCConn(t, plugins[1])
		defer client.Close()
		defer server.Stop()

		raw, err := client.Dispense("platform")
		require.NoError(err)
		require.Implements((*component.Platform)(nil), raw)

		_, ok := raw.(component.LogPlatform)
		require.False(ok, "does not implement LogPlatform")
	})
}

func TestPlatformDynamicFunc(t *testing.T) {
	testDynamicFunc(t, "platform", &mocks.Platform{}, func(v, f interface{}) {
		v.(*mocks.Platform).On("DeployFunc").Return(f)
	}, func(raw interface{}) interface{} {
		return raw.(component.Platform).DeployFunc()
	})
}

func TestPlatformConfig(t *testing.T) {
	mockV := &mockPlatformConfigurable{}
	testConfigurable(t, "platform", mockV, &mockV.Configurable)
}

type mockPlatformConfigurable struct {
	mocks.Platform
	mocks.Configurable
}

type mockPlatformLog struct {
	mocks.Platform
	mocks.LogPlatform
}
