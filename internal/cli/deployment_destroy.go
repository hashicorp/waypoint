package cli

import (
	"context"
	"strconv"

	"github.com/posener/complete"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

type DeploymentDestroyCommand struct {
	*baseCommand

	flagAll   bool
	flagForce bool
}

func (c *DeploymentDestroyCommand) Run(args []string) int {
	ctx := c.Ctx
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

	// Determine the deployments to delete
	var deployments []*pb.Deployment

	var err error
	if len(args) > 0 {
		// If we have arguments, we only delete the deployments specified.
		deployments, err = c.getDeployments(ctx, args)
		if err != nil {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}
	} else {
		// No arguments, get ALL deployments that are still physically created.
		deployments, err = c.allDeployments(ctx)
		if err != nil {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}
	}

	// Destroy each deployment
	c.ui.Output("%d deployments will be destroyed.", len(deployments), terminal.WithHeaderStyle())
	for _, deployment := range deployments {
		// Can't destroy a deployment that was not successful
		if deployment.Status.GetState() != pb.Status_SUCCESS {
			continue
		}

		// Get our app client
		app := c.project.App(deployment.Application.Application)

		c.ui.Output("Destroying deployment: %s", deployment.Id, terminal.WithInfoStyle())
		if err := app.Destroy(ctx, &pb.Job_DestroyOp{
			Target: &pb.Job_DestroyOp_Deployment{
				Deployment: deployment,
			},
		}); err != nil {
			c.ui.Output("Error destroying the deployment: %s", err.Error(), terminal.WithErrorStyle())
			return 1
		}
	}

	return 0
}

func (c *DeploymentDestroyCommand) getDeployments(ctx context.Context, ids []string) ([]*pb.Deployment, error) {
	var result []*pb.Deployment

	// Get each deployment
	client := c.project.Client()
	for _, id := range ids {
		ref := &pb.Ref_Operation{
			Target: &pb.Ref_Operation_Id{Id: id},
		}
		if v, err := strconv.ParseInt(id, 10, 64); err == nil {
			ref.Target = &pb.Ref_Operation_Sequence{
				Sequence: &pb.Ref_OperationSeq{
					Application: c.refApp,
					Number:      uint64(v),
				},
			}
		}

		deployment, err := client.GetDeployment(ctx, &pb.GetDeploymentRequest{
			Ref: ref,
		})
		if err != nil {
			return nil, err
		}

		result = append(result, deployment)
	}

	return result, nil
}

func (c *DeploymentDestroyCommand) allDeployments(ctx context.Context) ([]*pb.Deployment, error) {
	L := c.Log

	var result []*pb.Deployment

	client := c.project.Client()
	err := c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
		resp, err := client.ListDeployments(ctx, &pb.ListDeploymentsRequest{
			Application:   app.Ref(),
			Workspace:     c.project.WorkspaceRef(),
			PhysicalState: pb.Operation_CREATED,
			Order: &pb.OperationOrder{
				Order: pb.OperationOrder_COMPLETE_TIME,
				Desc:  true,
			},
		})
		if err != nil {
			return err
		}

		// If we aren't deploying all, then we have to find the released
		// deployment and NOT delete that.
		if !c.flagAll {
			release, err := client.GetLatestRelease(ctx, &pb.GetLatestReleaseRequest{
				Application: app.Ref(),
				Workspace:   c.project.WorkspaceRef(),
			})
			if status.Code(err) == codes.NotFound {
				L.Debug("no release found to exclude any deployments")
				err = nil
				release = nil
			}
			if err != nil {
				return nil
			}

			if release != nil {
				for i := 0; i < len(resp.Deployments); i++ {
					d := resp.Deployments[i]
					if d.Id == release.DeploymentId {
						L.Info("not destroying deployment that is released", "id", d.Id)
						resp.Deployments[len(resp.Deployments)-1], resp.Deployments[i] =
							resp.Deployments[i], resp.Deployments[len(resp.Deployments)-1]
						resp.Deployments = resp.Deployments[:len(resp.Deployments)-1]
						i--
					}
				}
			}
		}

		result = append(result, resp.Deployments...)
		return nil
	})

	return result, err
}

func (c *DeploymentDestroyCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:    "all",
			Target:  &c.flagAll,
			Usage:   "Delete ALL deployments, including released.",
			Default: false,
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "force",
			Target:  &c.flagForce,
			Usage:   "Yes to all confirmations.",
			Default: false,
		})
	})
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
	return formatHelp(`
Usage: waypoint deployment destroy [options] [id...]

  Destroy one or more deployments. This will "undeploy" this specific
  instance of an application.

  When no arguments are given, this will default to destroying ALL
  unreleased deployments. This will require interactive confirmation
  by the user unless the force flag (-force) is specified.

` + c.Flags().Help())
}
