package singleprocess

import (
	"context"
	"testing"

	"github.com/hashicorp/cap/oidc"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

func TestGetOIDCAuthURL(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(err)
	client := server.TestServer(t, impl)

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

	// Create
	{
		resp, err := client.UpsertAuthMethod(ctx, &pb.UpsertAuthMethodRequest{
			AuthMethod: serverptypes.TestAuthMethod(t, &pb.AuthMethod{
				Name: "TEST",
				Method: &pb.AuthMethod_Oidc{
					Oidc: amOIDC,
				},
			}),
		})
		require.NoError(err)
		require.NotNil(resp)
	}

	// Get our URL
	resp, err := client.GetOIDCAuthURL(ctx, &pb.GetOIDCAuthURLRequest{
		AuthMethod:  &pb.Ref_AuthMethod{Name: "TEST"},
		RedirectUri: "http://example.com",
		ClientNonce: "test",
	})
	require.NoError(err)
	require.NotNil(resp)
	require.NotEmpty(resp.Url)
	t.Logf("auth url: %s", resp.Url)
}
