package cli

import (
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

type TriggerCreateCommand struct {
	*baseCommand
}

func (c *TriggerCreateCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	flagSet := c.Flags()
	if err := c.Init(
		WithArgs(args),
		WithFlags(flagSet),
		WithNoConfig(),
	); err != nil {
		return 1
	}

	return 0
}

func (c *TriggerCreateCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, nil)
}

func (c *TriggerCreateCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *TriggerCreateCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *TriggerCreateCommand) Synopsis() string {
	return "Generate a trigger URL and register it to Waypoint server"
}

func (c *TriggerCreateCommand) Help() string {
	return formatHelp(`
Usage: waypoint trigger create [options]

  Create and register a trigger URL to Waypoint Server.

` + c.Flags().Help())
}
