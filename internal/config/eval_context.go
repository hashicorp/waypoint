package config

import (
	"github.com/hashicorp/hcl/v2"

	"github.com/hashicorp/waypoint/internal/config/funcs"
)

// EvalContext returns the common eval context to use for parsing all
// configurations. This should always be available for all config types.
//
// The pwd param is the directory to use as a working directory
// for determining things like relative paths. This should be considered
// the pwd over the actual process pwd.
func EvalContext(pwd string) *hcl.EvalContext {
	var result hcl.EvalContext

	// Start with our HCL stdlib
	result.Functions = funcs.Stdlib()

	// Add some of our functions

	return &result
}
