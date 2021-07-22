package cli

import (
	"context"
	"errors"
	stdflag "flag"
	"io"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/adrg/xdg"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clicontext"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/config/variables"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/grpcmetadata"
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

	// clientContext is set to the context information for the current
	// connection. This might not exist in the contextStorage yet if this
	// is from an env var or flags.
	clientContext *clicontext.Config

	// contextStorage is for CLI contexts.
	contextStorage *clicontext.Storage

	// refProject and refWorkspace the references for this CLI invocation.
	refProject   *pb.Ref_Project
	refApp       *pb.Ref_Application
	refWorkspace *pb.Ref_Workspace

	// variables hold the values set via flags and local env vars
	variables []*pb.Variable

	//---------------------------------------------------------------
	// Internal fields that should not be accessed directly

	// flagPlain is whether the output should be in plain mode.
	flagPlain bool

	// flagLabels are set via -label if flagSetOperation is set.
	flagLabels map[string]string

	// flagVars sets values for defined input variables
	flagVars map[string]string

	// flagVarFile is a HCL or JSON file setting one or more values
	// for defined input variables
	flagVarFile []string

	// flagRemote is whether to execute using a remote runner or use
	// a local runner.
	flagRemote bool

	// flagRemoteSource are the remote data source overrides for jobs.
	flagRemoteSource map[string]string

	// flagApp is the app to target.
	flagApp string

	// flagWorkspace is the workspace to work in.
	flagWorkspace string

	// flagConnection contains manual flag-based connection info.
	flagConnection clicontext.Config

	// args that were present after parsing flags
	args []string

	// options passed in at the global level
	globalOptions []Option

	// autoServer will be set to true if an automatic in-memory server
	// is allowd.
	autoServer bool

	// The home directory that we loaded the waypoint config from
	homeConfigPath string
}

// Close cleans up any resources that the command created. This should be
// defered by any CLI command that embeds baseCommand in the Run command.
func (c *baseCommand) Close() error {
	// Close our UI if it implements it. The glint-based UI does for example
	// to finish up all the CLI output.
	if closer, ok := c.ui.(io.Closer); ok && closer != nil {
		closer.Close()
	}

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

	for _, opt := range c.globalOptions {
		opt(&baseCfg)
	}

	for _, opt := range opts {
		opt(&baseCfg)
	}

	// Set some basic internal fields
	c.autoServer = !baseCfg.NoAutoServer

	// Init our UI first so we can write output to the user immediately.
	ui := baseCfg.UI
	if ui == nil {
		ui = terminal.ConsoleUI(c.Ctx)
	}

	c.ui = ui

	// Parse flags
	if err := baseCfg.Flags.Parse(baseCfg.Args); err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return err
	}
	c.args = baseCfg.Flags.Args()

	// Check for flags after args
	if err := checkFlagsAfterArgs(c.args, baseCfg.Flags); err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return err
	}

	// Reset the UI to plain if that was set
	if c.flagPlain {
		c.ui = terminal.NonInteractiveUI(c.Ctx)
	}

	// If we're parsing the connection from the arg, then use that.
	if baseCfg.ConnArg && len(c.args) > 0 {
		if err := c.flagConnection.FromURL(c.args[0]); err != nil {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return err
		}

		c.args = c.args[1:]
	}

	// With the flags we now know what workspace we're targeting
	c.refWorkspace = &pb.Ref_Workspace{Workspace: c.flagWorkspace}

	// Setup our base config path
	homeConfigPath, err := xdg.ConfigFile("waypoint/.ignore")
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return err
	}
	homeConfigPath = filepath.Dir(homeConfigPath)
	c.homeConfigPath = homeConfigPath
	c.Log.Debug("home configuration directory", "path", homeConfigPath)

	// Setup our base directory for context management
	contextStorage, err := clicontext.NewStorage(
		clicontext.WithDir(filepath.Join(homeConfigPath, "context")))
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return err
	}
	c.contextStorage = contextStorage

	// Parse the configuration
	c.cfg = &config.Config{}

	// If we have an app target requirement, we have to get it from the args
	// or the config.
	if baseCfg.AppTargetRequired {
		// If we have args, attempt to extract there first.
		if len(c.args) > 0 {
			match := reAppTarget.FindStringSubmatch(c.args[0])
			if match != nil {
				// Set our refs
				c.refProject = &pb.Ref_Project{Project: match[1]}
				c.refApp = &pb.Ref_Application{
					Project:     match[1],
					Application: match[2],
				}

				// Shift the args
				c.args = c.args[1:]

				// Explicitly set remote
				c.flagRemote = true
			}
		}

		if c.refApp == nil {
			baseCfg.Config = true
		}
	}

	// TODO: not sure if it's ok to move this. This is in support
	// of getting a project ref with an app target isn't a hard requirement
	// If we didn't get our ref, then we need to load config
	// Add an additional option?
	// if baseCfg.MaybeAppTarget {}
	if c.refApp == nil && baseCfg.MaybeApp {
		baseCfg.Config = true
	}

	// If we're loading the config, then get it.
	if baseCfg.Config {
		cfg, err := c.initConfig("", baseCfg.ConfigOptional)
		if err != nil {
			c.logError(c.Log, "failed to load config", err)
			return err
		}

		c.cfg = cfg
		if cfg != nil {
			c.refProject = &pb.Ref_Project{Project: cfg.Project}

			// If we require an app target and we still haven't set it,
			// and the user provided it via the CLI, set it now. This code
			// path is only reached if it wasn't set via the args either
			// above.
			if baseCfg.AppTargetRequired &&
				c.refApp == nil &&
				c.flagApp != "" {
				c.refApp = &pb.Ref_Application{
					Project:     cfg.Project,
					Application: c.flagApp,
				}
			}
		}
	}

	// Collect variable values from -var and -varfile flags,
	// and env vars set with WP_VAR_* and set them on the job
	vars, diags := variables.LoadVariableValues(c.flagVars, c.flagVarFile)
	if diags.HasErrors() {
		// we only return errors for file parsing, so we are specific
		// in the error log here
		c.logError(c.Log, "failed to load wpvars file", errors.New(diags.Error()))
		return diags
	}
	c.variables = vars

	// Create our client
	if baseCfg.Client {
		c.project, err = c.initClient(nil)
		if err != nil {
			c.logError(c.Log, "failed to create client", err)
			return err
		}
	}

	// Validate remote vs. local operations.
	if c.flagRemote && c.refApp == nil {
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
		if c.refApp == nil {
			if len(c.cfg.Apps()) != 1 {
				c.ui.Output(errAppModeSingle, terminal.WithErrorStyle())
				return ErrSentinel
			}

			c.refApp = &pb.Ref_Application{
				Project:     c.cfg.Project,
				Application: c.cfg.Apps()[0],
			}
		}
	}

	return nil
}

