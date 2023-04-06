// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package nullify

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNullify(t *testing.T) {
	intVal := 3

	cases := []struct {
		Name     string
		Input    interface{}
		Types    []interface{}
		Expected interface{}
	}{
		{
			"basic struct",
			&struct {
				A *int
				B int
			}{
				A: &intVal,
				B: 42,
			},
			[]interface{}{(*int)(nil)},
			&struct {
				A *int
				B int
			}{
				A: nil,
				B: 42,
			},
		},

		{
			"struct with no types",
			&struct {
				A *int
				B int
			}{
				A: &intVal,
				B: 42,
			},
			[]interface{}{(*string)(nil)},
			&struct {
				A *int
				B int
			}{
				A: &intVal,
				B: 42,
			},
		},
		{
			"unexported fields are ignored",
			&struct {
				a *int
				B int
			}{
				a: &intVal,
				B: 42,
			},
			[]interface{}{(*int)(nil)},
			&struct {
				a *int
				B int
			}{
				a: &intVal,
				B: 42,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require.NoError(t, Nullify(tt.Input, tt.Types...))
			require.Equal(t, tt.Input, tt.Expected)
		})
	}
}
