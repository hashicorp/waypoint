// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/posener/complete"
)

type UserTokenCommand struct {
	*baseCommand

	flagDuration        time.Duration
	flagUsername        string
	flagTriggerUrlToken bool
}

func (c *UserTokenCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
		WithNoLocalServer(), // local mode has no need for tokens
	); err != nil {
		return 1
	}

	var refUser *pb.Ref_User
	if c.flagUsername != "" {
		refUser = &pb.Ref_User{
			Ref: &pb.Ref_User_Username{
				Username: &pb.Ref_UserUsername{
					Username: c.flagUsername,
				},
			},
		}
	}

	// Generate the token
	client := c.project.Client()
	resp, err := client.GenerateLoginToken(c.Ctx, &pb.LoginTokenRequest{
		Duration: c.flagDuration.String(),
		User:     refUser,
		Trigger:  c.flagTriggerUrlToken,
	})
	if err != nil {
		c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	// We use fmt here and not the UI helpers because UI helpers will
	// trim tokens horizontally on terminals that are narrow.
	fmt.Println(resp.Token)
	return 0
}

func (c *UserTokenCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.DurationVar(&flag.DurationVar{
			Name:    "expires-in",
			Target:  &c.flagDuration,
			Usage:   "The duration until the token expires. i.e. '5m'.",
			Default: 720 * time.Hour, // 30 days
		})

		f.StringVar(&flag.StringVar{
			Name:   "username",
			Target: &c.flagUsername,
			Usage: "Username to generate the login token for. This defaults " +
				"to the currently logged in user.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:   "trigger-url-token",
			Target: &c.flagTriggerUrlToken,
			Usage: "Will generate a trigger auth token. This token can only be used " +
				"for trigger URL actions.",
			Default: false,
		})
	})
}

func (c *UserTokenCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *UserTokenCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *UserTokenCommand) Synopsis() string {
	return "Request a new token to access the server"
}

func (c *UserTokenCommand) Help() string {
	helpText := `
Usage: waypoint user token [options]

  Request a new login token for a user.

  You must be logged in already to generate a token. If you need to
  log in, use the "waypoint login" command.

  This generates a new token that can be used to authenticate directly
  to the Waypoint server. If you're inviting a new user to Waypoint,
  its recommended to generate an invite token with "waypoint user invite"
  or share the UI URL for logging in.
` + warnTokenDeprecated + "\n" + c.Flags().Help()

	return strings.TrimSpace(helpText)
}
