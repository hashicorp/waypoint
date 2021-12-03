package cli

import (
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

type TriggerApplyCommand struct {
	*baseCommand
}

func (c *TriggerApplyCommand) Run(args []string) int {
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

func (c *TriggerApplyCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, nil)
}

func (c *TriggerApplyCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *TriggerApplyCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *TriggerApplyCommand) Synopsis() string {
	return "Update a trigger URL configuration on Waypoint server"
}

func (c *TriggerApplyCommand) Help() string {
	return formatHelp(`
Usage: waypoint trigger apply [options]

  Update a trigger URL configuration on Waypoint Server.

` + c.Flags().Help())
}
