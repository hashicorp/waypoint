// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"context"
	"os"

	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/server/execclient"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type ExecCommand struct {
	*baseCommand

	flagInstanceId string
}

func (c *ExecCommand) targeted(
	ctx context.Context,
	app *clientpkg.App,
	client pb.WaypointClient,
	ec *execclient.Client,
) error {
	resp, err := client.ListInstances(ctx, &pb.ListInstancesRequest{
		Scope: &pb.ListInstancesRequest_Application_{
			Application: &pb.ListInstancesRequest_Application{
				Application: app.Ref(),
			},
		},
	})
	if err != nil {
		app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return ErrSentinel
	}

	var found bool
	for _, i := range resp.Instances {
		if i.Id == c.flagInstanceId {
			found = true
		}
	}

	if !found {
		app.UI.Output("Unable to find instance: %s", c.flagInstanceId, terminal.WithErrorStyle())
		return ErrSentinel
	}

	// Ok, the instance id is fine, let's go ahead and have execclient do its thing.
	ec.InstanceId = c.flagInstanceId

	return nil
}

func (c *ExecCommand) searchDeployments(
	ctx context.Context,
	app *clientpkg.App,
	client pb.WaypointClient,
	ec *execclient.Client,
) error {
	// Get the latest deployment
	c.Log.Debug("looking up deployments to use for app", "app", app.Ref().String())

	// Get the latest deployment
	resp, err := client.ListDeployments(ctx, &pb.ListDeploymentsRequest{
		Application: app.Ref(),
		Workspace:   c.project.WorkspaceRef(),
		Order: &pb.OperationOrder{
			Limit: 1,
			Order: pb.OperationOrder_COMPLETE_TIME,
			Desc:  true,
		},
		PhysicalState: pb.Operation_CREATED,
	})
	if err != nil {
		app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return ErrSentinel
	}
	if len(resp.Deployments) == 0 {
		app.UI.Output("No successful deployments found.", terminal.WithErrorStyle())
		return ErrSentinel
	}

	c.Log.Debug("found deployment", "deployment-id", resp.Deployments[0].Id)

	deployment := resp.Deployments[0]

	ec.DeploymentId = deployment.Id
	ec.DeploymentSeq = deployment.Sequence

	return nil
}

func (c *ExecCommand) Run(args []string) int {
	flagSet := c.Flags()

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(flagSet),
		WithSingleAppTarget(),
	); err != nil {
		return 1
	}

	args = flagSet.Args()
	if len(args) == 0 {
		c.ui.Output(
			"At least one argument expected.\n\n"+c.Help(),
			terminal.WithErrorStyle(),
		)

		return 1
	}

	var exitCode int
	client := c.project.Client()
	err := c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
		ec := &execclient.Client{
			Logger:  c.Log,
			UI:      c.ui,
			Context: ctx,
			Client:  client,
			Args:    args,
			Stdin:   os.Stdin,
			Stdout:  os.Stdout,
			Stderr:  os.Stderr,
		}

		var (
			err error
		)

		if c.flagInstanceId != "" {
			err = c.targeted(ctx, app, client, ec)
		} else {
			err = c.searchDeployments(ctx, app, client, ec)
		}
		if err != nil {
			app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}

		exitCode, err = app.Exec(ctx, ec)
		if err != nil {
			app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}

		return nil
	})
	if err != nil {
		return 1
	}

	return exitCode
}

func (c *ExecCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(s *flag.Sets) {
		f := s.NewSet("Command Options")
		f.StringVar(&flag.StringVar{
			Name:   "instance",
			Usage:  "Start an exec session on this specific instance",
			Target: &c.flagInstanceId,
		})
	})
}

func (c *ExecCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ExecCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ExecCommand) Synopsis() string {
	return "Execute a command in the context of a running application instance"
}

func (c *ExecCommand) Help() string {
	return formatHelp(`
Usage: waypoint exec [options] cmd

  Execute a command in the context of a running application instance.

  For example, you could run one of the following commands:

    waypoint exec bash
    waypoint exec rake db:migrate

` + c.Flags().Help())
}
