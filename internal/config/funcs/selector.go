package funcs

import (
	"github.com/hashicorp/go-bexpr"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/gocty"
)

func Selector() map[string]function.Function {
	return map[string]function.Function{
		"selectormatch": SelectorMatchFunc,
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

// SelectorMatch applies a selector to a map and returns true if the selector
// matches. The selector should be in go-bexpr format.
func SelectorMatch(m, selector cty.Value) (cty.Value, error) {
	return SelectorMatchFunc.Call([]cty.Value{m, selector})
}
