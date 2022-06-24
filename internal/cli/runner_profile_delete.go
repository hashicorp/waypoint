package cli

import (
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/posener/complete"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RunnerProfileDeleteCommand struct {
	*baseCommand
}

func (c *RunnerProfileDeleteCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoLocalServer(), // no auth in local mode
		WithNoConfig(),
	); err != nil {
		return 1
	}

	if len(c.args) == 0 {
		c.ui.Output("Runner profile id required.", terminal.WithErrorStyle())
		return 1
	}
	id := c.args[0]

	_, err := c.project.Client().DeleteOnDemandRunnerConfig(c.Ctx, &pb.DeleteOnDemandRunnerConfigRequest{
		Config: &pb.Ref_OnDemandRunnerConfig{
			Id: id,
		},
	})
	if err != nil && status.Code(err) != codes.NotFound {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	if status.Code(err) == codes.NotFound {
		c.ui.Output("runner profile not found", terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Runner profile deleted", terminal.WithHeaderStyle())

	return 0
}

func (c *RunnerProfileDeleteCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(sets *flag.Sets) {})
}

func (c *RunnerProfileDeleteCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *RunnerProfileDeleteCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *RunnerProfileDeleteCommand) Synopsis() string {
	return "Delete a runner profile."
}

func (c *RunnerProfileDeleteCommand) Help() string {
	return formatHelp(`
Usage: waypoint runner profile delete <id>

  Delete the specified runner profile.

`)
}
