package config

import (
	"github.com/hashicorp/hcl/v2"
)

type hclStage struct {
	Body hcl.Body `hcl:",remain"`
}

type hclBuild struct {
	Registry *hclStage `hcl:"registry,block"`
	Body     hcl.Body  `hcl:",remain"`
}

// Build are the build settings.
type Build struct {
	Labels map[string]string `hcl:"labels,optional"`
	Hooks  []*Hook           `hcl:"hook,block"`
	Use    *Use              `hcl:"use,block"`
}

// Registry are the registry settings.
type Registry struct {
	Labels map[string]string `hcl:"labels,optional"`
	Hooks  []*Hook           `hcl:"hook,block"`
	Use    *Use              `hcl:"use,block"`
}

// Deploy are the deploy settings.
type Deploy struct {
	Labels map[string]string `hcl:"labels,optional"`
	Hooks  []*Hook           `hcl:"hook,block"`
	Use    *Use              `hcl:"use,block"`
}

// Release are the release settings.
type Release struct {
	Labels map[string]string `hcl:"labels,optional"`
	Hooks  []*Hook           `hcl:"hook,block"`
	Use    *Use              `hcl:"use,block"`
}

// Use is something in the Waypoint configuration that is executed
// using some underlying plugin. This is a general shared structure that is
// used by internal/core to initialize all the proper plugins.
type Use struct {
	Type string   `hcl:",label"`
	Body hcl.Body `hcl:",remain"`
}

// Hook is the configuration for a hook that runs at specified times.
type Hook struct {
	When      string   `hcl:"when,attr"`
	Command   []string `hcl:"command,attr"`
	OnFailure string   `hcl:"on_failure,optional"`
}

func (h *Hook) ContinueOnFailure() bool {
	return h.OnFailure == "continue"
}
