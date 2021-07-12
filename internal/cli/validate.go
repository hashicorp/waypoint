package cli

import (
	"strings"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

type ValidateCommand struct {
	*baseCommand
}

func (c *ValidateCommand) Run(args []string) int {
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
		WithClient(false),
	); err != nil {
		return 1
	}

	if c.validateConfig() {
		return 0
	}

	return 1
}

func (c *ValidateCommand) validateConfig() bool {
	sg := c.ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Validating configuration file...")
	cfg, err := c.initConfig("", false)
	if err != nil {
		c.stepError(s, initStepConfig, err)
		return false
	}
	c.cfg = cfg
	c.refProject = &pb.Ref_Project{Project: cfg.Project}

	s.Update("Configuration file appears valid")
	s.Status(terminal.StatusOK)
	s.Done()

	return true
}

func (c *ValidateCommand) stepError(s terminal.Step, step initStepType, err error) {
	stepStrings := initStepStrings[step]

	s.Status(terminal.StatusError)
	s.Update(stepStrings.Error)
	s.Done()
	c.ui.Output("")
	if v := stepStrings.ErrorDetails; v != "" {
		c.ui.Output(strings.TrimSpace(v), terminal.WithErrorStyle())
		c.ui.Output("")
	}
	c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
}

func (c *ValidateCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(sets *flag.Sets) {
	})
}

func (c *ValidateCommand) Synopsis() string {
	return "Validate waypoint.hcl configuration"
}

func (c *ValidateCommand) Help() string {
	return formatHelp(`
Usage: waypoint validate [FILE]

  Validates a waypoint.hcl file.

`)
}
