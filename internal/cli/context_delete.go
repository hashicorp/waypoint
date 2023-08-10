// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cli

import (
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

type ContextDeleteCommand struct {
	*baseCommand

	flagAll bool
}

func (c *ContextDeleteCommand) Run(args []string) int {
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

	if c.flagAll {
		return c.runDeleteAll(args)
	}

	if len(args) != 1 {
		c.ui.Output(c.Help(), terminal.WithErrorStyle())
		return 1
	}

	name := args[0]

	// Get our contexts
	if err := c.contextStorage.Delete(name); err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Context %q deleted.", name, terminal.WithSuccessStyle())
	return 0
}

func (c *ContextDeleteCommand) runDeleteAll(args []string) int {
	if len(args) > 0 {
		c.ui.Output(c.Flags().Help(), terminal.WithErrorStyle())
		return 1
	}

	// Delete all
	list, err := c.contextStorage.List()
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	for _, name := range list {
		if err := c.contextStorage.Delete(name); err != nil {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}
	}

	c.ui.Output("%d context(s) deleted.", len(list), terminal.WithSuccessStyle())
	return 0
}

func (c *ContextDeleteCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:   "all",
			Target: &c.flagAll,
			Usage: "Delete all contexts. If this is specified, NAME should " +
				"not be specified in the command arguments.",
		})
	})
}

func (c *ContextDeleteCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ContextDeleteCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ContextDeleteCommand) Synopsis() string {
	return "Delete a context."
}

func (c *ContextDeleteCommand) Help() string {
	return formatHelp(`
Usage: waypoint context delete [options] NAME

  Deletes a context. This will succeed if the context is already deleted.

` + c.Flags().Help())
}
