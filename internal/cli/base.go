package cli

import (
	"context"
	"errors"
	"io"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"

	"github.com/hashicorp/waypoint/internal/clicontext"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/sdk/terminal"
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

	// LogOutput is the writer that Log points to. You SHOULD NOT use
	// this directly. We have access to this so you can use
	// hclog.OutputResettable if necessary.
	LogOutput io.Writer

	//---------------------------------------------------------------
	// The fields below are only available after calling Init.

	// cfg is the parsed configuration
	cfg *config.Config

	// UI is used to write to the CLI.
	ui terminal.UI

	// client for performing operations
	project *clientpkg.Project

	// contextStorage is for CLI contexts.
	contextStorage *clicontext.Storage

	// refProject and refWorkspace the references for this CLI invocation.
	refProject   *pb.Ref_Project
	refWorkspace *pb.Ref_Workspace

	//---------------------------------------------------------------
	// Internal fields that should not be accessed directly

	// app is the targeted application. This is only set if you use the
	// WithSingleApp option. You should not access this directly
	// though and use the DoApp function.
	app string

	// flagLabels are set via -label if flagSetOperation is set.
	flagLabels map[string]string

	// flagRemote is whether to execute using a remote runner or use
	// a local runner.
	flagRemote bool

	// flagWorkspace is the workspace to work in.
	flagWorkspace string

	// args that were present after parsing flags
	args []string
}

// Close cleans up any resources that the command created. This should be
// defered by any CLI command that embeds baseCommand in the Run command.
func (c *baseCommand) Close() error {
	return nil
}

// Init initializes the command by parsing flags, parsing the configuration,
// setting up the project, etc. You can control what is done by using the
// options.
//
// Init should be called FIRST within the Run function implementation. Many
// options will affect behavior of other functions that can be called later.
func (c *baseCommand) Init(opts ...Option) error {
	baseCfg := baseConfig{
		Config: true,
		Client: true,
	}
	for _, opt := range opts {
		opt(&baseCfg)
	}

	// Init our UI first so we can write output to the user immediately.
	c.ui = terminal.ConsoleUI(c.Ctx)

	// Parse flags
	if err := baseCfg.Flags.Parse(baseCfg.Args); err != nil {
		c.ui.Output(err.Error(), terminal.WithErrorStyle())
		return err
	}
	c.args = baseCfg.Flags.Args()

	// With the flags we now know what workspace we're targeting
	c.refWorkspace = &pb.Ref_Workspace{Workspace: c.flagWorkspace}

	// Setup our base config path
	homeConfigPath, err := xdg.ConfigFile("waypoint/.ignore")
	if err != nil {
		c.ui.Output(err.Error(), terminal.WithErrorStyle())
		return err
	}
	homeConfigPath = filepath.Dir(homeConfigPath)
	c.Log.Debug("home configuration directory", "path", homeConfigPath)

	// Setup our base directory for context management
	contextStorage, err := clicontext.NewStorage(
		clicontext.WithDir(filepath.Join(homeConfigPath, "context")))
	if err != nil {
		c.ui.Output(err.Error(), terminal.WithErrorStyle())
		return err
	}
	c.contextStorage = contextStorage

	// Parse the configuration
	c.cfg = &config.Config{}

	// If we're loading the config, then get it.
	if baseCfg.Config {
		cfg, err := c.initConfig(baseCfg.ConfigOptional)
		if err != nil {
			return err
		}

		c.cfg = cfg
		if cfg != nil {
			c.refProject = &pb.Ref_Project{Project: cfg.Project}
		}
	}

	// Create our client
	if baseCfg.Client {
		c.project, err = c.initClient()
		if err != nil {
			c.logError(c.Log, "failed to create client", err)
			return err
		}
	}

	// Validate remote vs. local operations.
	if c.flagRemote {
		if c.cfg == nil || c.cfg.Runner == nil || !c.cfg.Runner.Enabled {
			err := errors.New(
				"The `-remote` flag was specified but remote operations are not supported\n" +
					"for this project.\n\n" +
					"Remote operations must be manually enabled by using setting the 'runner.enabled'\n" +
					"setting in your Waypoint configuration file. Please see the documentation\n" +
					"on this setting for more information.")
			c.logError(c.Log, "", err)
			return err
		}
	}

	// If this is a single app mode then make sure that we only have
	// one app or that we have an app target.
	if baseCfg.AppTargetRequired {
		// TODO(mitchellh): when we support app targeting we can have more
		// than one as long as its targeted.
		if len(c.cfg.Apps) != 1 {
			c.ui.Output(errAppModeSingle, terminal.WithErrorStyle())
			return ErrSentinel
		}

		// Set our targeted app
		c.app = c.cfg.Apps[0].Name
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
func (c *baseCommand) DoApp(ctx context.Context, f func(context.Context, *clientpkg.App) error) error {
	var apps []*clientpkg.App
	for _, appCfg := range c.cfg.Apps {
		// If we're doing single targeting and this app isn't what we
		// want then continue. In practice we don't need to loop at all
		// for single targeting but it simplifies the implementation and
		// the performance here doesn't matter currently.
		if c.app != "" && appCfg.Name != c.app {
			continue
		}

		app := c.project.App(appCfg.Name)
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

// Retrieve the app config for the named application
func (c *baseCommand) AppConfig(name string) (*config.App, bool) {
	for _, appCfg := range c.cfg.Apps {
		if appCfg.Name == name {
			return appCfg, true
		}
	}

	return nil, false
}

// logError logs an error and outputs it to the UI.
func (c *baseCommand) logError(log hclog.Logger, prefix string, err error) {
	log.Error(prefix, "error", err)

	if prefix != "" {
		prefix += ": "
	}
	c.ui.Output("%s%s", prefix, err, terminal.WithErrorStyle())
}

// flagSet creates the flags for this command. The callback should be used
// to configure the set with your own custom options.
func (c *baseCommand) flagSet(bit flagSetBit, f func(*flag.Sets)) *flag.Sets {
	set := flag.NewSets()
	{
		f := set.NewSet("Global Options")
		f.StringVar(&flag.StringVar{
			Name:    "workspace",
			Target:  &c.flagWorkspace,
			Default: "default",
			Usage:   "Workspace to operate in. Defaults to 'default'.",
		})
	}

	if bit&flagSetOperation != 0 {
		f := set.NewSet("Operation Options")
		f.StringMapVar(&flag.StringMapVar{
			Name:   "label",
			Target: &c.flagLabels,
			Usage:  "Labels to set for this operation. Can be specified multiple times.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "remote",
			Target:  &c.flagRemote,
			Default: false,
			Usage: "True to use a remote runner to execute. This defaults to false \n" +
				"unless 'runner.default' is set in your configuration.",
		})
	}

	if f != nil {
		// Configure our values
		f(set)
	}

	return set
}

// flagSetBit is used with baseCommand.flagSet
type flagSetBit uint

const (
	flagSetNone      flagSetBit = 1 << iota
	flagSetOperation            // shared flags for operations (build, deploy, etc)
)

var (
	// ErrSentinel is a sentinel value that we can return from Init to force an exit.
	ErrSentinel = errors.New("error sentinel")

	errAppModeSingle = strings.TrimSpace(`
This command requires a single targeted app. You have multiple apps defined
so you can specify the app to target using the "-app" flag.
`)
)
