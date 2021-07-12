package cli

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/cap/util"
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	wpoidc "github.com/hashicorp/waypoint/internal/auth/oidc"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/serverclient"
)

type LoginCommand struct {
	*baseCommand

	flagAuthMethod string
}

func (c *LoginCommand) Run(args []string) int {
	// TODO server addr arg

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
		WithNoAutoServer(), // no need to login for local mode
	); err != nil {
		// This error specifically comes if we attempt to run this without
		// a server address configured.
		if errors.Is(err, serverclient.ErrNoServerConfig) {
			c.ui.Output(strings.TrimSpace(errLoginServerAddress), terminal.WithErrorStyle())
		}

		return 1
	}

	// Login with OIDC
	token, exitCode := c.loginOIDC()
	if exitCode > 0 {
		return exitCode
	}

	// Save our context and set it as the default. We copy flagConnection
	// which is already configured with the basic server connection stuff
	// from this command.
	newContext := c.flagConnection
	newContext.Server.AuthToken = token
	newContext.Server.RequireAuth = true

	// Set our contexts
	contextName := fmt.Sprintf("login_%d", time.Now().Unix())
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

func (c *LoginCommand) loginOIDC() (string, int) {
	// Get our OIDC auth methods
	respList, err := c.project.Client().ListOIDCAuthMethods(c.Ctx, &empty.Empty{})
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
	respURL, err := c.project.Client().GetOIDCAuthURL(c.Ctx, &pb.GetOIDCAuthURLRequest{
		AuthMethod:  refAM,
		RedirectUri: callbackSrv.RedirectUri(),
		Nonce:       callbackSrv.Nonce(),
	})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return "", 1
	}

	// Open the auth URL in the user browser or ask them to visit it.
	fmt.Printf(strings.TrimSpace(outVisitURL)+"\n\n", respURL.Url)
	if err := util.OpenURL(respURL.Url); err != nil {
		c.Log.Warn("error opening auth url", "err", err)
	}

	// Wait
	var req *pb.CompleteOIDCAuthRequest
	select {
	case <-c.Ctx.Done():
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
	respToken, err := c.project.Client().CompleteOIDCAuth(c.Ctx, req)
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return "", 1
	}

	return respToken.Token, 0
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

  If the server address is not specified and you have an active
  context (see "waypoint context"), then this command will reauthenticate
  to the currently active server.

` + c.Flags().Help())
}

const (
	errNoAuthMethods = `
Only token-based authentication is allowed by this server. To login using
a token, use the "waypoint context create" command.
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

	outVisitURL = `
Complete the authentication by visiting your authentication provider.
Opening your browser window now. If the browser window does not launch,
please visit the URL below:

%s
`
)
