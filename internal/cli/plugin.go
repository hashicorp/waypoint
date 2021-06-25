package cli

import (
	"context"
	"fmt"
	"github.com/hashicorp/waypoint-plugin-sdk"
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
		WithClient(false),
		WithNoAutoServer(),
		WithNoConfig(),
	); err != nil {
		return 1
	}
	args = flags.Args()

	pluginName := args[0]

	plugin, ok := plugin.Builtins[pluginName]
	if !ok {
		panic("no such plugin: " + pluginName)
	}

	// Run the plugin
	if !c.debugMode {
		sdk.Main(plugin...)
	} else {
		err := sdk.Debug(context.Background(), pluginName, plugin...)
		if err != nil {
			panic(fmt.Sprintf("Failed to launch plugin in debug mode: %s", err))
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
			Usage:   "set to true to run the plugin with support for debuggers like delve",
		})
	})
}

func (c *PluginCommand) Synopsis() string {
	return "Execute a built-in plugin."
}

func (c *PluginCommand) Help() string {
	return ""
}
