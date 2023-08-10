// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package flag

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringSlice(t *testing.T) {
	require := require.New(t)

	var valA, valB []string
	sets := NewSets()
	{
		set := sets.NewSet("A")
		set.StringSliceVar(&StringSliceVar{
			Name:   "a",
			Target: &valA,
		})
	}

	{
		set := sets.NewSet("B")
		set.StringSliceVar(&StringSliceVar{
			Name:   "b",
			Target: &valB,
		})
	}

	err := sets.Parse([]string{
		"-b", "somevalueB",
		"-a", "somevalueA,somevalueB",
	})
	require.NoError(err)

	require.Equal([]string{"somevalueB"}, valB)
	require.Equal([]string{"somevalueA", "somevalueB"}, valA)
}
