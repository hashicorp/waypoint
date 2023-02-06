package cli

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/posener/complete"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	empty "google.golang.org/protobuf/types/known/emptypb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/pkg/k8sauth"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverclient"
	"github.com/hashicorp/waypoint/pkg/serverconfig"
)

type K8SBootstrapCommand struct {
	*baseCommand

	flagRootTokenSecret        string
	flagRunnerTokenSecret      string
	flagODRImage               string
	flagODRServiceAccount      string
	flagODRImagePullPolicy     string
	flagAdvertiseService       string
	flagAdvertiseTLS           bool
	flagAdvertiseTLSSkipVerify bool
}

func (c *K8SBootstrapCommand) Run(args []string) int {
	ctx := c.Ctx
	log := c.Log

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
		WithNoLocalServer(),

		// Don't initialize the client because this is called before
		// the client is ready.
		WithNoClient(),
	); err != nil {
		return 1
	}

	// Get our Kubernetes client
	clientset, ns, _, err := k8sauth.ClientsetInCluster()
	if err != nil {
		c.ui.Output(
			"Error initializing Kubernetes client: %s",
			clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}
	serviceClient := clientset.CoreV1().Services(ns)
	secretClient := clientset.CoreV1().Secrets(ns)

	// Get the service
	var advertiseAddr string
	err = wait.PollImmediate(5*time.Second, 15*time.Minute, func() (bool, error) {
		c.ui.Output("Checking for service readiness every 5 seconds...")
		service, err := serviceClient.Get(ctx, c.flagAdvertiseService, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		if service.Spec.Type == "LoadBalancer" {
			if ig := service.Status.LoadBalancer.Ingress; len(ig) > 0 {
				// Prefer hostname over the IP
				if v := ig[0].Hostname; v != "" {
					advertiseAddr = v
				} else {
					advertiseAddr = ig[0].IP
				}

				return true, nil
			}
		} else {
			if ip := service.Spec.ClusterIP; ip != "" {
				advertiseAddr = ip
				return true, nil
			}
		}

		return false, nil
	})
	if err != nil {
		c.ui.Output(
			"Error waiting for Waypoint service: %s",
			clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	if advertiseAddr == "" {
		c.ui.Output(
			"Failed to detect waypoint-ui service address.",
			terminal.WithErrorStyle(),
		)
		return 1
	}

	// The advertise addr always needs the gRPC port. For our Helm chart
	// this isn't configurable so this is always correct.
	advertiseAddr += ":" + serverconfig.DefaultGRPCPort
	log.Info("service ready", "advertise_addr", advertiseAddr)

	// The service is ready so we should also be ready to connect. We
	// set a slightly longer timeout on the initial connection in case the
	// service came up quicker than Waypoint itself. A more robust check would
	// be checking the pods for readiness.
	log.Info("initializing server connection")
	proj, err := c.initClient(ctx, serverclient.Timeout(120*time.Second))
	if err != nil {
		c.ui.Output(
			"Error reconnecting with token: %s",
			clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	// Waypoint bootstrap
	log.Info("bootstrapping the server")
	client := proj.Client()
	resp, err := client.BootstrapToken(ctx, &empty.Empty{})
	if status.Code(err) == codes.PermissionDenied {
		// This is not an error, since our Helm chart will run this
		// bootstrap job on every upgrade as well and we just want to ignore it.
		c.ui.Output("Waypoint already bootstrapped. Doing nothing.")
		return 0
	}
	if err != nil {
		c.ui.Output(
			"Error bootstrapping the server: %s",
			clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}
	bootstrapTokenB64 := base64.StdEncoding.EncodeToString([]byte(resp.Token))
	log.Info("bootstrapping complete")

	// Set our token. We just do this as an env var because that's the easiest
	// way to get this to trigger. Env will override all other token sources.
	if err := os.Setenv(serverclient.EnvServerToken, resp.Token); err != nil {
		c.ui.Output(
			"Error setting token: %s",
			clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	// Reconnect
	log.Info("reconnecting to the server with the bootstrap token")
	proj, err = c.initClient(ctx)
	if err != nil {
		c.ui.Output(
			"Error reconnecting with token: %s",
			clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}
	client = proj.Client()

	// Set our server configuration
	log.Info("setting server configuration")
	_, err = client.SetServerConfig(ctx, &pb.SetServerConfigRequest{
		Config: &pb.ServerConfig{
			AdvertiseAddrs: []*pb.ServerConfig_AdvertiseAddr{
				{
					Addr:          advertiseAddr,
					Tls:           c.flagAdvertiseTLS,
					TlsSkipVerify: c.flagAdvertiseTLSSkipVerify,
				},
			},
			Platform: "kubernetes",
		},
	})
	if err != nil {
		c.ui.Output(
			"Error setting server configuration: %s",
			clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	// Get our runner token
	log.Info("generating login token for runner")
	tokenResp, err := client.GenerateLoginToken(c.Ctx, &pb.LoginTokenRequest{
		Duration: "", // Never expire for the static runner
	})
	if err != nil {
		c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}
	runnerTokenB64 := base64.StdEncoding.EncodeToString([]byte(tokenResp.Token))

	// Persist our root token
	log.Info("persisting root token", "secret", c.flagRootTokenSecret)
	_, err = secretClient.Patch(ctx, c.flagRootTokenSecret, types.JSONPatchType, []byte(
		fmt.Sprintf(`[{"op":"replace", "path": "/data/token", "value": "%s"}]`, bootstrapTokenB64)),
		metav1.PatchOptions{})
	if err != nil {
		c.ui.Output(
			"Error patching root token secret: %s",
			clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	// Persist our runner token
	log.Info("persisting runner token", "secret", c.flagRunnerTokenSecret)
	_, err = secretClient.Patch(ctx, c.flagRunnerTokenSecret, types.JSONPatchType, []byte(
		fmt.Sprintf(`[{"op":"replace", "path": "/data/token", "value": "%s"}]`, runnerTokenB64)),
		metav1.PatchOptions{})
	if err != nil {
		c.ui.Output(
			"Error patching runner token secret: %s",
			clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	// Configure our on-demand runner
	log.Info("storing on-demand runner configuration for Kubernetes")
	_, err = client.UpsertOnDemandRunnerConfig(ctx, &pb.UpsertOnDemandRunnerConfigRequest{
		Config: &pb.OnDemandRunnerConfig{
			OciUrl:     c.flagODRImage,
			PluginType: "kubernetes",
			PluginConfig: []byte(fmt.Sprintf(`{
	"service_account": "%s",
	"image_pull_policy": "%s"
}`, c.flagODRServiceAccount, c.flagODRImagePullPolicy)),
			ConfigFormat: pb.Hcl_JSON,
			Default:      true,
		},
	})
	if err != nil {
		c.ui.Output(
			"Error storing runner config on server: %s",
			clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	log.Info("bootstrap complete")
	return 0
}

func (c *K8SBootstrapCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetConnection, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.StringVar(&flag.StringVar{
			Name:   "root-token-secret",
			Target: &c.flagRootTokenSecret,
			Usage:  "The name of the Kubernetes secret to write the root token to.",
		})

		f.StringVar(&flag.StringVar{
			Name:   "runner-token-secret",
			Target: &c.flagRunnerTokenSecret,
			Usage:  "The name of the Kubernetes secret to write the runner token to.",
		})

		f.StringVar(&flag.StringVar{
			Name:   "odr-image",
			Target: &c.flagODRImage,
			Usage:  "The name and label of the container image to use for ODR.",
		})

		f.StringVar(&flag.StringVar{
			Name:   "odr-service-account",
			Target: &c.flagODRServiceAccount,
			Usage:  "The name of the Kubernetes service account to use for ODR.",
		})

		f.StringVar(&flag.StringVar{
			Name:   "odr-image-pull-policy",
			Target: &c.flagODRImagePullPolicy,
			Usage:  "The pull policy to use for the container image.",
		})

		f.StringVar(&flag.StringVar{
			Name:   "advertise-service",
			Target: &c.flagAdvertiseService,
			Usage:  "The name of the service to advertise.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "advertise-tls",
			Target:  &c.flagAdvertiseTLS,
			Usage:   "True if the advertise addr supports TLS.",
			Default: true,
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "advertise-tls-skip-verify",
			Target:  &c.flagAdvertiseTLSSkipVerify,
			Usage:   "True if the advertise addr TLS shouldn't be verified.",
			Default: false,
		})
	})
}

func (c *K8SBootstrapCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *K8SBootstrapCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *K8SBootstrapCommand) Synopsis() string {
	return "Post-Helm-install bootstrapping"
}

func (c *K8SBootstrapCommand) Help() string {
	return formatHelp(`
Usage: waypoint k8s bootstrap [options]

  Bootstrap a Waypoint installation from the Waypoint Helm chart.
  This is an internal command and not expected to be manually executed.
  This command only works with in-cluster Kubernetes authentication and
  will not work with out-of-cluster kubectl configuration.

  This command will do a number of things:

  1. Equivalent of "waypoint server bootstrap"
  2. Write a bootstrap token to the given Kubernetes secret
  3. Create a token for a static runner and write it to the configured
     Kubernetes secret.
  4. Configure Kubernetes on-demand runners.

  This command will only run if the server hasn't already been bootstrapped.
  If the server is bootstrapped, this will not run again. This doesn't handle
  partial failures well: if the server bootstrap succeeds but writing the
  secret fails, then the Waypoint installation should be fully uninstalled
  and then reinstalled. This is only use for fresh installations so there
  should be no concern of data loss in the event of a bootstrap failure.

` + c.Flags().Help())
}
