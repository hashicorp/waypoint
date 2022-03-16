package cli

import (
	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

type TaskListCommand struct {
	*baseCommand

	flagJson  bool
	flagLimit int
	flagDesc  bool
}

func (c *TaskListCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
	); err != nil {
		return 1
	}

	return 0
}

func (c *TaskListCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:    "desc",
			Target:  &c.flagDesc,
			Default: false,
			Usage:   "Output the list of Tasks from newest to oldest.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "json",
			Target:  &c.flagJson,
			Default: false,
			Usage:   "Output the list of Tasks as json.",
		})

		f.IntVar(&flag.IntVar{
			Name:    "limit",
			Target:  &c.flagLimit,
			Default: 0,
			Usage:   "If set, will limit the number of Tasks to list.",
		})
	})
}

func (c *TaskListCommand) Synopsis() string {
	return "List all On-Demand Runner Tasks in Waypoint"
}

func (c *TaskListCommand) Help() string {
	return formatHelp(`
Usage: waypoint task list [options]

  List all known On-Demand Runner Tasks from Waypoint server.

` + c.Flags().Help())
}
