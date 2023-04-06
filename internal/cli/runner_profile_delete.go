// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"github.com/posener/complete"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type RunnerProfileDeleteCommand struct {
	*baseCommand

	// NOTE(izaak): prior to this change, runner profiles were only ever deleted by name, and
	// we don't have a great record of enforcing uniqueness of runner profile names.
	// Ideally, you would only ever need to delete by name, but some existing users may
	// have more than one profile with the same name, and might need this functionality.
	// They'd have to do some work to discover the profile id though - we don't surface it in the CLI.
	// We can likely deprecate this flag in the future.
	flagId string
}

func (c *RunnerProfileDeleteCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoLocalServer(), // no auth in local mode
		WithNoConfig(),
	); err != nil {
		return 1
	}

	if len(c.args) == 0 && c.flagId == "" {
		c.ui.Output("Runner profile name required, or id flag required.", terminal.WithErrorStyle())
		return 1
	}
	var name string
	if len(c.args) > 0 {
		name = c.args[0]
	}

	_, err := c.project.Client().DeleteOnDemandRunnerConfig(c.Ctx, &pb.DeleteOnDemandRunnerConfigRequest{
		Config: &pb.Ref_OnDemandRunnerConfig{
			Name: name,
			Id:   c.flagId,
		},
	})
	if err != nil && status.Code(err) != codes.NotFound {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	if status.Code(err) == codes.NotFound {
		c.ui.Output("runner profile not found", terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Runner profile deleted", terminal.WithHeaderStyle())

	return 0
}

func (c *RunnerProfileDeleteCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(sets *flag.Sets) {
		f := sets.NewSet("Command Options")

		f.StringVar(&flag.StringVar{
			Name:    "id",
			Target:  &c.flagId,
			Default: "",
			Usage:   "The id of the runner profile to delete.",
		})
	})
}

func (c *RunnerProfileDeleteCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *RunnerProfileDeleteCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *RunnerProfileDeleteCommand) Synopsis() string {
	return "Delete a runner profile."
}

func (c *RunnerProfileDeleteCommand) Help() string {
	return formatHelp(`
Usage: waypoint runner profile delete <name>

  Delete the specified runner profile.

`)
}
