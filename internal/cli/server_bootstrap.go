package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/posener/complete"
	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

type ServerBootstrapCommand struct {
	*baseCommand

	flagContext        string
	flagContextDefault bool
}

func (c *ServerBootstrapCommand) Run(args []string) int {
	ctx := c.Ctx

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
		WithNoLocalServer(),
	); err != nil {
		return 1
	}

	// If we're running a local in-memory server, bootstrapping is not useful.
	if c.project.Local() {
		c.ui.Output(
			errBootstrapLocal,
			terminal.WithErrorStyle(),
		)
		return 1
	}

	client := c.project.Client()
	resp, err := client.BootstrapToken(ctx, &empty.Empty{})
	if err != nil {
		c.ui.Output(
			"Error bootstrapping the server: %s",
			clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	// Output our token
	c.ui.Output(resp.Token)

	// If we aren't storing a context, we're done
	if c.flagContext == "" {
		// NOTE(briancain): I don't think this will ever happen unless the user
		// specifically sets the flag to empty string. Our flag package will set
		// the default here, which will never be emptry string
		return 0
	} else if c.flagContext == "bootstrap-timestamp" {
		c.flagContext = fmt.Sprintf("bootstrap-%d", time.Now().Unix())
	}

	// Get our current context config and set our new token
	config := *c.clientContext
	config.Server.RequireAuth = true
	config.Server.AuthToken = resp.Token

	// Store it
	if err := c.contextStorage.Set(c.flagContext, &config); err != nil {
		c.ui.Output(
			"Error setting the CLI context: %s\n\n%s",
			clierrors.Humanize(err),
			errBootstrapContext,
			terminal.WithErrorStyle(),
		)
		return 1
	}
	if c.flagContextDefault {
		if err := c.contextStorage.SetDefault(c.flagContext); err != nil {
			c.ui.Output(
				"Error setting the CLI context: %s\n\n%s",
				clierrors.Humanize(err),
				errBootstrapContext,
				terminal.WithErrorStyle(),
			)
			return 1
		}
	}

	return 0
}

func (c *ServerBootstrapCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetConnection, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.StringVar(&flag.StringVar{
			Name:   "context-create",
			Target: &c.flagContext,
			Usage: "Create a CLI context for this bootstrapped server. The context name " +
				"will be the value of this flag. If this is an empty string, a context will " +
				"not be created",
			Default: "bootstrap-timestamp",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "context-set-default",
			Target:  &c.flagContextDefault,
			Default: true,
			Usage: "Set the newly bootstrapped server as the default CLI context. This " +
				"only has an effect if -context-create is non-empty.",
		})
	})
}

func (c *ServerBootstrapCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ServerBootstrapCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ServerBootstrapCommand) Synopsis() string {
	return "Bootstrap the server and retrieve the initial auth token"
}

func (c *ServerBootstrapCommand) Help() string {
	return formatHelp(`
Usage: waypoint server bootstrap [options]

  Bootstrap a new server and retrieve the initial auth token.

  When a server is started for the first time against an empty database,
  it is able to be bootstrapped. The bootstrap process retrieves the initial
  auth token for the server. After the auth token is retrieved, it can never
  be bootstrapped again.

  This command is only required for manually run servers. For servers
  installed with "waypoint install", the bootstrap is done automatically
  during the install process.

  The easiest way to run this command against a new server is by using
  flags to specify server connection information. This command will setup
  a CLI context by default.

` + c.Flags().Help())
}

var (
	errBootstrapContext = strings.TrimSpace(`
The Waypoint server successfully bootstrapped, but creating the context failed.

The bootstrap token is available above. The context could not be created
so the CLI is not configured to connect to the server. Please try to manually
recreate the context.
`)

	errBootstrapLocal = strings.TrimSpace(`
No running server detected.

Bootstrapping is only required for running servers. This error may happen
if you didn't specify a "-server-addr" or the server at that address has shut
down. Please start a server with "waypoint server run" and try again.
`)
)
