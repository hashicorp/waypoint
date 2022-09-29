package cli

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clicontext"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	configpkg "github.com/hashicorp/waypoint/internal/config"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverclient"
	"github.com/mr-tron/base58"
	"google.golang.org/protobuf/proto"
)

// This file contains the various methods that are used to perform
// the Init call on baseCommand. They are broken down into individual
// smaller methods for readability but more importantly to power the
// "init" subcommand. This allows us to share as much logic as possible
// between Init and "init" to help ensure that "init" succeeding means that
// other commands will succeed as well.

// initConfig initializes the configuration with the specified filename from the CLI.
// If filename is empty, it will default to configpkg.Filename.
// Not finding config is not an error case - config will return nil.
func (c *baseCommand) initConfig(filename string) (*configpkg.Config, configpkg.ValidationResults, error) {
	path, err := c.initConfigPath(filename)
	if err != nil {
		return nil, nil, err
	}

	if path == "" {
		return nil, nil, nil
	}

	return c.initConfigLoad(path)
}

// initConfigPath returns the path for the configuration file with the
// specified filename.
func (c *baseCommand) initConfigPath(filename string) (string, error) {
	path, err := configpkg.FindPath("", filename, true)
	if err != nil {
		return "", fmt.Errorf("Error looking for a Waypoint configuration: %s", err)
	}

	return path, nil
}

// initConfigLoad loads the configuration at the given path.
func (c *baseCommand) initConfigLoad(path string) (*configpkg.Config, configpkg.ValidationResults, error) {
	cfg, err := configpkg.Load(path, &configpkg.LoadOptions{
		Pwd:       filepath.Dir(path),
		Workspace: c.refWorkspace.Workspace,
	})
	if err != nil {
		return nil, []configpkg.ValidationResult{{Error: err}}, err
	}

	// Validate
	results, err := cfg.Validate()
	if err != nil {
		return nil, results, err
	}
	return cfg, results, nil
}

// Given a duration, check the clientContext's auth token to determine if the
// token expires within the duration.
// Emits a warning on the cli if token expires within the duration.
// Returns any errors which may have occurred parsing the token.
//
// This function provides no guarantee that the token itself is valid as
// validation checks can only happen on the server.
func (c *baseCommand) checkTokenExpiry(duration time.Duration) error {
	token := c.clientContext.Server.AuthToken
	tokenMagic := "wp24"
	data, err := base58.Decode(token)
	if err != nil {
		return err
	}

	var tt pb.TokenTransport
	err = proto.Unmarshal(data[len(tokenMagic):], &tt)
	if err != nil {
		return err
	}

	var body pb.Token
	err = proto.Unmarshal(tt.Body, &body)
	if err != nil {
		return err
	}

	if body.ValidUntil != nil {
		te := time.Unix(body.ValidUntil.Seconds, int64(body.ValidUntil.Nanos))
		expireWarnPeriod := time.Now().Add(duration)
		if te.Before(expireWarnPeriod) {
			c.ui.Output(fmt.Sprintf("The token used to authenticate with Waypoint "+
				"will be expiring at %s. Please reauthenticate with Waypoint soon.", te), terminal.WithWarningStyle())
		}
	}
	return nil
}

// initClient initializes the client.
//
// If ctx is nil, c.Ctx will be used. If ctx is non-nil, that context will be
// used and c.Ctx will be ignored.
func (c *baseCommand) initClient(
	ctx context.Context,
	connectOpts ...serverclient.ConnectOption,
) (*clientpkg.Project, error) {
	// We use our flag-based connection info if the user set an addr.
	var flagConnection *clicontext.Config
	if v := c.flagConnection; v.Server.Address != "" {
		flagConnection = &v
	}

	// Get the context we'll use. The ordering here is purposeful and creates
	// the following precedence: (1) context (2) env (3) flags where the
	// later values override the former.
	var err error
	connectOpts = append([]serverclient.ConnectOption{
		serverclient.FromContext(c.contextStorage, ""),
		serverclient.FromEnv(),
		serverclient.FromContextConfig(flagConnection),
		serverclient.Logger(c.Log.Named("serverclient")),
	}, connectOpts...)
	c.clientContext, err = serverclient.ContextConfig(connectOpts...)
	if err != nil {
		return nil, err
	}

	if err := c.checkTokenExpiry(time.Hour * 24 * 7); err != nil {
		c.Log.Debug("Unable to decode token when checking token expiry.")
	}

	// Start building our client options
	opts := []clientpkg.Option{
		clientpkg.WithLogger(c.Log),
		clientpkg.WithClientConnect(connectOpts...),
		clientpkg.WithProjectRef(c.refProject),
		clientpkg.WithWorkspaceRef(c.refWorkspace),
		clientpkg.WithVariables(c.variables),
		clientpkg.WithLabels(c.flagLabels),
		clientpkg.WithSourceOverrides(c.flagRemoteSource),
	}
	if c.noLocalServer {
		opts = append(opts, clientpkg.WithNoLocalServer())
	}

	if c.flagLocal != nil {
		opts = append(opts, clientpkg.WithUseLocalRunner(*c.flagLocal))
	}

	if c.ui != nil {
		opts = append(opts, clientpkg.WithUI(c.ui))
	}

	if c.cfg != nil {
		opts = append(opts, clientpkg.WithConfig(c.cfg))
	}

	if ctx == nil {
		ctx = c.Ctx
	}

	// Create our client
	return clientpkg.New(ctx, opts...)
}
