package config

import (
	"github.com/hashicorp/waypoint/internal/pkg/defaults"
)

// Default sets the default values where values are unset on this config.
// This will modify the config in place.
func (c *Config) Default() error {
	return defaults.Set(c)
}

func boolPtr(v bool) *bool {
	return &v
}
