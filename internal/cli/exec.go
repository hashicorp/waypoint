package cli

import (
	"context"
	"os"

	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/server/execclient"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

type ExecCommand struct {
	*baseCommand
}

func (c *ExecCommand) Run(args []string) int {
	flagSet := c.Flags()

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(flagSet),
		WithSingleApp(),
	); err != nil {
		return 1
	}

	args = flagSet.Args()

	var exitCode int
	client := c.project.Client()
	err := c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
		// Get the latest deployment
		resp, err := client.ListDeployments(ctx, &pb.ListDeploymentsRequest{
			Application: app.Ref(),
			Order: &pb.OperationOrder{
				Limit: 1,
				Order: pb.OperationOrder_COMPLETE_TIME,
				Desc:  true,
			},
			PhysicalState: pb.Operation_CREATED,
		})
		if err != nil {
			app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}
		if len(resp.Deployments) == 0 {
			app.UI.Output("No successful deployments found.", terminal.WithErrorStyle())
			return ErrSentinel
		}

		client := &execclient.Client{
			Logger:        c.Log,
			UI:            c.ui,
			Context:       ctx,
			Client:        client,
			DeploymentId:  resp.Deployments[0].Id,
			DeploymentSeq: resp.Deployments[0].Sequence,
			Args:          args,
			Stdin:         os.Stdin,
			Stdout:        os.Stdout,
			Stderr:        os.Stderr,
		}

		exitCode, err = client.Run()
		if err != nil {
			app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}

		return nil
	})
	if err != nil {
		return 1
	}

	return exitCode
}

func (c *ExecCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *ExecCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ExecCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ExecCommand) Synopsis() string {
	return "Execute a command in the context of a running application instance"
}

func (c *ExecCommand) Help() string {
	return formatHelp(`
Usage: waypoint exec [options]

  Execute a command in the context of a running application instance.

` + c.Flags().Help())
}
