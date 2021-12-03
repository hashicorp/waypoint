package cli

import (
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

type TriggerDeleteCommand struct {
	*baseCommand
}

func (c *TriggerDeleteCommand) Run(args []string) int {
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

func (c *TriggerDeleteCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, nil)
}

func (c *TriggerDeleteCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *TriggerDeleteCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *TriggerDeleteCommand) Synopsis() string {
	return "Delete a trigger URL from Waypoint server"
}

func (c *TriggerDeleteCommand) Help() string {
	return formatHelp(`
Usage: waypoint trigger delete [options]

  Delete a trigger URL from Waypoint Server.

` + c.Flags().Help())
}
