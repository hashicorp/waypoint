// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"context"
	"fmt"

	sdk "github.com/hashicorp/waypoint-plugin-sdk"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/plugin"
)

type PluginCommand struct {
	*baseCommand
	debugMode bool
}

func (c *PluginCommand) Run(args []string) int {
	flags := c.Flags()

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(flags),
		WithNoClient(),
		WithNoLocalServer(),
		WithNoConfig(),
	); err != nil {
		return 1
	}
	args = flags.Args()

	if len(args) == 0 {
		c.ui.Output("No argument specified.\n"+c.Flags().Help(), terminal.WithErrorStyle())
		return 1
	}

	pluginName := args[0]

	plugin, ok := plugin.Builtins[pluginName]
	if !ok {
		c.ui.Output("No such plugin: "+pluginName, terminal.WithErrorStyle())
		return 1
	}

	// Run the plugin
	if !c.debugMode {
		sdk.Main(plugin...)
	} else {
		err := sdk.Debug(context.Background(), pluginName, plugin...)
		if err != nil {
			c.ui.Output(fmt.Sprintf("Failed to launch plugin in debug mode: %s", err), terminal.WithErrorStyle())
			return 1
		}
	}

	return 0
}

func (c *PluginCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:    "debug",
			Default: false,
			Target:  &c.debugMode,
			Usage:   "Set to true to run the plugin with support for debuggers like delve.",
		})
	})
}

func (c *PluginCommand) Synopsis() string {
	return "Execute a built-in plugin."
}

func (c *PluginCommand) Help() string {
	return formatHelp(`
Usage: waypoint plugin [options] <plugin>

  Runs a specified plugin directly.

` + c.Flags().Help())
}
