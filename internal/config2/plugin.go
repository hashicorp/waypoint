package config

import (
	"github.com/hashicorp/waypoint-plugin-sdk/component"
)

// Plugin configures a plugin.
type Plugin struct {
	// Name of the plugin. This is expected to match the plugin binary
	// "waypoint-plugin-<name>" including casing.
	Name string `hcl:",label"`

	// Type is the type of plugin this is. This can be multiple.
	Type struct {
		Mapper   bool `hcl:"mapper,optional"`
		Builder  bool `hcl:"build,optional"`
		Registry bool `hcl:"registry,optional"`
		Platform bool `hcl:"deploy,optional"`
		Releaser bool `hcl:"release,optional"`
	} `hcl:"type,block"`

	// Checksum is the SHA256 checksum to validate this plugin.
	Checksum string `hcl:"checksum,optional"`
}

// Types returns the list of types that this plugin implements.
func (p *Plugin) Types() []component.Type {
	var result []component.Type
	for t, b := range p.typeMap() {
		if *b {
			result = append(result, t)
		}
	}

	return result
}

// markType marks that the given component type is supported by this plugin.
// This will panic if an unsupported plugin type is given.
func (p *Plugin) markType(typ component.Type) {
	m := p.typeMap()
	b, ok := m[typ]
	if !ok {
		panic("unknown type: " + typ.String())
	}

	*b = true
}

func (p *Plugin) typeMap() map[component.Type]*bool {
	return map[component.Type]*bool{
		component.MapperType:         &p.Type.Mapper,
		component.BuilderType:        &p.Type.Builder,
		component.RegistryType:       &p.Type.Registry,
		component.PlatformType:       &p.Type.Platform,
		component.ReleaseManagerType: &p.Type.Releaser,
	}
}
