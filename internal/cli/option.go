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

type baseConfig struct {
	Args  []string
	Flags *flag.Sets
}
