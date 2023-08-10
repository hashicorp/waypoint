// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cli

import (
	"github.com/posener/complete"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type TriggerDeleteCommand struct {
	*baseCommand

	flagTriggerName string
	flagTriggerId   string
}

func (c *TriggerDeleteCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	flagSet := c.Flags()
	if err := c.Init(
		WithArgs(args),
		WithFlags(flagSet),
		WithNoConfig(),
	); err != nil {
		return 1
	}

	if len(c.args) == 0 && c.flagTriggerId == "" {
		c.ui.Output("Trigger ID required.\n\n%s", c.Help(), terminal.WithErrorStyle())
		return 1
	}
	c.flagTriggerId = c.args[0]

	ctx := c.Ctx

	_, err := c.project.Client().DeleteTrigger(ctx, &pb.DeleteTriggerRequest{
		Ref: &pb.Ref_Trigger{
			Id: c.flagTriggerId,
		},
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.ui.Output("Trigger configuration for %q not found", c.flagTriggerId, clierrors.Humanize(err),
				terminal.WithErrorStyle())
			return 1
		}

		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Trigger %q deleted", c.flagTriggerId, terminal.WithSuccessStyle())

	return 0
}

func (c *TriggerDeleteCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")

		f.StringVar(&flag.StringVar{
			Name:   "id",
			Target: &c.flagTriggerId,
			Usage:  "The id of the trigger URL to delete.",
		})
	})
}

func (c *TriggerDeleteCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *TriggerDeleteCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *TriggerDeleteCommand) Synopsis() string {
	return "Delete a registered trigger URL."
}

func (c *TriggerDeleteCommand) Help() string {
	return formatHelp(`
Usage: waypoint trigger delete [options] <trigger-id>

  Delete a trigger URL from Waypoint Server.

` + c.Flags().Help())
}
