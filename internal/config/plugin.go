// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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

// Plugins returns all the plugins defined by this configuration. This
// will include the implicitly defined plugins via `use` statements.
func (c *Config) Plugins() []*Plugin {
	result := make([]*Plugin, len(c.Plugin))
	copy(result, c.Plugin)
	known := map[string]*Plugin{}
	for _, p := range result {
		known[p.Name] = p
	}

	// Helper to also track all the plugins defined within scoped stanzas
	trackStages := func(typ component.Type, ws, label []*scopedStage) {
		for _, s := range ws {
			result = trackPlugin(result, known, s.Use, typ)
		}
		for _, s := range label {
			result = trackPlugin(result, known, s.Use, typ)
		}
	}

	// Collect all the plugins used by all the apps.
	for _, app := range c.hclConfig.Apps {
		// Get all the implied stage plugins: build, deploy, etc.
		if v := app.BuildRaw; v != nil {
			result = trackPlugin(result, known, v.Use, component.BuilderType)
			trackStages(component.BuilderType, v.WorkspaceScoped, v.LabelScoped)
			if v := v.Registry; v != nil {
				result = trackPlugin(result, known, v.Use, component.RegistryType)
				trackStages(component.RegistryType, v.WorkspaceScoped, v.LabelScoped)
			}
		}
		if v := app.DeployRaw; v != nil {
			result = trackPlugin(result, known, v.Use, component.PlatformType)
			trackStages(component.PlatformType, v.WorkspaceScoped, v.LabelScoped)
		}
		if v := app.ReleaseRaw; v != nil {
			result = trackPlugin(result, known, v.Use, component.ReleaseManagerType)
			trackStages(component.ReleaseManagerType, v.WorkspaceScoped, v.LabelScoped)
		}
	}

	return result
}

// trackPlugin adds the plugin implied by the use statements to result if
// it hasn't been seen before (known via the "known" variable). This will
// modify the values result and known in-place.
func trackPlugin(
	result []*Plugin,
	known map[string]*Plugin,
	use *Use,
	typ component.Type,
) []*Plugin {
	// If we don't have a plugin defined then just return, we'll validate
	// this during validation phases.
	if use == nil || use.Type == "" {
		return result
	}

	p, ok := known[use.Type]
	if !ok {
		p = &Plugin{Name: use.Type}
		result = append(result, p)
		known[use.Type] = p
	}

	// Mark the type this is used for
	p.markType(typ)

	// Mark the plugin for mappers. Anything that is in a `use` anywhere
	// is always loaded for mappers.
	p.markType(component.MapperType)

	return result
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
