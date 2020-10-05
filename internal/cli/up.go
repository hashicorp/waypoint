package cli

import (
	"context"
	"strings"

	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
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
			app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}

		// Get the most recent pushed artifact
		push, err := client.GetLatestPushedArtifact(ctx, &pb.GetLatestPushedArtifactRequest{
			Application: app.Ref(),
			Workspace:   c.project.WorkspaceRef(),
		})
		if err != nil {
			app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}

		// Push it
		app.UI.Output("Deploying...", terminal.WithHeaderStyle())

		result, err := app.Deploy(ctx, &pb.Job_DeployOp{
			Artifact: push,
		})
		if err != nil {
			app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}
		deployUrl := result.Deployment.Preload.DeployUrl

		// Try to get the hostname
		var hostname *pb.Hostname
		hostnamesResp, err := client.ListHostnames(ctx, &pb.ListHostnamesRequest{
			Target: &pb.Hostname_Target{
				Target: &pb.Hostname_Target_Application{
					Application: &pb.Hostname_TargetApp{
						Application: result.Deployment.Application,
						Workspace:   result.Deployment.Workspace,
					},
				},
			},
		})
		if err == nil && len(hostnamesResp.Hostnames) > 0 {
			hostname = hostnamesResp.Hostnames[0]
		}

		// We're releasing, do that too.
		app.UI.Output("Releasing...", terminal.WithHeaderStyle())
		releaseResult, err := app.Release(ctx, &pb.Job_ReleaseOp{
			Deployment: result.Deployment,
			Prune:      true,
		})
		if err != nil {
			app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}

		releaseUrl := releaseResult.Release.Url

		// Output
		app.UI.Output("")
		switch {
		case releaseUrl != "":
			app.UI.Output(strings.TrimSpace(deployURLService)+"\n", terminal.WithSuccessStyle())
			app.UI.Output("   Release URL: %s", releaseUrl, terminal.WithSuccessStyle())
			app.UI.Output("Deployment URL: https://%s", deployUrl, terminal.WithSuccessStyle())

		case hostname != nil:
			app.UI.Output(strings.TrimSpace(deployURLService)+"\n", terminal.WithSuccessStyle())
			app.UI.Output("           URL: https://%s", hostname.Fqdn, terminal.WithSuccessStyle())
			app.UI.Output("Deployment URL: https://%s", deployUrl, terminal.WithSuccessStyle())

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
