package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/posener/complete"
	"google.golang.org/grpc"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clicontext"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/serverclient"
	"github.com/hashicorp/waypoint/internal/serverconfig"
	"github.com/hashicorp/waypoint/internal/serverinstall"
)

type InstallCommand struct {
	*baseCommand

	platform       string
	contextName    string
	contextDefault bool

	flagAcceptTOS bool
	flagRunner    bool
}

func (c *InstallCommand) Run(args []string) int {
	ctx := c.Ctx
	log := c.Log.Named("install")
	defer c.Close()

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
		WithClient(false),
	); err != nil {
		return 1
	}

	if !c.flagAcceptTOS {
		c.ui.Output(strings.TrimSpace(tosStatement), terminal.WithErrorStyle())
		return 1
	}

	var (
		contextConfig *clicontext.Config
		advertiseAddr *pb.ServerConfig_AdvertiseAddr
	)

	var err error
	var httpAddr string

	p, ok := serverinstall.Platforms[strings.ToLower(c.platform)]
	if !ok {
		c.ui.Output(
			"Error installing server into %s: invalid platform",
			c.platform,
			terminal.WithErrorStyle(),
		)

		return 1
	}

	result, err := p.Install(ctx, &serverinstall.InstallOpts{
		Log: log,
		UI:  c.ui,
	})
	if err != nil {
		c.ui.Output(
			"Error installing server into %s: %s", c.platform, clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)

		return 1
	}

	contextConfig = result.Context
	advertiseAddr = result.AdvertiseAddr
	httpAddr = result.HTTPAddr

	sg := c.ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Connecting to: %s", contextConfig.Server.Address)
	defer func() { s.Abort() }()

	// Connect
	log.Info("connecting to the server so we can set the server config", "addr", contextConfig.Server.Address)
	conn, err := serverclient.Connect(ctx,
		serverclient.FromContextConfig(contextConfig),
		serverclient.Timeout(5*time.Minute),
	)
	if err != nil {
		c.ui.Output(
			"Error connecting to server: %s\n\n%s",
			clierrors.Humanize(err),
			errInstallRunning,
			terminal.WithErrorStyle(),
		)
		return 1
	}
	client := pb.NewWaypointClient(conn)

	s.Update("Retrieving initial auth token...")

	// We need our bootstrap token immediately
	var callOpts []grpc.CallOption
	tokenResp, err := client.BootstrapToken(ctx, &empty.Empty{})
	if err != nil {
		c.ui.Output(
			"Error getting the initial token: %s\n\n%s",
			clierrors.Humanize(err),
			errInstallRunning,
			terminal.WithErrorStyle(),
		)
		return 1
	}
	if tokenResp != nil {
		log.Debug("token received, setting on context")
		contextConfig.Server.RequireAuth = true
		contextConfig.Server.AuthToken = tokenResp.Token

		callOpts = append(callOpts, grpc.PerRPCCredentials(
			serverclient.StaticToken(tokenResp.Token)))
	}

	// If we connected successfully, lets immediately setup our context.
	if c.contextName != "" {
		if err := c.contextStorage.Set(c.contextName, contextConfig); err != nil {
			c.ui.Output(
				"Error setting the CLI context: %s\n\n%s",
				clierrors.Humanize(err),
				errInstallRunning,
				terminal.WithErrorStyle(),
			)
			return 1
		}
		if c.contextDefault {
			if err := c.contextStorage.SetDefault(c.contextName); err != nil {
				c.ui.Output(
					"Error setting the CLI context: %s\n\n%s",
					clierrors.Humanize(err),
					errInstallRunning,
					terminal.WithErrorStyle(),
				)
				return 1
			}
		}
	}

	// Reconnect with the token set. The `contextConfig` has the token set on
	// it now so we can just reconnect with the same context.
	log.Info("reconnecting with our bootstrap token", "addr", contextConfig.Server.Address)
	conn.Close()
	conn, err = serverclient.Connect(ctx,
		serverclient.FromContextConfig(contextConfig),
		serverclient.Timeout(5*time.Minute),
	)
	if err != nil {
		c.ui.Output(
			"Error connecting to server with bootstrap token: %s\n\n%s",
			clierrors.Humanize(err),
			errInstallRunning,
			terminal.WithErrorStyle(),
		)
		return 1
	}
	client = pb.NewWaypointClient(conn)

	// Set the config
	s.Update("Configuring server...")
	log.Debug("setting the advertise address", "addr", fmt.Sprintf("%#v", advertiseAddr))
	_, err = client.SetServerConfig(ctx, &pb.SetServerConfigRequest{
		Config: &pb.ServerConfig{
			AdvertiseAddrs: []*pb.ServerConfig_AdvertiseAddr{
				advertiseAddr,
			},
		},
	}, callOpts...)
	if err != nil {
		c.ui.Output(
			"Error setting the advertise address: %s\n\n%s",
			clierrors.Humanize(err),
			errInstallRunning,
			terminal.WithErrorStyle(),
		)
		return 1
	}

	s.Update("Server installed and configured!")
	s.Done()

	if c.flagRunner {
		s = sg.Add("")

		// We need a new auth token for the runner so that the runner
		// can connect to the server. We don't want to reuse the bootstrap
		// token that is shared with the CLI cause that can be revoked.
		s.Update("Retrieving new auth token for runner...")
		resp, err := client.GenerateLoginToken(c.Ctx, &empty.Empty{})
		if err != nil {
			c.ui.Output(
				"Error retrieving auth token for runner: %s\n\n%s",
				clierrors.Humanize(err),
				errInstallRunner,
				terminal.WithErrorStyle(),
			)
			return 1
		}

		// Build a serverconfig that uses the advertise addr and includes
		// the token we just requested.
		connConfig := &serverconfig.Client{
			Address:       advertiseAddr.Addr,
			Tls:           advertiseAddr.Tls,
			TlsSkipVerify: advertiseAddr.TlsSkipVerify,
			RequireAuth:   true,
			AuthToken:     resp.Token,
		}

		// Install!
		s.Update("Installing runner...")
		err = p.InstallRunner(ctx, &serverinstall.InstallRunnerOpts{
			Log:             log,
			UI:              c.ui,
			AuthToken:       resp.Token,
			AdvertiseAddr:   advertiseAddr,
			AdvertiseClient: connConfig,
		})
		if err != nil {
			c.ui.Output(
				"Error installing the runner: %s\n\n%s",
				clierrors.Humanize(err),
				errInstallRunner,
				terminal.WithErrorStyle(),
			)
			return 1
		}

		s.Done()
	}

	// Close and success
	c.ui.Output(outInstallSuccess,
		c.contextName,
		advertiseAddr.Addr,
		"https://"+httpAddr,
		terminal.WithSuccessStyle(),
	)
	return 0
}

