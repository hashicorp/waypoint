package cli

import (
	"strings"

	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/sdk/terminal"
	"github.com/posener/complete"
)

type InitCommand struct {
	*baseCommand

	project *clientpkg.Project
}

func (c *InitCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
		WithClient(false),
	); err != nil {
		return 1
	}

	// Steps to run
	steps := []func(terminal.StepGroup) bool{
		c.validateConfig,
		c.validateServer,
	}

	sg := c.ui.StepGroup()
	for _, step := range steps {
		if !step(sg) {
			return 1
		}
	}
	sg.Wait()

	return 0
}

func (c *InitCommand) validateConfig(sg terminal.StepGroup) bool {
	s := sg.Add("Validating configuration file...")
	cfg, err := c.initConfig(false)
	if err != nil {
		c.stepError(s, initStepConfig, err)
		return false
	}
	var _ = cfg

	s.Update("Configuration file appears valid")
	s.Status(terminal.StatusOK)
	s.Done()

	return true
}

func (c *InitCommand) validateServer(sg terminal.StepGroup) bool {
	s := sg.Add("Validating server credentials...")
	client, err := c.initClient()
	if err != nil {
		c.stepError(s, initStepConnect, err)
		return false
	}
	c.project = client

	s.Update("Connection to server successful.")
	s.Status(terminal.StatusOK)
	s.Done()
	return true
}

func (c *InitCommand) stepError(s terminal.Step, step initStepType, err error) {
	stepStrings := initStepStrings[step]

	s.Status(terminal.StatusError)
	s.Update(stepStrings.Error)
	s.Done()
	c.ui.Output("")
	if v := stepStrings.ErrorDetails; v != "" {
		c.ui.Output(strings.TrimSpace(v), terminal.WithErrorStyle())
		c.ui.Output("")
	}
	c.ui.Output(err.Error(), terminal.WithErrorStyle())
}

func (c *InitCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *InitCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *InitCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *InitCommand) Synopsis() string {
	return "Initialize and validate a project."
}

func (c *InitCommand) Help() string {
	helpText := `
Usage: waypoint init [options]

  Initialize and validate a project.

  This is the first command that should be run for any new or existing
  Waypoint project per machine. This sets up the project if required and
  also validates that operations such as "up" will most likely work.

  This command is always safe to run multiple times. This command will never
  delete your configuration or any data in the server.

`

	return strings.TrimSpace(helpText)
}

type initStepType uint

const (
	initStepInvalid initStepType = iota
	initStepConfig
	initStepConnect
)

var initStepStrings = map[initStepType]struct {
	Error        string
	ErrorDetails string
}{
	initStepConfig: {
		Error: "Error loading configuration!",
	},

	initStepConnect: {
		Error: "Failed to initialize client for Waypoint server.",
		ErrorDetails: `
The Waypoint client validation step validates that we can connect to the
configured Waypoint server. If this is a local-only operation (no Waypoint
server is configured), then we validate that we can initialize local writes.
The error for this failure is shown below.
			`,
	},
}
