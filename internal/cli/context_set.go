// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"errors"

	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

type ContextSetCommand struct {
	*baseCommand
}

func (c *ContextSetCommand) Run(args []string) int {
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

	if c.flagWorkspace == "" {
		// Require one argument
		c.ui.Output(c.Help(), terminal.WithErrorStyle())
		return 1
	}

	contextName, err := c.contextStorage.Default()
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	if contextName == "" || contextName == "-" {
		c.ui.Output(
			clierrors.Humanize(errors.New("no default context exists")),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	// Load it and set it.
	cfg, err := c.contextStorage.Load(contextName)
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	// set new workspace
	cfg.Workspace = c.flagWorkspace

	// store updated context
	if err := c.contextStorage.Set(contextName, cfg); err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Context %q (%s) updated to use %s workspace.", contextName, cfg.Server.Address, cfg.Workspace, terminal.WithSuccessStyle())
	return 0
}

func (c *ContextSetCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *ContextSetCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ContextSetCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ContextSetCommand) Synopsis() string {
	return "Set a property of the current context."
}

func (c *ContextSetCommand) Help() string {
	return formatHelp(`
Usage: waypoint context set [options]

  Sets a property of the current context. The only property supported at this
  time is -workspace.

  To use this command, use the global -workspace flag to set the default
  workspace for the current context.

  To restore this CLI context to use the default workspace, use
  -workspace=default

` + c.Flags().Help())
}
