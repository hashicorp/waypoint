package cli

import (
	"context"
	"fmt"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/posener/complete"
)

type StatusCheckCommand struct {
	*baseCommand
}

func (c *StatusCheckCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
		WithSingleApp(),
	); err != nil {
		return 1
	}

	err := c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
		deployments, err := c.project.Client().ListDeployments(c.Ctx, &pb.ListDeploymentsRequest{
			Application:   c.refApp,
			Workspace:     c.refWorkspace,
			PhysicalState: pb.Operation_CREATED,
		})
		if err != nil {
			return fmt.Errorf("unable to get deployments for this project")
		}

		if len(deployments.Deployments) == 0 {
			return fmt.Errorf("this project has no deployments in this workspace")
		}
		// Query a check for each deployment
		for _, dep := range deployments.Deployments {
			resp, err := app.StatusReport(c.Ctx, &pb.Job_StatusReportOp{
				Target: &pb.Job_StatusReportOp_Deployment{
					Deployment: dep,
				},
			})
			if err != nil {
				return err
			}
			c.ui.Output(
				fmt.Sprintf("%s v%d is %s: %s",
					resp.StatusReport.Application.Application,
					dep.Sequence,
					resp.StatusReport.Health.HealthStatus,
					resp.StatusReport.Health.HealthMessage,
				), terminal.WithInfoStyle())
		}

		return nil
	})

	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Status report successfully requested",
		terminal.WithSuccessStyle(),
	)
	return 0
}

func (c *StatusCheckCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *StatusCheckCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *StatusCheckCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *StatusCheckCommand) Synopsis() string {
	return "Performs an health check for the current project."
}

func (c *StatusCheckCommand) Help() string {
	return formatHelp(`
Usage: waypoint status check

  Trigger a status check on the project.

  This command triggers a status check on the project,
  so that you can on a subsequent run of "waypoint status"
  see the updated status.

`)
}
