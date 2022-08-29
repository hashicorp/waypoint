package cli

import (
	"sort"
	"strings"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/runnerinstall"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/posener/complete"
)

type RunnerUninstallCommand struct {
	*baseCommand

	platform     string
	mode         string
	serverUrl    string
	serverCookie string
	id           string
}

func (c *RunnerUninstallCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *RunnerUninstallCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *RunnerUninstallCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")

		f.StringVar(&flag.StringVar{
			Name:   "platform",
			Usage:  "Platform to uninstall the Waypoint runner from.",
			Target: &c.platform,
		})

		f.StringVar(&flag.StringVar{
			Name:   "server-addr",
			Usage:  "Address of the Waypoint server.",
			EnvVar: "WAYPOINT_ADDR",
			Target: &c.serverUrl,
		})

		f.StringVar(&flag.StringVar{
			Name:   "id",
			Usage:  "ID of the Waypoint runner.",
			Target: &c.id,
		})

		// Add platforms in alphabetical order. A consistent order is important for repeatable doc generation.
		var sortedPlatformNames []string
		for name := range runnerinstall.Platforms {
			sortedPlatformNames = append(sortedPlatformNames, name)
		}
		sort.Strings(sortedPlatformNames)

		for _, name := range sortedPlatformNames {
			platform := runnerinstall.Platforms[name]
			platformSet := set.NewSet(name + " Options")
			platform.UninstallFlags(platformSet)
		}
	})
}

func (c *RunnerUninstallCommand) Help() string {
	return formatHelp(`
Usage: waypoint runner uninstall [options]

  Uninstall a Waypoint runner from server given a platform and Waypoint runner
  id. The platform should be specified as kubernetes, nomad, ecs, or docker.

  This will forget the runner on the server and remove any of the resources
  created by a runner installation.
` + c.Flags().Help())
}

func (c *RunnerUninstallCommand) Run(args []string) int {
	ctx := c.Ctx
	log := c.Log.Named("uninstall")
	defer c.Close()

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoLocalServer(), // no auth in local mode
		WithNoConfig(),
	); err != nil {
		return 1
	}

	if c.platform == "" {
		c.ui.Output(
			"The -platform flag is required.",
			terminal.WithErrorStyle(),
		)

		return 1
	}

	if c.id == "" {
		c.ui.Output(
			"The -id flag is required",
			terminal.WithErrorStyle(),
		)

		return 1
	}

	p, ok := runnerinstall.Platforms[strings.ToLower(c.platform)]
	if !ok {
		c.ui.Output(
			"Error installing server into %q: unsupported platform",
			c.platform,
			terminal.WithErrorStyle(),
		)

		return 1
	}

	sg := c.ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Uninstalling runner...")
	defer func() { s.Abort() }()

	err := p.Uninstall(ctx, &runnerinstall.InstallOpts{
		Log:        log,
		UI:         c.ui,
		ServerAddr: c.serverUrl,
		Id:         c.id,
	})
	if err != nil {
		s.Update("Error uninstalling runner")
		s.Status(terminal.StatusError)
		s.Done()
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}
	s.Update("Runner %q uninstalled successfully", c.id)
	s.Status(terminal.StatusOK)
	s.Done()

	s = sg.Add("Forgetting runner %q on server...", c.id)
	defer func() { s.Abort() }()

	_, err = c.project.Client().ForgetRunner(ctx, &pb.ForgetRunnerRequest{
		RunnerId: c.id,
	})
	if err != nil {
		s.Update("Couldn't forget runner: %s", clierrors.Humanize(err))
		s.Status(terminal.StatusWarn)
	} else {
		s.Update("Runner %q forgotten on server", c.id)
		s.Status(terminal.StatusOK)
	}
	s.Done()

	// TODO: Remove runner profiles associated solely with this runner

	return 0
}

func (c *RunnerUninstallCommand) Synopsis() string {
	return "Uninstall a Waypoint runner from Kubernetes, Nomad, ECS, or Docker"
}
