// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
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
	ctx := c.Ctx

	var taskId string
	if len(c.args) == 0 && c.flagRunJobId == "" {
		c.ui.Output("Task ID or Run Job Id required.\n\n%s", c.Help(), terminal.WithErrorStyle())
		return 1
	} else {
		taskId = c.args[0]
	}

	if c.flagRunJobId != "" && taskId != "" {
		c.ui.Output("Both Run Job Id and Task Id was supplied, will look up by Task Id", terminal.WithWarningStyle())
	}

	var (
		taskReq *pb.GetTaskRequest
	)

	if taskId != "" {
		taskReq = &pb.GetTaskRequest{
			Ref: &pb.Ref_Task{
				Ref: &pb.Ref_Task_Id{
					Id: taskId,
				},
			},
		}
	} else if c.flagRunJobId != "" {
		taskReq = &pb.GetTaskRequest{
			Ref: &pb.Ref_Task{
				Ref: &pb.Ref_Task_JobId{
					JobId: c.flagRunJobId,
				},
			},
		}
	}

	taskResp, err := c.project.Client().GetTask(ctx, taskReq)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.ui.Output("Task not found: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle())
			return 1
		}

		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}
	if taskResp == nil {
		c.ui.Output("The requested task was empty", terminal.WithWarningStyle())
		return 0
	}

	if c.flagJson {
		var m jsonpb.Marshaler
		m.Indent = "\t"
		str, err := m.MarshalToString(taskResp)
		if err != nil {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}

		fmt.Println(str)
		return 0
	}

	taskState, ok := pb.Task_State_name[int32(taskResp.Task.JobState)]
	if !ok {
		c.ui.Output("Unrecognized task state defined for task: ", taskResp.Task.JobState, terminal.WithErrorStyle())
		return 1
	}

	var workspace, project, application string
	if taskResp.TaskJob.Workspace != nil {
		workspace = taskResp.TaskJob.Workspace.Workspace
	}
	if taskResp.TaskJob.Application != nil {
		project = taskResp.TaskJob.Application.Project
		application = taskResp.TaskJob.Application.Application
	}

	c.ui.Output("On-Demand Runner Task Configuration", terminal.WithHeaderStyle())
	c.ui.NamedValues([]terminal.NamedValue{
		{
			Name: "Task ID", Value: taskResp.Task.Id,
		},
		{
			Name: "Task State", Value: taskState,
		},
		{
			Name: "Task Resource", Value: taskResp.Task.ResourceName,
		},
		{
			Name: "Workspace", Value: workspace,
		},
		{
			Name: "Project", Value: project,
		},
		{
			Name: "Application", Value: application,
		},
	}, terminal.WithInfoStyle())

	c.ui.Output("Run Job Configuration", terminal.WithHeaderStyle())
	if err := c.FormatJob(taskResp.TaskJob); err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Watch Job Configuration", terminal.WithHeaderStyle())
	if err := c.FormatJob(taskResp.WatchJob); err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Start Job Configuration", terminal.WithHeaderStyle())
	if err := c.FormatJob(taskResp.StartJob); err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Stop Job Configuration", terminal.WithHeaderStyle())
	if err := c.FormatJob(taskResp.StopJob); err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	return 0
}

// FormatJob takes a Job proto message and formats it into something nicer
// to read to the user.
func (c *TaskInspectCommand) FormatJob(job *pb.Job) error {
	if job == nil {
		return nil
	}

	var op string
	switch job.Operation.(type) {
	case *pb.Job_Noop_:
		op = "Noop"
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
	case *pb.Job_DestroyProject:
		op = "DestroyProject"
	default:
		op = "Unknown"
	}

	var jobState string
	switch job.State {
	case pb.Job_UNKNOWN:
		jobState = "Unknown"
	case pb.Job_QUEUED:
		jobState = "Queued"
	case pb.Job_WAITING:
		jobState = "Waiting"
	case pb.Job_RUNNING:
		jobState = "Running"
	case pb.Job_ERROR:
		jobState = "Error"
	case pb.Job_SUCCESS:
		jobState = "Success"
	default:
		jobState = "Unknown"
	}

	var targetRunner string
	switch target := job.TargetRunner.Target.(type) {
	case *pb.Ref_Runner_Any:
		targetRunner = "*"
	case *pb.Ref_Runner_Id:
		targetRunner = target.Id.Id
	}

	var completeTime string
	if time, err := ptypes.Timestamp(job.CompleteTime); err == nil {
		completeTime = humanize.Time(time)
	}
	var cancelTime string
	if time, err := ptypes.Timestamp(job.CancelTime); err == nil {
		cancelTime = humanize.Time(time)
	}

	// job had an error! Let's show the message
	var errMsg string
	if job.Error != nil {
		errMsg = job.Error.Message
	}

	c.ui.Output("Job Configuration:", terminal.WithInfoStyle())
	c.ui.NamedValues([]terminal.NamedValue{
		{
			Name: "Job ID", Value: job.Id,
		},
		{
			Name: "Singleton ID", Value: job.SingletonId,
		},
		{
			Name: "Operation", Value: op,
		},
		{
			Name: "Target Runner", Value: targetRunner,
		},
		{
			Name: "Workspace", Value: job.Workspace.Workspace,
		},
		{
			Name: "Project", Value: job.Application.Project,
		},
		{
			Name: "Application", Value: job.Application.Application,
		},
	}, terminal.WithInfoStyle())

	c.ui.Output("\nJob Results:", terminal.WithInfoStyle())
	c.ui.NamedValues([]terminal.NamedValue{
		{
			Name: "State", Value: jobState,
		},
		{
			Name: "Complete Time", Value: completeTime,
		},
		{
			Name: "Cancel Time", Value: cancelTime,
		},
		{
			Name: "Error Messsage", Value: errMsg,
		},
	}, terminal.WithInfoStyle())

	return nil
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

  Inspect an On-Demand Runner Tasks and all of the jobs associated with it.

` + c.Flags().Help())
}
