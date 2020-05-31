package plugin

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-argmapper"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// Not secret, just to avoid plugins being launched manually. The
	// cookie value is a random SHA256 via /dev/urandom. This cookie value
	// must NEVER be changed or plugins will stop working.
	MagicCookieKey:   "WAYPOINT_PLUGIN",
	MagicCookieValue: "be6c1928786a4df0222c13eef44ac846da2c0d461d99addc93f804601c6b7205",
}

// Plugins returns the list of available plugins and initializes them with
// the given components. This will panic if an invalid component is given.
func Plugins(opts ...Option) map[int]plugin.PluginSet {
	var c pluginConfig
	for _, opt := range opts {
		opt(&c)
	}

	// If we have no logger, we use the default
	if c.Logger == nil {
		c.Logger = hclog.L()
	}

	// Build our plugin types
	result := map[int]plugin.PluginSet{
		1: plugin.PluginSet{
			"mapper":         &MapperPlugin{},
			"builder":        &BuilderPlugin{},
			"platform":       &PlatformPlugin{},
			"log_platform":   &LogPlatformPlugin{},
			"registry":       &RegistryPlugin{},
			"releasemanager": &ReleaseManagerPlugin{},
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
	// Set the logger
	if err := setFieldValue(result, c.Logger); err != nil {
		panic(err)
	}

	return result
}

// pluginConfig is used to configure Plugins via Option calls.
type pluginConfig struct {
	Components []interface{}
	Mappers    []*argmapper.Func
	Logger     hclog.Logger
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
func WithMappers(ms ...*argmapper.Func) Option {
	return func(c *pluginConfig) {
		c.Mappers = append(c.Mappers, ms...)
	}
}

// WithLogger sets the logger for the plugins.
func WithLogger(log hclog.Logger) Option {
	return func(c *pluginConfig) { c.Logger = log }
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
	once := false
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
				// to it then we set the value directly. We continue setting
				// values because some values we set are available in multiple
				// plugins (loggers for example).
				if f.IsValid() && ct.AssignableTo(f.Type()) {
					f.Set(cv)
					once = true
				}
			}
		}
	}

	if !once {
		return fmt.Errorf("no plugin available for setting field of type %T", c)
	}

	return nil
}
