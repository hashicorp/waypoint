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

// WithSingleApp configures the CLI to expect a configuration with
// one or more apps defined but a single app targeted with `-app`.
// If only a single app exists, it is implicitly the target.
// Zero apps is an error.
func WithSingleApp() Option {
	return func(c *baseConfig) {
		c.AppTargetRequired = true
		c.Config = false
		c.Client = true
	}
}

func WithMaybeApp() Option {
	return func(c *baseConfig) {
		c.AppTargetRequired = false
		c.MaybeApp = true
		c.Config = false
		c.Client = true
	}
}

// WithNoConfig configures the CLI to not expect any project configuration.
// This will not read any configuration files.
func WithNoConfig() Option {
	return func(c *baseConfig) {
		c.Config = false
	}
}

// WithConfig configures the CLI to find and load any project configuration.
// If optional is true, no error will be shown if a config can't be found.
func WithConfig(optional bool) Option {
	return func(c *baseConfig) {
		c.Config = true
		c.ConfigOptional = optional
	}
}

// WithClient configures the CLI to initialize a client.
func WithClient(v bool) Option {
	return func(c *baseConfig) {
		c.Client = v
	}
}

// WithUI configures the CLI to use a specific UI implementation
func WithUI(ui terminal.UI) Option {
	return func(c *baseConfig) {
		c.UI = ui
	}
}

// WithNoAutoServer configures the CLI to not automatically spin up
// an in-memory server for this command.
func WithNoAutoServer() Option {
	return func(c *baseConfig) {
		c.NoAutoServer = true
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
	Args              []string
	Flags             *flag.Sets
	Config            bool
	ConfigOptional    bool
	Client            bool
	AppTargetRequired bool
	MaybeApp          bool
	UI                terminal.UI

	// NoAutoServer is true if an in-memory server is not allowed.
	NoAutoServer bool

	// ConnArg as true means we should parse the server address as an
	// argument (the first argument).
	ConnArg bool
}
