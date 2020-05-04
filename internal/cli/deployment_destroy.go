package cli

import (
	"context"
	"strings"

	"github.com/posener/complete"

	"github.com/hashicorp/waypoint/internal/core"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

type DeploymentDestroyCommand struct {
	*baseCommand
}

func (c *DeploymentDestroyCommand) Run(args []string) int {
	flags := c.Flags()

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(flags),
		WithSingleApp(),
	); err != nil {
		return 1
	}

	args = flags.Args()

	client := c.project.Client()
	err := c.DoApp(c.Ctx, func(ctx context.Context, app *core.App) error {
		// Get the latest deployment
		resp, err := client.ListDeployments(ctx, &pb.ListDeploymentsRequest{
			Limit:     1,
			Order:     pb.ListDeploymentsRequest_COMPLETE_TIME,
			OrderDesc: true,
		})
		if err != nil {
			app.UI.Output(err.Error(), terminal.WithErrorStyle())
			return ErrSentinel
		}
		if len(resp.Deployments) == 0 {
			app.UI.Output("No successful deployments found.", terminal.WithErrorStyle())
			return ErrSentinel
		}

		if err := app.DestroyDeploy(ctx, resp.Deployments[0]); err != nil {
			app.UI.Output("Error destroying the deployment: %s", err.Error(), terminal.WithErrorStyle())
			return ErrSentinel
		}

		return nil
	})
	if err != nil {
		return 1
	}

	return 0
}

func (c *DeploymentDestroyCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *DeploymentDestroyCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *DeploymentDestroyCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *DeploymentDestroyCommand) Synopsis() string {
	return "Destroy one or more deployments."
}

func (c *DeploymentDestroyCommand) Help() string {
	helpText := `
Usage: waypoint deployment destroy [options] [id...]

  Destroy one or more deployments. This will "undeploy" this specific
  instance of an application.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}
