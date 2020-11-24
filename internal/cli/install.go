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
	"github.com/hashicorp/waypoint/internal/serverinstall"
	"github.com/hashicorp/waypoint/internal/serverinstall/config"
)

type InstallCommand struct {
	*baseCommand

	Config         config.BaseConfig
	platform       string
	contextName    string
	contextDefault bool

	flagAcceptTOS bool
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

	p, err := serverinstall.NewServerPlatformInstaller(&c.Config, c.platform)
	if err != nil {
		c.ui.Output(
			"Error during server install: ", err,
			terminal.WithErrorStyle(),
		)
		return 1
	}

	contextConfig, advertiseAddr, httpAddr, err = p.Install(ctx, c.ui, log)
	if err != nil {
		c.ui.Output(
			"Error installing server into %s: %s", c.platform, clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)

		return 1
	}

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

	s.Done()

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
		f.StringVar(&flag.StringVar{
			Name:    "server-image",
			Target:  &c.Config.ServerImage,
			Usage:   "Docker image for the server image.",
			Default: "hashicorp/waypoint:latest",
		})

		f.StringMapVar(&flag.StringMapVar{
			Name:   "annotate-service",
			Target: &c.Config.ServiceAnnotations,
			Usage:  "Annotations for the Service generated.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "pull-policy",
			Target:  &c.Config.ImagePullPolicy,
			Usage:   "",
			Default: "Always",
		})

		f.BoolVar(&flag.BoolVar{
			Name:   "advertise-internal",
			Target: &c.Config.AdvertiseInternal,
			Usage: "Advertise the internal service address rather than the external. " +
				"This is useful if all your deployments will be able to access the private " +
				"service address. This will default to false but will be automatically set to " +
				"true if the external host is detected to be localhost.",
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
			Default: "kubernetes",
			Usage:   "Platform to install the Waypoint server into.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "accept-tos",
			Target:  &c.flagAcceptTOS,
			Usage:   acceptTOSHelp,
			Default: false,
		})

		f.StringVar(&flag.StringVar{
			Name:    "namespace",
			Target:  &c.Config.Namespace,
			Usage:   "Namespace to install the Waypoint server into for Nomad or Kubernetes.",
			Default: "",
		})

		f.StringVar(&flag.StringVar{
			Name:    "k8s-server-name",
			Target:  &c.Config.ServerName,
			Usage:   "Name of the Waypoint server for Kubernetes.",
			Default: "waypoint-server",
		})

		f.StringVar(&flag.StringVar{
			Name:    "k8s-service",
			Target:  &c.Config.ServiceName,
			Usage:   "Name of the Kubernetes service for the server.",
			Default: "waypoint",
		})

		f.StringVar(&flag.StringVar{
			Name:    "k8s-cpu-request",
			Target:  &c.Config.CPURequest,
			Usage:   "Configures the requested CPU amount for the Waypoint server in Kubernetes",
			Default: "100m",
		})

		f.StringVar(&flag.StringVar{
			Name:    "k8s-mem-request",
			Target:  &c.Config.MemRequest,
			Usage:   "Configures the requested memory amount for the Waypoint server in Kubernetes",
			Default: "256Mi",
		})

		f.StringVar(&flag.StringVar{
			Name:    "k8s-storage-request",
			Target:  &c.Config.StorageRequest,
			Usage:   "Configures the requested persistent volume size for the Waypoint server in Kubernetes.",
			Default: "1Gi",
		})

		f.BoolVar(&flag.BoolVar{
			Name:   "openshift",
			Target: &c.Config.OpenShift,
			Default: false,
			Usage:  "Enables installing the Waypoint server on Kubernetes on Red Hat OpenShift.",
		})

		f.StringVar(&flag.StringVar{
			Name:   "secret-file",
			Target: &c.Config.SecretFile,
			Usage:  "Use the Kubernetes Secret in the given path to access the Waypoint server image.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "pull-secret",
			Target:  &c.Config.ImagePullSecret,
			Usage:   "Secret to use to access the Waypoint server image on Kubernetes.",
			Default: "github",
		})

		f.StringVar(&flag.StringVar{
			Name:    "nomad-region",
			Target:  &c.Config.RegionF,
			Default: "global",
			Usage:   "Region to install the Waypoint server to on Nomad.",
		})
	
		f.StringSliceVar(&flag.StringSliceVar{
			Name:    "nomad-dc",
			Target:  &c.Config.DatacentersF,
			Default: []string{"dc1"},
			Usage:   "Datacenters to install to on Nomad platform.",
		})
	
		f.BoolVar(&flag.BoolVar{
			Name:    "nomad-policy-override",
			Target:  &c.Config.PolicyOverrideF,
			Default: false,
			Usage:   "Override the Nomad sentinel policy on enterprise Nomad platform.",
		})
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
