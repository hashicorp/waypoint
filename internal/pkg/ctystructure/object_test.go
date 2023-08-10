// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ctystructure

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

func TestObject(t *testing.T) {
	cases := []struct {
		Name   string
		Input  map[string]interface{}
		Output cty.Value
	}{
		{
			"primitives",
			map[string]interface{}{
				"string": "hello!",
			},
			cty.ObjectVal(map[string]cty.Value{
				"string": cty.StringVal("hello!"),
			}),
		},

		{
			"list of same types",
			map[string]interface{}{
				"ports": []interface{}{80, 100},
			},
			cty.ObjectVal(map[string]cty.Value{
				"ports": cty.ListVal([]cty.Value{
					cty.NumberIntVal(80),
					cty.NumberIntVal(100),
				}),
			}),
		},

		{
			"nested map",
			map[string]interface{}{
				"env": map[string]interface{}{
					"key":  "value",
					"port": 8080,
				},
			},
			cty.ObjectVal(map[string]cty.Value{
				"env": cty.ObjectVal(map[string]cty.Value{
					"key":  cty.StringVal("value"),
					"port": cty.NumberIntVal(8080),
				}),
			}),
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)

			val, err := Object(tt.Input)
			require.NoError(err)
			require.Equal(tt.Output, val)
		})
	}
}
