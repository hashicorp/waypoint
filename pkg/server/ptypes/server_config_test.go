// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ptypes

import (
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func TestValidateServerConfig(t *testing.T) {
	cases := []struct {
		Name   string
		Modify func(*pb.ServerConfig)
		Error  string
	}{
		{
			"valid",
			nil,
			"",
		},

		{
			"no advertise addrs",
			func(c *pb.ServerConfig) { c.AdvertiseAddrs = nil },
			"advertise_addrs: cannot be blank",
		},

		{
			"two advertise addrs",
			func(c *pb.ServerConfig) {
				c.AdvertiseAddrs = append(c.AdvertiseAddrs, nil)
			},
			"advertise_addrs: the length must be exactly 1",
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)

			cfg := TestServerConfig(t, nil)
			if f := tt.Modify; f != nil {
				f(cfg)
			}

			err := ValidateServerConfig(cfg)
			if tt.Error == "" {
				require.NoError(err)
				return
			}

			require.Error(err)
			require.Contains(err.Error(), tt.Error)
		})
	}
}
