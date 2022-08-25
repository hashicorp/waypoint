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

	platform              []string `hcl:"platform,optional"`
	skipAdopt             bool     `hcl:"skip_adopt,optional"`
	serverUrl             string   `hcl:"server_url,required"`
	id                    string   `hcl:"id,optional"`
	runnerProfileOdrImage string   `hcl:"odr_image,optional"`
	serverTls             bool     `hcl:"server_tls,optional"`
	serverTlsSkipVerify   bool     `hcl:"server_tls_skip_verify,optional"`
	serverRequireAuth     bool     `hcl:"server_require_auth,optional"`
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

		// Add platforms in alphabetical order. A consistent order is important for repeatable doc generation.
		var sortedPlatformNames []string
		for name := range runnerinstall.Platforms {
			sortedPlatformNames = append(sortedPlatformNames, name)
		}
		sort.Strings(sortedPlatformNames)

		f.EnumVar(&flag.EnumVar{
			Name:   "platform",
			Usage:  "Platform to install the Waypoint runner into. If unset, uses the platform of the local context.",
			Values: sortedPlatformNames,
			Target: &c.platform,
		})

		f.StringVar(&flag.StringVar{
			Name:   "server-addr",
			Usage:  "Address of the Waypoint server.",
			EnvVar: "WAYPOINT_ADDR",
			Target: &c.serverUrl,
		})

		f.StringVar(&flag.StringVar{
			Name:    "odr-image",
			Usage:   "Docker image for the on-demand runners.",
			Default: "hashicorp/waypoint-odr:latest",
			Target:  &c.runnerProfileOdrImage,
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "server-tls",
			Target:  &c.serverTls,
			Usage:   "Connect the runner to the server over TLS.",
			Default: true,
		})

		f.BoolVar(&flag.BoolVar{
			Name:   "server-tls-skip-verify",
			Target: &c.serverTlsSkipVerify,
			Usage:  "Skip TLS verification for runner connection to server.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:   "server-require-auth",
			Target: &c.serverRequireAuth,
			Usage:  "Send authentication details from runner to server.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "skip-adopt",
			Usage:   "Skip adoption of runner after it is installed.",
			Default: false,
			Target:  &c.skipAdopt,
		})

		f.StringVar(&flag.StringVar{
			Name:   "id",
			Usage:  "If this is set, the runner will use the specified id.",
			Target: &c.id,
		})

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

  Install a Waypoint runner to the specified platform: kubernetes, nomad, ecs, 
  or docker.

  This command will attempt to install a runner for the server configured in 
  the current Waypoint context. It will adopt the runner after installation, 
  unless the '-skip-adopt' flag is set to true.

` + c.Flags().Help())
}

func (c *RunnerInstallCommand) Synopsis() string {
	return "Install a Waypoint runner to Kubernetes, Nomad, ECS, or Docker"
}

func (c *RunnerInstallCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoLocalServer(), // no auth in local mode
		WithNoConfig(),
	); err != nil {
		return 1
	}

	ctx := c.Ctx
	log := c.Log.Named("install")
	defer c.Close()

	serverConfig, err := c.project.Client().GetServerConfig(ctx, &empty.Empty{})
	if err != nil {
		c.ui.Output(
			"Error getting server config.",
			clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	// If the user doesn't set a platform, set platform to the platform of the user's context
	platform := c.platform
	if len(platform) == 0 {
		platform = append(platform, serverConfig.Config.Platform)
	}

	p, ok := runnerinstall.Platforms[strings.ToLower(platform[0])]
	if !ok {
		c.ui.Output(
			"Error installing runner into %q: unsupported platform",
			platform[0],
			terminal.WithErrorStyle(),
		)
		return 1
	}

	sg := c.ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Connecting to: %s", c.serverUrl)
	defer func() { s.Abort() }()

	var cookie string
	if !c.skipAdopt {
		cookie = serverConfig.Config.Cookie
	}

	if c.serverUrl == "" {
		c.ui.Output(
			"-server-addr must be supplied for adoption.",
			terminal.WithErrorStyle(),
		)
		return 1
	}

	client := c.project.Client()
	s.Update("Finished connecting to: %s", c.serverUrl)
	s.Status(terminal.StatusOK)
	s.Done()

	// We generate the ID if the user doesn't provide one
	// This ID is used later to adopt the runner
	id := c.id
	if id == "" {
		id, err = server.Id()
		if err != nil {
			c.ui.Output("Error generating runner ID: %s", clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}
	}

	// collect any args after a `--` break to pass forward as secondary flags
	var secondaryArgs []string
	for i, f := range args {
		if f == "--" {
			secondaryArgs = args[(i + 1):]
			break
		}
	}

	s = sg.Add("Installing runner...")
	err = p.Install(ctx, &runnerinstall.InstallOpts{
		Log:        log,
		UI:         c.ui,
		Cookie:     cookie,
		ServerAddr: c.serverUrl,
		AdvertiseClient: &serverconfig.Client{
			Address:       c.serverUrl,
			Tls:           c.serverTls,
			TlsSkipVerify: c.serverTlsSkipVerify,
			RequireAuth:   c.serverRequireAuth,
		},
		Id:               id,
		RunnerAgentFlags: secondaryArgs,
	})
	if err != nil {
		c.ui.Output("Error installing runner: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}
	s.Update("Runner %q installed successfully to %s", id, platform[0])
	s.Status(terminal.StatusOK)
	s.Done()

	if c.skipAdopt {
		c.ui.Output(runnerInstalledButNotYetAdopted, terminal.WithInfoStyle())
	} else {
		err = installutil.AdoptRunner(ctx, c.ui, client, id, c.serverUrl)
		if err != nil {
			c.ui.Output("Error adopting runner: %s", clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}

		// Creating a new runner profile for the newly adopted runner
		var odrConfig *pb.OnDemandRunnerConfig
		s = sg.Add("Creating runner profile and targeting runner %s", strings.ToUpper(id))
		if odc, ok := p.(installutil.OnDemandRunnerConfigProvider); ok {
			odrConfig = odc.OnDemandRunnerConfig()
		} else {
			odrConfig = &pb.OnDemandRunnerConfig{
				Name: platform[0] + "-" + strings.ToUpper(id),
				TargetRunner: &pb.Ref_Runner{
					Target: &pb.Ref_Runner_Id{
						Id: &pb.Ref_RunnerId{
							Id: strings.ToUpper(id),
						},
					},
				},
				OciUrl:     c.runnerProfileOdrImage,
				PluginType: platform[0],
			}
		}
		runnerProfile, err := client.UpsertOnDemandRunnerConfig(ctx, &pb.UpsertOnDemandRunnerConfigRequest{Config: odrConfig})
		if err != nil {
			c.ui.Output("Error creating runner profile: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return 1
		}
		s.Update("Runner profile %q created successfully.", runnerProfile.Config.Name)
		s.Status(terminal.StatusOK)
		s.Done()
	}
	return 0
}

var (
	runnerInstalledButNotYetAdopted = strings.TrimSpace(`The installed runner must be adopted.
Please run "waypoint runner adopt" before the runner can start accepting jobs.
`)
)
