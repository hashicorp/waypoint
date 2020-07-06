package cli

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/posener/complete"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/hashicorp/waypoint/internal/clicontext"
	configpkg "github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/serverclient"
	"github.com/hashicorp/waypoint/internal/serverinstall"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

type InstallCommand struct {
	*baseCommand

	config         serverinstall.Config
	showYaml       bool
	contextName    string
	contextDefault bool
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
	); err != nil {
		return 1
	}

	// Decode our configuration
	output, err := serverinstall.Render(&c.config)
	if err != nil {
		c.ui.Output(
			"Error generating configuration: %s", err.Error(),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	stdout, stderr, err := c.ui.OutputWriters()
	if err != nil {
		panic(err)
	}

	if c.showYaml {
		fmt.Fprint(stdout, output)
		if output[:len(output)-1] != "\n" {
			fmt.Fprint(stdout, "\n")
		}

		return 0
	}

	cmd := exec.Command("kubectl", "create", "-f", "-")
	cmd.Stdin = strings.NewReader(output)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err = cmd.Run()
	if err != nil {
		c.ui.Output(
			"Error executing kubectl: %s", err.Error(),
			terminal.WithErrorStyle(),
		)

		return 1
	}

	st := c.ui.Status()
	defer st.Close()
	st.Update("Waiting for Kubernetes service to be ready...")

	// Build our K8S client.
	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
	clientconfig, err := config.ClientConfig()
	if err != nil {
		c.ui.Output(
			"Error initializing kubernetes client: %s", err.Error(),
			terminal.WithErrorStyle(),
		)
		return 1
	}
	clientset, err := kubernetes.NewForConfig(clientconfig)
	if err != nil {
		c.ui.Output(
			"Error initializing kubernetes client: %s", err.Error(),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	// Wait for our service to be ready
	log.Info("waiting for server service to become ready")
	var contextConfig clicontext.Config
	var advertiseAddr pb.ServerConfig_AdvertiseAddr
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

		// Get the port
		var port int32
		for _, spec := range svc.Spec.Ports {
			if spec.Name == "grpc" {
				port = spec.Port
				break
			}
		}
		if port == 0 {
			// If we didn't find the port, retry...
			log.Trace("no port found on service, retrying")
			return false, nil
		}

		// Set the address
		addr = fmt.Sprintf("%s:%d", addr, port)
		log.Info("server service ready", "addr", addr)

		// Set our advertise address
		advertiseAddr.Addr = addr
		advertiseAddr.Insecure = true

		// Set our connection information
		contextConfig = clicontext.Config{
			Server: configpkg.Server{
				Address:  addr,
				Insecure: true, // always for now
			},
		}

		return true, nil
	})
	if err != nil {
		c.ui.Output(
			"Error waiting for service ready: %s\n\n%s",
			err.Error(),
			errInstallRunning,
			terminal.WithErrorStyle(),
		)
		return 1
	}

	// Connect
	st.Update(fmt.Sprintf("Service ready. Connecting to: %s", contextConfig.Server.Address))
	log.Info("connecting to the server so we can set the server config", "addr", contextConfig.Server.Address)
	conn, err := serverclient.Connect(ctx,
		serverclient.FromContextConfig(&contextConfig),
		serverclient.Timeout(1*time.Minute),
	)
	if err != nil {
		c.ui.Output(
			"Error connecting to server: %s\n\n%s",
			err.Error(),
			errInstallRunning,
			terminal.WithErrorStyle(),
		)
		return 1
	}
	client := pb.NewWaypointClient(conn)

	// If we connected successfully, lets immediately setup our context.
	if c.contextName != "" {
		if err := c.contextStorage.Set(c.contextName, &contextConfig); err != nil {
			c.ui.Output(
				"Error setting the CLI context: %s\n\n%s",
				err.Error(),
				errInstallRunning,
				terminal.WithErrorStyle(),
			)
			return 1
		}
		if c.contextDefault {
			if err := c.contextStorage.SetDefault(c.contextName); err != nil {
				c.ui.Output(
					"Error setting the CLI context: %s\n\n%s",
					err.Error(),
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
				&advertiseAddr,
			},
		},
	})
	if err != nil {
		c.ui.Output(
			"Error setting the advertise address: %s\n\n%s",
			err.Error(),
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

		f.BoolVar(&flag.BoolVar{
			Name:   "show-yaml",
			Target: &c.showYaml,
			Usage:  "Show the YAML to be send to the cluster.",
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
	})
}

func (c *InstallCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *InstallCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *InstallCommand) Synopsis() string {
	return "Output Kubernetes configurations to run a self-hosted server."
}

func (c *InstallCommand) Help() string {
	helpText := `
Usage: waypoint install [options]

  Installs a Waypoint server to an existing Kubernetes cluster.

  By default, this will also automatically create a new default CLI context
  (see "waypoint context") so the CLI will be configured to use the newly
  installed server.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
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
`)
)
