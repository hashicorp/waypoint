package pluginclient

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/mitchellh/devflow/sdk/internal-shared/mapper"
	internalplugin "github.com/mitchellh/devflow/sdk/internal/plugin"
)

// ClientConfig returns the base client config to use when connecting
// to a plugin. This sets the handshake config, protocols, etc. Manually
// override any values you want to set.
func ClientConfig(log hclog.Logger) *plugin.ClientConfig {
	return &plugin.ClientConfig{
		HandshakeConfig:  internalplugin.Handshake,
		VersionedPlugins: internalplugin.Plugins(internalplugin.WithLogger(log)),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},

		// Syncing is very important to enable. This causes the stdout/err
		// of the plugin subprocesses to show up in our stdout/stderr.
		SyncStdout: os.Stdout,
		SyncStderr: os.Stderr,

		// We always set managed to true just in case we don't properly
		// call Kill so that CleanupClients gets it. If we do properly call
		// Kill, then it is a no-op to call it again so this is safe.
		Managed: true,

		// This is super important. There appears to be a bug with AutoMTLS
		// when using GRPCBroker and listening from the _client_ side. The
		// TLS fails to negotiate. For now we just disable this but we should
		// look into fixing that later.
		AutoMTLS: false,
	}
}

// Mappers returns the mappers supported by the plugin.
func Mappers(c *plugin.Client) ([]*mapper.Func, error) {
	rpcClient, err := c.Client()
	if err != nil {
		return nil, err
	}

	v, err := rpcClient.Dispense("mapper")
	if err != nil {
		return nil, err
	}

	client, ok := v.(*internalplugin.MapperClient)
	if !ok {
		return nil, fmt.Errorf("mapper service was unexpected type: %T", v)
	}

	return client.Mappers()
}
