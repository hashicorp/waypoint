package plugin

import (
	"github.com/hashicorp/go-plugin"
)

//go:generate sh -c "protoc -I proto/ proto/*.proto --go_out=plugins=grpc:proto/"

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// Not secret, just to avoid plugins being launched manually. The
	// cookie value is a random SHA256 via /dev/urandom
	MagicCookieKey:   "DEVFLOW_PLUGIN",
	MagicCookieValue: "be6c1928786a4df0222c13eef44ac846da2c0d461d99addc93f804601c6b7205",
}

// Plugins is the list of available versions.
var Plugins = map[int]plugin.PluginSet{
	1: plugin.PluginSet{
		"builder": &BuilderPlugin{},
	},
}
