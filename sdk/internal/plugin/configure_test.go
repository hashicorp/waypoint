package plugin

import (
	"testing"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/stretchr/testify/require"

	"github.com/mitchellh/devflow/sdk/component"
	"github.com/mitchellh/devflow/sdk/component/mocks"
)

// testConfigurable is a reusable helper that tests that a component implements
// the Configurable interfaces correctly.
func testConfigurable(
	t *testing.T,
	typ string, // plugin type
	impl interface{}, // full implementation
	mockC *mocks.Configurable,
) {
	require := require.New(t)

	var config struct {
		Name string `hcl:"name"`
	}
	mockC.On("Config").Return(&config, nil)

	plugins := Plugins(WithComponents(impl), WithMappers(testDefaultMappers(t)...))
	client, server := plugin.TestPluginGRPCConn(t, plugins[1])
	defer client.Close()
	defer server.Stop()

	raw, err := client.Dispense(typ)
	require.NoError(err)

	src := `name = "foo"`
	f, diag := hclparse.NewParser().ParseHCL([]byte(src), "test.hcl")
	require.False(diag.HasErrors())

	diag = component.Configure(raw, f.Body, nil)
	require.False(diag.HasErrors())
	require.Equal("foo", config.Name)
}
