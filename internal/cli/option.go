package cli

import (
	"github.com/mitchellh/devflow/internal/pkg/flag"
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

// WithSingleApp configures the CLI to expect a configuration with
// one or more apps defined but a single app targeted with `-app`.
// If only a single app exists, it is implicitly the target.
// Zero apps is an error.
func WithSingleApp() Option {
	return func(c *baseConfig) { c.AppMode = appModeSingle }
}

type baseConfig struct {
	Args    []string
	Flags   *flag.Sets
	AppMode appMode
}

// appMode is used with baseConfig to specify how we handle multiple
// apps in a configuration file. See the different Option functions more
// detailed documentation on each app mode.
type appMode uint8

const (
	appModeNone   appMode = iota // no apps required, no config required
	appModeSingle                // must target a single app
	appModeMulti                 // one or more apps, can target single
)
