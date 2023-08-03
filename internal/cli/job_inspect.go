// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type JobInspectCommand struct {
	*baseCommand

	flagJson bool
}

func (c *JobInspectCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
	); err != nil {
		return 1
	}
	ctx := c.Ctx

	var jobId string
	if len(c.args) == 0 {
		c.ui.Output("Job ID required.\n\n%s", c.Help(), terminal.WithErrorStyle())
		return 1
	} else {
		jobId = c.args[0]
	}

	resp, err := c.project.Client().GetJob(ctx, &pb.GetJobRequest{
		JobId: jobId,
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.ui.Output("Job id not found: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle())
			return 1
		}

		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}
	if resp == nil {
		c.ui.Output("The requested job id %q was empty", jobId, terminal.WithWarningStyle())
		return 0
	}

	if c.flagJson {
		data, err := protojson.MarshalOptions{
			Indent: "\t",
		}.Marshal(resp)
		if err != nil {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}

		fmt.Println(string(data))
		return 0
	}

	var op string
	switch resp.Operation.(type) {
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
	switch resp.State {
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
	switch target := resp.TargetRunner.Target.(type) {
	case *pb.Ref_Runner_Any:
		targetRunner = "*"
	case *pb.Ref_Runner_Id:
		targetRunner = target.Id.Id
	}

	var queueTime string
	if resp.QueueTime != nil {
		queueTime = humanize.Time(resp.QueueTime.AsTime())
	}

	var completeTime string
	if resp.CompleteTime != nil {
		completeTime = humanize.Time(resp.CompleteTime.AsTime())
	}
	var cancelTime string
	if resp.CancelTime != nil {
		cancelTime = humanize.Time(resp.CancelTime.AsTime())
	}

	// job had an error! Let's show the message
	var errMsg string
	if resp.Error != nil {
		errMsg = resp.Error.Message
	}

	c.ui.Output("Job Configuration", terminal.WithHeaderStyle())
	appProjectName := resp.GetApplication().GetProject()
	if appProjectName == "" {
		appProjectName = "deleted"
	}
	appName := resp.GetApplication().GetApplication()
	if appName == "" {
		appName = "deleted"
	}
	workspaceName := resp.GetWorkspace().GetWorkspace()
	if workspaceName == "" {
		workspaceName = "deleted"
	}
	c.ui.NamedValues([]terminal.NamedValue{
		{
			Name: "ID", Value: resp.Id,
		},
		{
			Name: "Singleton ID", Value: resp.SingletonId,
		},
		{
			Name: "Operation", Value: op,
		},
		{
			Name: "Target Runner", Value: targetRunner,
		},
		{
			Name: "Workspace", Value: workspaceName,
		},
		{
			Name: "Project", Value: appProjectName,
		},
		{
			Name: "Application", Value: appName,
		},
	}, terminal.WithInfoStyle())

	if resp.Pipeline != nil {
		c.ui.Output("Pipeline Info", terminal.WithHeaderStyle())
		c.ui.NamedValues([]terminal.NamedValue{
			{
				Name: "Name", Value: resp.Pipeline.PipelineName,
			},
			{
				Name: "ID", Value: resp.Pipeline.PipelineId,
			},
			{
				Name: "Step", Value: resp.Pipeline.Step,
			},
			{
				Name: "Run Sequence", Value: resp.Pipeline.RunSequence,
			},
		}, terminal.WithInfoStyle())
	}

	c.ui.Output("Job Results", terminal.WithHeaderStyle())
	c.ui.NamedValues([]terminal.NamedValue{
		{
			Name: "State", Value: jobState,
		},
		{
			Name: "Queue Time", Value: queueTime,
		},
		{
			Name: "Complete Time", Value: completeTime,
		},
		{
			Name: "Cancel Time", Value: cancelTime,
		},
		{
			Name: "Error Message", Value: errMsg,
		},
	}, terminal.WithInfoStyle())

	return 0
}

func (c *JobInspectCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:    "json",
			Target:  &c.flagJson,
			Default: false,
			Usage:   "Output the list of jobs as json.",
		})
	})
}

func (c *JobInspectCommand) Synopsis() string {
	return "Inspect the details of a job by id in Waypoint"
}

func (c *JobInspectCommand) Help() string {
	return formatHelp(`
Usage: waypoint job inspect [options] <job-id>

  Inspect the details of a job by id in Waypoint server.

` + c.Flags().Help())
}
