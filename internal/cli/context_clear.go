// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cli

import (
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

type ContextClearCommand struct {
	*baseCommand
}

func (c *ContextClearCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	flagSet := c.Flags()
	if err := c.Init(
		WithArgs(args),
		WithFlags(flagSet),
		WithNoConfig(),
		WithNoClient(),
	); err != nil {
		return 1
	}
	args = flagSet.Args()

	if len(args) != 0 {
		c.ui.Output(c.Flags().Help(), terminal.WithErrorStyle())
		return 1
	}

	// Get our contexts
	if err := c.contextStorage.UnsetDefault(); err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Default context cleared.", terminal.WithSuccessStyle())
	return 0
}

func (c *ContextClearCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *ContextClearCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ContextClearCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ContextClearCommand) Synopsis() string {
	return "Unset the default context."
}

func (c *ContextClearCommand) Help() string {
	return formatHelp(`
Usage: waypoint context clear

  This unsets any default context.

  This does not delete any contexts. There are two use cases to not have
  a default set: (1) forcing yourself to specify a context to use or
  (2) operating in local (non-server) mode.

  For (2), you may also set the value of the "WAYPOINT_CONTEXT" environment
  variable to "-" which will force a local mode operation if supported by
  the project.

` + c.Flags().Help())
}
