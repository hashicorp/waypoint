// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package funcs

import (
	"time"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

// Stdlib are the functions provided by the HCL stdlib.
func Datetime() map[string]function.Function {
	return map[string]function.Function{
		"timestamp": TimestampFunc,
	}
}

// TimestampFunc constructs a function that returns a string representation of the current date and time.
var TimestampFunc = function.New(&function.Spec{
	Params: []function.Parameter{},
	Type:   function.StaticReturnType(cty.String),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		return cty.StringVal(time.Now().UTC().Format(time.RFC3339)), nil
	},
})