// DoApp calls the callback for each app. This lets you execute logic
// in an app-specific context safely. This automatically handles any
// parallelization, waiting, and error handling. Your code should be
// thread-safe.
//
// If any error is returned, the caller should just exit. The error handling
// including messaging to the user is handled by this function call.
//
// If you want to early exit all the running functions, you should use
// the callback closure properties to cancel the passed in context. This
// will stop any remaining callbacks and exit early.
func (c *baseCommand) DoApp(ctx context.Context, f func(context.Context, *clientpkg.App) error) error {
	var appTargets []string
	if c.refApp != nil {
		appTargets = []string{c.refApp.Application}
	} else if c.cfg != nil {
		appTargets = append(appTargets, c.cfg.Apps()...)
	}

	var apps []*clientpkg.App
	for _, appName := range appTargets {
		app := c.project.App(appName)
		c.Log.Debug("will operate on app", "name", appName)
		apps = append(apps, app)
	}

	// Inject the metadata about the client, such as the runner id if it is running
	// a local runner.
	if id, ok := c.project.LocalRunnerId(); ok {
		ctx = grpcmetadata.AddRunner(ctx, id)
	}

	// Just a serialize loop for now, one day we'll parallelize.
	var finalErr error
	var didErrSentinel bool
	for _, app := range apps {
		// Support cancellation
		if err := ctx.Err(); err != nil {
			return err
		}

		if err := f(ctx, app); err != nil {
			if err != ErrSentinel {
				finalErr = multierror.Append(finalErr, err)
			} else {
				didErrSentinel = true
			}
		}
	}
	if finalErr == nil && didErrSentinel {
		finalErr = ErrSentinel
	}

	return finalErr
}

