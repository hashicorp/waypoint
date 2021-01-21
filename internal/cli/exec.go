package cli

import (
	"context"
	"os"
	"time"

	"github.com/posener/complete"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/server"
	"github.com/hashicorp/waypoint/internal/server/execclient"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
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

	// Ok, the instance id is fine, let's go ahead and have execclient do it's thing.
	ec.InstanceId = c.flagInstanceId

	return nil
}

func (c *ExecCommand) findDeployment(
	ctx context.Context,
	app *clientpkg.App,
	client pb.WaypointClient,
) (*pb.Deployment, error) {
	c.Log.Debug("looking up deployments to use for app", "app", app.Ref().String())

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
		return nil, ErrSentinel
	}
	if len(resp.Deployments) == 0 {
		app.UI.Output("No successful deployments found.", terminal.WithErrorStyle())
		return nil, ErrSentinel
	}

	c.Log.Debug("found deployment", "deployment-id", resp.Deployments[0].Id)
	return resp.Deployments[0], nil
}

func (c *ExecCommand) searchDeployments(
	ctx context.Context,
	app *clientpkg.App,
	client pb.WaypointClient,
	ec *execclient.Client,
) (chan error, error) {
	// Get the latest deployment
	deployment, err := c.findDeployment(ctx, app, client)
	if err != nil {
		return nil, err
	}

	c.Log.Debug("looking for instances of deployment to exec into")
	resp, err := client.FindExecInstance(ctx, &pb.FindExecInstanceRequest{
		DeploymentId: deployment.Id,
	})
	if err != nil {
		// If the server says there are no available instances, try to generate one
		// via the exec plugin. If the app's deployment plugin doesn't have an exec
		// plugin, this will generate an error that we will show to the user as meaning
		// there is no way to start an exec session.
		if st, ok := status.FromError(err); ok && st.Code() == codes.ResourceExhausted {
			c.Log.Debug("no instances found, trying exec via plugin")
			return c.viaPlugin(ctx, app, client, ec, deployment)
		}
		return nil, err
	}

	c.Log.Debug("found instance to exec into", "instance-id", resp.Instance.Id)

	ec.InstanceId = resp.Instance.Id
	ec.DeploymentId = deployment.Id
	ec.DeploymentSeq = deployment.Sequence

	return nil, nil
}

func (c *ExecCommand) viaPlugin(
	ctx context.Context,
	app *clientpkg.App,
	client pb.WaypointClient,
	ec *execclient.Client,
	deployment *pb.Deployment,
) (chan error, error) {
	instId, err := server.Id()
	if err != nil {
		return nil, err
	}

	ec.InstanceId = instId

	mon := make(chan pb.Job_State)
	done := make(chan error)

	// Start the plugin exec in the background so we can run the actual exec client
	// in the main goroutine.
	c.Log.Debug("launching exec job to start plugin")
	go func() {
		defer close(done)

		done <- app.Exec(ctx, &pb.Job_ExecOp{
			InstanceId: instId,
			Deployment: deployment,
		}, mon)
	}()

outer:
	for {
		select {

		// Someone or something canceled us
		case <-ctx.Done():
			return nil, ctx.Err()

		// The plugin has stopped
		case err := <-done:
			return nil, err

		// The runner job's status has changed
		case st, ok := <-mon:
			if !ok {
				return nil, status.Error(codes.Aborted, "job ended unexpectedly")
			}

			c.Log.Debug("received exec job status", "status", st.String())

			// We're waiting for just the initial job startup action here, then
			// we exit this loop so the execclient can run.
			switch st {

			// Error off the bat usually means the plugin crashed.
			case pb.Job_ERROR:
				return nil, status.Error(codes.Aborted, "job errored out before starting")

			// Queued by the server, all good.
			case pb.Job_WAITING:
				// ok

			// Plugin stopped, not good, unknown how the ordering would result here.
			case pb.Job_SUCCESS:
				return nil, status.Error(codes.Aborted, "job finished before we used it")

			// A runner has started the job, yay!
			case pb.Job_RUNNING:
				// super duper, let's check on the instance now.
				break outer
			}
		}
	}

	// Drain monitor so we don't hold up the job management code sending us updates.
	go func() {
		for {
			<-mon
		}
	}()

	// Look at the instances for the deployment and wait for our instance to pop up.

	// Only wait up to 30s for the instance to appear.
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	c.Log.Debug("looking for virtual instance to appear", "instance-id", instId)

	// Look at the instances and find our virtual instance.
	for {
		resp, err := client.ListInstances(ctx, &pb.ListInstancesRequest{
			Scope: &pb.ListInstancesRequest_DeploymentId{
				DeploymentId: deployment.Id,
			},
			WaitTimeout: "2s",
		})

		if err != nil {
			if err == context.Canceled {
				return nil, status.Error(codes.Aborted, "instance didn't appear after 30 seconds")
			}

			return nil, err
		}

		for _, inst := range resp.Instances {
			if inst.Id == instId {
				c.Log.Debug("virtual instance found in list of instances")
				// yeeehaw!
				return done, nil
			}
		}
	}
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
			err  error
			done chan error
		)

		if c.flagInstanceId != "" {
			err = c.targeted(ctx, app, client, ec)
		} else {
			done, err = c.searchDeployments(ctx, app, client, ec)
		}
		if err != nil {
			return err
		}

		exitCode, err = ec.Run()
		if err != nil {
			app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}

		// If there is a background task, wait for it to finish and report it's status.
		// Currently, this is the exec job if the plugin method was used.
		if done != nil {
			select {
			case err := <-done:
				if err != nil {
					app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
					return ErrSentinel
				}
			case <-ctx.Done():
				return ctx.Err()
			}
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
