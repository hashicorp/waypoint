package cli

import (
	"context"
	"fmt"
	"strings"

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

	flagPipelineId  string
	flagReattachRun bool
	flagRunSequence int
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

	if c.flagRunSequence > 0 && !c.flagReattachRun {
		c.ui.Output("The '-run' flag was specified, automatically assuming '-reattach'. "+
			"CLI will attempt to reattach to pipeline run %d",
			c.flagRunSequence,
			terminal.WithWarningStyle())

		c.flagReattachRun = true
	}

	err := c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
		// setup pipeline name to be used for UI printing
		pipelineIdent := pipelineName
		if c.flagPipelineId != "" {
			pipelineIdent = c.flagPipelineId
		}

		if c.flagReattachRun {
			app.UI.Output("Streaming pipeline %q run for application %q",
				pipelineIdent, app.Ref().Application, terminal.WithHeaderStyle())
		} else {
			app.UI.Output("Running pipeline %q for application %q",
				pipelineIdent, app.Ref().Application, terminal.WithHeaderStyle())
		}

		sg := app.UI.StepGroup()
		defer sg.Wait()

		step := sg.Add("")
		defer step.Abort()

		if !c.flagReattachRun {
			step.Update("Syncing pipeline configs...")

			_, err := app.ConfigSync(ctx, &pb.Job_ConfigSyncOp{})
			if err != nil {
				app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
				return ErrSentinel
			}
		}

		step.Update("Building pipeline execution request...")

		// build the initial job template for running the pipeline
		runJobTemplate := &pb.Job{
			Application:  app.Ref(),
			Workspace:    c.project.WorkspaceRef(),
			TargetRunner: &pb.Ref_Runner{Target: &pb.Ref_Runner_Any{}},
			Variables:    c.variables,
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

		var (
			resp       *pb.RunPipelineResponse
			respGet    *pb.GetPipelineRunResponse
			allRunJobs []string
			steps      int
			runSeq     uint64

			err error
		)
		if !c.flagReattachRun {
			step.Update("Requesting to queue run of pipeline %q...", pipelineIdent)

			// take pipeline id and queue a RunPipeline with a Job Template.
			resp, err = c.project.Client().RunPipeline(c.Ctx, runPipelineReq)
			if err != nil {
				return err
			}

			step.Update("Pipeline %q has started running. Attempting to read job stream sequentially in order", pipelineIdent)
			step.Done()

			steps = len(resp.JobMap)
			allRunJobs = resp.AllJobIds
			runSeq = resp.Sequence
		} else {
			getPipelineReq := &pb.GetPipelineRequest{
				Pipeline: runPipelineReq.Pipeline,
			}

			if c.flagRunSequence == 0 {
				// take pipeline id and queue a RunPipeline with a Job Template.
				respGet, err = c.project.Client().GetLatestPipelineRun(c.Ctx, getPipelineReq)
				if err != nil {
					return err
				}
			} else {
				// take pipeline id and queue a RunPipeline with a Job Template.
				respGet, err = c.project.Client().GetPipelineRun(c.Ctx, &pb.GetPipelineRunRequest{
					Pipeline: getPipelineReq.Pipeline,
					Sequence: uint64(c.flagRunSequence),
				})
				if err != nil {
					return err
				}
			}
			if respGet == nil {
				app.UI.Output("Getting a pipeline run returned a nil response", terminal.WithErrorStyle())
				return fmt.Errorf("Response was empty when requesting a pipeline run for pipeline %q", pipelineIdent)
			}

			step.Update("Attempting to read job stream sequentially in order for run %q", respGet.PipelineRun.Sequence)
			step.Done()

			steps = len(respGet.PipelineRun.Jobs)
			for _, j := range respGet.PipelineRun.Jobs {
				allRunJobs = append(allRunJobs, j.Id)
			}
			runSeq = respGet.PipelineRun.Sequence
		}
		// Receive job ids from running pipeline, use job client to attach to job stream
		// and stream here. First pass can be linear job streaming
		step = sg.Add("")
		defer step.Abort()

		step.Update("%d steps detected, run sequence %d", steps, runSeq)
		step.Done()

		var (
			deployUrl           string
			releaseUrl          string
			inplaceDeploy       bool
			finalVariableValues map[string]*pb.Variable_FinalValue
		)

		successful := steps
		for _, jobId := range allRunJobs {
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
			//app.UI.Output("Executing Step %q in workspace: %q", resp.JobMap[jobId].Step, ws, terminal.WithHeaderStyle())
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

			// Grab the deployment or release URL to display at the end of the pipeline run
			if job.Result.Up != nil {
				deployUrl = job.Result.Up.DeployUrl
				releaseUrl = job.Result.Up.ReleaseUrl
			} else if job.Result.Deploy != nil && job.Result.Deploy.Deployment.Preload != nil {
				deployUrl = job.Result.Deploy.Deployment.Preload.DeployUrl

				// inplace is true if this was an in-place deploy. We detect this
				// if we have a generation that uses a non-matching sequence number
				inplaceDeploy = job.Result.Deploy.Deployment.Generation != nil &&
					job.Result.Deploy.Deployment.Generation.Id != "" &&
					job.Result.Deploy.Deployment.Generation.InitialSequence != job.Result.Deploy.Deployment.Sequence
			} else if job.Result.Release != nil {
				releaseUrl = job.Result.Release.Release.Url
			}

			finalVariableValues = job.VariableFinalValues
		}

		// Show input variable values used in build
		// We do this here so that if the list is long, it doesn't
		// push the deploy/release URLs off the top of the terminal.
		// We also use the deploy result and not the release result,
		// because the data will be the same and this is the deployment command.
		app.UI.Output("Pipeline %q Run v%d Complete", pipelineIdent, runSeq, terminal.WithHeaderStyle())
		app.UI.Output("")
		tbl := fmtVariablesOutput(finalVariableValues)
		c.ui.Table(tbl)

		output := fmt.Sprintf("Pipeline %q (%s) finished! %d/%d steps successfully completed.", pipelineIdent, app.Ref().Project, successful, steps)
		if successful == 0 {
			app.UI.Output("✖ %s", output, terminal.WithErrorStyle())
		} else if successful < steps {
			app.UI.Output("● %s", output, terminal.WithWarningStyle())
		} else {
			app.UI.Output("✔ %s", output, terminal.WithSuccessStyle())
		}

		// Try to get the hostname
		var hostname *pb.Hostname
		hostnamesResp, err := c.project.Client().ListHostnames(ctx, &pb.ListHostnamesRequest{
			Target: &pb.Hostname_Target{
				Target: &pb.Hostname_Target_Application{
					Application: &pb.Hostname_TargetApp{
						Application: app.Ref(),
						Workspace:   c.project.WorkspaceRef(),
					},
				},
			},
		})
		if err == nil && len(hostnamesResp.Hostnames) > 0 {
			hostname = hostnamesResp.Hostnames[0]
		}

		// Output app URL
		app.UI.Output("")
		switch {
		case releaseUrl != "":
			printInplaceInfo(inplaceDeploy, app)
			app.UI.Output("   Release URL: %s", releaseUrl, terminal.WithSuccessStyle())
			if deployUrl != "" {
				app.UI.Output("Deployment URL: https://%s", deployUrl, terminal.WithSuccessStyle())
			} else {
				app.UI.Output(strings.TrimSpace(deployNoURL)+"\n", terminal.WithSuccessStyle())
			}
		case hostname != nil:
			printInplaceInfo(inplaceDeploy, app)
			app.UI.Output("           URL: https://%s", hostname.Fqdn, terminal.WithSuccessStyle())
			app.UI.Output("Deployment URL: https://%s", deployUrl, terminal.WithSuccessStyle())
		case deployUrl != "":
			printInplaceInfo(inplaceDeploy, app)
			app.UI.Output("Deployment URL: https://%s", deployUrl, terminal.WithSuccessStyle())
		default:
			app.UI.Output(strings.TrimSpace(deployNoURL)+"\n", terminal.WithSuccessStyle())
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

		f.BoolVar(&flag.BoolVar{
			Name:    "reattach",
			Target:  &c.flagReattachRun,
			Default: false,
			Usage: "If set, will replay or reattach to an existing pipeline run. If " +
				"'-run' is not specified, will attempt to read the latest run.",
		})

		f.IntVar(&flag.IntVar{
			Name:   "run",
			Target: &c.flagRunSequence,
			Usage:  "Replay or attach to a specific pipeline run by sequence number.",
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

	If '-reattach' is supplied, the CLI will attempt to reattach to an existing
	pipeline run. Defaults to latest, but if '-run' is specified, it will attach
	to that specific run by sequence number.

` + c.Flags().Help())
}
