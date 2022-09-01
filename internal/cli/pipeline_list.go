package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/golang/protobuf/jsonpb"
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type PipelineListCommand struct {
	*baseCommand

	flagJson bool
}

func (c *PipelineListCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
	); err != nil {
		return 1
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

	pipelinesResp, err := c.project.Client().ListPipelines(c.Ctx, &pb.ListPipelinesRequest{
		Project: projectRef,
	})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	if len(pipelinesResp.Pipelines) == 0 {
		return 0
	}
	pipelines := pipelinesResp.Pipelines

	if c.flagJson {
		var m jsonpb.Marshaler
		m.Indent = "\t"
		for _, p := range pipelines {
			str, err := m.MarshalToString(p)
			if err != nil {
				c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
				return 1
			}

			fmt.Println(str)
		}
		return 0
	}

	c.ui.Output("Waypoint Pipelines for %s", c.refProject.Project, terminal.WithHeaderStyle())

	tblHeaders := []string{"ID", "Name", "Owner", "Current Steps", "Last Run Started", "Last Run Completed", "State", "Total Runs"}
	tbl := terminal.NewTable(tblHeaders...)

	for _, pipeline := range pipelines {
		var owner string
		switch po := pipeline.Owner.(type) {
		case *pb.Pipeline_Project:
			owner = po.Project.Project
		default:
			owner = "???"
		}

		totalSteps := strconv.Itoa(len(pipeline.Steps))

		pipelineRunsResp, err := c.project.Client().ListPipelineRuns(c.Ctx, &pb.ListPipelineRunsRequest{
			Pipeline: &pb.Ref_Pipeline{
				Ref: &pb.Ref_Pipeline_Id{
					Id: pipeline.Id,
				},
			},
		})
		if err != nil {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}

		var totalRuns string
		var lastRunStart string
		var lastRunEnd string
		var state string

		runs := pipelineRunsResp.PipelineRuns
		if len(runs) > 0 {
			// Note(xx): This will be refactored in a future PR to use a pipeline bundle
			// that caches this information without having to make a call to the pipeline runs endpoint
			// for every pipeline and jobs for every pipeline run.
			lastRun := runs[len(runs)-1]
			totalRuns = strconv.FormatUint(lastRun.Sequence, 10)

			jobs := lastRun.Jobs
			j, err := c.project.Client().GetJob(c.Ctx, &pb.GetJobRequest{
				JobId: jobs[0].Id,
			})
			if err != nil {
				c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
				return 1
			}
			lastRunStart = humanize.Time(j.QueueTime.AsTime())

			j, err = c.project.Client().GetJob(c.Ctx, &pb.GetJobRequest{
				JobId: jobs[len(jobs)-1].Id,
			})
			if err != nil {
				c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
				return 1
			}

			lastRunEnd = humanize.Time(j.CompleteTime.AsTime())
			state = strings.ToLower(lastRun.State.String())
		}

		tblColumn := []string{
			pipeline.Id,
			pipeline.Name,
			owner,
			totalSteps,
			lastRunStart,
			lastRunEnd,
			state,
			totalRuns,
		}

		tbl.Rich(tblColumn, nil)
	}

	c.ui.Table(tbl)

	return 0
}

func (c *PipelineListCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:    "json",
			Target:  &c.flagJson,
			Default: false,
			Usage:   "Output the list of Pipelines as json.",
		})
	})
}

func (c *PipelineListCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *PipelineListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PipelineListCommand) Synopsis() string {
	return "List all pipelines for a project."
}

func (c *PipelineListCommand) Help() string {
	return formatHelp(`
Usage: waypoint pipeline list

  List all pipelines for a project.

` + c.Flags().Help())
}
