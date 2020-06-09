package config

import (
	"github.com/hashicorp/hcl/v2"
)

// Config is the configuration structure.
type Config struct {
	Server  *Server           `hcl:"server,block"`
	Project string            `hcl:"project,attr"`
	Apps    []*App            `hcl:"app,block"`
	Labels  map[string]string `hcl:"labels,optional"`
}

// App represents a single application.
type App struct {
	Name   string            `hcl:",label"`
	Path   string            `hcl:"path,optional"`
	Labels map[string]string `hcl:"labels,optional"`

	Build    *Component `hcl:"build,block"`
	Registry *Component `hcl:"registry,block"`
	Platform *Component `hcl:"deploy,block"`
	Release  *Component `hcl:"release,block"`
}

// Component is an internal name used to describe a single component such as
// build, deploy, releaser, etc.
type Component struct {
	Type   string            `hcl:",label"`
	Body   hcl.Body          `hcl:",remain"`
	Labels map[string]string `hcl:"labels,optional"`
}

// Server configures the remote server.
type Server struct {
	Address  string `hcl:"address,attr"`
	Insecure bool   `hcl:"insecure,optional"`

	// AddressInternal is a temporary config to work with local deployments
	// on platforms such as Docker for Mac. We need to discuss a more
	// long term approach to this.
	AddressInternal string `hcl:"address_internal,optional"`
}
