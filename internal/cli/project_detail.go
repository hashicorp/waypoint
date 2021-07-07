package cli

import (
	"fmt"
	"github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/posener/complete"
	"strings"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

type ProjectDetailCommand struct {
	*baseCommand
}

func (c *ProjectDetailCommand) Run(args []string) int {
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
	// Require one argument
	if len(args) != 1 {
		c.ui.Output(c.Help(), terminal.WithErrorStyle())
		return 1
	}
	name := args[0]

	resp, err := c.project.Client().GetProject(c.Ctx, &gen.GetProjectRequest{
		Project: &gen.Ref_Project{Project: name},
	})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	if resp == nil {
		c.ui.Output(fmt.Sprintf("Project \"%s\" not found.", name))
		return 0
	}

	headers := []string{
		"Key", "Value",
	}

	tbl := terminal.NewTable(headers...)

	addKV(tbl, "Name", resp.Project.Name)
	addKV(tbl, "Data Source", c.formatDataSourceName(resp.Project.DataSource))
	switch v := resp.Project.DataSource.Source.(type) {
	case *gen.Job_DataSource_Git:
		c.addGitDsTableEntries(tbl, v)
	case *gen.Job_DataSource_Local:
		// Nothing will be printed
	}

	workspace := c.refWorkspace.Workspace

	addKV(tbl,
		fmt.Sprintf("%s Scoped Settings", workspace),
		c.formatScopedSettings(resp.Project.ScopedSettings, workspace),
	)

	c.ui.Table(tbl)
	return 0
}

func (c *ProjectDetailCommand) formatScopedSettings(
	settings map[string]*gen.Project_ScopedProjectSettings,
	workspace string) string {
	if _, ok := settings[workspace]; !ok {
		return "N/A"
	}

	scopedSettings := settings[workspace]

	return scopedSettings.String()
}

func (c *ProjectDetailCommand) addGitDsTableEntries(tbl *terminal.Table, git *gen.Job_DataSource_Git) {
	addKV(tbl, "Git Url", git.Git.Url)
	addKV(tbl, "Git Ref", git.Git.Ref)
	addKV(tbl, "Git Path", git.Git.Path)
	addKV(tbl, "Git Auth Type", c.formatGitAuthTypeName(git.Git.Auth))
	c.addGitAuthData(tbl, git.Git.Auth)
}

func (c *ProjectDetailCommand) addGitAuthData(tbl *terminal.Table, auth interface{}) {
	switch v := auth.(type) {
	case *gen.Job_Git_Basic_:
		addKV(tbl, "Git Username", v.Basic.Username)
		addKV(tbl, "Git Password", strings.Repeat("*", len(v.Basic.Password)))
	case *gen.Job_Git_SSH:
		addKV(tbl, "Git SSH Private Key Password", strings.Repeat("*", len(v.Password)))
	}
}

func (c *ProjectDetailCommand) formatGitAuthTypeName(auth interface{}) string {
	switch auth.(type) {
	case *gen.Job_Git_Basic_:
		return "basic"
	case *gen.Job_Git_SSH:
		return "ssh"
	}
	return "unknown"
}

func (c *ProjectDetailCommand) formatDataSourceName(dataSource *gen.Job_DataSource) string {
	switch dataSource.Source.(type) {
	case *gen.Job_DataSource_Git:
		return "git"
	case *gen.Job_DataSource_Local:
		return "local"
	}
	return "unknown"
}

func addKV(tbl *terminal.Table, key string, val string) {
	tbl.Rich([]string{
		key, val,
	}, nil)
}

func (c *ProjectDetailCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *ProjectDetailCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ProjectDetailCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ProjectDetailCommand) Synopsis() string {
	return "Show details for the specified project."
}

func (c *ProjectDetailCommand) Help() string {
	return formatHelp(`
Usage: waypoint project details PROJECT-NAME

  This command lists information about a specified project.

`)
}
