package cli

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type JobCancelCommand struct {
	*baseCommand

	flagJson bool
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

	sg := c.ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Cancelling job %q", jobId)
	defer func() { s.Abort() }()

	_, err := c.project.Client().CancelJob(ctx, &pb.CancelJobRequest{
		JobId: jobId,
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

	s.Update("Marked job %q for cancellation", jobId)
	s.Done()

	return 0
}

func (c *JobCancelCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
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
