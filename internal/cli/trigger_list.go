package cli

import (
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

type TriggerListCommand struct {
	*baseCommand
}

func (c *TriggerListCommand) Run(args []string) int {
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

func (c *TriggerListCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, nil)
}

func (c *TriggerListCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *TriggerListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *TriggerListCommand) Synopsis() string {
	return "List trigger URL configurations on Waypoint server"
}

func (c *TriggerListCommand) Help() string {
	return formatHelp(`
Usage: waypoint trigger list [options]

  List trigger URL configurations on Waypoint Server.

` + c.Flags().Help())
}
