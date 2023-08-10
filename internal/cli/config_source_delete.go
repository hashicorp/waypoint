// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

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

type ConfigSourceDeleteCommand struct {
	*baseCommand

	flagType  string
	flagScope string
}

func (c *ConfigSourceDeleteCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
	); err != nil {
		return 1
	}

	// type is required if deleting
	if c.flagType == "" {
		c.ui.Output(c.Flags().Help(), terminal.WithErrorStyle())
		return 1
	}

	configSource := &pb.ConfigSource{
		Delete: true,
		Type:   c.flagType,
	}

	// If we have a workspace flag set, set that.
	if v := c.flagWorkspace; v != "" {
		configSource.Workspace = &pb.Ref_Workspace{
			Workspace: v,
		}
	}

	// Pre-calculate our project ref since we reuse this.
	projectRef := &pb.Ref_Project{Project: c.flagProject}

	// Depending on the scoping set our target
	switch c.flagScope {
	case "global":
		configSource.Scope = &pb.ConfigSource_Global{
			Global: &pb.Ref_Global{},
		}

	case "project":
		configSource.Scope = &pb.ConfigSource_Project{
			Project: projectRef,
		}

	case "app":
		if c.flagApp == "" {
			fmt.Fprintf(os.Stderr, "-scope requires -app set if scope is 'app'")
			return 1
		}
		configSource.Scope = &pb.ConfigSource_Application{
			Application: &pb.Ref_Application{
				Project:     projectRef.Project,
				Application: c.flagApp,
			},
		}

	default:
		err := fmt.Errorf("-scope needs to be one of 'global', 'project', or 'app'")
		c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	// Set our config
	client := c.project.Client()
	_, err := client.DeleteConfigSource(c.Ctx, &pb.DeleteConfigSourceRequest{
		ConfigSource: configSource,
	})
	if err != nil {
		c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Configuration deleted for dynamic source %q!", c.flagType, terminal.WithSuccessStyle())
	return 0
}

func (c *ConfigSourceDeleteCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.StringVar(&flag.StringVar{
			Name:   "type",
			Target: &c.flagType,
			Usage:  "Dynamic source type to delete, such as 'vault'.",
		})
		f.StringVar(&flag.StringVar{
			Name:   "scope",
			Target: &c.flagScope,
			Usage: "The scope for this configuration source. The configuration source will only " +
				"delete within this scope. This can be one of 'global', 'project', or " +
				"'app'.",
			Default: "global",
		})
	})
}

func (c *ConfigSourceDeleteCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ConfigSourceDeleteCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ConfigSourceDeleteCommand) Synopsis() string {
	return "Delete the configuration for a dynamic source plugin"
}

func (c *ConfigSourceDeleteCommand) Help() string {
	return formatHelp(`
Usage: waypoint config source-delete [options]

  Delete the configuration for a dynamic configuration source plugin.

  To use this command, you should specify a "-type" flag. Please see the
  documentation for the config source type you're configuring for details on
  what configuration fields are available.

  Configuration for this command is global. The "-app", "-project", and
  "-workspace" flags are ignored on this command.

` + c.Flags().Help())
}
