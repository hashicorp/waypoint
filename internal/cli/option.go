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

// WithProjectTargetRequired indicates that the command needs to operate against a project.
// Will require parsing config
func WithProjectTargetRequired() Option {
	return func(c *baseConfig) {
		// TODO(izaak): we should also parse config here
		c.ProjectTargetRequired = true
	}
}

// WithSingleApp configures the CLI to expect a configuration with
// one or more apps defined but a single app targeted with `-app`.
// If only a single app exists, it is implicitly the target.
// Zero apps is an error.
func WithSingleApp() Option {
	return func(c *baseConfig) {
		c.AppTargetRequired = true
		c.Config = false
	}
}

// WithMultipleApp configures the CLI to expect a configuration with
// one or more apps defined in a project. The option will prioritize a value
// provided to the -project flag.
func WithMultipleApp() Option {
	return func(c *baseConfig) {
		c.ProjectTargetRequired = true
		c.Config = false
	}
}

// TODO(izaak): delete
// WithOptionalApp configures the CLI to work with or without an explicit
// project config locally. It also allows for operations on multiple apps
// inside a project.
func WithOptionalApp() Option {
	return func(c *baseConfig) {
		c.AppTargetRequired = false
		c.AppOptional = true
		c.Config = false
		c.ConfigOptional = true
	}
}

// WithNoConfig configures the CLI to not expect any project configuration.
// This will not read any configuration files.
func WithNoConfig() Option {
	return func(c *baseConfig) {
		c.Config = false
	}
}

// TODO(izaak): delete
//// WithConfig configures the CLI to find and load any project configuration.
//// If optional is true, no error will be shown if a config can't be found.
//func WithConfig(optional bool) Option {
//	return func(c *baseConfig) {
//		c.Config = true
//		c.ConfigOptional = optional
//	}
//}

// WithNoClient prevents the CLI from instantiating a client automatically.
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

func WithRunnerRequired() Option {
	return func(c *baseConfig) {
		c.RunnerRequired = true
	}
}

type baseConfig struct {
	Args  []string
	Flags *flag.Sets

	RunnerRequired        bool
	ProjectTargetRequired bool
	NoClient              bool

	// TODO(izaak): how much of this can we delete? All of it?
	Config            bool
	ConfigOptional    bool
	AppOptional       bool
	AppTargetRequired bool

	UI terminal.UI

	// NoAutoServer is true if an in-memory server is not allowed.
	NoAutoServer bool

	// ConnArg as true means we should parse the server address as an
	// argument (the first argument).
	ConnArg bool
}
