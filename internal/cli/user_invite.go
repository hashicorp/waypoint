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

type UserInviteCommand struct {
	*baseCommand

	flagDuration time.Duration
	flagUsername string
}

func (c *UserInviteCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
		WithNoLocalServer(), // local mode has no need for tokens
	); err != nil {
		return 1
	}
	client := c.project.Client()

	req := &pb.InviteTokenRequest{
		Duration: c.flagDuration.String(),
	}

	if c.flagUsername != "" {
		req.Signup = &pb.Token_Invite_Signup{
			InitialUsername: c.flagUsername,
		}
	} else {
		// Invite ourselves, an existing user
		userResp, err := client.GetUser(c.Ctx, &pb.GetUserRequest{})
		if err != nil {
			c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}
		user := userResp.User

		req.Login = &pb.Token_Login{
			UserId: user.Id,
		}
	}

	// Generate the token
	resp, err := client.GenerateInviteToken(c.Ctx, req)
	if err != nil {
		c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	// We use fmt here and not the UI helpers because UI helpers will
	// trim tokens horizontally on terminals that are narrow.
	fmt.Println(resp.Token)
	return 0
}

func (c *UserInviteCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.DurationVar(&flag.DurationVar{
			Name:    "expires-in",
			Target:  &c.flagDuration,
			Usage:   "The duration until the token expires. i.e. '5m'.",
			Default: 24 * time.Hour,
		})

		f.StringVar(&flag.StringVar{
			Name:   "username",
			Target: &c.flagUsername,
			Usage: "Invite a new user and provide a username hint. The user " +
				"may still change their username after signing up. If this " +
				"isn't specified, an invite token for your current user will be " +
				"generated and no new user signup is performed.",
		})
	})
}

func (c *UserInviteCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *UserInviteCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *UserInviteCommand) Synopsis() string {
	return "Invite a user to join the Waypoint server"
}

func (c *UserInviteCommand) Help() string {
	helpText := `
Usage: waypoint user invite [options]

  Generate an invite token that can be used to log in.

  You must be logged in already to generate an invite token. If you need to
  log in, use the "waypoint login" command.

  This generates a new invite token. An invite token can be exchanged for
  a login token. If your Waypoint server has OIDC (non-token) auth enabled,
  it is recommended to instead invite users using your UI URL or directly
  via "waypoint login".

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}
