// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validationext

import (
	"fmt"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

func TestError(t *testing.T) {
	type Address struct {
		Number int
		Street string
	}
	type AddressJSON struct {
		Number int    `json:"number"`
		Street string `json:"street"`
	}

	s1 := &struct{ Name string }{Name: ""}
	s2 := &struct {
		Name    string
		Address *Address
	}{
		Name: "",
		Address: &Address{
			Number: 123,
			Street: "",
		},
	}
	s3 := &struct {
		Name    string       `json:"name"`
		Address *AddressJSON `json:"address"`
	}{
		Name: "",
		Address: &AddressJSON{
			Number: 123,
			Street: "",
		},
	}

	cases := []struct {
		Name     string
		Input    error
		Expected map[string]string
	}{
		{
			"nil",
			nil,
			nil,
		},

		{
			"non-validation error",
			fmt.Errorf("hello"),
			nil,
		},

		{
			"basic validation error",
			validation.ValidateStruct(s1,
				validation.Field(&s1.Name, validation.Required),
			),
			map[string]string{
				"Name": "cannot be blank",
			},
		},

		{
			"nested struct validation error",
			validation.ValidateStruct(s2,
				validation.Field(&s2.Name, validation.Required),
				validation.Field(&s2.Address, validation.Required),
				StructField(&s2.Address, func() []*validation.FieldRules {
					return []*validation.FieldRules{
						validation.Field(&s2.Address.Street, validation.Required),
					}
				}),
			),
			map[string]string{
				"Name":           "cannot be blank",
				"Address.Street": "cannot be blank",
			},
		},

		{
			"nested struct with json tags",
			validation.ValidateStruct(s3,
				validation.Field(&s3.Name, validation.Required),
				validation.Field(&s3.Address, validation.Required),
				StructField(&s3.Address, func() []*validation.FieldRules {
					return []*validation.FieldRules{
						validation.Field(&s3.Address.Street, validation.Required),
					}
				}),
			),
			map[string]string{
				"name":           "cannot be blank",
				"address.street": "cannot be blank",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)

			actual := Error(tt.Input)
			if len(tt.Expected) == 0 {
				require.Equal(actual, tt.Input)
				return
			}

			st, ok := status.FromError(actual)
			require.True(ok, fmt.Sprintf("%T", actual))

			for _, errMessage := range tt.Expected {
				require.Contains(st.Message(), errMessage)
			}

			require.Len(st.Details(), 1)
			br := st.Details()[0].(*errdetails.BadRequest)
			require.Len(br.FieldViolations, len(tt.Expected))
			for _, fv := range br.FieldViolations {
				require.Contains(tt.Expected, fv.Field)
				require.Contains(fv.Description, tt.Expected[fv.Field])
			}
		})
	}
}
