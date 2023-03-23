// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validationext

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsCidrBlock(t *testing.T) {
	cases := []struct {
		Input interface{}
		Valid bool
	}{
		{
			nil,
			false,
		},
		{
			int32(4),
			false,
		},
		{
			"192.168.0/24",
			false,
		},
		{
			"someString",
			false,
		},
		{
			"10.0.0.0/24",
			true,
		},
		{
			"10.255.255.255/24",
			true,
		},
		{
			"10.256.0.0/24",
			false,
		},
		{
			"10.0.256.0/24",
			false,
		},
		{
			"10.255.255asdfasdfaqsd.250",
			false,
		},
		{
			"192.168.0.0/24",
			true,
		},
		{
			"192.168.255.255/24",
			true,
		},
		{
			"192.168.256.0/24",
			false,
		},
		{
			"192.0.0.0/24",
			false,
		},
		{
			"172.16.0.0/24",
			true,
		},
		{
			"172.17.0.0/24",
			true,
		},
		{
			"172.18.0.0/24",
			true,
		},
		{
			"172.30.0.0/24",
			true,
		},
		{
			"172.20.0.0/24",
			true,
		},
		{
			"172.31.0.0/24",
			true,
		},
		{
			"172.15.0.0/24",
			false,
		},
		{
			"172.32.0.0/24",
			false,
		},
		{
			"172.192.0.0/24",
			false,
		},
		{
			"172.255.0.0/24",
			false,
		},
		{
			"10.0.0.0/7",
			false,
		},
		{
			"192.168.0.0/15",
			false,
		},
	}

	cases = append(cases, struct {
		Input interface{}
		Valid bool
	}{
		nil, false,
	})

	for _, tt := range cases {
		t.Run(fmt.Sprintf("%#v", tt.Input), func(t *testing.T) {
			err := IsPrivateCIDRBlock.Validate(tt.Input)
			require.Equal(t, tt.Valid, err == nil)
		})
	}
}
