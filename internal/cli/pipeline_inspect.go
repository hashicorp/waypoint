package cli

import (
	"fmt"

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

	pipelineRef := &pb.Ref_Pipeline{}
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
				Id: &pb.Ref_PipelineId{
					Id: c.args[0],
				},
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
	c.ui.NamedValues([]terminal.NamedValue{
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
	}, terminal.WithInfoStyle())

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

  Inspect the full details of a pipeline by id for a project.

` + c.Flags().Help())
}
