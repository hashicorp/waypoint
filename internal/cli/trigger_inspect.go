// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/posener/complete"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type TriggerInspectCommand struct {
	*baseCommand

	flagTriggerId string
	flagJson      bool
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

	if len(c.args) == 0 {
		c.ui.Output("Trigger ID required.\n\n%s", c.Help(), terminal.WithErrorStyle())
		return 1
	} else {
		c.flagTriggerId = c.args[0]
	}

	ctx := c.Ctx

	resp, err := c.project.Client().GetTrigger(ctx, &pb.GetTriggerRequest{
		Ref: &pb.Ref_Trigger{
			Id: c.flagTriggerId,
		},
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.ui.Output("Trigger configuration not found: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle())
			return 1
		}

		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	if c.flagJson {
		data, err := protojson.MarshalOptions{
			Indent: "\t",
		}.Marshal(resp.Trigger)
		if err != nil {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}

		fmt.Println(string(data))
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

	var opStr string
	switch triggerOpType := trigger.Operation.(type) {
	case *pb.Trigger_Build:
		opStr = "build"
	case *pb.Trigger_Push:
		opStr = "push"
	case *pb.Trigger_Deploy:
		opStr = "deploy"
	case *pb.Trigger_Destroy:
		switch triggerOpType.Destroy.Target.(type) {
		case *pb.Job_DestroyOp_Workspace:
			opStr = "destroy workspace"
		case *pb.Job_DestroyOp_Deployment:
			opStr = "destroy deployment"
		default:
			opStr = "unknown destroy operation target"
		}
	case *pb.Trigger_Release:
		opStr = "release"
	case *pb.Trigger_Up:
		opStr = "up"
	case *pb.Trigger_Init:
		opStr = "init"
	case *pb.Trigger_StatusReport:
		switch triggerOpType.StatusReport.Target.(type) {
		case *pb.Job_StatusReportOp_Deployment:
			opStr = "status report deployment"
		case *pb.Job_StatusReportOp_Release:
			opStr = "status report release"
		}
	default:
		opStr = fmt.Sprintf("unknown operation: %T", triggerOpType)
	}

	var tags string
	if len(trigger.Tags) > 0 {
		tags = strings.Join(trigger.Tags[:], ", ")
	}

	lastActiveTime := "n/a"
	if trigger.ActiveTime.IsValid() {
		lastActiveTime = humanize.Time(trigger.ActiveTime.AsTime())
	}

	c.ui.Output("Trigger URL config:", terminal.WithHeaderStyle())
	c.ui.NamedValues([]terminal.NamedValue{
		{
			Name: "Name", Value: trigger.Name,
		},
		{
			Name: "ID", Value: trigger.Id,
		},
		{
			Name: "Last Time Active", Value: lastActiveTime,
		},
		{
			Name: "Authenticated", Value: trigger.Authenticated,
		},
		{
			Name: "Operation", Value: opStr,
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
			Name: "Tags", Value: tags,
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
	return "Inspect a registered trigger URL configuration."
}

func (c *TriggerInspectCommand) Help() string {
	return formatHelp(`
Usage: waypoint trigger inspect [options] <trigger-id>

  Inspect a registered trigger URL configuration.

` + c.Flags().Help())
}
