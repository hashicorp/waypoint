package cli

import (
	"io/ioutil"
	"strings"

	"github.com/posener/complete"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/cli/datagen"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	configpkg "github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

type InitCommand struct {
	*baseCommand

	project *clientpkg.Project
	cfg     *configpkg.Config
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

	path, err := c.initConfigPath()
	if err != nil {
		c.ui.Output(err.Error(), terminal.WithErrorStyle())
		return 1
	}

	// If we have no config, initialize a new one.
	if path == "" {
		if !c.initNew() {
			return 1
		}

		return 0
	}

	// Steps to run
	steps := []func() bool{
		c.validateConfig,
		c.validateServer,
		c.validateProject,
		c.validatePlugins,
	}
	for _, step := range steps {
		if !step() {
			return 1
		}
	}

	c.ui.Output("")
	c.ui.Output("Project initialized!", terminal.WithStyle(terminal.SuccessBoldStyle))
	c.ui.Output(
		"You may now call 'waypoint up' to deploy your project or\n"+
			"commands such as 'waypoint build' to perform steps individually.",
		terminal.WithSuccessStyle(),
	)

	return 0
}

func (c *InitCommand) initNew() bool {
	data, err := datagen.Asset("init.tpl.hcl")
	if err != nil {
		// Should never happen because it is embedded.
		panic(err)
	}

	if err := ioutil.WriteFile(configpkg.Filename, data, 0644); err != nil {
		c.ui.Output(err.Error(), terminal.WithErrorStyle())
		return false
	}

	c.ui.Output("Initial Waypoint configuration created!", terminal.WithStyle(terminal.SuccessBoldStyle))
	c.ui.Output(strings.TrimSpace(`
No Waypoint configuration was found in this directory.

A sample configuration has been created in the file "waypoint.hcl". This
file is heavily commented to help you get started.

Once you've setup your initial configuration, run "waypoint init" again to
validate the configuration and initialize your project.
`),
		terminal.WithSuccessStyle(),
	)

	return true
}

func (c *InitCommand) validateConfig() bool {
	sg := c.ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Validating configuration file...")
	cfg, err := c.initConfig(false)
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

func (c *InitCommand) validateServer() bool {
	sg := c.ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Validating server credentials...")
	client, err := c.initClient()
	if err != nil {
		c.stepError(s, initStepConnect, err)
		return false
	}
	c.project = client

	if c.project.Local() {
		s.Update("Local mode initialized successfully")
	} else {
		s.Update("Connection to Waypoint server was successful")
	}

	s.Status(terminal.StatusOK)
	s.Done()
	return true
}

func (c *InitCommand) validateProject() bool {
	sg := c.ui.StepGroup()
	defer sg.Wait()

	ref := c.project.Ref()

	s := sg.Add("Checking if project %q is registered...", ref.Project)

	client := c.project.Client()
	resp, err := client.GetProject(c.Ctx, &pb.GetProjectRequest{Project: ref})
	if status.Code(err) == codes.NotFound {
		err = nil
		resp = nil
	}
	if err != nil {
		c.stepError(s, initStepProject, err)
		return false
	}

	var project *pb.Project
	if resp != nil {
		project = resp.Project
	}

	// If the project itself is missing, then register that.
	if project == nil {
		s.Status(terminal.StatusWarn)
		s.Update("Project %q is not registered with the server. Registering...", ref.Project)

		resp, err := client.UpsertProject(c.Ctx, &pb.UpsertProjectRequest{
			Project: &pb.Project{
				Name: ref.Project,
			},
		})
		if err != nil {
			c.stepError(s, initStepProject, err)
			return false
		}
		s.Status(terminal.StatusOK)

		project = resp.Project
	}

	pt := &serverptypes.Project{Project: project}
	for _, app := range c.cfg.Apps {
		if pt.App(app.Name) >= 0 {
			continue
		}

		// Missing an application, register it.
		s.Status(terminal.StatusWarn)
		s.Update("Application %q is not registered with the server. Registering...", app.Name)

		_, err := client.UpsertApplication(c.Ctx, &pb.UpsertApplicationRequest{
			Project: ref,
			Name:    app.Name,
		})
		if err != nil {
			c.stepError(s, initStepProject, err)
			return false
		}
		s.Status(terminal.StatusOK)
	}

	s.Update("Project %q and all apps are registered with the server.", ref.Project)
	s.Status(terminal.StatusOK)
	s.Done()
	return true
}

func (c *InitCommand) validatePlugins() bool {
	sg := c.ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Validating required plugins...")

	_, err := c.project.Validate(c.Ctx, &pb.Job_ValidateOp{})
	if err != nil {
		c.stepError(s, initStepPluginConfig, err)
		return false
	}

	s.Update("Plugins loaded and configured successfully")
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
	initStepPluginConfig
	initStepProject
)

var initStepStrings = map[initStepType]struct {
	Error        string
	ErrorDetails string
	Other        map[string]string
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

	initStepPluginConfig: {
		Error: "Failed to load and validate plugins!",
		ErrorDetails: `
This validation check ensures that you have all the required plugins available
and the configuration for each plugin (if it exists) is valid. The error message
below should tell you which plugin(s) failed.
		`,
	},

	initStepProject: {
		Error: "Error while checking for project registration.",
		ErrorDetails: `
There was an error while the checking if the project and applications
are registered with the Waypoint server. This error may be temporary and
you may retry to init. See the error message below.
		`,

		Other: map[string]string{
			"unregistered-desc": `
The project and apps must be registered prior to performing any operations.
This creates some metadata with the server. We require registration as a
verification that the project/app names are correct and that you're targeting
the correct server.
			`,
		},
	},
}
