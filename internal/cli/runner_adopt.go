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

type RunnerAdoptCommand struct {
	*baseCommand

	flagJson    bool
	flagPending bool
}

func (c *RunnerAdoptCommand) Run(args []string) int {
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

	_, err := c.project.Client().AdoptRunner(ctx, &pb.AdoptRunnerRequest{
		RunnerId: args[0],
		Adopt:    true,
	})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Runner %q adopted.", args[0], terminal.WithSuccessStyle())
	return 0
}

func (c *RunnerAdoptCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *RunnerAdoptCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *RunnerAdoptCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *RunnerAdoptCommand) Synopsis() string {
	return "Adopt a pending runner"
}

func (c *RunnerAdoptCommand) Help() string {
	return formatHelp(`
Usage: waypoint runner adopt [options] ID

  Adopt a pending runner with the given ID.

  The ID can be retrieved via the "waypoint runner list" command or
  the API. This ID must be the ID of a currently registered runner. The
  runner can be in any state: new, preadopted, adopted, or rejected. This
  will move the runner to the adopted state.

  Once a runner is adopted, that runner ID will remain adopted. Runners may
  restart and will accept jobs immediately, so long as they continue using
  the token they received during the adoption process.

` + c.Flags().Help())
}
