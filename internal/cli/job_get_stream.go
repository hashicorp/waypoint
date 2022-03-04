package cli

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type JobGetStreamCommand struct {
	*baseCommand

	flagJson bool
}

func (c *JobGetStreamCommand) Run(args []string) int {
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

	_, err := c.project.Client().GetJobStream(ctx, &pb.GetJobStreamRequest{
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

	// TODO(briancain): process and print terminal events like `internal/client/job.go`
	c.ui.Output("Job stream is not implemented yet!", terminal.WithWarningStyle())

	return 0
}

func (c *JobGetStreamCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		//f := set.NewSet("Command Options")
	})
}

func (c *JobGetStreamCommand) Synopsis() string {
	return "Attach a local CLI to a job stream by id"
}

func (c *JobGetStreamCommand) Help() string {
	return formatHelp(`
Usage: waypoint job get-stream [options] <job-id>

  Connects the local CLI to an active job stream.

` + c.Flags().Help())
}
