package cli

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

type TriggerListCommand struct {
	*baseCommand

	flagTriggerLabels []string

	flagJson bool
}

func (c *TriggerListCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	flagSet := c.Flags()
	if err := c.Init(
		WithArgs(args),
		WithFlags(flagSet),
		WithNoConfig(),
	); err != nil {
		return 1
	}
	ctx := c.Ctx

	req := &pb.ListTriggerRequest{
		Labels: c.flagTriggerLabels,
	}

	if c.flagWorkspace != "" {
		req.Workspace = &pb.Ref_Workspace{
			Workspace: c.flagWorkspace,
		}
	}

	if c.flagProject != "" {
		req.Project = &pb.Ref_Project{
			Project: c.flagProject,
		}
	}

	if c.flagApp != "" {
		req.Application = &pb.Ref_Application{
			Project:     c.flagProject,
			Application: c.flagApp,
		}
	}

	resp, err := c.project.Client().ListTriggers(ctx, req)
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	if len(resp.Triggers) == 0 {
		return 0
	}

	if c.flagJson {
		var m jsonpb.Marshaler
		m.Indent = "\t"
		for _, t := range resp.Triggers {
			str, err := m.MarshalToString(t)
			if err != nil {
				c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
				return 1
			}

			fmt.Println(str)
		}
		return 0
	}

	c.ui.Output("Trigger URL Configs", terminal.WithHeaderStyle())

	tbl := terminal.NewTable("ID", "Name", "Workspace", "Project", "Application", "Operation", "Description", "Labels", "Last Time Active")

	for _, t := range resp.Triggers {
		ws := "default"
		if t.Workspace != nil {
			ws = t.Workspace.Workspace
		}
		var proj, app, labels string
		if t.Project != nil {
			proj = t.Project.Project
		}
		if t.Application != nil {
			app = t.Application.Application
		}

		if len(t.Labels) > 0 {
			labels = strings.Join(t.Labels[:], ",")
		}

		var opStr string
		switch triggerOpType := t.Operation.(type) {
		case *pb.Trigger_Build:
			opStr = "build operation"
		case *pb.Trigger_Push:
			opStr = "push operation"
		case *pb.Trigger_Deploy:
			opStr = "deploy operation"
		case *pb.Trigger_Destroy:
			switch triggerOpType.Destroy.Target.(type) {
			case *pb.Job_DestroyOp_Workspace:
				opStr = "destroy workspace operation"
			case *pb.Job_DestroyOp_Deployment:
				opStr = "destroy deployment operation"
			default:
				opStr = "unknown destroy operation target"
			}
		case *pb.Trigger_Release:
			opStr = "release operation"
		case *pb.Trigger_Up:
			opStr = "up operation"
		case *pb.Trigger_Init:
			opStr = "init operation"
		default:
			opStr = fmt.Sprintf("unknown operation: %T", triggerOpType)
		}

		tbl.Rich([]string{
			t.Id,
			t.Name,
			ws,
			proj,
			app,
			opStr,
			t.Description,
			labels,
			t.ActiveTime.String(),
		}, nil)
	}

	c.ui.Table(tbl)

	return 0
}

func (c *TriggerListCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")

		f.StringSliceVar(&flag.StringSliceVar{
			Name:   "trigger-label",
			Target: &c.flagTriggerLabels,
			Usage: "A collection of labels to filter on. If the requested label does " +
				"not match any defined trigger URL, it will be omitted from the results. " +
				"Can be specified multiple times.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:   "json",
			Target: &c.flagJson,
			Usage: "Output trigger URL configuration list information as JSON. This includes " +
				"more fields since this is the complete API structure.",
		})

	})
}

func (c *TriggerListCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *TriggerListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *TriggerListCommand) Synopsis() string {
	return "List trigger URL configurations on Waypoint server"
}

func (c *TriggerListCommand) Help() string {
	return formatHelp(`
Usage: waypoint trigger list [options]

  List trigger URL configurations on Waypoint Server.

` + c.Flags().Help())
}
