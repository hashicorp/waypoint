package component

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/waypoint/sdk/docs"
)

// Configurable can be optionally implemented by any compontent to
// accept user configuration.
type Configurable interface {
	// Config should return a pointer to an allocated configuration
	// structure. This structure will be written to directly with the
	// decoded configuration. If this returns nil, then it is as if
	// Configurable was not implemented.
	Config() (interface{}, error)
}

// Documented can be optionally implemented by any component to
// return documentation about the component.
type Documented interface {
	// Documentation() returns a completed docs.Documentation struct
	// describing the components configuration.
	Documentation() (*docs.Documentation, error)
}

// ConfigurableNotify is an optional interface that can be implemented
// by any component to receive a notification that the configuration
// was decoded.
type ConfigurableNotify interface {
	Configurable

	// ConfigSet is called with the value of the configuration after
	// decoding is complete successfully.
	ConfigSet(interface{}) error
}

// Configure configures c with the provided configuration.
//
// If c does not implement Configurable AND body is non-empty, then it is
// an error. If body is empty in that case, it is not an error.
func Configure(c interface{}, body hcl.Body, ctx *hcl.EvalContext) hcl.Diagnostics {
	if c, ok := c.(Configurable); ok {
		// Get the configuration value
		v, err := c.Config()
		if err != nil {
			return hcl.Diagnostics{
				&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  err.Error(),
					Detail:   "",
				},
			}
		}

		// If the configuration structure is nil then we behave as if the
		// component is not configurable.
		if v == nil {
			return nil
		}

		// Decode
		if diag := gohcl.DecodeBody(body, ctx, v); len(diag) > 0 {
			return diag
		}

		// If decoding worked and we have a notification implementation, then
		// notify with the value.
		if cn, ok := c.(ConfigurableNotify); ok {
			if err := cn.ConfigSet(v); err != nil {
				return hcl.Diagnostics{
					&hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  err.Error(),
						Detail:   "",
					},
				}
			}
		}

		return nil
	}

	// If c doesn't implement Configurable, then we parse the content with
	// an empty schema which will error if there are any fields since its
	// non-conformant to the schema.
	_, diag := body.Content(&hcl.BodySchema{})
	return diag
}

// Documentation returns the documentation for the given component.
//
// If c does not implement Documented, nil is returned.
func Documentation(c interface{}) (*docs.Documentation, error) {
	if d, ok := c.(Documented); ok {
		return d.Documentation()
	}

	return nil, nil
}
