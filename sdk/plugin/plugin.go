package plugin

import (
	"fmt"
	"reflect"

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

// Plugins returns the list of available plugins and initializes them with
// the given components. This will panic if an invalid component is given.
func Plugins(components ...interface{}) map[int]plugin.PluginSet {
	// Build our plugin types
	result := map[int]plugin.PluginSet{
		1: plugin.PluginSet{
			"builder": &BuilderPlugin{},
		},
	}

	// Set the components
	for _, c := range components {
		if err := setComponent(result, c); err != nil {
			panic(err)
		}
	}

	return result
}

// setComponent sets the component on any public assignable fields
// of a plugin that match the type of c. An error is returned if
// c can't be assigned to any plugin type.
//
// preconditions:
//   - plugins in m are pointers to structs
func setComponent(m map[int]plugin.PluginSet, c interface{}) error {
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

	return fmt.Errorf("no plugin available for component of type %T", c)
}
