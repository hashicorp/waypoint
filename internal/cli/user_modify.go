// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cli

import (
	"strings"

	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type UserModifyCommand struct {
	*baseCommand

	flagUsername    string
	flagNewUsername string
	flagDisplay     string
}

func (c *UserModifyCommand) Run(args []string) int {
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

	// Get the user we're modifying.
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

	// Perform modifications
	willModify := false
	if v := c.flagNewUsername; v != "" {
		user.Username = v
		willModify = true
	}
	if v := c.flagDisplay; v != "" {
		user.Display = v
		willModify = true
	}

	if !willModify {
		c.ui.Output("at least one user modification flag must be specified"+
			c.Help(), terminal.WithErrorStyle())
		return 1
	}

	if _, err := client.UpdateUser(c.Ctx, &pb.UpdateUserRequest{User: user}); err != nil {
		c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	// This isn't the most user friendly output but it gets the job done at the moment.
	c.ui.Output("User modification successful.", terminal.WithSuccessStyle())
	return 0
}

func (c *UserModifyCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.StringVar(&flag.StringVar{
			Name:   "username",
			Target: &c.flagUsername,
			Usage:  "The user to modify. This defaults to the currently logged in user.",
		})

		f.StringVar(&flag.StringVar{
			Name:   "new-username",
			Target: &c.flagNewUsername,
			Usage:  "Set a new username for this user. This must be unique to the server.",
		})

		f.StringVar(&flag.StringVar{
			Name:   "display-name",
			Target: &c.flagDisplay,
			Usage: "The display name for a user. If this is set, this is used in some " +
				"places in the CLI and UI. This does not have to be unique.",
		})
	})
}

func (c *UserModifyCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *UserModifyCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *UserModifyCommand) Synopsis() string {
	return "Modify details about a user"
}

func (c *UserModifyCommand) Help() string {
	helpText := `
Usage: waypoint user modify [options]

  Modify details about a user.

  Some details such as username and display name may be updated.
  Use "waypoint user inspect" to see the current attributes for a user.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}
