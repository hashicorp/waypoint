package cli

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/golang/protobuf/jsonpb"
	"github.com/posener/complete"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type PipelineInspectCommand struct {
	*baseCommand

	flagJson         bool
	flagPipelineName string
	flagRunSequence  int
}

func (c *PipelineInspectCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
	); err != nil {
		return 1
	}

	if len(c.args) == 0 && c.flagPipelineName == "" {
		c.ui.Output("Pipeline ID or name required.\n\n%s", c.Help(), terminal.WithErrorStyle())
		return 1
	} else if len(c.args) != 0 && c.flagPipelineName != "" {
		c.ui.Output("Both pipeline name and ID were specified, using pipeline name", terminal.WithWarningStyle())
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

	var pipelineRef *pb.Ref_Pipeline
	if c.flagPipelineName != "" {
		pipelineRef = &pb.Ref_Pipeline{
			Ref: &pb.Ref_Pipeline_Owner{
				Owner: &pb.Ref_PipelineOwner{
					Project:      projectRef,
					PipelineName: c.flagPipelineName,
				},
			},
		}
	} else {
		pipelineRef = &pb.Ref_Pipeline{
			Ref: &pb.Ref_Pipeline_Id{
				Id: c.args[0],
			},
		}
	}
	resp, err := c.project.Client().GetPipeline(c.Ctx, &pb.GetPipelineRequest{
		Pipeline: pipelineRef,
	})

	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.ui.Output("Pipeline not found: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle())
			return 1
		}

		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}
	if resp == nil {
		c.ui.Output("The requested pipeline response was empty", terminal.WithWarningStyle())
		return 0
	}

	// rebuild the pipeline ref because ListPipelineRuns only takes ID
	runs, err := c.project.Client().ListPipelineRuns(c.Ctx, &pb.ListPipelineRunsRequest{
		Pipeline: &pb.Ref_Pipeline{
			Ref: &pb.Ref_Pipeline_Id{
				Id: resp.Pipeline.Id,
			},
		},
	})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	if c.flagJson {
		var m jsonpb.Marshaler
		m.Indent = "\t"
		str, err := m.MarshalToString(resp)
		if err != nil {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}

		fmt.Println(str)
		return 0
	}

	var owner string
	switch po := resp.Pipeline.Owner.(type) {
	case *pb.Pipeline_Project:
		owner = po.Project.Project
	default:
		owner = "???"
	}

	c.ui.Output("Pipeline Configuration", terminal.WithHeaderStyle())
	output := []terminal.NamedValue{
		{
			Name: "ID", Value: resp.Pipeline.Id,
		},
		{
			Name: "Name", Value: resp.Pipeline.Name,
		},
		{
			Name: "Owner", Value: owner,
		},
		{
			Name: "Root Step Name", Value: resp.RootStep,
		},
		{
			Name: "Total Steps", Value: len(resp.Pipeline.Steps),
		},
	}

	var (
		sha          string
		msg          string
		startJobTime string
		endJobTime   string
		pipelineRun  *pb.PipelineRun
	)

	seq := uint64(c.flagRunSequence)
	if seq > 0 {
		runResp, err := c.project.Client().GetPipelineRun(c.Ctx, &pb.GetPipelineRunRequest{
			Pipeline: &pb.Ref_Pipeline{
				Ref: &pb.Ref_Pipeline_Id{
					Id: resp.Pipeline.Id,
				},
			},
			Sequence: seq,
		})
		if err != nil {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}

		pipelineRun = runResp.PipelineRun
	} else if len(runs.PipelineRuns) > 0 {
		pipelineRun = runs.PipelineRuns[len(runs.PipelineRuns)-1]
	}

	// Determine start and end job times for requested or latest sequence
	if pipelineRun != nil && len(pipelineRun.Jobs) > 0 {
		startJob, err := c.project.Client().GetJob(c.Ctx, &pb.GetJobRequest{
			JobId: pipelineRun.Jobs[0].Id,
		})
		if err != nil {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}
		endJob, err := c.project.Client().GetJob(c.Ctx, &pb.GetJobRequest{
			JobId: pipelineRun.Jobs[len(pipelineRun.Jobs)-1].Id,
		})
		if err != nil {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}

		startJobTime = humanize.Time(startJob.QueueTime.AsTime())
		endJobTime = humanize.Time(endJob.QueueTime.AsTime())

		if startJob.DataSourceRef != nil {
			sha = startJob.DataSourceRef.Ref.(*pb.Job_DataSource_Ref_Git).Git.Commit
			msg = startJob.DataSourceRef.Ref.(*pb.Job_DataSource_Ref_Git).Git.CommitMessage
		}
	}

	if pipelineRun != nil {
		if seq > 0 {
			output = append(output, []terminal.NamedValue{
				{
					Name: "Run Sequence", Value: pipelineRun.Sequence,
				},
				{
					Name: "Jobs Queued", Value: pipelineRun.Jobs,
				},
				{
					Name: "Run Started", Value: startJobTime,
				},
				{
					Name: "Run Completed", Value: endJobTime,
				},
				{
					Name: "State", Value: pipelineRun.State,
				},
				{
					Name: "Git Commit SHA", Value: sha,
				},
				{
					Name: "Git Commit Message", Value: msg,
				},
			}...)
		} else {
			lastRun := runs.PipelineRuns[len(runs.PipelineRuns)-1]
			totalRuns := int(lastRun.Sequence)
			output = append(output, []terminal.NamedValue{
				{
					Name: "Total Runs", Value: totalRuns,
				},
				{
					Name: "Last Run Started", Value: startJobTime,
				},
				{
					Name: "Last Run Completed", Value: endJobTime,
				},
				{
					Name: "Last Run State", Value: lastRun.State,
				},
				{
					Name: "Last Run Git Commit SHA", Value: sha,
				},
				{
					Name: "Last Run Git Commit Message", Value: msg,
				},
			}...)
		}
	}

	c.ui.NamedValues(output, terminal.WithInfoStyle())

	// TODO(briancain): Use graphviz to build a pipeline graph and display in the terminal?
	return 0
}

func (c *PipelineInspectCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:    "json",
			Target:  &c.flagJson,
			Default: false,
			Usage:   "Output the Pipeline as json.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "name",
			Target:  &c.flagPipelineName,
			Default: "",
			Usage:   "Inspect a pipeline by name.",
		})

		f.IntVar(&flag.IntVar{
			Name:   "run",
			Target: &c.flagRunSequence,
			Usage:  "Inspect a specific run sequence.",
		})
	})
}

func (c *PipelineInspectCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *PipelineInspectCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PipelineInspectCommand) Synopsis() string {
	return "Inspect the full details of a pipeline by id"
}

func (c *PipelineInspectCommand) Help() string {
	return formatHelp(`
Usage: waypoint pipeline inspect [options] <pipeline-id>

  Inspect the full details of a pipeline by id or name for a project.
  Including the '-run' option will show full details of a specific run of this pipeline. 

` + c.Flags().Help())
}
