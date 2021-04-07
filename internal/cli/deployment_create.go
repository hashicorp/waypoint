package cli

import (
	"context"
	"strings"

	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

type DeploymentCreateCommand struct {
	*baseCommand

	flagRelease bool
}

func (c *DeploymentCreateCommand) Run(args []string) int {
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
		deployment := result.Deployment

		// Try to get the hostname
		var hostname *pb.Hostname
		hostnamesResp, err := client.ListHostnames(ctx, &pb.ListHostnamesRequest{
			Target: &pb.Hostname_Target{
				Target: &pb.Hostname_Target_Application{
					Application: &pb.Hostname_TargetApp{
						Application: deployment.Application,
						Workspace:   deployment.Workspace,
					},
				},
			},
		})
		if err == nil && len(hostnamesResp.Hostnames) > 0 {
			hostname = hostnamesResp.Hostnames[0]
		}

		// Release if we're releasing
		var releaseUrl string
		if c.flagRelease {
			// We're releasing, do that too.
			app.UI.Output("Releasing...", terminal.WithHeaderStyle())
			releaseResult, err := app.Release(ctx, &pb.Job_ReleaseOp{
				Deployment: deployment,
				Prune:      true,
			})
			if err != nil {
				app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
				return ErrSentinel
			}

			releaseUrl = releaseResult.Release.Url
		}

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
		return 1
	}

	return 0
}

func (c *DeploymentCreateCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:    "release",
			Target:  &c.flagRelease,
			Usage:   "Release this deployment immediately.",
			Default: true,
		})
	})
}

func (c *DeploymentCreateCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *DeploymentCreateCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *DeploymentCreateCommand) Synopsis() string {
	return "Deploy a pushed artifact"
}

func (c *DeploymentCreateCommand) Help() string {
	return formatHelp(`
Usage: waypoint deployment deploy [options]
Alias: waypoint deploy

  Deploy an application. This will deploy the most recent successful
  pushed artifact by default. You can view a list of recent artifacts
  using the "artifact list" command.

` + c.Flags().Help())
}

const (
	deployURLService = `
The deploy was successful! A Waypoint deployment URL is shown below. This
can be used internally to check your deployment and is not meant for external
traffic. You can manage this hostname using "waypoint hostname."
`

	deployNoURL = `
The deploy was successful!

The release did not provide a URL and the URL service is disabled on the
server, so no further URL information can be automatically provided. If
this is unexpected, please ensure the Waypoint server has both the URL service
enabled and advertise addresses set.
`
)
