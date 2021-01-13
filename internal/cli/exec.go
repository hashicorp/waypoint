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
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

type ExecCommand struct {
	*baseCommand

	flagInstanceId string
}

func (c *ExecCommand) targeted(args []string) int {
	var exitCode int

	client := c.project.Client()

	err := c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
		// We validate the users instance id first, since they're easy to mistype
		// and we can give them quicker and better feedback by doing this up front.

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

		// Ok, the instance id is fine, let's go ahead and have execclient do it's thing.

		client := &execclient.Client{
			Logger:     c.Log,
			UI:         c.ui,
			Context:    ctx,
			Client:     client,
			InstanceId: c.flagInstanceId,
			Args:       args,
			Stdin:      os.Stdin,
			Stdout:     os.Stdout,
			Stderr:     os.Stderr,
		}

		exitCode, err = client.Run()
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

func (c *ExecCommand) Run(args []string) int {
	flagSet := c.Flags()

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(flagSet),
		WithSingleApp(),
	); err != nil {
		return 1
	}

	args = flagSet.Args()

	if c.flagInstanceId != "" {
		return c.targeted(args)
	}

	var exitCode int
	client := c.project.Client()
	err := c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
		// Get the latest deployment
		resp, err := client.ListDeployments(ctx, &pb.ListDeploymentsRequest{
			Application: app.Ref(),
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

		client := &execclient.Client{
			Logger:        c.Log,
			UI:            c.ui,
			Context:       ctx,
			Client:        client,
			DeploymentId:  resp.Deployments[0].Id,
			DeploymentSeq: resp.Deployments[0].Sequence,
			Args:          args,
			Stdin:         os.Stdin,
			Stdout:        os.Stdout,
			Stderr:        os.Stderr,
		}

		exitCode, err = client.Run()
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
	return c.flagSet(0, func(s *flag.Sets) {
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
Usage: waypoint exec [options]

  Execute a command in the context of a running application instance.

` + c.Flags().Help())
}
