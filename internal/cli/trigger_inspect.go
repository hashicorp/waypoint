package cli

import (
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

type TriggerInspectCommand struct {
	*baseCommand

	flagTriggerName string
	flagTriggerId   string
}

func (c *TriggerInspectCommand) Run(args []string) int {
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

func (c *TriggerInspectCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")

		f.StringVar(&flag.StringVar{
			Name:   "name",
			Target: &c.flagTriggerName,
			Usage:  "The name of the trigger URL to inspect.",
		})

		f.StringVar(&flag.StringVar{
			Name:   "id",
			Target: &c.flagTriggerId,
			Usage:  "The id of the trigger URL to inspect.",
		})
	})
}

func (c *TriggerInspectCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *TriggerInspectCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *TriggerInspectCommand) Synopsis() string {
	return "Inspect a trigger URL from Waypoint server"
}

func (c *TriggerInspectCommand) Help() string {
	return formatHelp(`
Usage: waypoint trigger inspect [options]

  Inspect a trigger URL from Waypoint Server.

` + c.Flags().Help())
}
