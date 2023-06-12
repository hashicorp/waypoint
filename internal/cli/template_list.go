// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type ProjectTemplateListCommand struct {
	*baseCommand
}

func (c *ProjectTemplateListCommand) Run(args []string) int {
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
	); err != nil {
		return 1
	}
	ctx := c.Ctx

	out, _, err := c.ui.OutputWriters()
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	ptr, err := c.project.Client().ListProjectTemplates(ctx, &pb.ListProjectTemplatesRequest{})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}
	templates := ptr.ProjectTemplates

	if len(templates) == 0 {
		c.ui.Output("No project templates found.")
		return 0
	}

	table := tablewriter.NewWriter(out)
	table.SetHeader([]string{"Name", "Summary", "Terraform Module", "Terraform Module Version", "Tags"})
	table.SetBorder(false)

	for _, template := range templates {
		table.Rich([]string{
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
	}

	table.Render()

	return 0
}

func (c *ProjectTemplateListCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *ProjectTemplateListCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ProjectTemplateListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ProjectTemplateListCommand) Synopsis() string {
	return "List all project templates"
}

func (c *ProjectTemplateListCommand) Help() string {
	return formatHelp(`
Usage: waypoint template list [options]

  Lists all project templates stored on the Waypoint server.

` + c.Flags().Help())
}
