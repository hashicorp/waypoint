package cli

import (
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

type TriggerDeleteCommand struct {
	*baseCommand

	flagTriggerName string
	flagTriggerId   string
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
	ctx := c.Ctx

	_, err := c.project.Client().DeleteTrigger(ctx, &pb.DeleteTriggerRequest{
		Ref: &pb.Ref_Trigger{
			Name: c.flagTriggerName,
			Id:   c.flagTriggerId,
		},
	})
	if err != nil {
		c.ui.Output(
			"Error deleting trigger: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	c.ui.Output("Trigger deleted", terminal.WithSuccessStyle())

	return 0
}

func (c *TriggerDeleteCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")

		f.StringVar(&flag.StringVar{
			Name:   "name",
			Target: &c.flagTriggerName,
			Usage:  "The name of the trigger URL to delete.",
		})

		f.StringVar(&flag.StringVar{
			Name:   "id",
			Target: &c.flagTriggerId,
			Usage:  "The id of the trigger URL to delete.",
		})
	})
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
