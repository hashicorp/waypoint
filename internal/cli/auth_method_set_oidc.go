// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cli

import (
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/posener/complete"
)

type AuthMethodSetOIDCCommand struct {
	*baseCommand

	flagDisplayName    string
	flagDescription    string
	flagAccessSelector string
	flagMethod         pb.AuthMethod_OIDC
}

func (c *AuthMethodSetOIDCCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoLocalServer(), // no auth in local mode
		WithNoConfig(),
	); err != nil {
		return 1
	}

	if len(c.args) == 0 {
		c.ui.Output("name required for the auth method", terminal.WithErrorStyle())
		return 1
	}
	name := c.args[0]

	am := &pb.AuthMethod{
		Name:           name,
		DisplayName:    c.flagDisplayName,
		Description:    c.flagDescription,
		AccessSelector: c.flagAccessSelector,
		Method: &pb.AuthMethod_Oidc{
			Oidc: &c.flagMethod,
		},
	}

	_, err := c.project.Client().UpsertAuthMethod(c.Ctx, &pb.UpsertAuthMethodRequest{
		AuthMethod: am,
	})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Auth method configured.", terminal.WithSuccessStyle())
	return 0
}

func (c *AuthMethodSetOIDCCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.StringVar(&flag.StringVar{
			Name:   "display-name",
			Target: &c.flagDisplayName,
			Usage:  "Display name for the UI. Optional.",
		})

		f.StringVar(&flag.StringVar{
			Name:   "description",
			Target: &c.flagDescription,
			Usage:  "Short description of this auth method. Optional.",
		})

		f.StringVar(&flag.StringVar{
			Name:   "access-selector",
			Target: &c.flagAccessSelector,
			Usage: "Selector expression to control access based on claims. " +
				"See docs for more details.",
		})

		f = set.NewSet("OIDC Auth Method Options")
		f.StringVar(&flag.StringVar{
			Name:   "client-id",
			Target: &c.flagMethod.ClientId,
			Usage:  "The OAuth 2.0 Client Identifier.",
		})

		f.StringVar(&flag.StringVar{
			Name:   "client-secret",
			Target: &c.flagMethod.ClientSecret,
			Usage:  "The client secret corresponding with the client ID.",
		})

		f.StringSliceVar(&flag.StringSliceVar{
			Name:   "claim-scope",
			Target: &c.flagMethod.Scopes,
			Usage:  "The optional claims scope requested. May be specified multiple times.",
		})

		f.StringSliceVar(&flag.StringSliceVar{
			Name:   "signing-algorithm",
			Target: &c.flagMethod.SigningAlgs,
			Usage:  "The allowed signing algorithm. May be specified multiple times.",
		})

		f.StringVar(&flag.StringVar{
			Name:   "issuer",
			Target: &c.flagMethod.DiscoveryUrl,
			Usage: "Discovery URL of the OIDC provider that implements the " +
				".well-known/openid-configuration metadata endpoint.",
		})

		f.StringSliceVar(&flag.StringSliceVar{
			Name:   "issuer-ca-pem",
			Target: &c.flagMethod.DiscoveryCaPem,
			Usage:  "PEM-encoded certificates for connecting to the issuer. May be specified multiple times.",
		})

		f.StringSliceVar(&flag.StringSliceVar{
			Name:   "allowed-redirect-uri",
			Target: &c.flagMethod.AllowedRedirectUris,
			Usage: "Allowed URI for auth redirection. This automatically has " +
				"localhost (for CLI auth) and the server address configured. " +
				"If you have additional external addresses, you can specify them here. " +
				"May be specified multiple times.",
		})

		f.StringMapVar(&flag.StringMapVar{
			Name:   "claim-mapping",
			Target: &c.flagMethod.ClaimMappings,
			Usage: "Mapping of a claim to a variable value for the access selector. " +
				"This can be specified multiple times. Example value: " +
				"'http://example.com/key=key'",
		})

		f.StringMapVar(&flag.StringMapVar{
			Name:   "list-claim-mapping",
			Target: &c.flagMethod.ListClaimMappings,
			Usage: "Same as claim-mapping but for list values. " +
				"This can be repeated multiple times.",
		})
	})
}

func (c *AuthMethodSetOIDCCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *AuthMethodSetOIDCCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *AuthMethodSetOIDCCommand) Synopsis() string {
	return "Configure an OIDC auth method"
}

func (c *AuthMethodSetOIDCCommand) Help() string {
	return formatHelp(`
Usage: waypoint auth-method set oidc [options] NAME

  Configure an OIDC auth method.

` + c.Flags().Help())
}
