package plugin

import (
	"github.com/hashicorp/go-plugin"

	internalplugin "github.com/mitchellh/devflow/sdk/internal/plugin"
)

// ClientConfig returns the base client config to use when connecting
// to a plugin. This sets the handshake config, protocols, etc. Manually
// override any values you want to set.
func ClientConfig() *plugin.ClientConfig {
	return &plugin.ClientConfig{
		HandshakeConfig:  internalplugin.Handshake,
		VersionedPlugins: internalplugin.Plugins(),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	}
}
