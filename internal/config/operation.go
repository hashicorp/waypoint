package config

import (
	"github.com/hashicorp/hcl/v2"
)

// Operation is something in the Waypoint configuraiton that is executed
// using some underlying plugin. This is a general shared structure that is
// used by internal/core to initialize all the proper plugins.
type Operation struct {
	Type   string            `hcl:",label"`
	Body   hcl.Body          `hcl:",remain"`
	Labels map[string]string `hcl:"labels,optional"`
}
