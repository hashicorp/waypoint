package config

import (
	"github.com/hashicorp/hcl/v2"
)

// Use is something in the Waypoint configuration that is executed
// using some underlying plugin. This is a general shared structure that is
// used by internal/core to initialize all the proper plugins.
type Use struct {
	Type string   `hcl:",label"`
	Body hcl.Body `hcl:",remain"`
}

// hclUse as a minimal structure to extract the Use field of other configs.
type hclUse struct {
	Use *Use `hcl:"use,block"`
}

func (s *hclBuild) useContainer() hcl.Body {
	return s.Body
}

func (s *hclStage) useContainer() hcl.Body {
	return s.Body
}
