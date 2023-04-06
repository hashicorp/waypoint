// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/cap/util"
	"github.com/posener/complete"
	empty "google.golang.org/protobuf/types/known/emptypb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clicontext"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/pkg/k8sauth"
	wpoidc "github.com/hashicorp/waypoint/pkg/auth/oidc"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverclient"
	"github.com/hashicorp/waypoint/pkg/serverconfig"
)

type LoginCommand struct {
	*baseCommand

	flagAuthMethod     string
	flagToken          string
	flagK8S            bool
	flagK8SService     string
	flagK8STokenSecret string
	flagK8SNamespace   string
}

func (c *LoginCommand) Run(args []string) int {
	ctx := c.Ctx
	log := c.Log

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
		WithNoLocalServer(), // no need to login for local mode
		WithConnectionArg(),

		// Don't initialize the client automatically because if we have
		// -from-kubernetes set we may want to do more logic to detect the
		// server URL.
		WithNoClient(),
	); err != nil {
		// This error specifically comes if we attempt to run this without
		// a server address configured.
		if errors.Is(err, serverclient.ErrNoServerConfig) {
			c.ui.Output(strings.TrimSpace(errLoginServerAddress), terminal.WithErrorStyle())
		}

		return 1
	}

	// Get our default context. If the server address matches, then
	// we will simply overwrite that. We grab this early so that any errors
	// happen before we do the login loop.
	var contextDefault *clicontext.Config
	contextDefaultName, err := c.contextStorage.Default()
	if err == nil && contextDefaultName != "" {
		contextDefault, err = c.contextStorage.Load(contextDefaultName)
	}
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	// If we are using K8S and we don't have an address set, then
	// look it up using the service.
	if c.flagK8S && c.flagConnection.Server.Address == "" {
		log.Debug("-from-k8s with no address, detecting from Kubernetes service",
			"service", c.flagK8SService)
		addr, err := c.k8sServerAddr(ctx)
		if err != nil {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}

		c.flagConnection.Server.Address = addr

		// TLS is always skip verify for the current version of the helm
		// chart. We will add a way to detect TLS settings later.
		c.flagConnection.Server.Tls = true
		c.flagConnection.Server.TlsSkipVerify = true

		log.Debug("-from-kubernetes detected connection info",
			"address", addr,
			"tls", c.flagConnection.Server.Tls,
			"tls_skip_verify", c.flagConnection.Server.TlsSkipVerify,
		)
	}

	// Manually initialize our client. We have to do this because we have
	// WithClient(false) set above in case we populate the server addr with
	// Kubernetes info.
	c.project, err = c.initClient(ctx)
	if err != nil {
		c.ui.Output(
			"Error reconnecting with token: %s",
			clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	// Determine our auth func, which by default is OIDC
	var authFunc func(context.Context) (string, int)
	switch {
	case c.flagK8S:
		log.Info("login method", "method", "kubernetes")
		authFunc = c.loginK8S

	case c.flagToken != "":
		log.Info("login method", "method", "token")
		authFunc = c.loginToken

	default:
		log.Info("login method", "method", "OIDC")
		authFunc = c.loginOIDC
	}

	// Log in
	token, exitCode := authFunc(ctx)
	if exitCode > 0 {
		return exitCode
	}

	// Save our context and set it as the default. We copy flagConnection
	// which is already configured with the basic server connection stuff
	// from this command.
	newContext := c.flagConnection
	if c.clientContext != nil {
		// clientContext is always set to our actual context we used to
		// create our client. So this will accurately grab non-flag based
		// access i.e. loading our default context.
		newContext = *c.clientContext
	}
	newContext.Server.AuthToken = token
	newContext.Server.RequireAuth = true
	if c.flagK8S {
		newContext.Server.Platform = "kubernetes"
	}
	log.Debug("final login connection info",
		"address", newContext.Server.Address,
		"tls", newContext.Server.Tls,
		"tls_skip_verify", newContext.Server.TlsSkipVerify,
		"token", newContext.Server.AuthToken,
		"require_auth", newContext.Server.RequireAuth,
	)

	// Validate the connection
	_, err = c.initClient(ctx, serverclient.FromContextConfig(&newContext))
	if err != nil {
		c.ui.Output(fmt.Sprintf(
			strings.TrimSpace(errTokenValidation),
			clierrors.Humanize(err),
		), terminal.WithErrorStyle())
		return 1
	}

	// If the default context matches the server address, then we overwrite
	// that one. This prevents constant context sprawl as we reauth.
	contextName := fmt.Sprintf("login_%d", time.Now().Unix())
	if contextDefault != nil && newContext.Server.Address == contextDefault.Server.Address {
		contextName = contextDefaultName
	}

	// Set our contexts
	if err := c.contextStorage.Set(contextName, &newContext); err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}
	if err := c.contextStorage.SetDefault(contextName); err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	// Done!
	c.ui.Output("Authentication complete! You're now logged in.", terminal.WithSuccessStyle())

	return 0
}

