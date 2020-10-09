package cli

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/posener/complete"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clicontext"
	"github.com/hashicorp/waypoint/internal/clierrors"
	configpkg "github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/serverclient"
	"github.com/hashicorp/waypoint/internal/serverinstall"
)

type InstallCommand struct {
	*baseCommand

	config            serverinstall.Config
	showYaml          bool
	advertiseInternal bool
	contextName       string
	contextDefault    bool
	platform          string
	secretFile        string

	flagAcceptTOS bool
}

func (c *InstallCommand) InstallKubernetes(
	ctx context.Context, st terminal.Status, log hclog.Logger,
) (*clicontext.Config, *pb.ServerConfig_AdvertiseAddr, string, int) {
	stdout, stderr, err := c.ui.OutputWriters()
	if err != nil {
		panic(err)
	}

	if c.secretFile != "" {
		data, err := ioutil.ReadFile(c.secretFile)
		if err != nil {
			c.ui.Output(
				"Error reading Kubernetes secret file: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return nil, nil, "", 1
		}

		var secretData struct {
			Metadata struct {
				Name string `yaml:"name"`
			} `yaml:"metadata"`
		}

		err = yaml.Unmarshal(data, &secretData)
		if err != nil {
			c.ui.Output(
				"Error reading Kubernetes secret file: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return nil, nil, "", 1
		}

		if secretData.Metadata.Name == "" {
			c.ui.Output(
				"Invalid secret, no metadata.name",
				terminal.WithErrorStyle(),
			)
			return nil, nil, "", 1
		}

		c.config.ImagePullSecret = secretData.Metadata.Name

		c.ui.Output("Installing kubernetes secret...")

		cmd := exec.Command("kubectl", "create", "-f", "-")
		cmd.Stdin = bytes.NewReader(data)
		cmd.Stdout = stdout
		cmd.Stderr = stderr

		err = cmd.Run()
		if err != nil {
			c.ui.Output(
				"Error executing kubectl to install secret: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)

			return nil, nil, "", 1
		}
	}

	// Decode our configuration
	output, err := serverinstall.Render(&c.config)
	if err != nil {
		c.ui.Output(
			"Error generating configuration: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return nil, nil, "", 1
	}

	if c.showYaml {
		fmt.Fprint(stdout, output)
		if output[:len(output)-1] != "\n" {
			fmt.Fprint(stdout, "\n")
		}

		return nil, nil, "", 0
	}

	cmd := exec.Command("kubectl", "create", "-f", "-")
	cmd.Stdin = strings.NewReader(output)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err = cmd.Run()
	if err != nil {
		c.ui.Output(
			"Error executing kubectl: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)

		return nil, nil, "", 1
	}

	st.Update("Waiting for Kubernetes service to be ready...")

	// Build our K8S client.
	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
	clientconfig, err := config.ClientConfig()
	if err != nil {
		c.ui.Output(
			"Error initializing kubernetes client: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return nil, nil, "", 1
	}

	clientset, err := kubernetes.NewForConfig(clientconfig)
	if err != nil {
		c.ui.Output(
			"Error initializing kubernetes client: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return nil, nil, "", 1
	}

	// Wait for our service to be ready
	log.Info("waiting for server service to become ready")
	var contextConfig clicontext.Config
	var advertiseAddr pb.ServerConfig_AdvertiseAddr
	var httpAddr string
	var grpcAddr string

	err = wait.PollImmediate(2*time.Second, 10*time.Minute, func() (bool, error) {
		svc, err := clientset.CoreV1().Services(c.config.Namespace).Get(
			ctx, c.config.ServiceName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		ingress := svc.Status.LoadBalancer.Ingress
		if len(ingress) == 0 {
			log.Trace("ingress list is empty, waiting")
			return false, nil
		}

		addr := ingress[0].IP
		if addr == "" {
			addr = ingress[0].Hostname
		}

		// No address, still not ready
		if addr == "" {
			log.Trace("address is empty, waiting")
			return false, nil
		}

		// Get the ports
		var grpcPort int32
		var httpPort int32
		for _, spec := range svc.Spec.Ports {
			if spec.Name == "grpc" {
				grpcPort = spec.Port
			}

			if spec.Name == "http" {
				httpPort = spec.Port
			}

			if httpPort != 0 && grpcPort != 0 {
				break
			}
		}
		if grpcPort == 0 || httpPort == 0 {
			// If we didn't find the port, retry...
			log.Trace("no port found on service, retrying")
			return false, nil
		}

		// Set the grpc address
		grpcAddr = fmt.Sprintf("%s:%d", addr, grpcPort)
		log.Info("server service ready", "addr", addr)

		// HTTP address to return
		httpAddr = fmt.Sprintf("%s:%d", addr, httpPort)
		log.Info("http server ready", "httpAddr", addr)

		// Set our advertise address
		advertiseAddr.Addr = grpcAddr
		advertiseAddr.Tls = true
		advertiseAddr.TlsSkipVerify = true

		// If we want internal or we're a localhost address, we use the internal
		// address. The "localhost" check is specifically for Docker for Desktop
		// since pods can't reach this.
		if c.advertiseInternal || strings.HasPrefix(grpcAddr, "localhost:") {
			advertiseAddr.Addr = fmt.Sprintf("%s:%d",
				c.config.ServiceName,
				grpcPort,
			)
		}

		// Set our connection information
		contextConfig = clicontext.Config{
			Server: configpkg.Server{
				Address:       grpcAddr,
				Tls:           true,
				TlsSkipVerify: true, // always for now
			},
		}

		return true, nil
	})
	if err != nil {
		c.ui.Output(
			"Error waiting for service ready: %s\n\n%s",
			clierrors.Humanize(err),
			errInstallRunning,
			terminal.WithErrorStyle(),
		)
		return nil, nil, "", 1
	}

	return &contextConfig, &advertiseAddr, httpAddr, 0
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

	st := c.ui.Status()
	defer st.Close()

	var err error
	var httpAddr string

	switch c.platform {
	case "docker":
		contextConfig, advertiseAddr, httpAddr, err = serverinstall.InstallDocker(ctx, c.ui, st, &c.config)
		if err != nil {
			c.ui.Output(
				"Error installing server into docker: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)

			return 1
		}
	case "kubernetes":
		var code int
		contextConfig, advertiseAddr, httpAddr, code = c.InstallKubernetes(ctx, st, log)
		if code != 0 || c.showYaml {
			return code
		}

		// ok, inline below.
	case "nomad":
		contextConfig, advertiseAddr, httpAddr, err = serverinstall.InstallNomad(ctx, c.ui, st, &c.config)
		if err != nil {
			c.ui.Output(
				"Error installing server into Nomad: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)

			return 1
		}
	default:
		c.ui.Output(
			"Unknown server platform: %s", c.platform,
			terminal.WithErrorStyle(),
		)

		return 1
	}

	// Connect
	st.Update(fmt.Sprintf("Service ready. Connecting to: %s", contextConfig.Server.Address))
	log.Info("connecting to the server so we can set the server config", "addr", contextConfig.Server.Address)
	conn, err := serverclient.Connect(ctx,
		serverclient.FromContextConfig(contextConfig),
		serverclient.Timeout(2*time.Minute),
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

	// We need our bootstrap token immediately
	var callOpts []grpc.CallOption
	st.Update("Retrieving initial auth token...")
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
	st.Update("Configuring server...")
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

	// Close and success
	st.Close()
	c.ui.Output(outInstallSuccess,
		c.contextName,
		advertiseAddr.Addr,
		httpAddr,
		terminal.WithSuccessStyle(),
	)
	return 0
}

func (c *InstallCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.StringVar(&flag.StringVar{
			Name:    "namespace",
			Target:  &c.config.Namespace,
			Usage:   "Kubernetes namespace install into.",
			Default: "default",
		})

		f.StringVar(&flag.StringVar{
			Name:    "service",
			Target:  &c.config.ServiceName,
			Usage:   "Name of the Kubernetes service for the server.",
			Default: "waypoint",
		})

		f.StringVar(&flag.StringVar{
			Name:    "server-image",
			Target:  &c.config.ServerImage,
			Usage:   "Docker image for the server image.",
			Default: "docker.pkg.github.com/hashicorp/waypoint/alpha:latest",
		})

		f.StringMapVar(&flag.StringMapVar{
			Name:   "annotate-service",
			Target: &c.config.ServiceAnnotations,
			Usage:  "Annotations for the Service generated.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "pull-secret",
			Target:  &c.config.ImagePullSecret,
			Usage:   "Secret to use to access the waypoint server image",
			Default: "github",
		})

		f.StringVar(&flag.StringVar{
			Name:    "pull-policy",
			Target:  &c.config.ImagePullPolicy,
			Usage:   "",
			Default: "Always",
		})

		f.BoolVar(&flag.BoolVar{
			Name:   "show-yaml",
			Target: &c.showYaml,
			Usage:  "Show the YAML to be send to the cluster.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:   "advertise-internal",
			Target: &c.advertiseInternal,
			Usage: "Advertise the internal servivce address rather than the external. " +
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
			Usage:   "Platform to install the server into.",
		})

		f.StringVar(&flag.StringVar{
			Name:   "secret-file",
			Target: &c.secretFile,
			Usage:  "Use the Kubernetes Secret in the given path to access the waypoint server image",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "accept-tos",
			Target:  &c.flagAcceptTOS,
			Usage:   acceptTOSHelp,
			Default: false,
		})

		serverinstall.NomadFlags(f)
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

  Installs a Waypoint server to an existing Kubernetes cluster.

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
HTTP UI Address: %[3]s
`)
)
