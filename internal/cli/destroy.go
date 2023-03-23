// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"context"
	"strings"

	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type DestroyCommand struct {
	*baseCommand

	confirm bool
}

func (c *DestroyCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithMultiAppTargets(),
	); err != nil {
		return 1
	}

	err := c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
		if !c.confirm {
			proceed, err := c.ui.Input(&terminal.Input{
				Prompt: "Do you really want to destroy all resources for this app? Only 'yes' will be accepted to approve: ",
				Style:  "",
				Secret: false,
			})
			if err != nil {
				c.ui.Output(
					"Error destroying resources: %s",
					clierrors.Humanize(err),
					terminal.WithErrorStyle(),
				)
			} else if strings.ToLower(proceed) != "yes" {
				app.UI.Output("Destroying app %q requires confirmation.", app.Ref().GetApplication(), terminal.WithWarningStyle())
				return nil
			}
		}

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
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:    "auto-approve",
			Target:  &c.confirm,
			Default: false,
			Usage:   "Auto-approve destroying all resources. If unset, confirmation will be requested.",
		})
	})
}

func (c *DestroyCommand) Synopsis() string {
	return "Delete all the resources created"
}

func (c *DestroyCommand) Help() string {
	return formatHelp(`
Usage: waypoint destroy [options]

  Delete all resources created for all apps or project in the current workspace.
  Specify the -app to select a given app to delete resources for in a given 
  workspace.

  The workspace can continue to be used after this call, this just deletes
  all the resources created for all apps in the workspace up to this point.

  This functionality must be supported by the plugins in use and is dependent
  on their behavior. The expected behavior is that any physical resources created
  as part of deploys and releases are destroyed. For example, any load balancers,
  VMs, containers, etc.

  This targets apps in one workspace. You must call this for each workspace
  you've used if you want to destroy everything.

` + c.Flags().Help())
}
