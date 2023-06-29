// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"fmt"
	"os"

	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type ConfigSourceGetCommand struct {
	*baseCommand

	flagType  string
	flagScope string
}

func (c *ConfigSourceGetCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
	); err != nil {
		return 1
	}

	// type is required
	if c.flagScope != "all" && c.flagType == "" {
		c.ui.Output("A source type must be specified with '-type'.\n"+c.Help(), terminal.WithErrorStyle())
		return 1
	}

	getConfigSourceRequest := &pb.GetConfigSourceRequest{
		Type: c.flagType,
		Workspace: &pb.Ref_Workspace{
			Workspace: c.flagWorkspace,
		},
	}

	switch c.flagScope {
	case "global":
		getConfigSourceRequest.Scope = &pb.GetConfigSourceRequest_Global{
			Global: &pb.Ref_Global{},
		}

	case "project":
		getConfigSourceRequest.Scope = &pb.GetConfigSourceRequest_Project{
			Project: &pb.Ref_Project{
				Project: c.flagProject,
			},
		}

	case "app":
		if c.flagApp == "" {
			fmt.Fprintf(os.Stderr, "-scope requires -app set if scope is 'app'")
			return 1
		}

		getConfigSourceRequest.Scope = &pb.GetConfigSourceRequest_Application{
			Application: &pb.Ref_Application{
				Application: c.flagApp,
				Project:     c.flagProject,
			},
		}

	case "all":
		getConfigSourceRequest.Scope = &pb.GetConfigSourceRequest_All{All: true}

	default:
		err := fmt.Errorf("-scope needs to be one of 'global', 'project', or 'app'")
		c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	// Get our config source
	client := c.project.Client()
	resp, err := client.GetConfigSource(c.Ctx, getConfigSourceRequest)
	if err != nil {
		c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	var table *terminal.Table
	var scope, project, app, workspace string
	if c.flagScope == "all" {
		if len(resp.ConfigSources) == 0 {
			c.project.UI.Output(
				"No dynamic config sources are configured.\nUse the command "+
					"\"waypoint config source-set\" to add config sources.",
				terminal.WithWarningStyle())
			return 0
		}
		table = terminal.NewTable("Type", "Scope", "Project", "App", "Workspace")
		for _, cs := range resp.ConfigSources {
			switch ref := cs.Scope.(type) {
			case *pb.ConfigSource_Global:
				scope = "global"
			case *pb.ConfigSource_Project:
				scope = "project"
				project = ref.Project.Project
			case *pb.ConfigSource_Application:
				scope = "app"
				project = ref.Application.Project
				app = ref.Application.Application
			}

			if cs.Workspace != nil {
				workspace = cs.Workspace.Workspace
			}
			table.Rich([]string{
				cs.Type,
				scope,
				project,
				app,
				workspace,
			}, []string{
				"",
				"",
				"",
				"",
				"",
			})
		}
	} else {
		if len(resp.ConfigSources) == 0 {
			c.project.UI.Output(
				"Dynamic config source %q is not configured.\n\n"+
					"Note that this doesn't mean that this config source is not usable.\n"+
					"Many config sources work with no explicitly set configurations.",
				c.flagType, terminal.WithErrorStyle())
			return 1
		}

		// we use the first value because this will be the most specific since
		// we do a prefix search.
		cs := resp.ConfigSources[len(resp.ConfigSources)-1]
		switch ref := cs.Scope.(type) {
		case *pb.ConfigSource_Global:
			scope = "global"
		case *pb.ConfigSource_Project:
			scope = "project"
			project = ref.Project.Project
		case *pb.ConfigSource_Application:
			scope = "app"
			project = ref.Application.Project
			app = ref.Application.Application
		}

		if cs.Workspace != nil {
			workspace = cs.Workspace.Workspace
		}
		// Show config source info in a flat list where each project option
		//is its own row
		c.ui.Output("Config Source Info:", terminal.WithHeaderStyle())

		// Unset value strings will be omitted automatically
		c.ui.NamedValues([]terminal.NamedValue{
			{
				Name: "Type", Value: cs.Type,
			},
			{
				Name: "Scope", Value: scope,
			},
			{
				Name: "Project", Value: project,
			},
			{
				Name: "App", Value: app,
			},
			{
				Name: "Workspace", Value: workspace,
			},
		}, terminal.WithInfoStyle())
		c.ui.Output("")

		table = terminal.NewTable("Key", "Value")
		for k, v := range cs.Config {
			table.Rich([]string{
				k,
				v,
			}, []string{
				"",
				"",
			})
		}
	}

	c.ui.Table(table)
	return 0
}

func (c *ConfigSourceGetCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.StringVar(&flag.StringVar{
			Name:   "type",
			Target: &c.flagType,
			Usage:  "Dynamic source type to look up, such as 'vault'.",
		})

		f.StringVar(&flag.StringVar{
			Name: "scope",
			Usage: "The scope for this configuration source. The configuration source " +
				"will only appear within this scope. This can be one of 'all', " +
				"'global', 'project', or 'app'.",
			Default: "project",
			Target:  &c.flagScope,
		})
	})
}

func (c *ConfigSourceGetCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ConfigSourceGetCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ConfigSourceGetCommand) Synopsis() string {
	return "Get the configuration for a dynamic source plugin"
}

func (c *ConfigSourceGetCommand) Help() string {
	return formatHelp(`
Usage: waypoint config source-get [options]

  Get the configuration for a dynamic configuration source plugin.

  This does not list the dynamic configuration variables for an application.
  This command is for configuring the plugin that is used to fetch dynamic
  configurations globally for the server.

  To use this command, you must specify a "-type" flag.

  Configuration for this command is global. The "-app", "-project", and
  "-workspace" flags are ignored on this command.

` + c.Flags().Help())
}
