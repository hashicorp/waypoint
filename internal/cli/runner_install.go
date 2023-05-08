package cli

import (
	"sort"
	"strings"

	"github.com/posener/complete"
	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"

	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/installutil"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/runnerinstall"
	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverconfig"
)

type RunnerInstallCommand struct {
	*baseCommand

	platform              []string `hcl:"platform,optional"`
	skipAdopt             bool     `hcl:"skip_adopt,optional"`
	serverUrl             string   `hcl:"server_url"`
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
			Default: installutil.DefaultODRImage,
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

  To further customize the runner installation, you may pass advanced flag
  options specified in the documentation for the 'runner agent' command. To set
  these values, include a '--' after the full argument list for 'install',
  followed by these advanced flag options. As an example, to set a label k/v
  on the runner profile that is generated as part of adopting a runner during
  the install, the command would be:

    waypoint runner install -server-addr=localhost:9701 -server-tls-skip-verify -- -label=environment=primary

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
	var runnerPlatform string
	if len(c.platform) == 0 {
		runnerPlatform = serverConfig.Config.Platform
	}

	p, ok := runnerinstall.Platforms[strings.ToLower(runnerPlatform)]
	if !ok {
		c.ui.Output(
			"Error installing runner into %q: unsupported platform",
			runnerPlatform,
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

	// This loop is used to parse any arguments supplied after `--` for label flags so that we can
	// apply them to our runner profile as target labels later
	targetLabels := make(map[string]string)
	for i, arg := range secondaryArgs {
		// A label flag can be either `-label=key=value` or `-label key=value`
		// so we need to parse for both cases
		if strings.Contains(arg, "-label=") {
			kv := strings.Split(strings.TrimPrefix(arg, "-label="), "=")
			targetLabels[kv[0]] = kv[1]
		} else if strings.Contains(arg, "-label") {
			// If -label is the final argument and there is no KV pair following it, we don't attempt to parse
			if i+1 < len(secondaryArgs) {
				// We get the next argument because if it's space delimited, the KV pair
				// should be the next positional argument
				kvPair := secondaryArgs[i+1]
				// If there's no "=", then we skip the argument because it's not a KV pair
				if !strings.Contains(kvPair, "=") {
					continue
				} else {
					kv := strings.Split(strings.TrimPrefix(secondaryArgs[i+1], "-label="), "=")
					targetLabels[kv[0]] = kv[1]
				}
			}
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
		c.ui.Output(runnerInstallFailed, runnerPlatform, id, terminal.WithWarningStyle())
		return 1
	}
	s.Update("Runner %q installed successfully to %s", id, runnerPlatform)
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
		s = sg.Add("Creating runner profile and targeting runner %s", id)
		if odc, ok := p.(installutil.OnDemandRunnerConfigProvider); ok {
			odrConfig = odc.OnDemandRunnerConfig()
			odrConfig.Name = odrConfig.Name + "-" + id
			odrConfig.OciUrl = c.runnerProfileOdrImage // Use what we got from flags (or the default)
		} else {
			odrConfig = &pb.OnDemandRunnerConfig{
				Name:       runnerPlatform + "-" + id,
				OciUrl:     c.runnerProfileOdrImage,
				PluginType: runnerPlatform,
			}
		}
		if len(targetLabels) != 0 {
			odrConfig.TargetRunner = &pb.Ref_Runner{
				Target: &pb.Ref_Runner_Labels{
					Labels: &pb.Ref_RunnerLabels{
						Labels: targetLabels,
					},
				},
			}
		} else {
			odrConfig.TargetRunner = &pb.Ref_Runner{
				Target: &pb.Ref_Runner_Id{
					Id: &pb.Ref_RunnerId{Id: id},
				},
			}
		}

		// if we have no runner profiles, make this one the default
		profiles, err := client.ListOnDemandRunnerConfigs(ctx, &empty.Empty{})
		if err != nil {
			c.ui.Output("Error getting runner profiles: %s", clierrors.Humanize(err))
			return 1
		}
		if len(profiles.Configs) == 0 {
			odrConfig.Default = true
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

	runnerInstallFailed = strings.TrimSpace(`
Please run the following to clean up the resources from the unsuccessful runner installation,
specifying additional platform flags as needed:

waypoint runner uninstall -platform=%[1]s -id=%[2]s <additional_platform_flags>
`)
)
