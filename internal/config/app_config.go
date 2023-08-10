// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package config

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
	"github.com/zclconf/go-cty/cty/function"

	"github.com/hashicorp/waypoint/internal/config/dynamic"
	"github.com/hashicorp/waypoint/internal/pkg/partial"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serversort "github.com/hashicorp/waypoint/pkg/server/sort"
)

// genericConfig represents the `config` stanza that can be placed
// both in the app and at the project level.
type genericConfig struct {
	// internal are variables which can be seen for templating but are not exposed
	// by default to an application or runner.
	InternalRaw hcl.Expression `hcl:"internal,optional"`

	// env are variables that will be exported into the application or runners
	// environment.
	EnvRaw hcl.Expression `hcl:"env,optional"`

	// file are paths that will be written to disk in the context of the application
	// environment.
	FileRaw hcl.Expression `hcl:"file,optional"`

	// Indicates a signal to send the application when config files change.
	FileChangeSignal string `hcl:"file_change_signal,optional"`

	// WorkspaceScoped are workspace-scoped config variables.
	WorkspaceScoped []*scopedConfig `hcl:"workspace,block"`

	// LabelScoped are label-selector-scoped config variables.
	LabelScoped []*scopedConfig `hcl:"label,block"`

	// Runner are runner-scoped configurations.
	Runner *runnerConfig `hcl:"runner,block"`

	ctx       *hcl.EvalContext    // ctx is the context to use when evaluating
	scopeFunc func(*pb.ConfigVar) // scopeFunc should set the scope for the config var
}

// runnerConfig is used for the `runner` blocks within `config`.
type runnerConfig struct {
	// Same as genericConfig, see there for docs.
	InternalRaw     hcl.Expression  `hcl:"internal,optional"`
	EnvRaw          hcl.Expression  `hcl:"env,optional"`
	FileRaw         hcl.Expression  `hcl:"file,optional"`
	WorkspaceScoped []*scopedConfig `hcl:"workspace,block"`
	LabelScoped     []*scopedConfig `hcl:"label,block"`
}

func (c *runnerConfig) configVars(
	ctx *hcl.EvalContext,
	scopeFunc func(*pb.ConfigVar),
) ([]*pb.ConfigVar, error) {
	if c == nil {
		return nil, nil
	}

	// Setup our scope function
	runnerScopeFunc := func(v *pb.ConfigVar) {
		// Inherit our parent scope.
		scopeFunc(v)
		v.Target.Runner = &pb.Ref_Runner{
			Target: &pb.Ref_Runner_Any{
				Any: &pb.Ref_RunnerAny{},
			},
		}
	}

	// Build a generic config since we're most similar to that and that
	// lets us share logic across the config var functions without
	// copy and pasting.
	generic := &genericConfig{
		InternalRaw:     c.InternalRaw,
		EnvRaw:          c.EnvRaw,
		FileRaw:         c.FileRaw,
		WorkspaceScoped: c.WorkspaceScoped,
		LabelScoped:     c.LabelScoped,
		scopeFunc:       runnerScopeFunc,
	}

	return generic.configVars(ctx)
}

// scopedConfig is used for the `workspace` and `label`-scoped config blocks
// within genericConfig as a way to further scope configuration.
type scopedConfig struct {
	// Scope is the label for the block. This is reused for both workspace
	// and label scoped variables so this could be either of those.
	Scope string `hcl:",label"`

	// Same as genericConfig, see there for docs.
	InternalRaw hcl.Expression `hcl:"internal,optional"`
	EnvRaw      hcl.Expression `hcl:"env,optional"`
	FileRaw     hcl.Expression `hcl:"file,optional"`
}

// configVars returns the set of ConfigVars ready to be sent to the API server.
//
// scopeFunc must be provided to set the proper scoping on the rendered
// variables since this struct on its own doesn't know.
func (s *scopedConfig) configVars(
	ctx *hcl.EvalContext,
	scopeFunc func(*pb.ConfigVar),
) ([]*pb.ConfigVar, error) {
	// sortVars performs a topological sort of the variables via references, so
	// the pairs can be evaluated top to bottom safely.
	pairs, err := sortVars(ctx, []sortVarMap{
		{Expr: s.EnvRaw, Prefix: "config.env."},
		{Expr: s.InternalRaw, Prefix: "config.internal.", Internal: true},
		{Expr: s.FileRaw, Prefix: "config.file.", Path: true},
	})
	if err != nil {
		return nil, err
	}

	return configVars(ctx, pairs, scopeFunc)
}

// ConfigVars returns the set of config vars that can be sent to the API server.
// This will include all the scoped config vars such as workspaces, labels,
// and runners.
func (c *genericConfig) ConfigVars() ([]*pb.ConfigVar, error) {
	if c == nil {
		return nil, nil
	}

	// Build our evaluation context for the config vars
	ctx := c.ctx
	ctx = appendContext(ctx, &hcl.EvalContext{
		Functions: dynamic.Register(map[string]function.Function{}),
	})
	ctx = finalizeContext(ctx)

	return c.configVars(ctx)
}

