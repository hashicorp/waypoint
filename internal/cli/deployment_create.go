package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/posener/complete"

	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/sdk/terminal"
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

	c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
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
		deployment := result.Deployment

		// Try to get the hostname
		var hostname *pb.Hostname
		var deployUrl string
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

			deployUrl = fmt.Sprintf(
				"%s--%s%s",
				hostname.Hostname,
				deployment.Id,
				strings.TrimPrefix(hostname.Fqdn, hostname.Hostname),
			)
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
			app.UI.Output("URL: %s", releaseUrl, terminal.WithSuccessStyle())

		case hostname != nil:
			app.UI.Output(strings.TrimSpace(deployURLService)+"\n", terminal.WithSuccessStyle())
			app.UI.Output("           URL: https://%s", hostname.Fqdn, terminal.WithSuccessStyle())
			app.UI.Output("Deployment URL: https://%s", deployUrl, terminal.WithSuccessStyle())

		default:
			app.UI.Output(strings.TrimSpace(deployNoURL)+"\n", terminal.WithSuccessStyle())
		}

		return nil
	})

	return 0
}

func (c *DeploymentCreateCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:    "release",
			Target:  &c.flagRelease,
			Usage:   "Release this deployment immedately.",
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
	return "Deploy a pushed artifact."
}

func (c *DeploymentCreateCommand) Help() string {
	helpText := `
Usage: waypoint deployment deploy [options]

  Deploy an application. This will deploy the most recent successful
  pushed artifact by default. You can view a list of recent artifacts
  using the "artifact list" command.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

const (
	deployURLService = `
The deploy was successful! A Waypoint deployment URL is shown below. This
can be used internally to check your deployment and is not meant for external
traffic. You can manage this hostname using "waypoint hostname."
`

	deployNoURL = `
The deploy was successful!
`
)
