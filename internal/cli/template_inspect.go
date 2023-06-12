// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"fmt"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/posener/complete"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type ProjectTemplateInspectCommand struct {
	*baseCommand

	flagJson bool

	flagID string
}

func (c *ProjectTemplateInspectCommand) Run(args []string) int {
	flagSet := c.Flags()
	if err := c.Init(
		WithArgs(args),
		WithFlags(flagSet),
		WithNoConfig(),
	); err != nil {
		return 1
	}
	args = flagSet.Args()
	ctx := c.Ctx

	if len(args) > 1 {
		c.ui.Output("Only one project template may be specified at a time.\n\n"+c.Help(), terminal.WithErrorStyle())
		return 1
	}

	name := ""
	if len(args) == 1 {
		name = args[0]
	}

	if name != "" && c.flagID != "" {
		c.ui.Output("Name argument and id flag may not be specified together.\n\n"+c.Help(), terminal.WithErrorStyle())
		return 1
	}
	if name == "" && c.flagID == "" {
		c.ui.Output("Missing project template name or id.\n\n"+c.Help(), terminal.WithErrorStyle())
		return 1
	}

	out, _, err := c.ui.OutputWriters()
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	var tref pb.Ref_ProjectTemplate
	if name != "" {
		tref.Ref = &pb.Ref_ProjectTemplate_Name{
			Name: name,
		}
	}
	if c.flagID != "" {
		tref.Ref = &pb.Ref_ProjectTemplate_Id{
			Id: c.flagID,
		}
		name = c.flagID
	}

	tr, err := c.project.Client().GetProjectTemplate(ctx, &pb.GetProjectTemplateRequest{
		ProjectTemplate: &tref,
	})
	if err != nil {
		errMsg := clierrors.Humanize(err)
		if status.Code(err) == codes.NotFound || tr.ProjectTemplate == nil {
			errMsg = fmt.Sprintf("Project template %q does not exist", name)
		}
		c.ui.Output(errMsg, terminal.WithErrorStyle())
		return 1
	}
	template := tr.ProjectTemplate

	if c.flagJson {
		data, err := protojson.MarshalOptions{
			Indent: "\t",
		}.Marshal(template)
		if err != nil {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}

		c.ui.Output(string(data))
		return 0
	}

	table := tablewriter.NewWriter(out)
	table.SetHeader([]string{
		"ID",
		"Name",
		"Summary",
		"Terraform Module",
		"Terraform Module Version",
		"Tags",
	})
	table.SetBorder(false)

	table.Rich([]string{
		template.Id,
		template.Name,
		template.Summary,
		template.TerraformNocodeModule.Source,
		template.TerraformNocodeModule.Version,
		strings.Join(template.Tags, ", "),
	}, []tablewriter.Colors{
		{},
		{},
		{},
		{},
		{},
	})

	table.Render()

	return 0
}

func (c *ProjectTemplateInspectCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(sets *flag.Sets) {
		f := sets.NewSet("Command Options")

		f.BoolVar(&flag.BoolVar{
			Name:   "json",
			Target: &c.flagJson,
			Usage:  "Output project information as JSON.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "id",
			Target:  &c.flagID,
			Default: "",
			Usage:   "Id of project template. Mutually exclusive with name argument.",
		})
	})

}

func (c *ProjectTemplateInspectCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ProjectTemplateInspectCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ProjectTemplateInspectCommand) Synopsis() string {
	return "View a single project template"
}

func (c *ProjectTemplateInspectCommand) Help() string {
	return formatHelp(`
Usage: waypoint template inspect [options] [NAME]

  Show detailed information for a single project template given a name or ID.

` + c.Flags().Help())
}
