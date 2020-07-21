package cli

import (
	"context"
	"strings"

	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

type UpCommand struct {
	*baseCommand
}

func (c *UpCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithSingleApp(),
	); err != nil {
		return 1
	}

	client := c.project.Client()

	err := c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
		// Build it
		app.UI.Output("Building...", terminal.WithHeaderStyle())

		_, err := app.Build(ctx, &pb.Job_BuildOp{})
		if err != nil {
			app.UI.Output(err.Error(), terminal.WithErrorStyle())
			return ErrSentinel
		}

		// Get the most recent pushed artifact
		push, err := client.GetLatestPushedArtifact(ctx, &pb.GetLatestPushedArtifactRequest{
			Application: app.Ref(),
			Workspace:   c.project.WorkspaceRef(),
		})
		if err != nil {
			app.UI.Output(err.Error(), terminal.WithErrorStyle())
			return ErrSentinel
		}

		// Push it
		app.UI.Output("Deploying...", terminal.WithHeaderStyle())

		result, err := app.Deploy(ctx, &pb.Job_DeployOp{
			Artifact: push,
		})
		if err != nil {
			app.UI.Output(err.Error(), terminal.WithErrorStyle())
			return ErrSentinel
		}

		cfg, ok := c.cfg.AppConfig(c.app)
		if !ok {
			app.UI.Output(
				"Strangly the application configuration is unavailable",
				terminal.WithErrorStyle(),
			)
			return ErrSentinel
		}

		if cfg.Release == nil {
			app.UI.Output("App Ready!", terminal.WithHeaderStyle())
			app.UI.Output("The application did not provide a release step, so here is the deployment info:")

			app.UI.NamedValues([]terminal.NamedValue{
				{
					Name:  "Id",
					Value: result.Deployment.Id,
				},
			}, terminal.WithInfoStyle())

			return nil
		}

		// We're releasing, do that too.
		app.UI.Output("Releasing...", terminal.WithHeaderStyle())

		releaseResult, err := app.Release(ctx, &pb.Job_ReleaseOp{
			TrafficSplit: &pb.Release_Split{
				Targets: []*pb.Release_SplitTarget{
					{
						DeploymentId: result.Deployment.Id,
						Percent:      100,
					},
				},
			},
		})
		if err != nil {
			app.UI.Output(err.Error(), terminal.WithErrorStyle())
			return ErrSentinel
		}

		app.UI.Output("App Ready!", terminal.WithHeaderStyle())
		app.UI.Output("The application can be accessed using the following information:")

		app.UI.NamedValues([]terminal.NamedValue{
			{
				Name:  "URL",
				Value: releaseResult.Release.Url,
			},
		}, terminal.WithInfoStyle())

		return nil
	})

	if err != nil {
		if err != ErrSentinel {
			c.ui.Output(err.Error(), terminal.WithErrorStyle())
		}

		return 1
	}

	return 0
}

func (c *UpCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
	})
}

func (c *UpCommand) Synopsis() string {
	return "Perform the build, deploy, and release steps for the app."
}

func (c *UpCommand) Help() string {
	helpText := `
Usage: waypoint up [options]

  Perform the build, deploy, and release steps for the app.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}
