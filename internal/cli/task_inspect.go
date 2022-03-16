package cli

import (
	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

type TaskInspectCommand struct {
	*baseCommand

	flagJson     bool
	flagRunJobId string
}

func (c *TaskInspectCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
	); err != nil {
		return 1
	}

	return 0
}

func (c *TaskInspectCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:    "json",
			Target:  &c.flagJson,
			Default: false,
			Usage:   "Output the list of jobs as json.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "run-job-id",
			Target:  &c.flagRunJobId,
			Default: "",
			Usage:   "Look up a Task by Run Job Id.",
		})
	})
}

func (c *TaskInspectCommand) Synopsis() string {
	return "Inspect an On-Demand Runner Task from Waypoint"
}

func (c *TaskInspectCommand) Help() string {
	return formatHelp(`
Usage: waypoint task inspect [options] <task-id>

  List all known On-Demand Runner Tasks from Waypoint server.

` + c.Flags().Help())
}
