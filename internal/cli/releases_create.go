package cli

import (
	"context"
	"strconv"
	"strings"

	"github.com/posener/complete"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

type ReleaseCreateCommand struct {
	*baseCommand

	flagRepeat      bool
	flagDeployment  string
	flagPrune       bool
	flagPruneRetain int
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

		var deploy *pb.Deployment

		if c.flagDeployment == "" {
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
			deploy = resp.Deployments[0]
		} else if i, err := strconv.ParseUint(c.flagDeployment, 10, 64); err == nil {
			deploy, err = client.GetDeployment(ctx, &pb.GetDeploymentRequest{
				Ref: &pb.Ref_Operation{
					Target: &pb.Ref_Operation_Sequence{
						Sequence: &pb.Ref_OperationSeq{
							Application: app.Ref(),
							Number:      i,
						},
					},
				},
			})

			if err != nil {
				app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
				return ErrSentinel
			}
		} else {
			deploy, err = client.GetDeployment(ctx, &pb.GetDeploymentRequest{
				Ref: &pb.Ref_Operation{
					Target: &pb.Ref_Operation_Id{
						Id: c.flagDeployment,
					},
				},
			})

			if err != nil {
				app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
				return ErrSentinel
			}
		}

		// If the latest release already deployed this then we're done.
		if release != nil && release.DeploymentId == deploy.Id {
			if c.flagRepeat {
				c.Log.Warn("deployment already released but -repeat specified, will re-release")
			} else {
				c.Log.Warn("deployment already released")
				app.UI.Output(strings.TrimSpace(releaseUpToDate),
					deploy.Id,
					terminal.WithSuccessStyle())
				return nil
			}
		}

		if deploy.State != pb.Operation_CREATED {
			app.UI.Output("Deployment specified is not available (state=%s)", deploy.State,
				terminal.WithErrorStyle())
			return ErrSentinel
		}

		// UI
		app.UI.Output("Releasing...", terminal.WithHeaderStyle())

		// Release
		result, err := app.Release(ctx, &pb.Job_ReleaseOp{
			Deployment: deploy,

			Prune:               c.flagPrune,
			PruneRetain:         int32(c.flagPruneRetain),
			PruneRetainOverride: c.flagPruneRetain >= 0,
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

		app.UI.Output("\nRelease URL: %s", result.Release.Url, terminal.WithSuccessStyle())
		return nil
	})
	if err != nil {
		return 1
	}

	return 0
}

func (c *ReleaseCreateCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:    "repeat",
			Target:  &c.flagRepeat,
			Usage:   "Re-release if deploy is already released.",
			Default: false,
		})

		f.StringVar(&flag.StringVar{
			Name:    "deployment",
			Aliases: []string{"d"},
			Target:  &c.flagDeployment,
			Usage:   "Release the specified deployment.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "prune",
			Target:  &c.flagPrune,
			Usage:   "Prune old unreleased deployments.",
			Default: true,
		})

		f.IntVar(&flag.IntVar{
			Name:   "prune-retain",
			Target: &c.flagPruneRetain,
			Usage: "The number of unreleased deployments to keep. " +
				"If this isn't set or is set to any negative number, " +
				"then this will default to 1 on the server. If you want to prune " +
				"all unreleased deployments, set this to 0.",
			Default: -1,
		})
	})
}

func (c *ReleaseCreateCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ReleaseCreateCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ReleaseCreateCommand) Synopsis() string {
	return "Release a deployment"
}

func (c *ReleaseCreateCommand) Help() string {
	return formatHelp(`
Usage: waypoint release [options] [id]

  Open a deployment to traffic.

  This defaults to the latest deployment. Other deployments can be
  specified by sequence number or long ID.

` + c.Flags().Help())
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
