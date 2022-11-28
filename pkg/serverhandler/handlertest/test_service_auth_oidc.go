package handlertest

import (
	"context"
	"testing"

	"github.com/hashicorp/cap/oidc"
	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func init() {
	tests["auth_oidc"] = []testFunc{
		TestOIDCAuth,
		TestOIDCAuth_accessSelector,
	}
}

func TestOIDCAuth(t *testing.T, factory Factory) {
	require := require.New(t)
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

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
		AllowedRedirectUris: []string{"https://example.com"},
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
		RedirectUri: "https://example.com",
	})
	require.NoError(err)
	require.NotNil(resp)
	require.NotEmpty(resp.Url)
	t.Logf("auth url: %s", resp.Url)

	// Setup our test provider to auth
	oidcTP.SetExpectedState("state")
	oidcTP.SetExpectedAuthCode("hello")
	oidcTP.SetExpectedAuthNonce("nonce")

	// Complete our auth
	respAuth, err := client.CompleteOIDCAuth(ctx, &pb.CompleteOIDCAuthRequest{
		AuthMethod:  &pb.Ref_AuthMethod{Name: "TEST"},
		RedirectUri: "https://example.com",
		State:       "state",
		Code:        "hello",
		Nonce:       "nonce",
	})
	require.NoError(err)
	require.NotNil(respAuth)
	require.NotEmpty(respAuth.Token)
	user := respAuth.User
	require.NotNil(user)
	require.NotEmpty(user.Links)
	t.Logf("id claims %#v", respAuth.IdClaimsJson)
	t.Logf("user claims %#v", respAuth.UserClaimsJson)

	// Complete our auth again. We should get the same user.
	{
		respAuth, err := client.CompleteOIDCAuth(ctx, &pb.CompleteOIDCAuthRequest{
			AuthMethod:  &pb.Ref_AuthMethod{Name: "TEST"},
			RedirectUri: "https://example.com",
			State:       "state",
			Code:        "hello",
			Nonce:       "nonce",
		})
		require.NoError(err)
		require.NotNil(respAuth)

		user2 := respAuth.User
		require.NotNil(user2)
		require.Equal(user.Id, user2.Id)
	}
}

func TestOIDCAuth_accessSelector(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	// Create our OIDC test provider
	oidcTP := oidc.StartTestProvider(t)
	oidcTP.SetClientCreds("alice", "big-secret")
	_, _, tpAlg, _ := oidcTP.SigningKeys()

	cases := []struct {
		Name        string
		Selector    string
		Claims      map[string]interface{}
		Mapping     map[string]string
		ListMapping map[string]string
		Err         string
	}{
		{
			"list success",
			"hashicorp in list.g",
			map[string]interface{}{
				"groups": []string{"hashicorp", "dadgarcorp", "umbrellacorp"},
				"admin":  false,
			},
			nil,
			map[string]string{"groups": "g"},
			"",
		},

		{
			"list failure",
			"nopecorp in list.g",
			map[string]interface{}{
				"groups": []string{"hashicorp", "dadgarcorp", "umbrellacorp"},
				"admin":  false,
			},
			nil,
			map[string]string{"groups": "g"},
			"denied access",
		},

		{
			"key success",
			"value.admin != true",
			map[string]interface{}{
				"groups": []string{"hashicorp", "dadgarcorp", "umbrellacorp"},
				"admin":  false,
			},
			map[string]string{"admin": "admin"},
			nil,
			"",
		},

		{
			"key failure",
			"value.admin != false",
			map[string]interface{}{
				"groups": []string{"hashicorp", "dadgarcorp", "umbrellacorp"},
				"admin":  false,
			},
			map[string]string{"admin": "admin"},
			nil,
			"denied",
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)

			// Set some custom claims to write access selectors for
			oidcTP.SetCustomClaims(tt.Claims)

			// Create
			{
				amOIDC := &pb.AuthMethod_OIDC{
					ClientId:            "alice",
					ClientSecret:        "big-secret",
					DiscoveryUrl:        oidcTP.Addr(),
					DiscoveryCaPem:      []string{oidcTP.CACert()},
					SigningAlgs:         []string{string(tpAlg)},
					AllowedRedirectUris: []string{"https://example.com"},
					ClaimMappings:       tt.Mapping,
					ListClaimMappings:   tt.ListMapping,
				}

				resp, err := client.UpsertAuthMethod(ctx, &pb.UpsertAuthMethodRequest{
					AuthMethod: serverptypes.TestAuthMethod(t, &pb.AuthMethod{
						Name:           "TEST",
						AccessSelector: tt.Selector,
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
				RedirectUri: "https://example.com",
			})
			require.NoError(err)
			require.NotNil(resp)
			require.NotEmpty(resp.Url)

			// Setup our test provider to auth
			oidcTP.SetExpectedState("state")
			oidcTP.SetExpectedAuthCode("hello")
			oidcTP.SetExpectedAuthNonce("nonce")

			// Complete our auth
			respAuth, err := client.CompleteOIDCAuth(ctx, &pb.CompleteOIDCAuthRequest{
				AuthMethod:  &pb.Ref_AuthMethod{Name: "TEST"},
				RedirectUri: "https://example.com",
				State:       "state",
				Code:        "hello",
				Nonce:       "nonce",
			})
			if tt.Err == "" {
				require.NoError(err)
				require.NotNil(respAuth)
				require.NotEmpty(respAuth.Token)
				return
			}

			require.Error(err)
			require.Contains(err.Error(), tt.Err)
		})
	}
}
