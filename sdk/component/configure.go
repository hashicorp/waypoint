package component

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
)

// Configurable can be optionally implemented by any compontent to
// accept user configuration.
type Configurable interface {
	// Config should return a pointer to an allocated configuration
	// structure. This structure will be written to directly with the
	// decoded configuration.
	Config() interface{}
}

// Configure configures c with the provided configuration.
//
// If c does not implement Configurable AND body is non-empty, then it is
// an error. If body is empty in that case, it is not an error.
func Configure(c interface{}, body hcl.Body, ctx *hcl.EvalContext) hcl.Diagnostics {
	if c, ok := c.(Configurable); ok {
		return gohcl.DecodeBody(body, ctx, c.Config())
	}

	// If c doesn't implement Configurable, then we parse the content with
	// an empty schema which will error if there are any fields since its
	// non-conformant to the schema.
	_, diag := body.Content(&hcl.BodySchema{})
	return diag
}
