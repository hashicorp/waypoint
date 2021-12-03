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
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

type TriggerInspectCommand struct {
	*baseCommand

	flagTriggerName string
	flagTriggerId   string
	flagJson        bool
}

func (c *TriggerInspectCommand) Run(args []string) int {
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

	resp, err := c.project.Client().GetTrigger(ctx, &pb.GetTriggerRequest{
		Ref: &pb.Ref_Trigger{
			Name: c.flagTriggerName,
			Id:   c.flagTriggerId,
		},
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.ui.Output("Trigger configuration not found", clierrors.Humanize(err),
				terminal.WithErrorStyle())
			return 1
		}

		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	if c.flagJson {
		var m jsonpb.Marshaler
		m.Indent = "\t"
		str, err := m.MarshalToString(resp.Trigger)
		if err != nil {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}

		fmt.Println(str)
		return 0
	}

	trigger := resp.Trigger

	var ws, proj, app string
	if trigger.Workspace != nil {
		ws = trigger.Workspace.Workspace
	}
	if trigger.Project != nil {
		proj = trigger.Project.Project
	}
	if trigger.Application != nil {
		app = trigger.Application.Application
	}

	// TODO: might have to get a case statement for figuring out the operation

	c.ui.Output("Trigger URL config:", terminal.WithHeaderStyle())
	c.ui.NamedValues([]terminal.NamedValue{
		{
			Name: "Name", Value: trigger.Name,
		},
		{
			Name: "ID", Value: trigger.Id,
		},
		{
			Name: "Last Time Active", Value: trigger.ActiveTime.String(),
		},
		{
			Name: "Authenticated", Value: trigger.Authenticated,
		},
		{
			Name: "Operation", Value: trigger.Operation,
		},
		{
			Name: "Workspace", Value: ws,
		},
		{
			Name: "Project", Value: proj,
		},
		{
			Name: "Application", Value: app,
		},
		{
			Name: "Labels", Value: trigger.Labels,
		},
		{
			Name: "Description", Value: trigger.Description,
		},
	}, terminal.WithInfoStyle())

	return 0
}

func (c *TriggerInspectCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")

		f.StringVar(&flag.StringVar{
			Name:   "id",
			Target: &c.flagTriggerId,
			Usage:  "The id of the trigger URL to inspect.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:   "json",
			Target: &c.flagJson,
			Usage: "Output trigger URL configuration information as JSON. This includes " +
				"more fields since this is the complete API structure.",
		})
	})
}

func (c *TriggerInspectCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *TriggerInspectCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *TriggerInspectCommand) Synopsis() string {
	return "Inspect a trigger URL from Waypoint server"
}

func (c *TriggerInspectCommand) Help() string {
	return formatHelp(`
Usage: waypoint trigger inspect [options]

  Inspect a trigger URL from Waypoint Server.

` + c.Flags().Help())
}
