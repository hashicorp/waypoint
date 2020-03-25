package cli

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
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

	// app is the targeted application. This is only set if you use the
	// WithSingleApp option. You should not access this directly
	// though and use the DoApp function.
	app string
}

// Init initializes the command by parsing flags, parsing the configuration,
// setting up the project, etc. You can control what is done by using the
// options.
//
// Init should be called FIRST within the Run function implementation. Many
// options will affect behavior of other functions that can be called later.
func (c *baseCommand) Init(opts ...Option) error {
	var baseCfg baseConfig
	for _, opt := range opts {
		opt(&baseCfg)
	}

	// Init our UI first so we can write output to the user immediately.
	c.ui = &terminal.BasicUI{}

	// Parse flags
	if err := baseCfg.Flags.Parse(baseCfg.Args); err != nil {
		c.ui.Output(err.Error(), terminal.WithErrorStyle())
		return err
	}

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

	// If this is a single app mode then make sure that we only have
	// one app or that we have an app target.
	if baseCfg.AppMode == appModeSingle {
		// TODO(mitchellh): when we support app targeting we can have more
		// than one as long as its targeted.
		if len(cfg.Apps) != 1 {
			c.project.UI.Output(errAppModeSingle, terminal.WithErrorStyle())
			return ErrSentinel
		}

		// Set our targeted app
		c.app = cfg.Apps[0].Name
	}

	return nil
}

// DoApp calls the callback for each app. This lets you execute logic
// in an app-specific context safely. This automatically handles any
// parallelization, waiting, and error handling. Your code should be
// thread-safe.
//
// If any error is returned, the caller should just exit. The error handling
// including messaging to the user is handling by this function call.
//
// If you want to early exit all the running functions, you should use
// the callback closure properties to cancel the passed in context. This
// will stop any remaining callbacks and exit early.
func (c *baseCommand) DoApp(ctx context.Context, f func(context.Context, *core.App) error) error {
	var apps []*core.App
	for _, appCfg := range c.cfg.Apps {
		// If we're doing single targeting and this app isn't what we
		// want then continue. In practice we don't need to loop at all
		// for single targeting but it simplifies the implementation and
		// the performance here doesn't matter currently.
		if c.app != "" && appCfg.Name != c.app {
			continue
		}

		app, err := c.project.App(appCfg.Name)
		if err != nil {
			panic(err)
		}

		c.Log.Debug("will operate on app", "name", appCfg.Name)
		apps = append(apps, app)
	}

	// Just a serialize loop for now, one day we'll parallelize.
	var finalErr error
	for _, app := range apps {
		// Support cancellation
		if err := ctx.Err(); err != nil {
			return err
		}

		if err := f(ctx, app); err != nil {
			finalErr = multierror.Append(finalErr, err)
		}
	}

	return finalErr
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

		if f != nil {
			// Configure our values
			f(set)
		}

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

var (
	// ErrSentinel is a sentinel value that we can return from Init to force an exit.
	ErrSentinel = errors.New("error sentinel")

	errAppModeSingle = strings.TrimSpace(`
This command requires a single targeted app. You have multiple apps defined
so you can specify the app to target using the "-app" flag.
`)
)
