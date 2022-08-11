package tokenutil

import (
	"context"
	"net/url"

	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"

	"github.com/hashicorp/go-hclog"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func SetupExternalCreds(ctx context.Context, log hclog.Logger, token, aud string) (credentials.PerRPCCredentials, error) {
	transport, _, err := TokenDecode(token)
	if err != nil {
		return nil, err
	}

	if transport.ExternalCreds != nil {
		if oc, ok := transport.ExternalCreds.(*pb.TokenTransport_OauthCreds); ok {
			conf := &clientcredentials.Config{
				ClientID:       oc.OauthCreds.ClientId,
				ClientSecret:   oc.OauthCreds.ClientSecret,
				TokenURL:       oc.OauthCreds.Url,
				EndpointParams: url.Values{"audience": {aud}},
			}

			oauthToken, err := conf.Token(ctx)
			if err != nil {
				return nil, err
			}

			log.Info("utilizing credentials fetched via oauth for server auth",
				"oauth-url", oc.OauthCreds.Url,
				"oauth-client-id", oc.OauthCreds.ClientId,
			)

			// Remove the external creds from the token before we transmit it.
			minToken, err := StripCreds(transport)
			if err != nil {
				return nil, err
			}

			// We pass back the oauth access token and the waypoint token because
			// the waypoint server uses the token for identity as well.
			return &TokenAndAuth{
				PerRPCCredentials: oauth.NewOauthAccess(oauthToken),
				Token:             minToken,
			}, nil
		}
	}

	return nil, nil
}
