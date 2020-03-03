package sdk

import (
	"github.com/hashicorp/go-plugin"

	"github.com/mitchellh/devflow/sdk/internal-shared/mapper"
	"github.com/mitchellh/devflow/sdk/internal-shared/protomappers"
	sdkplugin "github.com/mitchellh/devflow/sdk/internal/plugin"
)

//go:generate sh -c "protoc -I`go list -m -f \"{{.Dir}}\" github.com/mitchellh/protostructure` -I proto/ proto/*.proto --go_out=plugins=grpc:proto/"

// Main is the primary entrypoint for plugins serving components. This
// function never returns; it blocks until the program is exited. This should
// be called immediately in main() in your plugin binaries, no prior setup
// should be done.
func Main(opts ...Option) {
	var c config

	// Default our mappers
	c.Mappers = append(c.Mappers, protomappers.All...)

	// Build config
	for _, opt := range opts {
		opt(&c)
	}

	// Build up our mappers
	var mappers []*mapper.Func
	for _, raw := range c.Mappers {
		// If the mapper is already a mapper.Func, then we let that through as-is
		m, ok := raw.(*mapper.Func)
		if !ok {
			var err error
			m, err = mapper.NewFunc(raw)
			if err != nil {
				panic(err)
			}
		}

		mappers = append(mappers, m)
	}

	// Serve
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: sdkplugin.Handshake,
		VersionedPlugins: sdkplugin.Plugins(
			sdkplugin.WithComponents(c.Components...),
			sdkplugin.WithMappers(mappers...),
		),
		GRPCServer: plugin.DefaultGRPCServer,
	})
}

// config is the configuration for Main. This can only be modified using
// Option implementations.
type config struct {
	// Components is the list of components to serve from the plugin.
	Components []interface{}

	// Mappers is the list of mapper functions.
	Mappers []interface{}
}

// Option modifies config. Zero or more can be passed to Main.
type Option func(*config)

// WithComponents specifies a list of components to serve from the plugin
// binary. This will append to the list of components to serve. You can
// currently only serve at most one of each type of plugin.
func WithComponents(cs ...interface{}) Option {
	return func(c *config) { c.Components = append(c.Components, cs...) }
}

// WithMappers specifies a list of mappers to apply to the plugin.
//
// Mappers are functions that take zero or more arguments and return
// one result (optionally with an error). These can be used to convert argument
// types as needed for your plugin functions. For example, you can convert a
// proto type to a richer Go struct.
//
// Mappers must take zero or more arguments and return exactly one or two
// values where the second return type must be an error. Example:
//
//   func() *Value
//   func() (*Value, error)
//   -- the above with any arguments
//
// This will append the mappers to the list of available mappers. A set of
// default mappers is always included to convert from SDK proto types to
// richer Go structs.
func WithMappers(ms ...interface{}) Option {
	return func(c *config) { c.Mappers = append(c.Mappers, ms...) }
}
