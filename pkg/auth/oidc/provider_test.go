// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package oidc

import (
	"context"
	"testing"

	"github.com/hashicorp/cap/oidc"
	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func TestProviderCache(t *testing.T) {
	cache := NewProviderCache()
	defer cache.Close()
	require := require.New(t)
	ctx := context.Background()

	// Create our OIDC test provider
	oidcTP := oidc.StartTestProvider(t)
	oidcTP.SetClientCreds("alice", "big-secret")
	_, _, tpAlg, _ := oidcTP.SigningKeys()

	// Create our auth method configuration
	amOIDC := &pb.AuthMethod_OIDC{
		ClientId:            "alice",
		ClientSecret:        "big-secret",
		DiscoveryUrl:        oidcTP.Addr(),
		DiscoveryCaPem:      []string{oidcTP.CACert()},
		SigningAlgs:         []string{string(tpAlg)},
		AllowedRedirectUris: []string{"http://example.com"},
	}

	// Create the auth method
	am := serverptypes.TestAuthMethod(t, &pb.AuthMethod{
		Name: "TEST",
		Method: &pb.AuthMethod_Oidc{
			Oidc: amOIDC,
		},
	})

	// Get our first value
	p, err := cache.Get(ctx, am, nil)
	require.NoError(err)
	require.NotNil(p)

	// Get a second value, they should be pointer equal
	p2, err := cache.Get(ctx, am, nil)
	require.NoError(err)
	require.Equal(p, p2)

	// Update the config
	amOIDC.AllowedRedirectUris = []string{"http://example.com/foo"}
	p2, err = cache.Get(ctx, am, nil)
	require.NoError(err)
	require.NotNil(p2)
	require.NotEqual(p, p2)
}
