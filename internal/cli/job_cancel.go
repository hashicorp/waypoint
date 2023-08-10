// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cli

import (
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type JobCancelCommand struct {
	*baseCommand

	flagForce bool
}

func (c *JobCancelCommand) Run(args []string) int {
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

	if c.flagForce {
		c.ui.Output("You requested to use force to cancel a job! Be aware that this "+
			"operation is dangerous and could result in some bad behavior or failure modes in Waypoint.",
			terminal.WithWarningStyle())
		c.ui.Output("If this is not your intention, ctrl-c now! The CLI will sleep for 3 seconds...",
			terminal.WithWarningStyle())
		time.Sleep(3 * time.Second)
	}

	sg := c.ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Cancelling job %q", jobId)
	defer func() { s.Abort() }()

	_, err := c.project.Client().CancelJob(ctx, &pb.CancelJobRequest{
		JobId: jobId,
		Force: c.flagForce,
	})
	if err != nil {
		s.Update("Failed to marked job %q for cancellation", jobId)
		s.Status(terminal.StatusError)
		s.Done()

		if status.Code(err) == codes.NotFound {
			c.ui.Output("Job id not found: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle())
			return 1
		}

		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	if !c.flagForce {
		job, err := c.project.Client().GetJob(ctx, &pb.GetJobRequest{JobId: jobId})
		if err != nil {
			c.ui.Output("Job not found: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle())
			return 1
		}
		s.Update("Marked job %q for cancellation", jobId)
		if job.State == pb.Job_RUNNING {
			c.ui.Output("Waypoint will gracefully cancel the requested job and wait for any\ndownstream listeners to close the connection. This could take a while.")
		}
	} else {
		s.Update("Forcefully marked job %q for cancellation", jobId)
		s.Status(terminal.StatusWarn)
	}
	s.Done()

	return 0
}

func (c *JobCancelCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:    "dangerously-force",
			Target:  &c.flagForce,
			Default: false,
			Usage: "Will forcefully cancel the job. This will immediately mark the " +
				"job as complete in the server, regardless of the real job status. This " +
				"may leave dangling resources or cause concurrency issues if the underlying " +
				"job doesn't gracefully cancel. USE WITH CAUTION.",
		})
	})
}

func (c *JobCancelCommand) Synopsis() string {
	return "Cancel a running a job by id"
}

func (c *JobCancelCommand) Help() string {
	return formatHelp(`
Usage: waypoint job cancel [options] <job-id>

  Cancel a running job by id from Waypoint server.

` + c.Flags().Help())
}
