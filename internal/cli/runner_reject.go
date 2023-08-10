// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cli

import (
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type RunnerRejectCommand struct {
	*baseCommand

	flagJson    bool
	flagPending bool
}

func (c *RunnerRejectCommand) Run(args []string) int {
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
		Adopt:    false,
	})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Runner %q rejected.", args[0], terminal.WithSuccessStyle())
	return 0
}

func (c *RunnerRejectCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *RunnerRejectCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *RunnerRejectCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *RunnerRejectCommand) Synopsis() string {
	return "Reject a pending runner"
}

func (c *RunnerRejectCommand) Help() string {
	return formatHelp(`
Usage: waypoint runner reject [options] ID

  Reject a pending or adopted runner with the given ID.

  The ID can be retrieved via the "waypoint runner list" command or
  the API. This ID must be the ID of a currently registered runner. The
  runner can be in any state: new, preadopted, adopted, or rejected. This
  will move the runner to the rejected state.

  A rejected runner will never be sent any configuration or jobs. Runners
  that were previously adopted will continue their currently running jobs
  and then will not receive any further jobs.

` + c.Flags().Help())
}
