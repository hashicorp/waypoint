// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package funcs

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty/function"
)

// This adds the functions that should be available to any HCL run in the
// entrypoint. It excludes VCS functions functions because that information
// is lost to be used in the entrypoint.
func AddEntrypointFunctions(ctx *hcl.EvalContext) {
	// Start with our HCL stdlib
	set := Stdlib()

	// add functions to our context
	addFuncs := func(fs map[string]function.Function) {
		for k, v := range fs {
			set[k] = v
		}
	}

	// Add some of our functions
	addFuncs(Filesystem())
	addFuncs(Encoding())
	addFuncs(Datetime())
	addFuncs(Jsonnet())
	addFuncs(Selector())

	ctx.Functions = set
}

// This adds the functions that should be able to any HCL run in the
// waypoint server and CLI context.
func AddStandardFunctions(ctx *hcl.EvalContext, pwd string) {
	// Start with our HCL stdlib
	set := Stdlib()

	// add functions to our context
	addFuncs := func(fs map[string]function.Function) {
		for k, v := range fs {
			set[k] = v
		}
	}

	// Add some of our functions
	addFuncs(VCSGitFuncs(pwd))
	addFuncs(Filesystem())
	addFuncs(Encoding())
	addFuncs(Datetime())
	addFuncs(Jsonnet())
	addFuncs(Selector())

	ctx.Functions = set
}
