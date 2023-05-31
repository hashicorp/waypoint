// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/posener/complete"
)

type ConfigDeleteCommand struct {
	*baseCommand
}

func (c *ConfigDeleteCommand) Run(args []string) int {
	initOpts := []Option{
		WithArgs(args),
		WithFlags(c.Flags()),

		// Don't allow a local in-mem server because configuration
		// makes no sense with the local server.
		WithNoLocalServer(),
		WithNoConfig(), // no waypoint.hcl
	}

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(initOpts...); err != nil {
		return 1
	}

	// If there are no command arguments, check if the command has
	// been invoked with a pipe like `cat .env | waypoint config set`.
	if len(c.args) == 0 {
		info, err := os.Stdin.Stat()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to get console mode for stdin")
			return 1
		}

		// If there's no pipe, there are no arguments. Fail.
		if info.Mode()&os.ModeNamedPipe == 0 {
			_, _ = fmt.Fprintf(os.Stderr, "config delete requires at least one key entry")
			return 1
		}

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			c.args = append(c.args, scanner.Text())
		}
	}

	// Get our API client
	client := c.project.Client()

	var req pb.ConfigDeleteRequest

	for _, arg := range c.args {
		// Build our initial config var
		configVar := &pb.ConfigVar{
			Target: &pb.ConfigVar_Target{},
			Name:   arg,
			Value:  &pb.ConfigVar_Unset{},
		}

		req.Variables = append(req.Variables, configVar)
	}

	_, err := client.DeleteConfig(c.Ctx, &req)
	if err != nil {
		c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	return 0
}

func (c *ConfigDeleteCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {})
}

func (c *ConfigDeleteCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ConfigDeleteCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ConfigDeleteCommand) Synopsis() string {
	return "Delete a config variable."
}

func (c *ConfigDeleteCommand) Help() string {
	return formatHelp(`
Usage: waypoint config delete <name> <name2> ...

  Delete a config variable from the system. This cannot be undone.

` + c.Flags().Help())
}
