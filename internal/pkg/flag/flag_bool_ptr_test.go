// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package flag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoolPtrVar(t *testing.T) {
	varTrue := true
	varFalse := false
	cases := map[string]struct {
		input     *string
		expected  *bool
		omitValue bool // e.g. -poll instead of -poll=true
		shouldErr bool
	}{
		"omitted": {
			expected: nil,
		},
		// support -flag for "true"
		"bool flag behavior": {
			omitValue: true,
			expected:  &varTrue,
		},
		"true": {
			input:    strPtr("TRUE"),
			expected: &varTrue,
		},
		"false": {
			input:    strPtr("False"),
			expected: &varFalse,
		},
		// empty is equivilant to a user supplying -flag="", and not a scenario
		// where the flag is simply omittied.
		"empty": {
			input:     strPtr(""),
			shouldErr: true,
		},
		"bad input": {
			input:     strPtr("a non-truthy value"),
			shouldErr: true,
		},
		// only accept values Go's ParseBool would accept
		"weird truthy input": {
			input:     strPtr("tRuE"),
			shouldErr: true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			var valA *bool
			var valB string
			sets := NewSets()
			{
				set := sets.NewSet("A")
				set.BoolPtrVar(&BoolPtrVar{
					Name:   "a",
					Target: &valA,
				})
			}
			{
				// borrowed from string_slice_test, just to have another input
				set := sets.NewSet("B")
				set.StringVar(&StringVar{
					Name:   "b",
					Target: &valB,
				})
			}

			var err error
			inputs := []string{"-b=somevalueB"}
			if c.input != nil {
				inputs = append(inputs, "-a="+*c.input)
			}
			if c.omitValue == true {
				inputs = append(inputs, "-a")
			}
			err = sets.Parse(inputs)
			if c.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, c.expected, valA)
		})
	}
}

func strPtr(s string) *string {
	return &s
}
