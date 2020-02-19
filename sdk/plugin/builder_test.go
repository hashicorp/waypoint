package plugin

import (
	"testing"

	"github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/require"

	"github.com/mitchellh/devflow/internal/component"
	"github.com/mitchellh/devflow/internal/mapper"
)

func TestBuilderBuild(t *testing.T) {
	require := require.New(t)

	client, server := plugin.TestPluginGRPCConn(t, Plugins[1])
	defer client.Close()
	defer server.Stop()

	raw, err := client.Dispense("builder")
	require.NoError(err)
	builder := raw.(component.Builder)
	f := builder.BuildFunc().(*mapper.Func)
	require.NotNil(f)
}
