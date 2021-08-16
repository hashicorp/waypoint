package cli

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

type ProjectInspectCommand struct {
	*baseCommand

	flagJson bool
}

func (c *ProjectInspectCommand) Run(args []string) int {
	flagSet := c.Flags()
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(flagSet),
		WithConfig(true), // optional config loading
	); err != nil {
		return 1
	}

	cmdArgs := flagSet.Args()
	var projectTarget string
	if len(cmdArgs) > 1 {
		c.ui.Output("No more than 1 argument required.\n\n"+c.Help(), terminal.WithErrorStyle())
		return 1
	} else if len(cmdArgs) == 0 {
		// If we're in a project dir, load the name. Otherwise we'll
		// try the arg passed in
		if c.project.Ref() != nil {
			projectTarget = c.project.Ref().Project
		} else {
			c.ui.Output("Project argument required, and not in a project directory..\n\n"+
				c.Help(), terminal.WithErrorStyle())
			return 1
		}
	} else if len(cmdArgs) == 1 {
		// project requested
		projectTarget = cmdArgs[0]
	}

	err := c.FormatProject(projectTarget)
	if err != nil {
		c.ui.Output("Failed to format project: %s"+
			clierrors.Humanize(err),
			terminal.WithErrorStyle())
		return 1
	}
	return 0
}

func (c *ProjectInspectCommand) FormatProject(projectTarget string) error {
	// Get our API client
	client := c.project.Client()

	resp, err := client.GetProject(c.Ctx, &pb.GetProjectRequest{
		Project: &pb.Ref_Project{
			Project: projectTarget,
		},
	})
	if err != nil {
		return err
	}
	project := resp.Project
	workspaces := resp.Workspaces

	projectHeaders := []string{
		"Project", "Applications", "Workspaces", "Remote Enabled", "Data Source", "Project Poll Enabled",
		"Project Poll Interval", "App Status Poll Enabled", "App Status Poll Interval",
	}

	projectTbl := terminal.NewTable(projectHeaders...)

	var appNames []string
	for _, app := range project.Applications {
		appNames = append(appNames, app.Name)
	}

	var workspaceNames []string
	for _, ws := range workspaces {
		workspaceNames = append(workspaceNames, ws.Workspace.Workspace)
	}

	dataSource := "Local" // if unset, assume local
	if project.DataSource != nil {
		switch project.DataSource.Source.(type) {
		case *pb.Job_DataSource_Local:
			dataSource = "Local"
		case *pb.Job_DataSource_Git:
			dataSource = "Git"
		}
	}

	var projectPollEnabled bool
	var projectPollInterval string
	if project.DataSourcePoll != nil {
		projectPollEnabled = project.DataSourcePoll.Enabled
		projectPollInterval = project.DataSourcePoll.Interval
	}

	var appPollEnabled bool
	var appPollInterval string
	if project.StatusReportPoll != nil {
		appPollEnabled = project.StatusReportPoll.Enabled
		appPollInterval = project.StatusReportPoll.Interval
	}

	statusColor := ""
	columns := []string{
		project.Name,
		strings.Join(appNames, "\n"),
		strings.Join(workspaceNames, "\n"),
		strconv.FormatBool(project.RemoteEnabled),
		dataSource,
		strconv.FormatBool(projectPollEnabled),
		projectPollInterval,
		strconv.FormatBool(appPollEnabled),
		appPollInterval,
	}

	projectTbl.Rich(
		columns,
		[]string{
			statusColor,
		},
	)

	if c.flagJson {
		err = c.outputProjectJson(projectTbl)
		if err != nil {
			return err
		}
	} else {
		c.ui.Table(projectTbl)
	}

	return nil
}

func (c *ProjectInspectCommand) outputProjectJson(t *terminal.Table) error {
	project := c.formatJsonMap(t)

	data, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return err
	}

	c.ui.Output(string(data))

	return nil
}

// Takes a terminal Table and formats it into a map of key values to be used
// for formatting a JSON output response
func (c *ProjectInspectCommand) formatJsonMap(t *terminal.Table) []map[string]interface{} {
	result := []map[string]interface{}{}
	for _, row := range t.Rows {
		c := map[string]interface{}{}

		for j, r := range row {
			// Remove any whitespacess in key
			header := strings.ReplaceAll(t.Headers[j], " ", "")
			c[header] = r.Value
		}
		result = append(result, c)
	}

	return result
}

func (c *ProjectInspectCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(sets *flag.Sets) {
		f := sets.NewSet("Command Options")

		f.BoolVar(&flag.BoolVar{
			Name:   "json",
			Target: &c.flagJson,
			Usage:  "Output project information as JSON.",
		})
	})
}

func (c *ProjectInspectCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ProjectInspectCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ProjectInspectCommand) Synopsis() string {
	return "Inspect the details of a project."
}

func (c *ProjectInspectCommand) Help() string {
	return formatHelp(`
Usage: waypoint project inspect [project-name]

  Inspect the details of a given project.

  Projects usually map to a single version control repository and contain
  exactly one "waypoint.hcl" configuration. A project may contain multiple
  applications.

  A project is registered via the web UI, "waypoint project apply",
  or "waypoint init".

` + c.Flags().Help())
}
