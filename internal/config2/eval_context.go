package config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/gocty"

	"github.com/hashicorp/waypoint/internal/config/funcs"
)

// EvalContext returns the common eval context to use for parsing all
// configurations. This should always be available for all config types.
//
// The pwd param is the directory to use as a working directory
// for determining things like relative paths. This should be considered
// the pwd over the actual process pwd.
func EvalContext(parent *hcl.EvalContext, pwd string) *hcl.EvalContext {
	// NewChild works even with parent == nil so this is valid
	result := parent.NewChild()

	// Start with our HCL stdlib
	result.Functions = funcs.Stdlib()

	// add functions to our context
	addFuncs := func(fs map[string]function.Function) {
		for k, v := range fs {
			result.Functions[k] = v
		}
	}

	// Add some of our functions
	addFuncs(funcs.VCSGitFuncs(pwd))
	addFuncs(funcs.Filesystem(pwd))
	addFuncs(funcs.Encoding())

	return result
}

// addPathValue adds the "path" variable to the context.
func addPathValue(ctx *hcl.EvalContext, v map[string]string) {
	value, err := gocty.ToCtyValue(v, cty.Map(cty.String))
	if err != nil {
		// map[string]string conversion should never fail
		panic(err)
	}

	if ctx.Variables == nil {
		ctx.Variables = map[string]cty.Value{}
	}

	ctx.Variables["path"] = value
}
