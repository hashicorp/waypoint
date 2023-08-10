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

type HostnameDeleteCommand struct {
	*baseCommand
}

func (c *HostnameDeleteCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
	); err != nil {
		return 1
	}

	if len(c.args) == 0 {
		c.ui.Output("hostname required for deletion", terminal.WithErrorStyle())
		return 1
	}
	hostname := c.args[0]

	_, err := c.project.Client().DeleteHostname(c.Ctx, &pb.DeleteHostnameRequest{
		Hostname: hostname,
	})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Hostname deleted.", terminal.WithSuccessStyle())
	return 0
}

func (c *HostnameDeleteCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *HostnameDeleteCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *HostnameDeleteCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *HostnameDeleteCommand) Synopsis() string {
	return "Delete a previously registered hostname."
}

func (c *HostnameDeleteCommand) Help() string {
	return formatHelp(`
Usage: waypoint hostname delete HOSTNAME

  Delete a previously registered hostname.

` + c.Flags().Help())
}
