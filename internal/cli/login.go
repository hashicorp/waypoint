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
	"github.com/hashicorp/waypoint/internal/clicontext"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/serverclient"
)

type LoginCommand struct {
	*baseCommand

	flagAuthMethod string
	flagToken      string
}

func (c *LoginCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
		WithNoAutoServer(), // no need to login for local mode
		WithConnectionArg(),
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

	// Determine our auth func, which by default is OIDC
	var authFunc func() (string, int)
	switch {
	case c.flagToken != "":
		authFunc = c.loginToken

	default:
		authFunc = c.loginOIDC
	}

	// Log in
	token, exitCode := authFunc()
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

func (c *LoginCommand) loginToken() (string, int) {
	// First we decode the token to ensure it is valid and also to figure
	// out if we have a login or invite token.
	decodeResp, err := c.project.Client().DecodeToken(c.Ctx, &pb.DecodeTokenRequest{
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
	convertResp, err := c.project.Client().ConvertInviteToken(c.Ctx, &pb.ConvertInviteTokenRequest{
		Token: c.flagToken,
	})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return "", 1
	}

	return convertResp.Token, 0
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
	// We purposely use fmt here and NOT c.ui because the ui will truncate
	// our URL (a known bug).
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

	// Output the claims in the debug log
	c.Log.Warn("OIDC authentication complete",
		"user_id", respToken.User.Id,
		"username", respToken.User.Username,
		"id_claims", respToken.IdClaimsJson,
		"user_claims", respToken.UserClaimsJson,
	)

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

		f.StringVar(&flag.StringVar{
			Name:   "token",
			Target: &c.flagToken,
			Usage: "Auth with a token. This will force auth-method to 'token'. " +
				"This works with both login and invite tokens.",
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

	outVisitURL = `
Complete the authentication by visiting your authentication provider.
Opening your browser window now. If the browser window does not launch,
please visit the URL below:

%s
`
)
