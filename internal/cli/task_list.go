package cli

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/dustin/go-humanize"
	"github.com/golang/protobuf/jsonpb"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type TaskListCommand struct {
	*baseCommand

	flagJson       bool
	flagLimit      int
	flagDesc       bool
	flagTaskStates []string
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

	listTaskReq := &pb.ListTaskRequest{}

	if len(c.flagTaskStates) > 0 {
		var ts []pb.Task_State
		// convert to int32 const from string value
		for _, state := range c.flagTaskStates {
			s, ok := pb.Task_State_value[state]
			if !ok {
				// this shouldn't happen given the State flag is an enum var, but protect
				// against it anyway
				c.ui.Output("Undefined task job state value: "+state, terminal.WithErrorStyle())
				return 1
			} else {
				ts = append(ts, pb.Task_State(s))
			}
		}

		listTaskReq.TaskState = ts
	}

	resp, err := c.project.Client().ListTask(ctx, listTaskReq)
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
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].StartJob.QueueTime.AsTime().Before(tasks[j].StartJob.QueueTime.AsTime())
		})
	} else {
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].StartJob.QueueTime.AsTime().After(tasks[j].StartJob.QueueTime.AsTime())
		})
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

	tblHeaders := []string{"ID", "Run Job Operation", "Pipeline", "Task State", "Project", "Time Created", "Time Completed"}
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
		case *pb.Job_WatchTask:
			op = "WatchTask"
		case *pb.Job_Init:
			op = "Init"
		case *pb.Job_PipelineStep:
			op = "PipelineStep"
		default:
			op = "Unknown"
			c.Log.Debug("encountered unsupported task operation", "op", t.TaskJob.Operation)
		}

		var project string
		if t.TaskJob.Application != nil {
			project = t.TaskJob.Application.Project
		}

		var queueTime string
		if t.StartJob.QueueTime != nil {
			queueTime = humanize.Time(t.StartJob.QueueTime.AsTime())
		}

		var completeTime string
		if t.StopJob.CompleteTime != nil {
			completeTime = humanize.Time(t.StopJob.CompleteTime.AsTime())
		}

		pipeline := ""
		if t.TaskJob.Pipeline != nil {
			pipeline = t.TaskJob.Pipeline.Pipeline + "[run: " + strconv.FormatUint(t.TaskJob.Pipeline.RunSequence, 10) + "]" + "[step: " + t.TaskJob.Pipeline.Step + "]"
		}

		tblColumn := []string{
			t.Task.Id,
			op,
			pipeline,
			pb.Task_State_name[int32(t.Task.JobState)],
			project,
			queueTime,
			completeTime,
		}

		tbl.Rich(tblColumn, nil)
	}

	c.ui.Table(tbl)

	return 0
}

var taskStateValues = []string{pb.Task_State_name[0],
	pb.Task_State_name[1],
	pb.Task_State_name[2],
	pb.Task_State_name[3],
	pb.Task_State_name[4],
	pb.Task_State_name[5],
	pb.Task_State_name[6],
	pb.Task_State_name[7],
	pb.Task_State_name[8],
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

		f.EnumVar(&flag.EnumVar{
			Name:   "state",
			Target: &c.flagTaskStates,
			Values: taskStateValues,
			Usage:  "List Tasks that only match the requested state. Can be repeated multiple times.",
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
