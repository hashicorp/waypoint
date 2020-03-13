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
