package cli

import (
	"context"
	"strings"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
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

	err := c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
		result, err := app.Up(ctx, &pb.Job_UpOp{})
		if err != nil {
			app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}

		// Common reused values
		releaseUrl := result.Up.ReleaseUrl
		appUrl := result.Up.AppUrl
		deployUrl := result.Up.DeployUrl

		// Output
		app.UI.Output("")
		switch {
		case releaseUrl != "":
			app.UI.Output(strings.TrimSpace(deployURLService)+"\n", terminal.WithSuccessStyle())
			app.UI.Output("   Release URL: %s", releaseUrl, terminal.WithSuccessStyle())
			if deployUrl != "" {
				app.UI.Output("Deployment URL: %s", deployUrl, terminal.WithSuccessStyle())
			}

		case appUrl != "" && deployUrl != "":
			app.UI.Output(strings.TrimSpace(deployURLService)+"\n", terminal.WithSuccessStyle())
			app.UI.Output("           URL: %s", appUrl, terminal.WithSuccessStyle())
			app.UI.Output("Deployment URL: %s", deployUrl, terminal.WithSuccessStyle())

		default:
			app.UI.Output(strings.TrimSpace(deployNoURL)+"\n", terminal.WithSuccessStyle())
		}

		return nil
	})

	if err != nil {
		if err != ErrSentinel {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
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
	return "Perform the build, deploy, and release steps for the app"
}

func (c *UpCommand) Help() string {
	return formatHelp(`
Usage: waypoint up [options]

  Perform the build, deploy, and release steps for the app.

` + c.Flags().Help())
}
