// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/posener/complete"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type TriggerListCommand struct {
	*baseCommand

	flagTriggerTags []string

	flagJson bool
	flagFull bool
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
		Tags: c.flagTriggerTags,
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
		m := protojson.MarshalOptions{
			Indent: "\t",
		}
		for _, t := range resp.Triggers {
			data, err := m.Marshal(t)
			if err != nil {
				c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
				return 1
			}

			fmt.Println(string(data))
		}
		return 0
	}

	c.ui.Output("Trigger URL Configs", terminal.WithHeaderStyle())

	tblHeaders := []string{"ID", "Name", "Workspace", "Project", "Application", "Operation"}
	if c.flagFull {
		tblHeaders = append(tblHeaders, "Authenticated", "Description", "Tags", "Last Time Active")
	}
	tbl := terminal.NewTable(tblHeaders...)

	for _, t := range resp.Triggers {
		ws := "default"
		if t.Workspace != nil && t.Workspace.Workspace != "" {
			ws = t.Workspace.Workspace
		}

		var proj, app, tags string
		if t.Project != nil {
			proj = t.Project.Project
		}
		if t.Application != nil {
			app = t.Application.Application
		}

		if len(t.Tags) > 0 {
			tags = strings.Join(t.Tags[:], ",")
		}

		var opStr string
		switch triggerOpType := t.Operation.(type) {
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

		var lastActiveTime string
		if t.ActiveTime != nil {
			lastActiveTime = humanize.Time(t.ActiveTime.AsTime())
		}

		tblColumn := []string{
			t.Id,
			t.Name,
			ws,
			proj,
			app,
			opStr,
		}

		if c.flagFull {
			tblColumn = append(tblColumn, strconv.FormatBool(t.Authenticated), t.Description, tags, lastActiveTime)
		}

		tbl.Rich(tblColumn, nil)
	}

	c.ui.Table(tbl)

	return 0
}

func (c *TriggerListCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")

		f.StringSliceVar(&flag.StringSliceVar{
			Name:   "trigger-tag",
			Target: &c.flagTriggerTags,
			Usage: "A collection of tags to filter on. If the requested tag does " +
				"not match any defined trigger URL, it will be omitted from the results. " +
				"Can be specified multiple times.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:   "full",
			Target: &c.flagFull,
			Usage:  "Output the full list of options for a trigger configuration.",
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
	return "List registered trigger URL configurations."
}

func (c *TriggerListCommand) Help() string {
	return formatHelp(`
Usage: waypoint trigger list [options]

  List trigger URL configurations on Waypoint Server.

` + c.Flags().Help())
}
