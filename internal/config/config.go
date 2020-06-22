package config

import (
	"github.com/hashicorp/hcl/v2"
)

// Filename is the default filename for the Waypoint configuration.
const Filename = "waypoint.hcl"

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

	Build    *Build     `hcl:"build,block"`
	Platform *Operation `hcl:"deploy,block"`
	Release  *Operation `hcl:"release,block"`
}

// Server configures the remote server.
type Server struct {
	Address  string `hcl:"address,attr"`
	Insecure bool   `hcl:"insecure,optional"`

	// AddressInternal is a temporary config to work with local deployments
	// on platforms such as Docker for Mac. We need to discuss a more
	// long term approach to this.
	AddressInternal string `hcl:"address_internal,optional"`

	// Indicates that we need to present a token to connect to this server.
	// We don't allow the token to be hardcoded into the config though, we
	// always read that out of an env var later.
	RequireAuth bool `hcl:"require_auth,optional"`
}

// Build are the build settings.
type Build struct {
	Type string   `hcl:",label"`
	Body hcl.Body `hcl:",remain"`

	Labels   map[string]string `hcl:"labels,optional"`
	Registry *Operation        `hcl:"registry,block"`
}

func (b *Build) Operation() *Operation {
	if b == nil {
		return nil
	}

	return &Operation{
		Type:   b.Type,
		Body:   b.Body,
		Labels: b.Labels,
	}
}

func (b *Build) RegistryOperation() *Operation {
	if b == nil {
		return nil
	}

	return b.Registry
}
