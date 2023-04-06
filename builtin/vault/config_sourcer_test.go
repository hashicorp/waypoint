// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	pb "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	"github.com/hashicorp/waypoint/builtin/vault/testvault"
)

func TestConfigSourcer(t *testing.T) {
	ctx := context.Background()
	log := hclog.L()
	require := require.New(t)

	client, closer := testvault.TestVault(t)
	defer closer()

	// Write our secret
	_, err := client.Logical().Write("secret/data/my-secret", map[string]interface{}{
		"data": map[string]interface{}{
			"value": "world",
		},
	})
	require.NoError(err)

	cs := &ConfigSourcer{Client: client}
	defer cs.stop()

	// Read
	result, err := cs.read(ctx, log, []*component.ConfigRequest{
		{
			Name: "HELLO",
			Config: map[string]string{
				"path": "secret/data/my-secret",
				"key":  "/data/value",
			},
		},
	})
	require.NoError(err)
	require.NotNil(result)
	require.Len(result, 1)

	v := result[0]
	require.Equal("HELLO", v.Name)
	require.Equal("world", v.Result.(*pb.ConfigSource_Value_Value).Value)
}
