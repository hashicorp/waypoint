package config

import (
	"github.com/hashicorp/hcl/v2"
)

// Config is the configuration structure.
type Config struct {
	App []*App `hcl:"app,block"`
}

// App represents a single application.
type App struct {
	Name string `hcl:",label"`
	Path string `hcl:"path,optional"`

	Build    *Component `hcl:"build,block"`
	Registry *Component `hcl:"registry,block"`
	Deploy   *Component `hcl:"deploy,block"`
}

// Component is an internal name used to describe a single component such as
// build, deploy, releaser, etc.
type Component struct {
	Type string   `hcl:",label"`
	Body hcl.Body `hcl:",remain"`
}
