// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ptypes

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func TestServerConfig(t testing.T, src *pb.ServerConfig) *pb.ServerConfig {
	t.Helper()

	if src == nil {
		src = &pb.ServerConfig{}
	}

	require.NoError(t, mergo.Merge(src, &pb.ServerConfig{
		AdvertiseAddrs: []*pb.ServerConfig_AdvertiseAddr{
			{
				Addr: "127.0.0.1",
			},
		},
	}))

	return src
}

// ValidateServerConfig validates the server config structure.
// TODO: This still panics if the server config is nil
func ValidateServerConfig(c *pb.ServerConfig) error {
	return validation.ValidateStruct(c,
		validation.Field(&c.AdvertiseAddrs, validation.Required, validation.Length(1, 1)),
	)
}
