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
	addFuncs(funcs.Filesystem())
	addFuncs(funcs.Encoding())
	addFuncs(funcs.Datetime())

	return result
}

// appendContext makes child a child of parent and returns the new context.
// If child is nil this returns parent.
func appendContext(parent, child *hcl.EvalContext) *hcl.EvalContext {
	if child == nil {
		return parent
	}

	// We need to get the full tree of contexts since we need to go
	// parent => child traversal but HCL only supports child => parent.
	var tree []*hcl.EvalContext
	for current := child; current != nil; current = current.Parent() {
		tree = append(tree, current)
	}

	// Go backward through the tree (parent => child) order to ensure
	// that we merge all the context trees properly.
	for i := len(tree) - 1; i >= 0; i-- {
		current := tree[i]
		parent = parent.NewChild()
		parent.Variables = current.Variables
		parent.Functions = current.Functions
	}

	return parent
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

// finalizeContext should be called whenever an HCL context is being used
// as the final call.
func finalizeContext(ctx *hcl.EvalContext) *hcl.EvalContext {
	ctx = ctx.NewChild()
	ctx.Functions = funcs.MakeTemplateFuncs(ctx)
	return ctx
}

// hclContextContainer is an interface that config structs that have an HCL
// context may implement. We use this for certain things such as mapoperation()
// to set the proper context.
type hclContextContainer interface {
	hclContext() *hcl.EvalContext
}
