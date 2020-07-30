package cli

import (
	"context"
	"strings"

	"github.com/posener/complete"

	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
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
	if len(args) != 1 {
		c.ui.Output("you must supply a single deployment ID", terminal.WithErrorStyle)
		return 1
	}

	client := c.project.Client()
	err := c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
		// Get the deployment
		deployment, err := client.GetDeployment(ctx, &pb.GetDeploymentRequest{
			DeploymentId: args[0],
		})
		if err != nil {
			app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}

		// Can't destroy a deployment that was not successful
		if deployment.Status.GetState() != pb.Status_SUCCESS {
			app.UI.Output("Cannot destroy deployment that is not successful", terminal.WithErrorStyle())
			return ErrSentinel
		}

		if err := app.DestroyDeploy(ctx, &pb.Job_DestroyDeployOp{
			Deployment: deployment,
		}); err != nil {
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
