// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package validationext

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_IsVersion(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// This is weird but forced by ozzo-validation, so users should
			// pair that with validation.Required.
			"",
			false,
		},

		{
			"bob",
			false,
		},

		{
			"1",
			true,
		},

		{
			"1.0",
			true,
		},

		{
			"1.0.1",
			true,
		},

		{
			"v1",
			true,
		},

		{
			"v1.0",
			true,
		},

		{
			"v1.0.1",
			true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.Input, func(t *testing.T) {
			err := IsVersion.Validate(tt.Input)
			require.Equal(t, tt.Valid, err == nil)
		})
	}
}

func Test_MeetsConstraints(t *testing.T) {
	cases := []struct {
		Input      string
		Constraint []string
		Valid      bool
	}{
		{
			// This is weird but forced by ozzo-validation, so users should
			// pair that with validation.Required.
			"",
			[]string{"< 1"},
			false,
		},

		{
			"bob",
			[]string{"< 1"},
			false,
		},

		{
			"1",
			[]string{"< 2"},
			true,
		},

		{
			"1.0",
			[]string{"<1"},
			false,
		},

		{
			"1.0.1",
			[]string{">0", "<2"},
			true,
		},

		{
			"v1.0.1",
			[]string{">0", "<2"},
			true,
		},

		{
			"1.0.2",
			[]string{">v0", "<v2"},
			true,
		},

		{
			"v1.0.3",
			[]string{">v2"},
			false,
		},

		{
			"v1.0.4",
			[]string{">2"},
			false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.Input, func(t *testing.T) {
			f := MeetsConstraints(tt.Constraint...)
			err := f.Validate(tt.Input)
			require.Equal(t, tt.Valid, err == nil)
		})
	}

	testPanic := func() { MeetsConstraints("not a constraint") }
	// If a bad constraint is passed, we explode.
	t.Run("check we panic on bad constraint", func(t *testing.T) {
		require.Panics(t, testPanic)
	})
}
