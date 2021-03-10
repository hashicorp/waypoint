package cli

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/hashicorp/waypoint/internal/clicontext"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	configpkg "github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/serverclient"
)

// This file contains the various methods that are used to perform
// the Init call on baseCommand. They are broken down into individual
// smaller methods for readability but more importantly to power the
// "init" subcommand. This allows us to share as much logic as possible
// between Init and "init" to help ensure that "init" succeeding means that
// other commands will succeed as well.

// initConfig initializes the configuration with the specified filename from the CLI.
// If filename is empty, it will default to configpkg.Filename.
func (c *baseCommand) initConfig(filename string, optional bool) (*configpkg.Config, error) {
	path, err := c.initConfigPath(filename)
	if err != nil {
		return nil, err
	}

	if path == "" {
		if optional {
			return nil, nil
		}

		return nil, errors.New("A Waypoint configuration file is required but wasn't found.")
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
func (c *baseCommand) initConfigLoad(path string) (*configpkg.Config, error) {
	cfg, err := configpkg.Load(path, filepath.Dir(path))
	if err != nil {
		return nil, err
	}

	// Validate
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// initClient initializes the client.
func (c *baseCommand) initClient() (*clientpkg.Project, error) {
	// We use our flag-based connection info if the user set an addr.
	var flagConnection *clicontext.Config
	if v := c.flagConnection; v.Server.Address != "" {
		flagConnection = &v
	}

	// Get the context we'll use. The ordering here is purposeful and creates
	// the following precedence: (1) context (2) env (3) flags where the
	// later values override the former.
	var err error
	connectOpts := []serverclient.ConnectOption{
		serverclient.FromContext(c.contextStorage, ""),
		serverclient.FromEnv(),
		serverclient.FromContextConfig(flagConnection),
	}
	c.clientContext, err = serverclient.ContextConfig(connectOpts...)
	if err != nil {
		return nil, err
	}

	// Start building our client options
	opts := []clientpkg.Option{
		clientpkg.WithLogger(c.Log),
		clientpkg.WithClientConnect(connectOpts...),
		clientpkg.WithProjectRef(c.refProject),
		clientpkg.WithWorkspaceRef(c.refWorkspace),
		clientpkg.WithLabels(c.flagLabels),
		clientpkg.WithSourceOverrides(c.flagRemoteSource),
	}
	if !c.flagRemote && c.autoServer {
		opts = append(opts, clientpkg.WithLocal())
	}

	if c.ui != nil {
		opts = append(opts, clientpkg.WithUI(c.ui))
	}

	// Create our client
	return clientpkg.New(c.Ctx, opts...)
}
