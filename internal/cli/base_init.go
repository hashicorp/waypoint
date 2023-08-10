// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cli

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/hashicorp/waypoint/internal/clicontext"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	configpkg "github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/pkg/serverclient"
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
