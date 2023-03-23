// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// Simple wrapper around the hclwrite call to turn a traversal into a string.
func renderTraversal(t hcl.Traversal) (string, error) {
	tokens := hclwrite.TokensForTraversal(t)
	var buf bytes.Buffer

	_, err := tokens.WriteTo(&buf)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// Captures the information about a variable calculated during sort to be used
// without recalculation by the caller.
type analyzedPair struct {
	Pair hcl.KeyValuePair
	Refs []string

	// Name is the name of the key, which is either an env var name,
	// file path, or internal var name.
	Name string

	// Internal is true if this is an internal var (in the `internal =` map)
	Internal bool

	// Path is true if this is a file config, not an env var config.
	Path bool
}

// VariableLoopError is returned when, in the course of sorting the variables,
// a loop is detected. This means the variables can never be properly evaluated.
type VariableLoopError struct {
	LoopVars []string
}

func (v *VariableLoopError) Error() string {
	return fmt.Sprintf("loop detected amongst variables: %s", strings.Join(v.LoopVars, ", "))
}

// sortVars performs a topological sort on EnvRaw and InternalRaw. The sort
// yields the pairs in most referenced to least referenced order. Meaning
// that the if pair X references pair R, then R will be before X in the slice.
func (c *genericConfig) sortTopLevelVars(ctx *hcl.EvalContext) ([]*analyzedPair, error) {
	return sortVars(ctx, []sortVarMap{
		{Expr: c.EnvRaw, Prefix: "config.env."},
		{Expr: c.InternalRaw, Prefix: "config.internal.", Internal: true},
		{Expr: c.FileRaw, Prefix: "config.file.", Path: true},
	})
}

// sortVarMap is used as an input to sortVars to specify a map of variables.
type sortVarMap struct {
	Expr     hcl.Expression // HCL map of vars
	Prefix   string         // HCL prefix of the reference to this.
	Internal bool           // True if an internal var
	Path     bool           // True if a file var
}

// sortVars performs a topological sort on the given input maps and
// yields the pairs in most referenced to least referenced order. Meaning
// that the if pair X references pair R, then R will be before X in the slice.
func sortVars(ctx *hcl.EvalContext, maps []sortVarMap) ([]*analyzedPair, error) {
	// The algorithm used to perform the sort is Kahn's topological sorting algorithm.
	// https://www.geeksforgeeks.org/topological-sorting-indegree-based-solution/
	//
	// degrees tracks how many times a variable is referenced.
	// pairMap maps a variable's name to its data.
	degrees := map[string]int{}
	pairMap := map[string]*analyzedPair{}

	for _, m := range maps {
		pairs, diags := hcl.ExprMap(m.Expr)
		if diags.HasErrors() {
			continue
		}

		for _, pair := range pairs {
			// Decode the key. The key must be a string.
			val, diags := pair.Key.Value(ctx)
			if diags.HasErrors() {
				return nil, diags
			}
			if val.Type() != cty.String {
				rng := pair.Key.Range()
				return nil, &hcl.Diagnostic{
					Severity:    hcl.DiagError,
					Summary:     "key must be string",
					Subject:     &rng,
					Expression:  pair.Key,
					EvalContext: ctx,
				}
			}

			key := val.AsString()

			// We track the references using the "traversal" name, for instance
			// config.env.blah. So we need to create this long name as the referenced
			// name.
			pubName := m.Prefix + key

			var refs []string
			for _, ref := range pair.Value.Variables() {
				name, err := renderTraversal(ref)
				if err != nil {
					return nil, err
				}

				refs = append(refs, name)
			}

			// We initialize each variable to 0 to pick up later. This way, all variables,
			// even if unreferenced, are in degrees.
			degrees[pubName] = 0
			pairMap[pubName] = &analyzedPair{
				Pair:     pair,
				Name:     key,
				Refs:     refs,
				Internal: m.Internal,
				Path:     m.Path,
			}
		}
	}

	// Now we build up degrees by walking all the references on all the pairs.
	// This is the start of Kahn's algorithm.
	for _, pair := range pairMap {
		for _, ref := range pair.Refs {
			degrees[ref]++
		}
	}

	// toCheck is a work list of pairs that should now be checked.
	var toCheck []*analyzedPair

	// We initialize toCheck by finding all the pairs with no references.
	for name, deg := range degrees {
		if deg == 0 {
			toCheck = append(toCheck, pairMap[name])
		}
	}

	var order []*analyzedPair

	// This loop basically is walking already known good variables and
	// trying to make more variables good by reduced degree that we built
	// up earlier.
	for len(toCheck) > 0 {
		x := toCheck[len(toCheck)-1]
		toCheck = toCheck[:len(toCheck)-1]

		order = append(order, x)

		for _, ref := range x.Refs {
			deg := degrees[ref] - 1
			degrees[ref] = deg

			if deg == 0 {
				// The ref might be to a variable that isn't known atm
				if pair, ok := pairMap[ref]; ok {
					toCheck = append(toCheck, pair)
				}
			}
		}
	}

	// Now check that everything in degrees is 0, otherwise there is a loop!
	var loopVars []string

	for name, deg := range degrees {
		if deg != 0 {
			loopVars = append(loopVars, name)
		}
	}

	if len(loopVars) > 0 {
		sort.Strings(loopVars)
		return nil, &VariableLoopError{LoopVars: loopVars}
	}

	// Gotta reverse it because the order is least references to most and we want
	// to evaluate the most ref'd first.

	for i := 0; i < len(order)/2; i++ {
		x, y := order[i], order[len(order)-1-i]
		order[i], order[len(order)-1-i] = y, x
	}

	return order, nil
}
