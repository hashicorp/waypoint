package config

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
