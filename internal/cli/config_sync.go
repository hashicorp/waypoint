package cli

import (
	"context"

	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type ConfigSyncCommand struct {
	*baseCommand
}

func (c *ConfigSyncCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithMultiAppTargets(),
	); err != nil {
		return 1
	}

	err := c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
		sg := app.UI.StepGroup()
		defer sg.Wait()

		step := sg.Add("Synchronizing configuration variables and pipeline configs...")
		defer step.Abort()

		jobResult, err := app.ConfigSync(ctx, &pb.Job_ConfigSyncOp{})
		if err != nil {
			app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}

		step.Update("Configuration variables synchronized successfully!")
		step.Done()

		if jobResult.PipelineConfigSync != nil && len(jobResult.PipelineConfigSync.SyncedPipelines) > 0 {
			step := sg.Add("Configuration for pipelines synchronized successfully!")
			step.Done()

			// Extra space
			app.UI.Output("")
			for name, ref := range jobResult.PipelineConfigSync.SyncedPipelines {
				pipelineRef, ok := ref.Ref.(*pb.Ref_Pipeline_Id)
				if !ok {
					app.UI.Output("failed to convert pipeline ref", terminal.WithErrorStyle())
					return ErrSentinel
				}

				app.UI.Output("Pipeline %q (%s) synchronized!", name, pipelineRef.Id, terminal.WithInfoStyle())
			}
		}

		return nil
	})
	if err != nil {
		return 1
	}

	return 0
}

func (c *ConfigSyncCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *ConfigSyncCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ConfigSyncCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ConfigSyncCommand) Synopsis() string {
	return "Synchronize declared variables and pipeline configs in a waypoint.hcl"
}

func (c *ConfigSyncCommand) Help() string {
	return formatHelp(`
Usage: waypoint config sync [options]

  Synchronize declared application configuration in the waypoint.hcl file
  for existing and new deployments.

  Conflicting configuration keys will be overwritten. Configuration keys
  that do not exist in the waypoint.hcl file but exist on the server will not
  be deleted.

` + c.Flags().Help())
}
