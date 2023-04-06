// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type TaskCancelCommand struct {
	*baseCommand

	flagRunJobId string
}

func (c *TaskCancelCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
	); err != nil {
		return 1
	}
	ctx := c.Ctx

	var taskId string
	if len(c.args) > 0 {
		taskId = c.args[0]
	}

	if taskId == "" && c.flagRunJobId == "" {
		c.ui.Output("Task Id or Run Job Id required.\n\n%s", c.Help(), terminal.WithErrorStyle())
		return 1
	} else if c.flagRunJobId != "" && taskId != "" {
		c.ui.Output("Both Run Job Id and Task Id was supplied, will look up by Task Id", terminal.WithWarningStyle())
	}

	taskReq := &pb.CancelTaskRequest{Ref: &pb.Ref_Task{}}

	if taskId != "" {
		taskReq.Ref.Ref = &pb.Ref_Task_Id{
			Id: taskId,
		}
	} else if c.flagRunJobId != "" {
		taskReq.Ref.Ref = &pb.Ref_Task_JobId{
			JobId: c.flagRunJobId,
		}
	}

	_, err := c.project.Client().CancelTask(ctx, taskReq)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.ui.Output("Task not found: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle())
			return 1
		}

		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Task %q has been requested to be cancelled", taskId, terminal.WithWarningStyle())
	return 0
}

func (c *TaskCancelCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")

		f.StringVar(&flag.StringVar{
			Name:    "run-job-id",
			Target:  &c.flagRunJobId,
			Default: "",
			Usage:   "Cancel a Task by Run Job Id.",
		})
	})
}

func (c *TaskCancelCommand) Synopsis() string {
	return "Cancel an On-Demand Runner Task running in Waypoint"
}

func (c *TaskCancelCommand) Help() string {
	return formatHelp(`
Usage: waypoint task cancel [options] <task-id>

  Cancel an On-Demand Runner Task and all of the jobs associated with it.

` + c.Flags().Help())
}
