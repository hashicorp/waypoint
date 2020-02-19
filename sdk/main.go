package sdk

import (
	"github.com/hashicorp/go-plugin"

	sdkplugin "github.com/mitchellh/devflow/sdk/plugin"
)

//go:generate sh -c "protoc -I proto/ proto/*.proto --go_out=plugins=grpc:proto/"

// Main is the primary entrypoint for plugins serving components. This
// function never returns; it blocks until the program is exited. This should
// be called immediately in main() in your plugin binaries, no prior setup
// should be done.
func Main(opts ...Option) {
	// Build config
	var c config
	for _, opt := range opts {
		opt(&c)
	}

	// Serve
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig:  sdkplugin.Handshake,
		VersionedPlugins: sdkplugin.Plugins(c.Components...),
		GRPCServer:       plugin.DefaultGRPCServer,
	})
}

// config is the configuration for Main. This can only be modified using
// Option implementations.
type config struct {
	// Components is the list of components to serve from the plugin.
	Components []interface{}
}

// Option modifies config. Zero or more can be passed to Main.
type Option func(*config)

// WithComponents specifies a list of components to serve from the plugin
// binary. This will append to the list of components to serve. You can
// currently only serve at most one of each type of plugin.
func WithComponents(cs ...interface{}) Option {
	return func(c *config) { c.Components = append(c.Components, cs...) }
}
