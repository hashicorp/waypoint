// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cli

import (
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

// Option is used to configure Init on baseCommand.
type Option func(c *baseConfig)

// WithArgs sets the arguments to the command that are used for parsing.
// Remaining arguments can be accessed using your flag set and asking for Args.
// Example: c.Flags().Args().
func WithArgs(args []string) Option {
	return func(c *baseConfig) { c.Args = args }
}

// WithFlags sets the flags that are supported by this command. This MUST
// be set otherwise a panic will happen. This is usually set by just calling
// the Flags function on your command implementation.
func WithFlags(f *flag.Sets) Option {
	return func(c *baseConfig) { c.Flags = f }
}

// WithProjectTarget configures the CLI to expect a configuration
// with a project specified, either through a waypoint.hcl file or
// the -project flag
func WithProjectTarget() Option {
	return func(c *baseConfig) {
		c.ProjectTargetRequired = true
	}
}

// WithSingleAppTarget configures the CLI to expect a configuration with
// one or more apps defined but a single app targeted with `-app`.
// If only a single app exists, it is implicitly the target.
// Zero apps is an error.
func WithSingleAppTarget() Option {
	return func(c *baseConfig) {
		c.SingleAppTarget = true
		// Projects are implicitly required for app targets
		c.ProjectTargetRequired = true
	}
}

// WithMultiAppTargets configures the CLI to allow this command to run against
// every app specified by the user - either each application in the project,
// or the single app specified with the -app flag.
func WithMultiAppTargets() Option {
	return func(c *baseConfig) {
		c.MultiAppTarget = true
		// Projects are implicitly required for app targets
		c.ProjectTargetRequired = true
	}
}

// WithNoConfig configures the CLI to not expect any project configuration.
// This will not read any configuration files.
func WithNoConfig() Option {
	return func(c *baseConfig) {
		c.NoConfig = true
	}
}

// WithNoClient configures the CLI to not use a client
func WithNoClient() Option {
	return func(c *baseConfig) {
		c.NoClient = true
	}
}

// WithUI configures the CLI to use a specific UI implementation
func WithUI(ui terminal.UI) Option {
	return func(c *baseConfig) {
		c.UI = ui
	}
}

// WithNoLocalServer configures the CLI to not automatically spin up
// an in-memory server for this command.
func WithNoLocalServer() Option {
	return func(c *baseConfig) {
		c.NoLocalServer = true
	}
}

// WithConnectionArg parses the first argument in the CLI as connection
// info if it exists. This parses it according to the clicontext.Config.FromURL
// method.
func WithConnectionArg() Option {
	return func(c *baseConfig) {
		c.ConnArg = true
	}
}

type baseConfig struct {
	Args  []string
	Flags *flag.Sets

	ProjectTargetRequired bool
	NoClient              bool
	SingleAppTarget       bool
	MultiAppTarget        bool
	NoConfig              bool

	UI terminal.UI

	// NoLocalServer is true if an in-memory server is not allowed.
	NoLocalServer bool

	// ConnArg as true means we should parse the server address as an
	// argument (the first argument).
	ConnArg bool
}
