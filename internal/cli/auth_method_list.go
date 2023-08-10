// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cli

import (
	"fmt"

	"github.com/posener/complete"
	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type AuthMethodListCommand struct {
	*baseCommand
}

func (c *AuthMethodListCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoLocalServer(), // no auth in local mode
		WithNoConfig(),
	); err != nil {
		return 1
	}

	resp, err := c.project.Client().ListAuthMethods(c.Ctx, &empty.Empty{})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	table := terminal.NewTable("Name", "Method")

	// We always add "token" even though that's technically not a type
	// of auth method. This just makes it more obvious to users that
	// token auth is always available.
	table.Rich([]string{
		"token",
		"token",
	}, nil)

	for _, am := range resp.AuthMethods {
		method := fmt.Sprintf("%T", am.Method)
		switch am.Method.(type) {
		case *pb.AuthMethod_Oidc:
			method = "OIDC"
		}

		table.Rich([]string{
			am.Name,
			method,
		}, nil)
	}

	c.ui.Table(table)
	return 0
}

func (c *AuthMethodListCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *AuthMethodListCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *AuthMethodListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *AuthMethodListCommand) Synopsis() string {
	return "List all configured auth methods"
}

func (c *AuthMethodListCommand) Help() string {
	return formatHelp(`
Usage: waypoint auth-method list

  List all the auth methods configured with the Waypoint server.

  This will list all the ways that a user can log in to the Waypoint server.
  For most day-to-day Waypoint users, this doesn't provide much value. You
  can use the results of this command with "waypoint login" to target a
  specific auth method. However, if there is only one auth method other than "token", then
  "waypoint login" automatically uses that method.

`)
}
