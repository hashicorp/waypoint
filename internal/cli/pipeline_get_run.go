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

type PipelineGetRunCommand struct {
	*baseCommand

	flagPipelineId  string
	flagRunSequence int
}

func (c *PipelineGetRunCommand) Run(args []string) int {
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

		app.UI.Output("Streaming pipeline %q run for application %q",
			pipelineIdent, app.Ref().Application, terminal.WithHeaderStyle())

		sg := app.UI.StepGroup()
		defer sg.Wait()

		step := sg.Add("Reading pipeline run...")
		defer step.Abort()

		// build the base api request
		getPipelineReq := &pb.GetPipelineRequest{}

		// Reference by ID if set, otherwise by name
		if c.flagPipelineId != "" {
			getPipelineReq.Pipeline = &pb.Ref_Pipeline{
				Ref: &pb.Ref_Pipeline_Id{
					Id: c.flagPipelineId,
				},
			}
		} else {
			getPipelineReq.Pipeline = &pb.Ref_Pipeline{
				Ref: &pb.Ref_Pipeline_Owner{
					Owner: &pb.Ref_PipelineOwner{
						Project:      &pb.Ref_Project{Project: app.Ref().Project},
						PipelineName: pipelineName,
					},
				},
			}
		}

		step.Update("Requesting to stream run of pipeline %q...", pipelineIdent)

		var (
			resp *pb.GetPipelineRunResponse
			err  error
		)
		if c.flagRunSequence == 0 {
			// take pipeline id and queue a RunPipeline with a Job Template.
			resp, err = c.project.Client().GetLatestPipelineRun(c.Ctx, getPipelineReq)
			if err != nil {
				return err
			}
		} else {
			// take pipeline id and queue a RunPipeline with a Job Template.
			resp, err = c.project.Client().GetPipelineRun(c.Ctx, &pb.GetPipelineRunRequest{
				Pipeline: getPipelineReq.Pipeline,
				Sequence: uint64(c.flagRunSequence),
			})
			if err != nil {
				return err
			}
		}
		if resp == nil {
			app.UI.Output("Getting a pipeline run returned a nil response", terminal.WithErrorStyle())
			return fmt.Errorf("Response was empty when requesting a pipeline run for pipeline %q", pipelineIdent)
		}

		step.Update("Attempting to read job stream sequentially in order for run %q", resp.PipelineRun.Sequence)
		step.Done()

		// Receive job ids from running pipeline, use job client to attach to job stream
		// and stream here. First pass can be linear job streaming
		step = sg.Add("")
		defer step.Abort()

		steps := len(resp.PipelineRun.Jobs)
		step.Update("%d steps detected, run sequence %d, pipeline run status %q",
			steps, resp.PipelineRun.Sequence, pb.PipelineRun_State_name[int32(resp.PipelineRun.State)])
		step.Done()

		successful := steps
		for _, jobIdRef := range resp.PipelineRun.Jobs {
			jobId := jobIdRef.Id
			job, err := c.project.Client().GetJob(c.Ctx, &pb.GetJobRequest{
				JobId: jobId,
			})
			if err != nil {
				return err
			}
			// NOTE(briancain): We intentionally skip Noop type jobs because currently
			// we make step Refs for pipelines run a Noop job to make dependency tracking
			// for pipeline step refs easier. We don't stream a noop output job because
			// there's nothing to stream.
			if _, ok := job.Operation.(*pb.Job_Noop_); ok {
				continue
			}

			ws := "default"
			if job.Workspace != nil {
				ws = job.Workspace.Workspace
			}
			stepName := job.Pipeline.Step
			app.UI.Output("Executing Step %q in workspace: %q", stepName, ws, terminal.WithHeaderStyle())
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

func (c *PipelineGetRunCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.StringVar(&flag.StringVar{
			Name:    "id",
			Target:  &c.flagPipelineId,
			Default: "",
			Usage:   "Run a pipeline by ID.",
		})

		f.IntVar(&flag.IntVar{
			Name:   "run",
			Target: &c.flagRunSequence,
			Usage:  "Inspect a specific run sequence.",
		})
	})
}

func (c *PipelineGetRunCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *PipelineGetRunCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PipelineGetRunCommand) Synopsis() string {
	return "Attach to a pipeline runs job stream."
}

func (c *PipelineGetRunCommand) Help() string {
	return formatHelp(`
Usage: waypoint pipeline get-run [options] <pipeline-name>

  Attempts to reattach the CLI to an existing pipeline run. Defaults to latest,
	but if '-run' is specified, it will attach to that specific run by sequence number.

` + c.Flags().Help())
}
