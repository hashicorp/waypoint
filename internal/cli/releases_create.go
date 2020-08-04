package cli

import (
	"context"
	"strings"

	"github.com/posener/complete"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

type ReleaseCreateCommand struct {
	*baseCommand
}

func (c *ReleaseCreateCommand) Run(args []string) int {
	defer c.Close()

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
		// Get the latest release
		release, err := client.GetLatestRelease(ctx, &pb.GetLatestReleaseRequest{
			Application: app.Ref(),
			Workspace:   c.project.WorkspaceRef(),
		})
		if status.Code(err) == codes.NotFound {
			err = nil
			release = nil
		}
		if err != nil {
			app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}

		// Get the latest deployment
		resp, err := client.ListDeployments(ctx, &pb.ListDeploymentsRequest{
			Application:   app.Ref(),
			Workspace:     c.project.WorkspaceRef(),
			PhysicalState: pb.Operation_CREATED,
			Order: &pb.OperationOrder{
				Limit: 1,
				Order: pb.OperationOrder_COMPLETE_TIME,
				Desc:  true,
			},
		})
		if err != nil {
			app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}
		if len(resp.Deployments) == 0 {
			app.UI.Output(strings.TrimSpace(releaseNoDeploys), terminal.WithErrorStyle())
			return ErrSentinel
		}
		deploy := resp.Deployments[0]

		// If the latest release already deployed this then we're done.
		if release != nil && release.DeploymentId == deploy.Id {
			app.UI.Output(strings.TrimSpace(releaseUpToDate),
				deploy.Id,
				terminal.WithSuccessStyle())
			return nil
		}

		// UI
		app.UI.Output("Releasing...", terminal.WithHeaderStyle())

		// Release
		result, err := app.Release(ctx, &pb.Job_ReleaseOp{
			Deployment: deploy,
			Prune:      true,
		})
		if err != nil {
			app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}

		if result.Release.Url == "" {
			app.UI.Output("\n"+strings.TrimSpace(releaseNoUrl),
				deploy.Id,
				terminal.WithSuccessStyle())
			return nil
		}

		app.UI.Output("\nURL: https://%s", result.Release.Url, terminal.WithSuccessStyle())
		return nil
	})
	if err != nil {
		return 1
	}

	return 0
}

func (c *ReleaseCreateCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, nil)
}

func (c *ReleaseCreateCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ReleaseCreateCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ReleaseCreateCommand) Synopsis() string {
	return "Release a deployment."
}

func (c *ReleaseCreateCommand) Help() string {
	helpText := `
Usage: waypoint release [options]

  Open a deployment to traffic. This will default to shifting traffic
  100% to the latest deployment. You can specify multiple percentages to
  split traffic between releases.

Examples:

  "waypoint release" - will send 100% of traffic to the latest deployment.

  "waypoint release 90" - will send 90% of traffic to the latest deployment
  and 10% of traffic to the prior deployment.

  "waypoint release '+10'" - will send 10% more traffic to the latest deployment.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

const (
	releaseNoUrl = `
Deployment %s marked as released.

No release manager was configured and the configured platform doesn't
natively support releases. This means that releasing doesn't generate any
public URL. Waypoint marked the deployment aboved as "released" for server
history and to prevent commands such as "waypoint destroy" from deleting
the deployment by default.
`

	releaseNoDeploys = `
No deployments were found.

This may be because this application has never deployed before or it may be
because any previous deploys have been destroyed. Create a new deployment
using "waypoint deploy" and try again.
`

	releaseUpToDate = `
The deployment %q is already the released deployment. Nothing to be done.
`
)
