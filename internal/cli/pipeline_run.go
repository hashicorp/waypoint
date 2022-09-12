package cli

import (
	"context"
	"fmt"

	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	jobstream "github.com/hashicorp/waypoint/internal/jobstream"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type PipelineRunCommand struct {
	*baseCommand

	flagPipelineId string
}

func (c *PipelineRunCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
	); err != nil {
		return 1
	}

	var pipelineName string
	if len(c.args) == 0 && c.flagPipelineId == "" {
		c.ui.Output("Pipeline name or ID required.\n\n%s", c.Help(), terminal.WithErrorStyle())
		return 1
	} else if c.flagPipelineId == "" {
		pipelineName = c.args[0]
	} else {
		c.ui.Output("Both pipeline name and ID were specified, using pipeline ID", terminal.WithWarningStyle())
	}

	if c.flagLocal != nil && *c.flagLocal {
		// TODO(briancain): Remove this warning when local support for Pipelines is introduced.
		// GitHub: https://github.com/hashicorp/waypoint/issues/3813
		c.ui.Output("At the moment, the initial Tech Preview of Custom Pipelines does not allow "+
			"for executing pipelines with a local runner. The CLI will attempt to run the "+
			"requested pipeline but it will likely fail if the project was not configured "+
			"to run remotely.",
			terminal.WithWarningStyle())
	}

	err := c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
		// setup pipeline name to be used for UI printing
		pipelineIdent := pipelineName
		if c.flagPipelineId != "" {
			pipelineIdent = c.flagPipelineId
		}

		app.UI.Output("Running pipeline %q for application %q",
			pipelineIdent, app.Ref().Application, terminal.WithHeaderStyle())

		sg := app.UI.StepGroup()
		defer sg.Wait()

		step := sg.Add("Syncing pipeline configs...")
		defer step.Abort()

		_, err := app.ConfigSync(ctx, &pb.Job_ConfigSyncOp{})
		if err != nil {
			app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}

		step.Update("Building pipeline execution request...")

		// build the initial job template for running the pipeline
		runJobTemplate := &pb.Job{
			Application:  app.Ref(),
			Workspace:    c.project.WorkspaceRef(),
			TargetRunner: &pb.Ref_Runner{Target: &pb.Ref_Runner_Any{}},
		}

		// build the base api request
		runPipelineReq := &pb.RunPipelineRequest{
			JobTemplate: runJobTemplate,
		}

		// Reference by ID if set, otherwise by name
		if c.flagPipelineId != "" {
			runPipelineReq.Pipeline = &pb.Ref_Pipeline{
				Ref: &pb.Ref_Pipeline_Id{
					Id: c.flagPipelineId,
				},
			}
		} else {
			runPipelineReq.Pipeline = &pb.Ref_Pipeline{
				Ref: &pb.Ref_Pipeline_Owner{
					Owner: &pb.Ref_PipelineOwner{
						Project:      &pb.Ref_Project{Project: app.Ref().Project},
						PipelineName: pipelineName,
					},
				},
			}
		}

		step.Update("Requesting to queue run of pipeline %q...", pipelineIdent)

		// take pipeline id and queue a RunPipeline with a Job Template.
		resp, err := c.project.Client().RunPipeline(c.Ctx, runPipelineReq)
		if err != nil {
			return err
		}

		step.Update("Pipeline %q has started running. Attempting to read job stream sequentially in order", pipelineIdent)
		step.Done()

		// Receive job ids from running pipeline, use job client to attach to job stream
		// and stream here. First pass can be linear job streaming
		step = sg.Add("")
		defer step.Abort()

		steps := len(resp.JobMap)
		step.Update("%d steps detected, run sequence %d", steps, resp.Sequence)
		step.Done()

		successful := steps
		for _, jobId := range resp.AllJobIds {
			job, err := c.project.Client().GetJob(c.Ctx, &pb.GetJobRequest{
				JobId: jobId,
			})
			if err != nil {
				return err
			}
			ws := "default"
			if job.Workspace != nil {
				ws = job.Workspace.Workspace
			}
			app.UI.Output("Executing Step %q in workspace: %q", resp.JobMap[jobId].Step, ws, terminal.WithHeaderStyle())
			app.UI.Output("Reading job stream (jobId: %s)...", jobId, terminal.WithInfoStyle())
			app.UI.Output("")

			_, err = jobstream.Stream(c.Ctx, jobId,
				jobstream.WithClient(c.project.Client()),
				jobstream.WithUI(app.UI))
			if err != nil {
				return err
			}

			job, err = c.project.Client().GetJob(c.Ctx, &pb.GetJobRequest{
				JobId: jobId,
			})
			if err != nil {
				return err
			}
			if job.State != pb.Job_SUCCESS {
				successful--
			}
		}

		output := fmt.Sprintf("Pipeline %q (%s) finished! %d/%d steps successfully completed.", pipelineIdent, app.Ref().Project, successful, steps)
		if successful == 0 {
			app.UI.Output("✖ %s", output, terminal.WithErrorStyle())
		} else if successful < steps {
			app.UI.Output("● %s", output, terminal.WithWarningStyle())
		} else {
			app.UI.Output("✔ %s", output, terminal.WithSuccessStyle())
		}

		return nil
	})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	return 0
}

func (c *PipelineRunCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.StringVar(&flag.StringVar{
			Name:    "id",
			Target:  &c.flagPipelineId,
			Default: "",
			Usage:   "Run a pipeline by ID.",
		})
	})
}

func (c *PipelineRunCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *PipelineRunCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PipelineRunCommand) Synopsis() string {
	return "Manually execute a pipeline by name."
}

func (c *PipelineRunCommand) Help() string {
	return formatHelp(`
Usage: waypoint pipeline run [options] <pipeline-name>

  Run a pipeline by name. If run outside of a project dir, a '-project' flag is
	required. Before running a requested pipeline, this command will sync
	pipeline configs so it runs the most up to date configuration version for a
	pipeline.

` + c.Flags().Help())
}
