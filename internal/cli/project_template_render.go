package cli

import (
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/jobstream"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type ProjectTemplateRenderCommand struct {
	*baseCommand

	flagFromTemplate string

	flagProjectName string

	flagProjectDescription string

	flagGithubRepoOwner string
}

func (c *ProjectTemplateRenderCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
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

	sg := c.ui.StepGroup()
	defer sg.Wait()

	// TODO: check for existing project and fail if one exists

	// TODO: validate flagFromProject set
	resp, err := c.project.Client().GetProjectTemplate(ctx, &pb.GetProjectTemplateRequest{
		ProjectTemplate: &pb.Ref_ProjectTemplate{
			Name: c.flagFromTemplate,
		},
	})
	if err != nil {
		c.ui.Output(
			"Error getting project template: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	// Default to template repo owner
	githubRepoOwner := c.flagGithubRepoOwner
	if githubRepoOwner == "" {
		githubRepoOwner = resp.ProjectTemplate.SourceCodePlatform.(*pb.ProjectTemplate_Github).Github.Source.Owner
	}

	s := sg.Add("Initiating creation of project %s from template %s", c.flagProjectName, c.flagFromTemplate)
	defer func() { s.Abort() }()

	upsertResp, err := c.project.Client().UpsertProjectFromTemplate(ctx, &pb.UpsertProjectFromTemplateRequest{
		ProjectName: c.flagProjectName,
		Description: c.flagProjectDescription,
		SourceCodePlatformDestinationOptions: &pb.UpsertProjectFromTemplateRequest_Github{
			Github: &pb.ProjectTemplate_SourceCodePlatformGithub_Destination_Options{
				Owner: githubRepoOwner,
			},
		},
		Template: &pb.Ref_ProjectTemplate{
			Name: c.flagFromTemplate,
		},
	})
	if err != nil {
		c.ui.Output(
			"Error upserting project template: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	// Ignore the job result for now

	s.Update("Streaming template job %q", upsertResp.JobId)
	_, err = jobstream.Stream(c.Ctx, upsertResp.JobId,
		jobstream.WithClient(c.project.Client()),
		jobstream.WithUI(c.ui))

	if err != nil {
		c.ui.Output("job %q to upsert project failed", upsertResp.JobId, clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	s.Update("Getting project %q", c.flagProjectName)
	// Check for an existing project of the same name.
	getProjResp, err := c.project.Client().GetProject(ctx, &pb.GetProjectRequest{
		Project: &pb.Ref_Project{
			Project: c.flagProjectName,
		},
	})
	if err != nil {
		c.ui.Output("failed to get newly created project %q", c.flagProjectName, clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	// TODO: run `git clone` right here
	// TODO: recommend running `waypoint up`

	s.Update("Project template created! %q", getProjResp.Project.DataSource.Source.(*pb.Job_DataSource_Git).Git.Url)
	s.Done()

	upCommand := &UpCommand{
		baseCommand: c.baseCommand,
	}

	upExitCode := upCommand.Run([]string{"-project", c.flagProjectName, "-app", c.flagProjectName})

	c.ui.Output("Your new project %s is deployed!", c.flagProjectName)
	c.ui.Output("To clone it, run:")
	c.ui.Output("git clone %s", getProjResp.Project.DataSource.Source.(*pb.Job_DataSource_Git).Git.Url)

	c.ui.Output("Then make your change, commit and push, and run:")
	c.ui.Output("waypoint up")

	return upExitCode
}

func (c *ProjectTemplateRenderCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(sets *flag.Sets) {
		f := sets.NewSet("Command Options")

		f.StringVar(&flag.StringVar{
			Name:    "template",
			Target:  &c.flagFromTemplate,
			Default: "",
			Usage:   "todo",
		})

		f.StringVar(&flag.StringVar{
			Name:    "project-name",
			Target:  &c.flagProjectName,
			Default: "",
			Usage:   "todo",
		})

		f.StringVar(&flag.StringVar{
			Name:    "project-description",
			Target:  &c.flagProjectDescription,
			Default: "",
			Usage:   "todo",
		})

		f.StringVar(&flag.StringVar{
			Name:    "github-repo-owner",
			Target:  &c.flagGithubRepoOwner,
			Default: "",
			Usage:   "todo",
		})
	})
}

func (c *ProjectTemplateRenderCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ProjectTemplateRenderCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ProjectTemplateRenderCommand) Synopsis() string {
	return "Create or update a project template."
}

func (c *ProjectTemplateRenderCommand) Help() string {
	return formatHelp(`
Usage: waypoint project template render [options]

  Create or update a project template.

  This will create a new project template with the given options. If a 
  project template with the same name already exists, this will update 
  the existing project template using the fields that are set.

` + c.Flags().Help())
}