// logError logs an error and outputs it to the UI.
func (c *baseCommand) logError(log hclog.Logger, prefix string, err error) {
	if err == ErrSentinel {
		return
	}

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

		f.BoolVar(&flag.BoolVar{
			Name:    "plain",
			Target:  &c.flagPlain,
			Default: false,
			Usage:   "Plain output: no colors, no animation.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "app",
			Target:  &c.flagApp,
			Default: "",
			Usage: "App to target. Certain commands require a single app target for " +
				"Waypoint configurations with multiple apps. If you have a single app, " +
				"then this can be ignored.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "workspace",
			Target:  &c.flagWorkspace,
			Default: "default",
			Usage:   "Workspace to operate in.",
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

		f.StringMapVar(&flag.StringMapVar{
			Name:   "remote-source",
			Target: &c.flagRemoteSource,
			Usage: "Override configurations for how remote runners source data. " +
				"This is specified to the data source type being used in your configuration. " +
				"This is used for example to set a specific Git ref to run against.",
		})

		f.StringMapVar(&flag.StringMapVar{
			Name:   "var",
			Target: &c.flagVars,
			Usage:  "Variable value to set for this operation. Can be specified multiple times.",
		})

		f.StringSliceVar(&flag.StringSliceVar{
			Name:   "var-file",
			Target: &c.flagVarFile,
			Usage: "HCL or JSON file containing variable values to set for this " +
				"operation. If any \"*.auto.wpvars\" or \"*.auto.wpvars.json\" " +
				"files are present, they will be automatically loaded.",
		})
	}

	if bit&flagSetConnection != 0 {
		f := set.NewSet("Connection Options")
		f.StringVar(&flag.StringVar{
			Name:   "server-addr",
			Target: &c.flagConnection.Server.Address,
			Usage:  "Address for the server.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "server-tls",
			Target:  &c.flagConnection.Server.Tls,
			Default: true,
			Usage:   "True if the server should be connected to via TLS.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "server-tls-skip-verify",
			Target:  &c.flagConnection.Server.TlsSkipVerify,
			Default: false,
			Usage:   "True to skip verification of the TLS certificate advertised by the server.",
		})
	}

	if f != nil {
		// Configure our values
		f(set)
	}

	return set
}

// checkFlagsAfterArgs checks for a very common user error scenario where
// CLI flags are specified after positional arguments. Since we use the
// stdlib flag package, this is not allowed. However, we can detect this
// scenario, and notify a user. We can't easily automatically fix it because
// its hard to tell positional vs intentional flags.
func checkFlagsAfterArgs(args []string, set *flag.Sets) error {
	if len(args) == 0 {
		return nil
	}

	// Build up our arg map for easy searching.
	flagMap := map[string]struct{}{}
	for _, v := range args {
		// If we reach a "--" we're done. This is a common designator
		// in CLIs (such as exec) that everything following is fair game.
		if v == "--" {
			break
		}

		// There is always at least 2 chars in a flag "-v" example.
		if len(v) < 2 {
			continue
		}

		// Flags start with a hyphen
		if v[0] != '-' {
			continue
		}

		// Detect double hyphen flags too
		if v[1] == '-' {
			v = v[1:]
		}

		// More than double hyphen, ignore. note this looks like we can
		// go out of bounds and panic cause this is the 3rd char if we have
		// a double hyphen and we only protect on 2, but since we check first
		// against plain "--" we know that its not exactly "--" AND the length
		// is at least 2, meaning we can safely imply we have length 3+ for
		// double-hyphen prefixed values.
		if v[1] == '-' {
			continue
		}

		// If we have = for "-foo=bar", trim out the =.
		if idx := strings.Index(v, "="); idx >= 0 {
			v = v[:idx]
		}

		flagMap[v[1:]] = struct{}{}
	}

	// Now look for anything that looks like a flag we accept. We only
	// look for flags we accept because that is the most common error and
	// limits the false positives we'll get on arguments that want to be
	// hyphen-prefixed.
	didIt := false
	set.VisitSets(func(name string, s *flag.Set) {
		s.VisitAll(func(f *stdflag.Flag) {
			if _, ok := flagMap[f.Name]; ok {
				// Uh oh, we done it. We put a flag after an arg.
				didIt = true
			}
		})
	})

	if didIt {
		return errFlagAfterArgs
	}

	return nil
}

// flagSetBit is used with baseCommand.flagSet
type flagSetBit uint

const (
	flagSetNone       flagSetBit = 1 << iota
	flagSetOperation             // shared flags for operations (build, deploy, etc)
	flagSetConnection            // shared flags for server connections
)

var (
	// ErrSentinel is a sentinel value that we can return from Init to force an exit.
	ErrSentinel = errors.New("error sentinel")

	errFlagAfterArgs = errors.New(strings.TrimSpace(`
Flags must be specified before positional arguments in the CLI command.
For example "waypoint up -example project" not "waypoint up project -example".
Please reorder your arguments and try again.

Note: we can't automatically fix this or allow this since we can't safely
detect what you want as flag arguments and what you want as positional arguments.
The underlying library we use for flag parsing (the Go standard library)
enforces this requirement. Sorry!
`))

	errAppModeSingle = strings.TrimSpace(`
This command requires a single targeted app. You have multiple apps defined
so you can specify the app to target using the "-app" flag.
`)

	reAppTarget = regexp.MustCompile(`^(?P<project>[-0-9A-Za-z_]+)/(?P<app>[-0-9A-Za-z_]+)$`)

	snapshotUnimplementedErr = strings.TrimSpace(`
The current Waypoint server does not support snapshots. Rerunning the command
with '-snapshot=false' is required, and there will be no automatic data backups
for the server.
`)
)