func (c *LoginCommand) loginToken(ctx context.Context) (string, int) {
	// First we decode the token to ensure it is valid and also to figure
	// out if we have a login or invite token.
	decodeResp, err := c.project.Client().DecodeToken(ctx, &pb.DecodeTokenRequest{
		Token: c.flagToken,
	})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return "", 1
	}
	token := decodeResp.Token

	// If we have a login token, then we're just done cause that can be stored directly.
	if _, ok := token.Kind.(*pb.Token_Login_); ok {
		return c.flagToken, 0
	}

	// Then it must be an invite token
	if _, ok := token.Kind.(*pb.Token_Invite_); !ok {
		c.ui.Output(strings.TrimSpace(errTokenInvalid), terminal.WithErrorStyle())
		return "", 1
	}

	// Convert it
	convertResp, err := c.project.Client().ConvertInviteToken(ctx, &pb.ConvertInviteTokenRequest{
		Token: c.flagToken,
	})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return "", 1
	}

	return convertResp.Token, 0
}

func (c *LoginCommand) loginOIDC(ctx context.Context) (string, int) {
	// Get our OIDC auth methods
	respList, err := c.project.Client().ListOIDCAuthMethods(ctx, &empty.Empty{})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return "", 1
	}
	if len(respList.AuthMethods) == 0 {
		c.ui.Output(strings.TrimSpace(errNoAuthMethods), terminal.WithErrorStyle())
		return "", 1
	}

	// If no auth method is specified, we use the only one. If more than
	// one is configured, it is an error.
	if c.flagAuthMethod == "" {
		if len(respList.AuthMethods) > 1 {
			var names []string
			for _, m := range respList.AuthMethods {
				names = append(names, m.Name)
			}
			sort.Strings(names)

			c.ui.Output(
				strings.TrimSpace(errManyAuthMethods),
				strings.Join(names, "\n"),
				terminal.WithErrorStyle(),
			)
			return "", 1
		}

		c.flagAuthMethod = respList.AuthMethods[0].Name
	}
	refAM := &pb.Ref_AuthMethod{Name: c.flagAuthMethod}

	// Start our callback server
	callbackSrv, err := wpoidc.NewCallbackServer()
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return "", 1
	}

	// Get the auth URL
	respURL, err := c.project.Client().GetOIDCAuthURL(ctx, &pb.GetOIDCAuthURLRequest{
		AuthMethod:  refAM,
		RedirectUri: callbackSrv.RedirectUri(),
		Nonce:       callbackSrv.Nonce(),
	})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return "", 1
	}

	// Open the auth URL in the user browser or ask them to visit it.
	// We purposely use fmt here and NOT c.ui because the ui will truncate
	// our URL (a known bug).
	fmt.Printf(strings.TrimSpace(outVisitURL)+"\n\n", respURL.Url)
	if err := util.OpenURL(respURL.Url); err != nil {
		c.Log.Warn("error opening auth url", "err", err)
	}

	// Wait
	var req *pb.CompleteOIDCAuthRequest
	select {
	case <-ctx.Done():
		// User cancelled
		return "", 1

	case err := <-callbackSrv.ErrorCh():
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return "", 1

	case req = <-callbackSrv.SuccessCh():
		// We got our data!
	}

	// Complete the auth
	req.AuthMethod = refAM
	respToken, err := c.project.Client().CompleteOIDCAuth(ctx, req)
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return "", 1
	}

	// Output the claims in the debug log
	c.Log.Warn("OIDC authentication complete",
		"user_id", respToken.User.Id,
		"username", respToken.User.Username,
		"id_claims", respToken.IdClaimsJson,
		"user_claims", respToken.UserClaimsJson,
	)

	return respToken.Token, 0
}

func (c *LoginCommand) loginK8S(ctx context.Context) (string, int) {
	// Get our Kubernetes client
	clientset, ns, _, err := k8sauth.Clientset("", "")
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return "", 1
	}

	if c.flagK8SNamespace != "" {
		ns = c.flagK8SNamespace
	}

	secretClient := clientset.CoreV1().Secrets(ns)

	// Get the secret
	secret, err := secretClient.Get(ctx, c.flagK8STokenSecret, metav1.GetOptions{})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return "", 1
	}

	// Get our token
	tokenB64 := secret.Data["token"]
	if len(tokenB64) == 0 {
		c.ui.Output(strings.TrimSpace(errK8STokenEmpty), terminal.WithErrorStyle())
		return "", 1
	}

	return string(tokenB64), 0
}

