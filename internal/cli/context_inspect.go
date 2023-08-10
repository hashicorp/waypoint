// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cli

import (
	"encoding/json"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/posener/complete"
)

type ContextInspectCommand struct {
	*baseCommand

	flagJson bool
}

func (c *ContextInspectCommand) Run(args []string) int {
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

	if len(c.args) >= 1 {
		cc, err := c.contextStorage.Load(c.args[0])
		if err != nil {
			c.ui.Output("Error loading context '%s': %s", c.args[0], err)
			return 1
		}

		if c.flagJson {
			data, err := json.MarshalIndent(cc.Server, "", "  ")
			if err != nil {
				c.ui.Output("Error rendering json: %s", err)
				return 1
			}

			c.ui.Output(string(data))
			return 0
		}

		c.ui.Output("Context Info:", terminal.WithHeaderStyle())

		c.ui.NamedValues([]terminal.NamedValue{
			{
				Name: "address", Value: cc.Server.Address,
			},
			{
				Name: "address internal", Value: cc.Server.AddressInternal,
			},
			{
				Name: "tls", Value: cc.Server.Tls,
			},
			{
				Name: "tls skip verify", Value: cc.Server.TlsSkipVerify,
			},
			{
				Name: "require auth", Value: cc.Server.RequireAuth,
			},
			{
				Name: "platform", Value: cc.Server.Platform,
			},
		}, terminal.WithInfoStyle())

		c.ui.Output("Workspace Info:", terminal.WithHeaderStyle())
		workspace := cc.Workspace
		if workspace == "" {
			workspace = "default"
		}

		c.ui.NamedValues([]terminal.NamedValue{
			{
				Name: "Name", Value: workspace,
			},
		}, terminal.WithInfoStyle())

		return 0
	}

	def, err := c.contextStorage.Default()
	if err != nil {
		def = "<unknown>"
	}

	if c.flagJson {
		data, err := json.MarshalIndent(map[string]interface{}{
			"config_path":     c.homeConfigPath,
			"default_context": def,
		}, "", "  ")
		if err != nil {
			c.ui.Output("Error rendering json: %s", err)
			return 1
		}

		c.ui.Output(string(data))
		return 0
	}

	c.ui.Output("Context Settings:", terminal.WithHeaderStyle())

	if def == "" {
		def = "<unset>"
	}

	c.ui.NamedValues([]terminal.NamedValue{
		{
			Name: "config path", Value: c.homeConfigPath,
		},
		{
			Name: "default context", Value: def,
		},
	}, terminal.WithInfoStyle())

	return 0
}

func (c *ContextInspectCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")

		f.BoolVar(&flag.BoolVar{
			Name:    "json",
			Target:  &c.flagJson,
			Usage:   "Output information in JSON format",
			Default: false,
		})
	})
}

func (c *ContextInspectCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ContextInspectCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ContextInspectCommand) Synopsis() string {
	return "Output context info."
}

func (c *ContextInspectCommand) Help() string {
	return formatHelp(`
Usage: waypoint context inspect [<name>]

  Output information about a waypoint context or general context info.

` + c.Flags().Help())
}
