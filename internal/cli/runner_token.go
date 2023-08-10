// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type RunnerTokenCommand struct {
	*baseCommand

	flagDuration time.Duration
	flagId       string
	flagLabels   map[string]string
}

func (c *RunnerTokenCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
		WithNoLocalServer(), // local mode has no need for tokens
	); err != nil {
		return 1
	}

	// Generate the token
	client := c.project.Client()
	resp, err := client.GenerateRunnerToken(c.Ctx, &pb.GenerateRunnerTokenRequest{
		Duration: c.flagDuration.String(),
		Id:       c.flagId,
		Labels:   c.flagLabels,
	})
	if err != nil {
		c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	// We use fmt here and not the UI helpers because UI helpers will
	// trim tokens horizontally on terminals that are narrow.
	fmt.Println(resp.Token)
	return 0
}

func (c *RunnerTokenCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.DurationVar(&flag.DurationVar{
			Name:    "expires-in",
			Target:  &c.flagDuration,
			Usage:   "The duration until the token expires. i.e. '5m'.",
			Default: 720 * time.Hour, // 30 days
		})

		f.StringVar(&flag.StringVar{
			Name:   "id",
			Target: &c.flagId,
			Usage:  "Id to restrict this token to. If empty, all runner IDs are valid.",
		})

		f.StringMapVar(&flag.StringMapVar{
			Name:   "label",
			Target: &c.flagLabels,
			Usage: "Labels that must match the runner for this token to be valid. " +
				"These are set in 'k=v' format and this flag can be repeated to " +
				"set multiple labels. If no labels are set, runners with any labels " +
				"are valid.",
		})
	})
}

func (c *RunnerTokenCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *RunnerTokenCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *RunnerTokenCommand) Synopsis() string {
	return "Request a new token to run a runner"
}

func (c *RunnerTokenCommand) Help() string {
	helpText := `
Usage: waypoint runner token [options]

  Generate a new runner token used for runners.

  This generates a new token that can be used for "waypoint runner agent".
  Generating a token in advance enables the "pre-adoption" mode where a runner
  avoids the manual adoption process and begins accepting work immediately.

  While this makes it possible to do things such as automate runners all the
  way to running, it is generally not recommended since it requires a higher
  level of security to transfer a token around.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}
