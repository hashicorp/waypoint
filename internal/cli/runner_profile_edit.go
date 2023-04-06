// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"github.com/posener/complete"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"

	cli "github.com/hashicorp/waypoint/internal/cli/editor"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type RunnerProfileEditCommand struct {
	*baseCommand
	flagName string
}

func (c *RunnerProfileEditCommand) Run(args []string) int {
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

	// Setup flag name if argument to command is given
	if c.flagName == "" && len(args) == 0 {
		c.ui.Output("Must provide a runner profile name either by '-name' or argument.\n\n%s",
			c.Help(), terminal.WithErrorStyle())
		return 1
	} else if c.flagName != "" && len(args) > 0 {
		c.ui.Output("Cannot set name both via argument and '-name'. Pick one and run the command again.\n\n%s",
			c.Help(), terminal.WithErrorStyle())
		return 1
	} else if c.flagName == "" && len(args) > 0 {
		c.flagName = args[0]
	}

	var (
		od      *pb.OnDemandRunnerConfig
		updated bool
	)

	// NOTE(briancain): !!! IMPORTANT !!!
	// Don't use StepGroups for this CLI package. It will overwrite lines in
	// the terminal editor making it really difficult to read what you are typing.
	// We default to the classic Output class instead so there's no hanging text.

	if c.flagName != "" {
		c.ui.Output("Checking for an existing runner profile: %s", c.flagName)
		// Check for an existing project of the same name.
		resp, err := c.project.Client().GetOnDemandRunnerConfig(ctx, &pb.GetOnDemandRunnerConfigRequest{
			Config: &pb.Ref_OnDemandRunnerConfig{
				Name: c.flagName,
			},
		})
		if status.Code(err) == codes.NotFound {
			// If the error is a not found error, act as though there is no error
			// and the project is nil so that we can handle that later.
			resp = nil
			err = nil
		}
		if err != nil {
			c.ui.Output(
				"Error checking for project: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return 1
		}

		if resp != nil {
			od = resp.Config
			c.ui.Output("Updating runner profile %q (%q)...", od.Name, od.Id)
			updated = true
		} else {
			c.ui.Output("No existing runner profile found for id %q...command will create a new profile", c.flagName)
			od = &pb.OnDemandRunnerConfig{
				Name: c.flagName,
			}
		}
	} else {
		c.ui.Output("Creating new runner profile named %q", c.flagName)
		od = &pb.OnDemandRunnerConfig{
			Name: c.flagName,
		}
	}

	edited, _, err := cli.Run(od.PluginConfig)
	if err != nil {
		c.ui.Output(
			"Error editing runner profile: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	od.PluginConfig = edited

	// Upsert
	resp, err := c.project.Client().UpsertOnDemandRunnerConfig(ctx, &pb.UpsertOnDemandRunnerConfigRequest{
		Config: od,
	})
	if err != nil {
		c.ui.Output(
			"Error upserting runner profile: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	if updated {
		c.ui.Output("Runner profile %q updated", resp.Config.Name, terminal.WithSuccessStyle())
	} else {
		c.ui.Output("Runner profile %q created", resp.Config.Name, terminal.WithSuccessStyle())
	}

	return 0
}

func (c *RunnerProfileEditCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(sets *flag.Sets) {
		f := sets.NewSet("Command Options")

		f.StringVar(&flag.StringVar{
			Name:    "name",
			Target:  &c.flagName,
			Default: "",
			Usage:   "The name of an existing runner profile to update.",
		})
	})
}

func (c *RunnerProfileEditCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *RunnerProfileEditCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *RunnerProfileEditCommand) Synopsis() string {
	return "Edit an existing runner profile."
}

func (c *RunnerProfileEditCommand) Help() string {
	return formatHelp(`
Usage: waypoint runner profile edit [OPTIONS] <profile-name>

  Edit an existing runner profile.
` + c.Flags().Help())
}