func (c *LoginCommand) k8sServerAddr(ctx context.Context) (string, error) {
	// Get our Kubernetes client
	clientset, ns, _, err := k8sauth.Clientset("", "")
	if err != nil {
		return "", err
	}

	if c.flagK8SNamespace != "" {
		ns = c.flagK8SNamespace
	}

	serviceClient := clientset.CoreV1().Services(ns)

	// Get the service
	var advertiseAddr string
	waitOut := false
	err = wait.PollImmediate(5*time.Second, 15*time.Minute, func() (bool, error) {
		service, err := serviceClient.Get(ctx, c.flagK8SService, metav1.GetOptions{})
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

		// Only show this once.
		if !waitOut {
			c.ui.Output("Waiting for the Waypoint service to become ready...")
			waitOut = true
		}

		return false, nil
	})
	if err != nil {
		return "", err
	}

	if advertiseAddr == "" {
		return "", fmt.Errorf("Failed to detect waypoint-ui service address.")
	}

	// The advertise addr always needs the gRPC port. For our Helm chart
	// this isn't configurable so this is always correct.
	advertiseAddr += ":" + serverconfig.DefaultGRPCPort

	c.ui.Output("Waypoint server URL: %s", advertiseAddr)
	return advertiseAddr, nil
}

func (c *LoginCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetConnection, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.StringVar(&flag.StringVar{
			Name:   "auth-method",
			Target: &c.flagAuthMethod,
			Usage: "Auth method to use for login. This will default to " +
				"the only available auth method if only one exists.",
		})

		f.StringVar(&flag.StringVar{
			Name:   "token",
			Target: &c.flagToken,
			Usage: "Auth with a token. This will force auth-method to 'token'. " +
				"This works with both login and invite tokens.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:   "from-kubernetes",
			Target: &c.flagK8S,
			Usage: "Perform the initial authentication after Waypoint is installed " +
				"to Kubernetes. This requires kubectl to be configured with access to the " +
				"Kubernetes cluster. The primary use case of this is to get the first " +
				"token from a Waypoint installation. After that, future users should use " +
				"a configured auth method or request a token from an administrator.",
		})

		f.StringVar(&flag.StringVar{
			Name:   "from-kubernetes-service",
			Target: &c.flagK8SService,
			Usage: "The name of the Kubernetes service to get the server address from " +
				"when using the -from-kubernetes flag.",
			Default: "waypoint-ui",
		})

		f.StringVar(&flag.StringVar{
			Name:   "from-kubernetes-secret",
			Target: &c.flagK8STokenSecret,
			Usage: "The name of the Kubernetes secret that has the Waypoint token " +
				"when using the -from-kubernetes flag.",
			Default: "waypoint-server-token",
		})

		f.StringVar(&flag.StringVar{
			Name:   "from-kubernetes-namespace",
			Target: &c.flagK8SNamespace,
			Usage: "The name of the Kubernetes namespace that has the Waypoint token " +
				"when using the -from-kubernetes flag.",
		})
	})
}

func (c *LoginCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *LoginCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *LoginCommand) Synopsis() string {
	return "Log in to a Waypoint server"
}

func (c *LoginCommand) Help() string {
	return formatHelp(`
Usage: waypoint login [server address]

  Log in to a Waypoint server.

  This is usually the first command a new user runs to gain CLI access to
  an existing Waypoint server.

  If the server address is not specified and you have an active
  context (see "waypoint context"), then this command will reauthenticate
  to the currently active server.

  This command can be used for token-based authentication as well as
  other forms such as OIDC. You can use "-token" to specify a login or
  invite token and configure the CLI to access the server.

  The "-from-kubernetes" flag can be used after a fresh Waypoint installation
  on Kubernetes to log in using the bootstrap token. This requires local
  Kubernetes connection configuration via a KUBECONFIG environment variable
  or file.

` + c.Flags().Help())
}

const (
	errNoAuthMethods = `
Only token-based authentication is allowed by this server. To login using
a token, use the "waypoint login" command with the "-token" flag.
`

	errManyAuthMethods = `
The Waypoint server has multiple auth methods configured. You must specify
which auth method you want to use using the "-auth-method" flag. The list
of available auth methods are:

%s
`

	errLoginServerAddress = `
This error usually is because you forgot to specify an address for
a server as an argument. Please use "waypoint login <address>" where
"<address>" is the address to your Waypoint server. For example:

waypoint login example.com

or

waypoint login https://example.com
`

	errTokenInvalid = `
The specified token is not a valid login or invite token. Please
double-check the token and try again.
`

	errK8STokenEmpty = `
The Waypoint token in the Kubernetes secret is empty. This is usually
for one of two reasons. First, the Waypoint server may not be bootstrapped.
After installing Waypoint on Kubernetes, it takes a few minutes for Waypoint
to bootstrap itself.

If Waypoint is already bootstrapped, it's possible the server administrator
already deleted the secret. Future users should not use this authentication
method and should instead ask another Waypoint user to generate an invite token
for them.
`

	errTokenValidation = `
Error while validating the login token. This generally shouldn't happen.
Waypoint performs a final verification that the login is valid and this failed.
Please see the error message below:

%s
`

	outVisitURL = `
Complete the authentication by visiting your authentication provider.
Opening your browser window now. If the browser window does not launch,
please visit the URL below:

%s
`
)