func (c *InstallCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:    "accept-tos",
			Target:  &c.flagAcceptTOS,
			Usage:   acceptTOSHelp,
			Default: false,
		})

		f.StringVar(&flag.StringVar{
			Name:    "context-create",
			Target:  &c.contextName,
			Default: fmt.Sprintf("install-%d", time.Now().Unix()),
			Usage: "Create a context with connection information for this installation. " +
				"The default value will be suffixed with a timestamp at the time the command is executed.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "context-set-default",
			Target:  &c.contextDefault,
			Default: true,
			Usage:   "Set the newly installed server as the default CLI context.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "platform",
			Target:  &c.platform,
			Default: "",
			Usage:   "Platform to install the Waypoint server into.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "x-runner",
			Target:  &c.flagRunner,
			Usage:   "Install a runner in addition to the server",
			Default: false,
			Hidden:  true,
		})

		for name, platform := range serverinstall.Platforms {
			platformSet := set.NewSet(name + " Options")
			platform.InstallFlags(platformSet)
		}
	})
}

func (c *InstallCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *InstallCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *InstallCommand) Synopsis() string {
	return "Install the Waypoint server to Kubernetes, Nomad, or Docker"
}

func (c *InstallCommand) Help() string {
	return formatHelp(`
Usage: waypoint server install [options]
Alias: waypoint install

	Installs a Waypoint server to an existing platform. The platform should be
	specified as kubernetes, nomad, or docker.

  By default, this will also automatically create a new default CLI context
  (see "waypoint context") so the CLI will be configured to use the newly
  installed server.

  This command will require you to accept the Waypoint Terms of Service
  and Privacy Policy for the Waypoint URL service by specifying the "-accept-tos"
  flag. This only applies to the Waypoint URL service. You may disable the
  URL service by manually running the server. If you disable the URL service,
  you do not need to accept any terms.

` + c.Flags().Help())
}

var (
	errInstallRunning = strings.TrimSpace(`
The Waypoint server has been deployed, but due to this error we were
unable to automatically configure the local CLI or the Waypoint server
advertise address. You must do this manually using "waypoint context"
and "waypoint server config-set".
`)

	errInstallRunner = strings.TrimSpace(`
The Waypoint runner failed to install. This error occurred after the
Waypoint server was successfully installed. Your CLI is configured to
use the installed server. If you want to retry, you must uninstall the
server first.
`)

	outInstallSuccess = strings.TrimSpace(`
Waypoint server successfully installed and configured!

The CLI has been configured to connect to the server automatically. This
connection information is saved in the CLI context named %[1]q.
Use the "waypoint context" CLI to manage CLI contexts.

The server has been configured to advertise the following address for
entrypoint communications. This must be a reachable address for all your
deployments. If this is incorrect, manually set it using the CLI command
"waypoint server config-set".

Advertise Address: %[2]s
Web UI Address: %[3]s
`)
)
