// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"

	"github.com/hashicorp/waypoint/internal/config/variables"
	"github.com/hashicorp/waypoint/pkg/config/funcs"
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

	funcs.AddStandardFunctions(result, pwd)

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

// defineContextVarsIfNeeded will set an empty map to ctx.Variables if needed
func defineContextVarsIfNeeded(ctx *hcl.EvalContext) {
	if ctx.Variables == nil {
		ctx.Variables = map[string]cty.Value{}
	}
}

// addWorkspaceValue adds the workspace values to the context. This
// adds the `workspace` map and currently only supports the `workspace.name`
// value.
func addWorkspaceValue(ctx *hcl.EvalContext, v string) {
	addMapVariable(ctx, "workspace", map[string]string{
		"name": v,
	})
}

// addPathValue adds the "path" variable to the context.
func addPathValue(ctx *hcl.EvalContext, v map[string]string) {
	addMapVariable(ctx, "path", v)
}

// addMapVariable adds a map[string]string to the context
func addMapVariable(ctx *hcl.EvalContext, varName string, v map[string]string) {
	value, err := gocty.ToCtyValue(v, cty.Map(cty.String))
	if err != nil {
		// map[string]string conversion should never fail
		panic(err)
	}

	addCtyVariable(ctx, varName, value)
}

// addCtyVariable adds a cty variable to the context
func addCtyVariable(ctx *hcl.EvalContext, varName string, value cty.Value) {
	defineContextVarsIfNeeded(ctx)
	ctx.Variables[varName] = value
}

// finalizeContext should be called whenever an HCL context is being used
// as the final call.
func finalizeContext(ctx *hcl.EvalContext) *hcl.EvalContext {
	ctx = ctx.NewChild()
	ctx.Functions = funcs.MakeTemplateFuncs(ctx)
	return ctx
}

// AddVariables uses the final map of InputValues to add all input variables
// to the given hcl EvalContext.
func AddVariables(ctx *hcl.EvalContext, vs variables.Values) *hcl.EvalContext {
	variables.AddInputVariables(ctx, vs)
	return ctx
}

// hclContextContainer is an interface that config structs that have an HCL
// context may implement. We use this for certain things such as mapoperation()
// to set the proper context.
type hclContextContainer interface {
	hclContext() *hcl.EvalContext
}
