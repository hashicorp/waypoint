package cli

import (
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/installutil"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/runnerinstall"
	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverconfig"
	"github.com/posener/complete"
	empty "google.golang.org/protobuf/types/known/emptypb"
	"sort"
	"strings"
)

type RunnerInstallCommand struct {
	*baseCommand

	platform            string
	adopt               bool
	serverUrl           string
	serverCookie        string
	id                  string
	serverTls           bool
	serverTlsSkipVerify bool
	serverRequireAuth   bool
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

		f.BoolVar(&flag.BoolVar{
			Name:    "server-tls",
			Target:  &c.serverTls,
			Usage:   "If true, the Waypoint runner will connect to the server over TLS.",
			Default: true,
		})

		f.BoolVar(&flag.BoolVar{
			Name:   "server-tls-skip-verify",
			Target: &c.serverTlsSkipVerify,
			Usage:  "If true, will not validate TLS cert presented by the server.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:   "server-require-auth",
			Target: &c.serverRequireAuth,
			Usage:  "If true, will send authentication details.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "adopt",
			Usage:   "Adopt the runner after it is installed.",
			Default: true,
			Target:  &c.adopt,
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

func (c *RunnerInstallCommand) Help() string {
	return formatHelp(`
Usage: waypoint runner install [options]

  Installs a Waypoint runner to an existing platform. The platform should be
  specified as kubernetes, nomad, ecs, or docker.

  By default, this will adopt the runner after it is installed. The install will
  attempt to install a runner for the server configured in the current Waypoint 
  context.

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

	sg := c.ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Connecting to: %s", c.serverUrl)
	defer func() { s.Abort() }()

	if c.adopt && c.serverCookie == "" {
		c.ui.Output(
			"Server cookie must be supplied for adoption.",
			terminal.WithErrorStyle(),
		)
		return 1
	}

	client := c.project.Client()
	conn, err := client.GetServerConfig(ctx, &empty.Empty{})
	if err != nil {
		c.ui.Output("Error getting server config: %s", clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}
	s.Update("Retrieved server config")
	s.Status(terminal.StatusOK)
	s.Done()

	// TODO: Evaluate if generating a token for non-adoption mode is necessary
	token := &pb.NewTokenResponse{}
	if !c.adopt {
		s = sg.Add("Generating runner token...")
		token, err = client.GenerateRunnerToken(ctx, &pb.GenerateRunnerTokenRequest{
			Duration: "",
			Id:       "",
			Labels:   nil,
		})
		if err != nil {
			c.ui.Output("Error generating runner token: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
		}
		s.Update("Runner token generated")
		s.Status(terminal.StatusOK)
		s.Done()
	}

	// We generate the ID if the user doesn't provide one
	// This ID is used later to adopt the runner
	var id string
	if c.id == "" {
		id, err = server.Id()
		if err != nil {
			c.ui.Output("Error generating runner ID: %s", clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}
	} else {
		id = c.id
	}

	s = sg.Add("Installing runner...")
	err = p.Install(ctx, &runnerinstall.InstallOpts{
		Log:        log,
		UI:         c.ui,
		Cookie:     c.serverCookie,
		ServerAddr: c.serverUrl,
		AdvertiseClient: &serverconfig.Client{
			Address:       c.serverUrl,
			Tls:           c.serverTls,
			TlsSkipVerify: c.serverTlsSkipVerify,
			RequireAuth:   c.serverRequireAuth,
			AuthToken:     token.Token,
		},
		Id: id,
	})
	if err != nil {
		c.ui.Output("Error installing runner: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}
	s.Update("Runner %s installed successfully", id)
	s.Status(terminal.StatusOK)
	s.Done()

	if c.adopt {
		err = installutil.AdoptRunner(ctx, c.ui, client, id)
		if err != nil {
			c.ui.Output("Error adopting runner: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return 1
		}

		s = sg.Add("Creating runner profile and targeting runner %s", strings.ToUpper(id))
		// Creating a new runner profile for the newly adopted runner
		runnerProfile, err := client.UpsertOnDemandRunnerConfig(ctx, &pb.UpsertOnDemandRunnerConfigRequest{
			Config: &pb.OnDemandRunnerConfig{
				Name: c.platform + "-" + strings.ToUpper(id),
				TargetRunner: &pb.Ref_Runner{
					Target: &pb.Ref_Runner_Id{
						Id: &pb.Ref_RunnerId{
							Id: strings.ToUpper(id),
						},
					},
				},
				OciUrl:               "hashicorp/waypoint-odr:latest",
				EnvironmentVariables: nil,
				PluginType:           c.platform,
				PluginConfig:         nil,
				ConfigFormat:         0,
				Default:              false,
			},
		})
		if err != nil {
			c.ui.Output("Error creating runner profile: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return 1
		}
		s.Update("Runner profile %s created successfully.", runnerProfile.Config.Name)
		s.Status(terminal.StatusOK)
		s.Done()
		return 0
	}
	return 0
}

func (c *RunnerInstallCommand) Synopsis() string {
	return "Installs a Waypoint runner to Kubernetes, Nomad, ECS, or Docker"
}
