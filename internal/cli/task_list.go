package cli

import (
	"fmt"

	"github.com/golang/protobuf/jsonpb"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
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
	ctx := c.Ctx

	resp, err := c.project.Client().ListTask(ctx, &pb.ListTaskRequest{})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	if len(resp.Tasks) == 0 {
		return 0
	}
	tasks := resp.Tasks

	// Reverse task list if requested
	if c.flagDesc {
		// reverse in place
		for i, j := 0, len(tasks)-1; i < j; i, j = i+1, j-1 {
			tasks[i], tasks[j] = tasks[j], tasks[i]
		}
	}

	// limit to the first n jobs
	if c.flagLimit > 0 && c.flagLimit <= len(tasks) {
		tasks = tasks[:c.flagLimit]
	}

	if c.flagJson {
		var m jsonpb.Marshaler
		m.Indent = "\t"
		for _, t := range resp.Tasks {
			str, err := m.MarshalToString(t)
			if err != nil {
				c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
				return 1
			}

			fmt.Println(str)
		}
		return 0
	}

	c.ui.Output("Waypoint On-Demand Runner Tasks", terminal.WithHeaderStyle())

	tblHeaders := []string{"ID", "Start Job Id", "Run Job Id", "Stop Job Id"}
	tbl := terminal.NewTable(tblHeaders...)

	for _, t := range tasks {
		var taskJobId, startJobId, stopJobId string
		if t.TaskJob != nil {
			taskJobId = t.TaskJob.Id
		}
		if t.StartJob != nil {
			startJobId = t.StartJob.Id
		}
		if t.StopJob != nil {
			stopJobId = t.StopJob.Id
		}

		tblColumn := []string{
			t.Task.Id,
			startJobId,
			taskJobId,
			stopJobId,
		}

		tbl.Rich(tblColumn, nil)
	}

	c.ui.Table(tbl)

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

  List all known On-Demand Runner Tasks from Waypoint server. Each task is a
  Waypoint Job tuple made up of a StartTask, RunTask, and StopTask.

` + c.Flags().Help())
}