func (c *genericConfig) configVars(
	ctx *hcl.EvalContext,
) ([]*pb.ConfigVar, error) {
	// We copy ourselves to a scopedConfig so we can share the configVars
	// function. Otherwise, the two functions are nearly identical.
	rootScope := &scopedConfig{
		InternalRaw: c.InternalRaw,
		EnvRaw:      c.EnvRaw,
		FileRaw:     c.FileRaw,
	}
	result, err := rootScope.configVars(ctx, c.scopeFunc)
	if err != nil {
		return nil, err
	}

	// Do runners
	if c.Runner != nil {
		next, err := c.Runner.configVars(ctx, c.scopeFunc)
		if err != nil {
			return nil, err
		}

		result = append(result, next...)
	}

	// Build up our workspace-scoped configs.
	for _, wsScope := range c.WorkspaceScoped {
		next, err := wsScope.configVars(ctx, func(v *pb.ConfigVar) {
			// Always apply our root scope so that if this is a workspace-scoped
			// var WITHIN an app-scoped genericConfig, then it gets that target
			// too.
			c.scopeFunc(v)

			// Apply our own filters.
			v.Target.Workspace = &pb.Ref_Workspace{Workspace: wsScope.Scope}
		})
		if err != nil {
			return nil, err
		}

		result = append(result, next...)
	}

	// Build up our label-scoped configs.
	for _, scoped := range c.LabelScoped {
		next, err := scoped.configVars(ctx, func(v *pb.ConfigVar) {
			// Always apply our root scope so that if this is a label-scoped
			// var WITHIN an app-scoped genericConfig, then it gets that target
			// too.
			c.scopeFunc(v)

			// Apply our own filters.
			v.Target.LabelSelector = scoped.Scope
		})
		if err != nil {
			return nil, err
		}

		result = append(result, next...)
	}

	// Sort our results by name. This helps with deterministic behavior
	// in API calls, user output, etc. without forcing all callers to worry
	// about sorting.
	sort.Sort(serversort.ConfigName(result))

	return result, nil
}

// configVars returns the "rendered" list of config vars that are ready to
// be sent to the API server. As inputs, this requires the topologically
// sorted set of config vars (from sortVars) so that ordering is already
// pre-determined.
//
// The scopeFunc can be used to modify the config var and set proper
// targeting and other values. This is called before the value is set.
func configVars(
	ctx *hcl.EvalContext,
	sortedVars []*analyzedPair,
	scopeFunc func(*pb.ConfigVar),
) ([]*pb.ConfigVar, error) {
	// We're going to build up the variables as we go along using these 4 maps.
	ctx.Variables = map[string]cty.Value{}

	var (
		env      = map[string]cty.Value{}
		internal = map[string]cty.Value{}
		file     = map[string]cty.Value{}
		config   = map[string]cty.Value{}
	)

	var result []*pb.ConfigVar
	for _, pair := range sortedVars {
		key := pair.Name

		// Start building our var
		var newVar pb.ConfigVar
		newVar.Target = &pb.ConfigVar_Target{}
		newVar.Name = key
		newVar.Internal = pair.Internal
		newVar.NameIsPath = pair.Path
		scopeFunc(&newVar)

		// Decode the value
		val, diags := pair.Pair.Value.Value(ctx)
		if diags.HasErrors() {
			// Ok, we can't read its value right now. Let's do a partial evaluation then.
			str, err := partial.EvalExpression(ctx, pair.Pair.Value)
			if err != nil {
				return nil, err
			}

			newVar.Value = &pb.ConfigVar_Static{
				Static: str,
			}

			// We don't advertise these variables in the eval context because
			// we don't want them to be substituted as strings into other variables.
			// If the current variable is referenced by a later variable, we want
			// that to be a normal HCL template expansion of the variable reference,
			// not the contents. Quick example:
			//
			// a = "${g} ${s}"
			// b = "more: ${a}"
			// g = unknown()
			// s = "ok"
			//
			// After running the algorithm, we want b to still be 'more: ${a}', NOT
			// 'more: ${g} ok'. The reason being the 2nd one confuses the escaping
			// as it appears like it might be data that was returned from a file or
			// something.
		} else {
			switch val.Type() {
			case dynamic.Type:
				newVar.Value = &pb.ConfigVar_Dynamic{
					Dynamic: val.EncapsulatedValue().(*pb.ConfigVar_DynamicVal),
				}

			default:
				// For non-config val types we try to convert it to a string
				// as a static value.
				var err error
				val, err = convert.Convert(val, cty.String)
				if err != nil {
					return nil, err
				}

				// We have to escape any HCL we find in the string so that we don't
				// evaluate it down-stream.
				// First, we need to check if the value is not null, since we allow
				// `null` defaults for input variables, and a user may forget to
				// provide a value to an input variable
				if val.IsNull() {
					return nil, fmt.Errorf("could not evaluate %q in app config with `null` value", newVar.Name)
				}
				newVar.Value = &pb.ConfigVar_Static{
					Static: hclEscaper.Replace(val.AsString()),
				}

				if pair.Internal {
					internal[pair.Name] = val

					// Because of the nature of the hcl map type, we have to rebuild these
					// each time we modify them.
					config["internal"] = cty.MapVal(internal)
				} else if pair.Path {
					file[pair.Name] = val

					// Because of the nature of the hcl map type, we have to rebuild these
					// each time we modify them.
					config["file"] = cty.MapVal(file)
				} else {
					env[pair.Name] = val

					// Because of the nature of the hcl map type, we have to rebuild these
					// each time we modify them.
					config["env"] = cty.MapVal(env)
				}

				ctx.Variables["config"] = cty.MapVal(config)
			}
		}

		result = append(result, &newVar)
	}

	return result, nil
}

var (
	// hclEscaper is used to escape HCL in our config values.
	hclEscaper = strings.NewReplacer("${", "$${", "%{", "%%{")
)
