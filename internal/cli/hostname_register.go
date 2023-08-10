// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cli

import (
	"context"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/posener/complete"
)

type HostnameRegisterCommand struct {
	*baseCommand
}

func (c *HostnameRegisterCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithMultiAppTargets(),
	); err != nil {
		return 1
	}

	hostname := ""
	if len(c.args) > 0 {
		hostname = c.args[0]
	}
	if hostname != "" && c.flagApp == "" {
		c.ui.Output("A target app is required when providing a specific hostname to register.",
			terminal.WithErrorStyle(),
		)
		return 1
	}

	client := c.project.Client()
	err := c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
		app.UI.Output("Registering hostname for %s...", app.Ref().Application, terminal.WithHeaderStyle())
		resp, err := client.CreateHostname(ctx, &pb.CreateHostnameRequest{
			Hostname: hostname,
			Target: &pb.Hostname_Target{
				Target: &pb.Hostname_Target_Application{
					Application: &pb.Hostname_TargetApp{
						Application: app.Ref(),
						Workspace:   c.project.WorkspaceRef(),
					},
				},
			},
		})
		if err != nil {
			app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}

		c.ui.Output(resp.Hostname.Fqdn, terminal.WithSuccessStyle())
		return nil
	})
	if err != nil {
		return 1
	}

	return 0
}

func (c *HostnameRegisterCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *HostnameRegisterCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *HostnameRegisterCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *HostnameRegisterCommand) Synopsis() string {
	return "Register a hostname to route to your apps."
}

func (c *HostnameRegisterCommand) Help() string {
	return formatHelp(`
Usage: waypoint hostname register [hostname]

  Register a hostname with the URL service to route to your apps.

  The URL service must be enabled and configured with the Waypoint server.
  This will output the fully qualified domain name that should begin
  routing immediately.

` + c.Flags().Help())
}
