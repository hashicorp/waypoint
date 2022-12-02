package cli

import (
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
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
	); err != nil {
		return 1
	}

	resp, err := c.project.Client().ListProjectTemplates(c.Ctx, &pb.ListProjectTemplatesRequest{})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	if len(resp.ProjectTemplates) == 0 {
		return 0
	}

	c.ui.Output("Project templates")

	tbl := terminal.NewTable("Name", "Repo", "Description")

	for _, p := range resp.ProjectTemplates {

		tbl.Rich([]string{
			p.Name,
			p.SourceCodePlatform.(*pb.ProjectTemplate_Github).Github.Source.Repo + "/" + p.SourceCodePlatform.(*pb.ProjectTemplate_Github).Github.Source.Repo,
			p.Description,
		}, nil)
	}

	c.ui.Table(tbl)

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
	return "List all registered project templates."
}

func (c *ProjectTemplateListCommand) Help() string {
	return formatHelp(`
Usage: waypoint project template list

  List project templates.

  Project templates are used as templates for creating projects.
`)
}
