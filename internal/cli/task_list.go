package cli

import (
	"fmt"

	"github.com/dustin/go-humanize"
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

	tblHeaders := []string{"ID", "Run Job Operation", "Task State", "Time Completed"}
	tbl := terminal.NewTable(tblHeaders...)

	for _, t := range tasks {
		var op string
		// Job_Noop seems to be missing the isJob_operation method
		switch t.TaskJob.Operation.(type) {
		case *pb.Job_Build:
			op = "Build"
		case *pb.Job_Push:
			op = "Push"
		case *pb.Job_Deploy:
			op = "Deploy"
		case *pb.Job_Destroy:
			op = "Destroy"
		case *pb.Job_Release:
			op = "Release"
		case *pb.Job_Validate:
			op = "Validate"
		case *pb.Job_Auth:
			op = "Auth"
		case *pb.Job_Docs:
			op = "Docs"
		case *pb.Job_ConfigSync:
			op = "ConfigSync"
		case *pb.Job_Exec:
			op = "Exec"
		case *pb.Job_Up:
			op = "Up"
		case *pb.Job_Logs:
			op = "Logs"
		case *pb.Job_QueueProject:
			op = "QueueProject"
		case *pb.Job_Poll:
			op = "Poll"
		case *pb.Job_StatusReport:
			op = "StatusReport"
		case *pb.Job_StartTask:
			op = "StartTask"
		case *pb.Job_StopTask:
			op = "StopTask"
		case *pb.Job_Init:
			op = "Init"
		default:
			op = "Unknown"
		}

		var completeTime string
		if t.StopJob.CompleteTime != nil {
			completeTime = humanize.Time(t.StopJob.CompleteTime.AsTime())
		}

		tblColumn := []string{
			t.Task.Id,
			op,
			pb.Task_State_name[int32(t.Task.JobState)],
			completeTime,
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
