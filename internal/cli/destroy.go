package cli

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

type DestroyCommand struct {
	*baseCommand
}

func (c *DestroyCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithSingleApp(),
	); err != nil {
		return 1
	}

	err := c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
		if err := app.Destroy(ctx, &pb.Job_DestroyOp{
			Target: &pb.Job_DestroyOp_Workspace{
				Workspace: &empty.Empty{},
			},
		}); err != nil {
			c.ui.Output("Error destroying: %s", err.Error(), terminal.WithErrorStyle())
			return ErrSentinel
		}

		app.UI.Output("Destroy successful!", terminal.WithSuccessStyle())
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

func (c *DestroyCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
	})
}

func (c *DestroyCommand) Synopsis() string {
	return "Delete all the resources created for an app"
}

func (c *DestroyCommand) Help() string {
	return formatHelp(`
Usage: waypoint destroy [options]

  Delete all resources created for an app in the current workspace.

  The workspace can continue to be used after this call, this just deletes
  all the resources created for this app up to this point.

  This functionality must be supported by the plugins in use and is dependent
  on their behavior. The expect behavior is that any physical resources created
  as part of deploys and releases are destroyed. For example, any load balancers,
  VMs, containers, etc.

  This targets one app in one workspace. You must call this for each workspace
  you've used if you want to destroy everything.

` + c.Flags().Help())
}
