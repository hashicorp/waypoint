// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type RunnerForgetCommand struct {
	*baseCommand

	flagJson    bool
	flagPending bool
}

func (c *RunnerForgetCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	flagSet := c.Flags()
	if err := c.Init(
		WithArgs(args),
		WithFlags(flagSet),
		WithNoConfig(),
	); err != nil {
		return 1
	}
	ctx := c.Ctx
	args = flagSet.Args()

	// Require one argument
	if len(args) != 1 {
		c.ui.Output(c.Help(), terminal.WithErrorStyle())
		return 1
	}

	_, err := c.project.Client().ForgetRunner(ctx, &pb.ForgetRunnerRequest{
		RunnerId: args[0],
	})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Runner %q forgotten.", args[0], terminal.WithSuccessStyle())
	return 0
}

func (c *RunnerForgetCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *RunnerForgetCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *RunnerForgetCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *RunnerForgetCommand) Synopsis() string {
	return "Forget a previously registered runner"
}

func (c *RunnerForgetCommand) Help() string {
	return formatHelp(`
Usage: waypoint runner forget [options] ID

  Forget a previously registered runner.

  This will delete any records of a previously registered runner. If the
  runner is currently running, it will begin to error on the next job request.
  On subsequent registrations, the Waypoint server behaves as if it has never
  seen this runner before and triggers the full adoption process again.

` + c.Flags().Help())
}
