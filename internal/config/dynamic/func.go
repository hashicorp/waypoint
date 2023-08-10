// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

// Package dynamic contains the HCL function, types, and logic for
// implementing dynamic config sourcing HCL configuration. This
// doesn't implement the actual logic behind configuration fetching
// (see `internal/appconfig`), only the declaration in an HCL file.
package dynamic

import (
	"reflect"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/gocty"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// Register registers the configdynamic function calls to the given
// map of HCL functions that can be used for an hcl.EvalContext.
// This function ensures that we can change the function name easily
// and consistently or add new dynamic-related functions.
func Register(m map[string]function.Function) map[string]function.Function {
	m["dynamic"] = Func

	// This is deprecated, but we keep this form working cause it costs nothing.
	m["configdynamic"] = Func

	return m
}

var (
	// Func implements the configdynamic() HCL function.
	Func = function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name: "from",
				Type: cty.String,
			},

			{
				Name: "config",
				Type: cty.Map(cty.String),
			},
		},
		Type: function.StaticReturnType(Type),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			var config map[string]string
			if err := gocty.FromCtyValue(args[1], &config); err != nil {
				return cty.NilVal, err
			}

			return cty.CapsuleVal(Type, &pb.ConfigVar_DynamicVal{
				From:   args[0].AsString(),
				Config: config,
			}), nil
		},
	})

	// Type is the encapsulated type that is returned by `Func`.
	// Encapsulated types are opaque in HCL and can be decoded
	// in Go to native Go types.
	Type = cty.Capsule("configval",
		reflect.TypeOf((*pb.ConfigVar_DynamicVal)(nil)).Elem())
)
