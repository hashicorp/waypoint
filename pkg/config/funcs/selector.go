// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package funcs

import (
	"errors"

	"github.com/hashicorp/go-bexpr"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/gocty"
)

func Selector() map[string]function.Function {
	return map[string]function.Function{
		"selectormatch":  SelectorMatchFunc,
		"selectorlookup": SelectorLookupFunc,
	}
}

// SelectorMatchFunc constructs a function that applies a label selector
// to a map and returns true/false if there is a match.
var SelectorMatchFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name: "map",
			Type: cty.Map(cty.String),
		},
		{
			Name: "selector",
			Type: cty.String,
		},
	},
	Type: function.StaticReturnType(cty.Bool),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		s := args[1].AsString()
		eval, err := bexpr.CreateEvaluator(s)
		if err != nil {
			return cty.UnknownVal(cty.String), err
		}

		var m map[string]string
		if err := gocty.FromCtyValue(args[0], &m); err != nil {
			return cty.UnknownVal(cty.String), err
		}

		result, err := eval.Evaluate(m)
		if err != nil {
			return cty.UnknownVal(cty.String), err
		}

		return cty.BoolVal(result), nil
	},
})

// SelectorLookupFunc constructs a function that applies a label selector
// to a map and returns true/false if there is a match.
var SelectorLookupFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name: "map",
			Type: cty.Map(cty.String),
		},
		{
			Name: "selectormap",
			Type: cty.Map(cty.DynamicPseudoType),
		},
		{
			Name: "default",
			Type: cty.DynamicPseudoType,
		},
	},
	Type: func(args []cty.Value) (ret cty.Type, err error) {
		expected := args[1].Type().ElementType()
		if !args[2].Type().Equals(expected) {
			return cty.NilType, errors.New("default value type must match types of selector map")
		}

		return expected, nil
	},
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		var m map[string]string
		if err := gocty.FromCtyValue(args[0], &m); err != nil {
			return cty.UnknownVal(cty.String), err
		}

		// Go through the selector map and find one that matches.
		for it := args[1].ElementIterator(); it.Next(); {
			key, val := it.Element()

			s := key.AsString()
			eval, err := bexpr.CreateEvaluator(s)
			if err != nil {
				return cty.UnknownVal(cty.String), err
			}

			result, err := eval.Evaluate(m)
			if err != nil {
				return cty.UnknownVal(cty.String), err
			}

			if result {
				return val, nil
			}
		}

		return args[2], nil
	},
})

// SelectorMatch applies a selector to a map and returns true if the selector
// matches. The selector should be in go-bexpr format.
func SelectorMatch(m, selector cty.Value) (cty.Value, error) {
	return SelectorMatchFunc.Call([]cty.Value{m, selector})
}

// SelectorLookup applies a selector to a map and returns true if the selector
// matches. The selector should be in go-bexpr format.
func SelectorLookup(m, selector, def cty.Value) (cty.Value, error) {
	return SelectorLookupFunc.Call([]cty.Value{m, selector, def})
}
