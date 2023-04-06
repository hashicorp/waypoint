// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	jobstream "github.com/hashicorp/waypoint/internal/jobstream"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

type JobGetStreamCommand struct {
	*baseCommand

	flagJson bool
}

func (c *JobGetStreamCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
	); err != nil {
		return 1
	}

	var jobId string
	if len(c.args) == 0 {
		c.ui.Output("Job ID required.\n\n%s", c.Help(), terminal.WithErrorStyle())
		return 1
	} else {
		jobId = c.args[0]
	}

	sg := c.ui.StepGroup()
	defer sg.Wait()

	step := sg.Add("Reading job stream (jobId: %s) ...", jobId)
	defer step.Abort()

	// Ignore the job result for now
	_, err := jobstream.Stream(c.Ctx, jobId,
		jobstream.WithClient(c.project.Client()),
		jobstream.WithUI(c.ui))

	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	step.Update("Finished reading job stream")
	step.Done()

	return 0
}

func (c *JobGetStreamCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		//f := set.NewSet("Command Options")
	})
}

func (c *JobGetStreamCommand) Synopsis() string {
	return "Attach a local CLI to a job stream by id"
}

func (c *JobGetStreamCommand) Help() string {
	return formatHelp(`
Usage: waypoint job get-stream [options] <job-id>

  Connects the local CLI to an active job stream.

` + c.Flags().Help())
}
