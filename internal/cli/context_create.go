// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cli

import (
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clicontext"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

type ContextCreateCommand struct {
	*baseCommand

	flagConfig     clicontext.Config
	flagSetDefault bool
}

func (c *ContextCreateCommand) Run(args []string) int {
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

	// Require one argument
	if len(args) != 1 {
		c.ui.Output(c.Help(), terminal.WithErrorStyle())
		return 1
	}

	name := args[0]

	// Set the context
	if err := c.contextStorage.Set(name, &c.flagConfig); err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	if c.flagSetDefault {
		if err := c.contextStorage.SetDefault(name); err != nil {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}
	}

	c.ui.Output("Context %q created.", name, terminal.WithSuccessStyle())
	return 0
}

func (c *ContextCreateCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:    "set-default",
			Target:  &c.flagSetDefault,
			Default: true,
			Usage:   "Set this context as the new default for the CLI.",
		})
		f.StringVar(&flag.StringVar{
			Name:   "server-addr",
			Target: &c.flagConfig.Server.Address,
			Usage:  "Address for the server.",
		})
		f.StringVar(&flag.StringVar{
			Name:   "server-auth-token",
			Target: &c.flagConfig.Server.AuthToken,
			Usage:  "Authentication token to use to connect to the server.",
		})
		f.StringVar(&flag.StringVar{
			Name:    "server-platform",
			Target:  &c.flagConfig.Server.Platform,
			Default: "n/a",
			Usage:   "The current platform that Waypoint server is running on.",
		})
		f.BoolVar(&flag.BoolVar{
			Name:    "server-tls",
			Target:  &c.flagConfig.Server.Tls,
			Usage:   "If true, will connect to the server over TLS.",
			Default: true,
		})
		f.BoolVar(&flag.BoolVar{
			Name:   "server-tls-skip-verify",
			Target: &c.flagConfig.Server.TlsSkipVerify,
			Usage:  "If true, will not validate TLS cert presented by the server.",
		})
		f.BoolVar(&flag.BoolVar{
			Name:   "server-require-auth",
			Target: &c.flagConfig.Server.RequireAuth,
			Usage:  "If true, will send authentication details.",
		})
	})
}

func (c *ContextCreateCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ContextCreateCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ContextCreateCommand) Synopsis() string {
	return "Create a context."
}

func (c *ContextCreateCommand) Help() string {
	return formatHelp(`
Usage: waypoint context create [options] NAME

  Creates a new context.

` + c.Flags().Help())
}
