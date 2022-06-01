package cli

import (
	"context"

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
	if len(c.args) == 0 {
		c.ui.Output("Pipeline Name required.\n\n%s", c.Help(), terminal.WithErrorStyle())
		return 1
	} else {
		pipelineName = c.args[0]
	}

	// Pre-calculate our project ref
	projectRef := &pb.Ref_Project{Project: c.flagProject}
	if c.flagProject == "" {
		if c.project != nil {
			projectRef = c.project.Ref()
		}

		if projectRef == nil {
			c.ui.Output("You must specify a project with -project or be inside an existing project directory.\n"+c.Help(),
				terminal.WithErrorStyle())
			return 1
		}
	}

	err := c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
		app.UI.Output("Running pipeline %q for application %q",
			pipelineName, app.Ref().Application, terminal.WithHeaderStyle())

		sg := app.UI.StepGroup()
		defer sg.Wait()

		step := sg.Add("Building pipeline execution request...")
		defer step.Abort()

		runJobTemplate := &pb.Job{
			Application: app.Ref(),
			Workspace:   &pb.Ref_Workspace{Workspace: c.flagWorkspace},

			TargetRunner: &pb.Ref_Runner{Target: &pb.Ref_Runner_Any{}},
		}

		runPipelineReq := &pb.RunPipelineRequest{
			JobTemplate: runJobTemplate,

			Pipeline: &pb.Ref_Pipeline{
				Ref: &pb.Ref_Pipeline_Owner{
					Owner: &pb.Ref_PipelineOwner{
						Project:      projectRef,
						PipelineName: pipelineName,
					},
				},
			},
		}

		step.Update("Requesting to queue run of pipeline %q...", pipelineName)

		// take pipeline id and queue a RunPipeline with a Job Template.
		resp, err := c.project.Client().RunPipeline(c.Ctx, runPipelineReq)
		if err != nil {
			return err
		}

		step.Update("Pipeline %q has started running. Attempting to read job stream sequentially in order", pipelineName)
		step.Done()

		// TODO: do we need a stream timeout?
		// Receieve job ids from running pipeline, use job client to attach to job stream
		// and stream here. First pass can be linear job streaming
		for _, jobId := range resp.AllJobIds {
			app.UI.Output("Reading job stream (jobId: %s)...", jobId, terminal.WithHeaderStyle())
			app.UI.Output("")

			// Throw away the job result for now. We could do something fancy with this
			// later.
			// TODO: unknown stream event: type=*gen.GetJobStreamResponse_Download_ ??
			_, err := jobstream.Stream(c.Ctx, jobId, jobstream.WithClient(c.project.Client()))
			if err != nil {
				return err
			}
		}

		app.UI.Output("✔ Pipeline %q (%s) finished!", pipelineName, projectRef, terminal.WithSuccessStyle())

		return nil
	})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	return 0
}

func (c *PipelineRunCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
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
	required.

` + c.Flags().Help())
}
