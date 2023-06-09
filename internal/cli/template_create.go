package cli

import (
	"os"

	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type ProjectTemplateCreateCommand struct {
	*baseCommand

	flagSummary                    string
	flagExpandedSummary            string
	flagReadmeMarkdownTemplatePath string
	flagWaypointHCLTemplatePath    string
	flagTFCNoCodeModuleSource      string
	flagTFCNoCodeModuleVersion     string
	flagTags                       []string
}

func (c *ProjectTemplateCreateCommand) Run(args []string) int {
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

	if len(args) != 1 {
		c.ui.Output("Single argument required.\n\n"+c.Help(), terminal.WithErrorStyle())
		return 1
	}

	name := args[0]

	var template pb.ProjectTemplate

	template.Name = name
	if c.flagSummary != "" {
		template.Summary = c.flagSummary
	}
	if c.flagExpandedSummary != "" {
		template.ExpandedSummary = c.flagExpandedSummary
	}

	if c.flagReadmeMarkdownTemplatePath != "" {
		rmt, err := os.ReadFile(c.flagReadmeMarkdownTemplatePath)
		if err != nil {
			errMsg := "Unable to read readme.md template file: %s"
			if err == os.ErrNotExist {
				errMsg = "Readme template file does not exist: %s"
			}

			c.ui.Output(errMsg, clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}
		template.ReadmeMarkdownTemplate = rmt
	}

	if c.flagWaypointHCLTemplatePath != "" {
		wpt, err := os.ReadFile(c.flagWaypointHCLTemplatePath)
		if err != nil {
			errMsg := "Unable to read waypoint.hcl template file: %s"
			if err == os.ErrNotExist {
				errMsg = "Waypoint.hcl template file does not exist: %s"
			}

			c.ui.Output(errMsg, clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}
		template.WaypointProject = &pb.ProjectTemplate_WaypointProject{
			WaypointHclTemplate: wpt,
		}
	}

	if c.flagTFCNoCodeModuleSource != "" && c.flagTFCNoCodeModuleVersion == "" {
		c.ui.Output("Terraform No Code module version required.\n\n"+c.Help(), terminal.WithErrorStyle())
		return 1
	}
	if c.flagTFCNoCodeModuleSource == "" && c.flagTFCNoCodeModuleVersion != "" {
		c.ui.Output("Terraform No Code module source required.\n\n"+c.Help(), terminal.WithErrorStyle())
		return 1
	}
	if c.flagTFCNoCodeModuleSource != "" && c.flagTFCNoCodeModuleVersion != "" {
		template.TerraformNocodeModule = &pb.ProjectTemplate_TerraformNocodeModule{
			Source:  c.flagTFCNoCodeModuleSource,
			Version: c.flagTFCNoCodeModuleVersion,
		}
	}

	template.Tags = c.flagTags

	_, err := c.project.Client().CreateProjectTemplate(ctx, &pb.CreateProjectTemplateRequest{
		ProjectTemplate: &template,
	})
	if err != nil {
		c.ui.Output("Error creating project template: %s", clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("template created!")

	return 0
}

func (c *ProjectTemplateCreateCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(sets *flag.Sets) {
		f := sets.NewSet("Command Options")

		f.StringVar(&flag.StringVar{
			Name:    "summary",
			Target:  &c.flagSummary,
			Default: "",
			Usage:   "Summary for the project template",
		})

		f.StringVar(&flag.StringVar{
			Name:    "expanded-summary",
			Target:  &c.flagExpandedSummary,
			Default: "",
			Usage:   "Expanded Summary for the project template",
		})

		f.StringVar(&flag.StringVar{
			Name:    "readme-markdown-template-path",
			Target:  &c.flagReadmeMarkdownTemplatePath,
			Default: "",
			Usage:   "Path to a markdown readme template for projects created from a project template",
		})

		f.StringVar(&flag.StringVar{
			Name:    "waypoint-hcl-template-path",
			Target:  &c.flagWaypointHCLTemplatePath,
			Default: "",
			Usage:   "Path to a templated waypoint.hcl file for projects created from a project template",
		})

		f.StringVar(&flag.StringVar{
			Name:    "tfc-nocode-module-source",
			Target:  &c.flagTFCNoCodeModuleSource,
			Default: "",
			Usage:   "The name of the Terraform no-code module from a Terraform registry that the template should use to provision infrastructure for Waypoint projects created from the template",
		})

		f.StringVar(&flag.StringVar{
			Name:    "tfc-nocode-module-version",
			Target:  &c.flagTFCNoCodeModuleVersion,
			Default: "",
			Usage:   "The version of the Terraform no-code module from a Terraform registry that the template should use to provision infrastructure for Waypoint projects created from the template",
		})

		f.StringSliceVar(&flag.StringSliceVar{
			Name:   "tag",
			Target: &c.flagTags,
			Usage:  "A tag to add to the project template",
		})
	})
}

func (c *ProjectTemplateCreateCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ProjectTemplateCreateCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ProjectTemplateCreateCommand) Synopsis() string {
	return "Create a project template."
}

func (c *ProjectTemplateCreateCommand) Help() string {
	return formatHelp(`
Usage: waypoint template create [options] NAME

  Create a project template.

  This will create a new project template with the given options.

  When running this command the -waypoint-hcl-template-path,
  -tfc-nocode-module-source, and -tfc-nocode-module-version flags are required.

` + c.Flags().Help())
}
