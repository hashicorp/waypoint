package cli

import (
	"context"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2/hclsimple"

	"github.com/mitchellh/devflow/internal/config"
	"github.com/mitchellh/devflow/internal/core"
	"github.com/mitchellh/devflow/internal/pkg/flag"
	"github.com/mitchellh/devflow/sdk/datadir"
	"github.com/mitchellh/devflow/sdk/terminal"
)

// baseCommand is embedded in all commands to provide common logic and data.
//
// The unexported values are not available until after Init is called. Some
// values are only available in certain circumstances, read the documentation
// for the field to determine if that is the case.
type baseCommand struct {
	// Ctx is the base context for the command. It is up to commands to
	// utilize this context so that cancellation works in a timely manner.
	Ctx context.Context

	// Log is the logger to use.
	Log hclog.Logger

	//---------------------------------------------------------------
	// The fields below are only available after calling Init.

	// cfg is the parsed configuration
	cfg *config.Config

	// dir is the project directory
	dir *datadir.Project

	// project is the main project for the configuration
	project *core.Project

	// UI is used to write to the CLI.
	ui terminal.UI

	//---------------------------------------------------------------
	// Internal fields that should not be accessed directly

	flags     *flag.Sets
	flagsOnce sync.Once
}

// Init initializes the command by parsing flags, parsing the configuration,
// setting up the project, etc. You can control what is done by using the
// options.
func (c *baseCommand) Init() error {
	// Init our UI first so we can write output to the user immediately.
	c.ui = &terminal.BasicUI{}

	// Parse the configuration
	c.Log.Debug("reading configuration", "path", "devflow.hcl")
	var cfg config.Config
	if err := hclsimple.DecodeFile("devflow.hcl", nil, &cfg); err != nil {
		c.logError(c.Log, "error decoding configuration", err)
		return err
	}
	c.cfg = &cfg

	// Setup our project data directory
	c.Log.Debug("preparing project directory", "path", ".devflow")
	projDir, err := datadir.NewProject(".devflow")
	if err != nil {
		c.logError(c.Log, "error preparing data directory", err)
		return err
	}
	c.dir = projDir

	// Create our project
	c.project, err = core.NewProject(c.Ctx,
		core.WithLogger(c.Log),
		core.WithConfig(&cfg),
		core.WithDataDir(projDir),
	)
	if err != nil {
		c.logError(c.Log, "failed to create project", err)
		return err
	}

	return nil
}

// logError logs an error and outputs it to the UI.
func (c *baseCommand) logError(log hclog.Logger, prefix string, err error) {
	log.Error(prefix, "error", err)
	c.ui.Output("%s: %s", prefix, err, terminal.WithErrorStyle())
}

// flagSet creates the flags for this command. The callback should be used
// to configure the set with your own custom options.
//
// The result is cached on the command to save performance on future calls
// since it may be called multiple times for help generation.
func (c *baseCommand) flagSet(bit flagSetBit, f func(*flag.Sets)) *flag.Sets {
	c.flagsOnce.Do(func() {
		set := flag.NewSets()

		// Configure our values
		f(set)

		// Cache
		c.flags = set
	})

	return c.flags
}

// flagSetBit is used with baseCommand.flagSet
type flagSetBit uint

const (
	flagSetNone flagSetBit = 1 << iota
	flagSetHTTP            // not used currently, should replace when we don't need
)
