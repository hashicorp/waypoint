// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"strings"

	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type UserInspectCommand struct {
	*baseCommand

	flagUsername string
}

func (c *UserInspectCommand) Run(args []string) int {
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

	userResp, err := client.GetUser(c.Ctx, &pb.GetUserRequest{User: refUser})
	if err != nil {
		c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}
	user := userResp.User

	// This isn't the most user friendly output but it gets the job done at the moment.
	c.ui.Output(user.String())
	return 0
}

func (c *UserInspectCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.StringVar(&flag.StringVar{
			Name:   "username",
			Target: &c.flagUsername,
			Usage:  "The user to lookup. This defaults to the currently logged in user.",
		})
	})
}

func (c *UserInspectCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *UserInspectCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *UserInspectCommand) Synopsis() string {
	return "Show details about a single user"
}

func (c *UserInspectCommand) Help() string {
	helpText := `
Usage: waypoint user inspect [options]

  Show details about a single user, defaulting to the currently logged in user.

  This shows details about the currently logged in user or any user that
  is specified via flags.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}
