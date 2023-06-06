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

type ConfigSourceSetCommand struct {
	*baseCommand

	flagType   string
	flagConfig map[string]string
	flagDelete bool
	flagScope  string
}

func (c *ConfigSourceSetCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
	); err != nil {
		return 1
	}

	// type is required and config is required if we're not deleting.
	if c.flagType == "" || (!c.flagDelete && len(c.flagConfig) == 0) {
		c.ui.Output(c.Flags().Help(), terminal.WithErrorStyle())
		return 1
	}

	configSource := &pb.ConfigSource{
		Delete: c.flagDelete,
		Type:   c.flagType,
		Config: c.flagConfig,
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
	_, err := client.SetConfigSource(c.Ctx, &pb.SetConfigSourceRequest{
		ConfigSource: configSource,
	})
	if err != nil {
		c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Configuration set for dynamic source %q!", c.flagType, terminal.WithSuccessStyle())
	return 0
}

func (c *ConfigSourceSetCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.StringVar(&flag.StringVar{
			Name:   "type",
			Target: &c.flagType,
			Usage:  "Dynamic source type to configure, such as 'vault'.",
		})
		f.StringMapVar(&flag.StringMapVar{
			Name:   "config",
			Target: &c.flagConfig,
			Usage: "Configuration for the dynamic source type. This may be repeated. " +
				"The fields available are dependent on the dynamic source type, so please " +
				"check the documentation for that specific type for more information.",
		})
		f.BoolVar(&flag.BoolVar{
			Name:   "delete",
			Target: &c.flagDelete,
			Usage: "Delete the configuration for this source type. If this is set " +
				"then the -config flag is ignored.",
		})
		f.StringVar(&flag.StringVar{
			Name:   "scope",
			Target: &c.flagScope,
			Usage: "The scope for this configuration source. The configuration source will only " +
				"appear within this scope. This can be one of 'global', 'project', or " +
				"'app'.",
			Default: "global",
		})
	})
}

func (c *ConfigSourceSetCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ConfigSourceSetCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ConfigSourceSetCommand) Synopsis() string {
	return "Set the configuration for a dynamic source plugin"
}

func (c *ConfigSourceSetCommand) Help() string {
	return formatHelp(`
Usage: waypoint config source-set [options]

  Set the configuration for a dynamic configuration source plugin.

  This does not add a dynamic configuration variable to your application.
  This command is for configuring the plugin that is used to fetch dynamic
  configurations globally. For example, configuring authentication information
  or server addresses and so on.

  To use this command, you should specify a "-type" flag along with one or more
  "-config" values. Please see the documentation for the config source type
  you're configuring for details on what configuration fields are available.

  This command overrides all configuration already set for a configuration
  source plugin. When modifying an existing configuration, all desired
  "-config" flags will need to be set each time the command is ran.

` + c.Flags().Help())
}
