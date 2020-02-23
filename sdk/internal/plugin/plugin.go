package plugin

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/go-plugin"
	"github.com/mitchellh/devflow/sdk/pkg/mapper"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// Not secret, just to avoid plugins being launched manually. The
	// cookie value is a random SHA256 via /dev/urandom
	MagicCookieKey:   "DEVFLOW_PLUGIN",
	MagicCookieValue: "be6c1928786a4df0222c13eef44ac846da2c0d461d99addc93f804601c6b7205",
}

// Plugins returns the list of available plugins and initializes them with
// the given components. This will panic if an invalid component is given.
func Plugins(opts ...Option) map[int]plugin.PluginSet {
	var c pluginConfig
	for _, opt := range opts {
		opt(&c)
	}

	// Build our plugin types
	result := map[int]plugin.PluginSet{
		1: plugin.PluginSet{
			"builder": &BuilderPlugin{},
		},
	}

	// Set the various field values
	for _, c := range c.Components {
		if err := setFieldValue(result, c); err != nil {
			panic(err)
		}
	}

	// Set the mappers
	if err := setFieldValue(result, c.Mappers); err != nil {
		panic(err)
	}

	return result
}

// pluginConfig is used to configure Plugins via Option calls.
type pluginConfig struct {
	Components []interface{}
	Mappers    []*mapper.Func
}

// Option configures Plugins
type Option func(*pluginConfig)

// WithComponents sets the components to configure for the plugins.
// This will append to the components.
func WithComponents(cs ...interface{}) Option {
	return func(c *pluginConfig) { c.Components = append(c.Components, cs...) }
}

// WithMappers sets the mappers to configure for the plugins. This will
// append to the existing mappers.
func WithMappers(ms ...*mapper.Func) Option {
	return func(c *pluginConfig) { c.Mappers = append(c.Mappers, ms...) }
}

// setFieldValue sets the given value c on any exported field of an available
// plugin that matches the type of c. An error is returned if c can't be
// assigned to ANY plugin type.
//
// preconditions:
//   - plugins in m are pointers to structs
func setFieldValue(m map[int]plugin.PluginSet, c interface{}) error {
	cv := reflect.ValueOf(c)
	ct := cv.Type()

	// Go through each pluginset
	for _, set := range m {
		// Go through each plugin
		for _, p := range set {
			// Get the value, dereferencing the pointer. We expect
			// the value to be &SomeStruct{} so we must deref once.
			v := reflect.ValueOf(p).Elem()

			// Go through all the fields
			for i := 0; i < v.NumField(); i++ {
				f := v.Field(i)

				// If the field is valid and our component can be assigned
				// to it then we set the value directly. We then return since
				// we expect a plugin to only represent one kind of plugin.
				if f.IsValid() && ct.AssignableTo(f.Type()) {
					f.Set(cv)
					return nil
				}
			}
		}
	}

	return fmt.Errorf("no plugin available for setting field of type %T", c)
}
