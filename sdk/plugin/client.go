package plugin

import (
	"os"

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

		// Syncing is very important to enable. This causes the stdout/err
		// of the plugin subprocesses to show up in our stdout/stderr.
		SyncStdout: os.Stdout,
		SyncStderr: os.Stderr,
	}
}
