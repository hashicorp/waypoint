package cli

import (
	"context"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/runnerinstall"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/posener/complete"
	"sort"
	"strings"
	"time"
)

type RunnerInstallCommand struct {
	*baseCommand

	platform     string
	mode         string
	serverUrl    string
	serverCookie string
	id           string
}

func (c *RunnerInstallCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *RunnerInstallCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *RunnerInstallCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")

		f.StringVar(&flag.StringVar{
			Name:   "platform",
			Usage:  "Platform to install the Waypoint runner into.",
			Target: &c.platform,
		})

		f.StringVar(&flag.StringVar{
			Name:   "server-addr",
			Usage:  "Address of the Waypoint server.",
			EnvVar: "WAYPOINT_ADDR",
			Target: &c.serverUrl,
		})

		// TODO: Determine if adoption or preadoption will be default
		f.StringVar(&flag.StringVar{
			Name:    "mode",
			Usage:   "Installation mode: adoption or preadoption.",
			Default: "adoption",
			Target:  &c.mode,
		})

		f.StringVar(&flag.StringVar{
			Name:   "server-cookie",
			Usage:  "Server cookie for the Waypoint cluster for which you're targeting this runner.",
			Target: &c.serverCookie,
		})

		f.StringVar(&flag.StringVar{
			Name:   "id",
			Usage:  "If this is set, the runner will use the specified id.",
			Target: &c.id,
		})

		// Add platforms in alphabetical order. A consistent order is important for repeatable doc generation.
		i := 0
		sortedPlatformNames := make([]string, len(runnerinstall.Platforms))
		for name := range runnerinstall.Platforms {
			sortedPlatformNames[i] = name
			i++
		}
		sort.Strings(sortedPlatformNames)

		for _, name := range sortedPlatformNames {
			platform := runnerinstall.Platforms[name]
			platformSet := set.NewSet(name + " Options")
			platform.InstallFlags(platformSet)
		}
	})
}

// TODO: Add description
func (c *RunnerInstallCommand) Help() string {
	return formatHelp(`
Usage: waypoint runner install [options]
` + c.Flags().Help())
}

func (c *RunnerInstallCommand) Run(args []string) int {
	ctx := c.Ctx
	log := c.Log.Named("install")
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

	p, ok := runnerinstall.Platforms[strings.ToLower(c.platform)]
	if !ok {
		if c.platform == "" {
			c.ui.Output(
				"The -platform flag is required.",
				terminal.WithErrorStyle(),
			)

			return 1
		}

		c.ui.Output(
			"Error installing server into %q: unsupported platform",
			c.platform,
			terminal.WithErrorStyle(),
		)

		return 1
	}

	if c.mode != "adoption" && c.mode != "preadoption" {
		c.ui.Output(
			"Unsupported runner install mode.",
			terminal.WithErrorStyle(),
		)
		return 1
	}

	if c.mode == "adoption" && c.serverCookie == "" {
		c.ui.Output(
			"Server cookie must be supplied for adoption mode.",
			terminal.WithErrorStyle(),
		)
		return 1
	}

	client := c.project.Client()

	token := &pb.NewTokenResponse{}
	if c.mode == "adoption" {
		log.Debug("Generating runner token.")
		token, err := client.GenerateRunnerToken(ctx, &pb.GenerateRunnerTokenRequest{
			Duration: "",
			Id:       "",
			Labels:   nil,
		})
		if err != nil {
			c.ui.Output("Error generating runner token: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
		}
		token.GetToken()
	}

	log.Debug("Installing runner.")
	err := p.Install(ctx, &runnerinstall.InstallOpts{
		Log:             log,
		UI:              c.ui,
		AuthToken:       token.Token,
		Cookie:          c.serverCookie,
		ServerAddr:      c.serverUrl,
		AdvertiseClient: nil,
		Id:              c.id,
	})
	if err != nil {
		c.ui.Output("Error installing runner: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}
	log.Debug("Runner installed.")

	// TODO: Only run the below in adoption mode
	var id string
	if c.id == "" {
		// TODO: Get the randomly generated ID of the newly installed runner
	} else {
		id = c.id
	}

	// Waits 5 minutes for the server to detect the new runner before timing out
	d := time.Now().Add(time.Minute * time.Duration(5))
	ctx, cancel := context.WithDeadline(ctx, d)
	defer cancel()
	ticker := time.NewTicker(5 * time.Second)
	for true {
		select {
		case <-ticker.C:
		case <-ctx.Done():
			c.ui.Output("Cancelled.",
				terminal.WithErrorStyle(),
			)
			return 1
		}
		runner, err := client.GetRunner(ctx, &pb.GetRunnerRequest{RunnerId: id})
		if err != nil {
			c.ui.Output("Error getting runner: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return 1
		}
		if runner == nil {
			continue
		}
		_, err = client.AdoptRunner(ctx, &pb.AdoptRunnerRequest{
			RunnerId: runner.Id,
			Adopt:    true,
		})
		if err != nil {
			c.ui.Output("Error adopting runner: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return 1
		}
		return 0
	}
	return 0
}

func (c *RunnerInstallCommand) Synopsis() string {
	return "Installs a Waypoint runner to Kubernetes, Nomad, ECS, or Docker"
}
