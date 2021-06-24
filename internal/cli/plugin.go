package cli

import (
	"context"
	"fmt"
	"github.com/hashicorp/waypoint-plugin-sdk"
	"github.com/hashicorp/waypoint/internal/plugin"
	"os"
)

type PluginCommand struct {
	*baseCommand
}

func (c *PluginCommand) Run(args []string) int {
	plugin, ok := plugin.Builtins[args[0]]
	if !ok {
		panic("no such plugin: " + args[0])
	}

	debug := os.Getenv("WAYPOINT_PLUGIN_DEBUG") != ""

	// Run the plugin
	if !debug {
		sdk.Main(plugin...)
	} else {
		err := sdk.Debug(context.Background(), "pack", plugin...)
		if err != nil {
			panic(fmt.Sprintf("Failed to launch plugin in debug mode: %v", err))
		}
	}

	return 0
}

func (c *PluginCommand) Synopsis() string {
	return "Execute a built-in plugin."
}

func (c *PluginCommand) Help() string {
	return ""
}
