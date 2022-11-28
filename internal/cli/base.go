package cli

import (
	"context"
	"errors"
	stdflag "flag"
	"fmt"
	"io"
	"os"
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
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

const (
	defaultWorkspace        = "default"
	defaultWorkspaceEnvName = "WAYPOINT_WORKSPACE"
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

	// refProject references the project for this CLI invocation. A project
	// reference will be looked for in config, or in the -project flag.
	refProject *pb.Ref_Project

	// refApps references the apps for this CLI invocation. An app
	// reference will be looked for in config, or in the -app flag.
	refApps []*pb.Ref_Application

	// refWorkspace referenced the workspace for this CLI invocation
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

	// flagLocal indicates that any operations performed must happen with a local runner
	// not a remote runner.
	flagLocal *bool

	// flagRemoteSource are the remote data source overrides for jobs.
	flagRemoteSource map[string]string

	// flagApp is the app to target.
	flagApp string

	// flagProject is the project to target.
	flagProject string

	// flagWorkspace is the workspace to work in.
	flagWorkspace string

	// flagConnection contains manual flag-based connection info.
	flagConnection clicontext.Config

	// args that were present after parsing flags
	args []string

	// options passed in at the global level
	globalOptions []Option

	// noLocalServer prevents the creation of a local in-memory server
	noLocalServer bool

	// The home directory that we loaded the waypoint config from
	homeConfigPath string

	// deprecatedFlagRemote is whether to execute using a remote runner or use
	// a local runner.
	deprecatedFlagRemote *bool
}

// Close cleans up any resources that the command created. This should be
// deferred by any CLI command that embeds baseCommand in the Run command.
func (c *baseCommand) Close() error {
	// Close the project client, which gracefully shuts down the local runner
	if c.project != nil {
		c.project.Close()
	}

	// Close our UI if it implements it. The glint-based UI does for example
	// to finish up all the CLI output.
	if closer, ok := c.ui.(io.Closer); ok && closer != nil {
		closer.Close()
	}

	return nil
}

// Checks for deprecated flags and args.
func (c *baseCommand) checkDeprecatedFlags() error {
	// Check for deprecated project/app syntax.
	// NOTE(izaak): we should remove this in the next major (v0.8.0) because it
	// collides with arguments that contain a single slash (i.e. `waypoint exec bin/bash`)
	if len(c.args) > 0 {
		match := reAppTarget.FindStringSubmatch(c.args[0])
		if match != nil {
			return errors.New(errDeprecatedProjectAppArg)
		}
	}

	// Check for deprecated remote flag
	if c.deprecatedFlagRemote != nil {
		return fmt.Errorf("The -remote flag has been deprecated. Use -local=%t instead", !*c.deprecatedFlagRemote)
	}
	return nil
}

func (c *baseCommand) showValidations(validationResults config.ValidationResults) {
	if len(validationResults) == 0 {
		return
	}

	c.ui.Output("The following validation issues were detected:", terminal.WithHeaderStyle())

	for _, vr := range validationResults {
		if vr.Error != nil {
			c.ui.Output(vr.Error.Error(), terminal.WithErrorStyle())
		} else if vr.Warning != "" {
			c.ui.Output(vr.Warning, terminal.WithWarningStyle())
		}
	}

	c.ui.Output("")
}

// Init initializes the command by parsing flags, parsing the configuration,
// setting up the project, etc. You can control what is done by using the
// options.
//
// Init should be called FIRST within the Run function implementation. Many
// options will affect behavior of other functions that can be called later.
//
// In broad strokes, Init populates fields on the baseCommand by doing the following:
// - Parse flags
// - Parse input variables
// - Creates a project client
// - Triggers creation of the in-memory server (if necessary)
// - Starts a local runner (if necessary)
// - Attempts to find a waypoint.hcl config file, and parse it
// - Determines which project/apps are being targeted, by looking at
//   the -project and -app flags, the local config, the waypoint server.
func (c *baseCommand) Init(opts ...Option) error {
	baseCfg := baseConfig{}

	for _, opt := range c.globalOptions {
		opt(&baseCfg)
	}

	for _, opt := range opts {
		opt(&baseCfg)
	}

	// Set some basic internal fields
	c.noLocalServer = baseCfg.NoLocalServer

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

	if err := c.checkDeprecatedFlags(); err != nil {
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

	// load workspace from cli/env/storage
	workspace, err := c.workspace()
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return err
	}

	c.refWorkspace = &pb.Ref_Workspace{Workspace: workspace}

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

	// Now we begin parsing project and/or app values, when they are required
	// by the command.
	// The goal is to set c.refProject and c.refApps for the following options:
	// WithSingleAppTarget:
	//   - 1 app in []c.refApps
	//   - c.refProject set
	// WithMultiAppTargets:
	//   - 1 or more apps in []c.refApps
	//   - c.refProject set
	// WithProjectTarget:
	//   - c.refProject set
	//   - value of []c.refApps doesn't matter; likely will be set when using
	//     a local waypoint.hcl or if someone also includes -app flag, but not
	//     required

	// 1. Parse the configuration

	if !baseCfg.NoConfig {
		var vr config.ValidationResults

		// Try parsing config
		c.cfg, vr, err = c.initConfig("")

		c.showValidations(vr)

		if err != nil {
			return err
		}

		// If that worked, set our refs
		if c.cfg != nil {
			// Warn if the project from config and the project from flags conflict
			if c.flagProject != "" && c.flagProject != c.cfg.Project {
				c.ui.Output(warnProjectFlagMismatch, c.cfg.Project, c.flagProject, terminal.WithWarningStyle())

				// NOTE(izaak): unless we force remoteness, we may spawn a local runner which will operate against
				// the current config (which isn't relevant)
				c.Log.Debug("Forcing any future operations to occur remotely because the relevant waypoint.hcl is not present.")
				flagLocal := false
				c.flagLocal = &flagLocal
			} else {
				// This config is good - use it to obtain our refs.
				c.refProject = &pb.Ref_Project{Project: c.cfg.Project}
				for _, app := range c.cfg.Apps() {
					c.refApps = append(c.refApps, &pb.Ref_Application{
						Project:     c.cfg.Project,
						Application: app,
					})
				}
			}
		}
	}

	// 2. Parse project flags; overwrite any c.refProject value set from
	// config parsing, as precedence order means we take the most specific value
	// which is the -project flag
	if c.flagProject != "" {
		c.refProject = &pb.Ref_Project{Project: c.flagProject}
	}

	// 2.a. if c.refProject is nil at this point but we know it's required, we fail out
	if baseCfg.ProjectTargetRequired && c.refProject == nil {
		// The user must not have specified a project flag, and config parsing didn't produce one either.

		// NOTE(izaak) The UX here will be refined in the next pass - it's ok that this is terse for now.
		err := errors.New("No project specified, and no waypoint.hcl found.")
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return err
	}

	// 3. Create our client
	// We must do this after the project ref and vars are set, and we need
	// the client to find the project for the specified app targets.
	if !baseCfg.NoClient {
		c.project, err = c.initClient(nil)
		if err != nil {
			c.logError(c.Log, "failed to create client", err)
			return err
		}
	}

	// 4. Parse app flags; overwrite any []c.refApps value set from
	// config parsing, as precedence order means we take the most specific value
	// which is the -app flag
	if c.flagApp != "" {
		// NOTE: we could allow app to be specified multiple times in the future
		c.refApps = []*pb.Ref_Application{{Application: c.flagApp}}
	}

	// 4.a. If app(s) are required but not set, we do a final check to the
	// Waypoint server to see if it knows what apps belong to the project. We
	// set ProjectTargetRequired to `true` for both AppTarget options, so at this
	// point, if an AppTarget option is set then we should have a c.refProject.
	if (baseCfg.SingleAppTarget || baseCfg.MultiAppTarget) && len(c.refApps) == 0 {
		// We must not have found an app from config or flags, so we need to resort to the API.
		c.Log.Debug("No apps found via CLI or API - listing them from the CLI.")
		resp, err := c.project.Client().GetProject(c.Ctx, &pb.GetProjectRequest{Project: c.refProject})
		if err != nil {
			c.logError(c.Log, fmt.Sprintf("Failed to get project %s", c.refProject.Project), err)
			return err
		}
		for _, app := range resp.Project.Applications {
			c.refApps = append(c.refApps, &pb.Ref_Application{
				Application: app.Name,
				Project:     app.Project.Project,
			})
		}

		// This should be a very-edge case; we have done everything we possibly
		// can to find apps, and we don't have them
		if len(c.refApps) == 0 {
			err = fmt.Errorf("This command requires an app to be targeted, but no apps were found in project %q.", c.refProject.Project)
			c.logError(c.Log, "", err)
			return err
		}
	}

	// 4.b. Check to ensure there is 1 and only 1 target for SingleAppTarget cmd
	if baseCfg.SingleAppTarget && len(c.refApps) > 1 {
		c.ui.Output(errAppModeSingle, terminal.WithErrorStyle())
		return ErrSentinel
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
	// c.refApps is set in c.Init(), based on the project
	// the command is running against
	var apps []*clientpkg.App
	for _, refApp := range c.refApps {
		app := c.project.App(refApp.Application)
		c.Log.Debug("will operate on app", "name", refApp.Application)

		apps = append(apps, app)
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
			Aliases: []string{"a"},
			Default: "",
			Usage: "App to target. Certain commands require a single app target for " +
				"Waypoint configurations with multiple apps. If you have a single app, " +
				"then this can be ignored.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "project",
			Target:  &c.flagProject,
			Aliases: []string{"p"},
			Default: "",
			Usage:   "Project to target.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "workspace",
			Target:  &c.flagWorkspace,
			Aliases: []string{"w"},
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

		f.BoolPtrVar(&flag.BoolPtrVar{
			Name:   "remote",
			Target: &c.deprecatedFlagRemote,
			Hidden: true,
			Usage:  "True to use a remote runner to execute the operation.",
		})

		f.BoolPtrVar(&flag.BoolPtrVar{
			Name:   "local",
			Target: &c.flagLocal,
			Usage: "True to use a local runner to execute the operation, false to use a remote runner. \n" +
				"If unset, Waypoint will automatically determine where the operation will occur, \n" +
				"defaulting to remote if possible.",
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
// it's hard to tell positional vs intentional flags.
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

// workspace computes the workspace based on available values, in this order of
// precedence (last value wins):
//
// - value stored in the CLI context
// - value from the environment variable WAYPOINT_WORKSPACE
// - value set in the CLI flag -workspace
//
// The default value is "default"
func (c *baseCommand) workspace() (string, error) {
	// load env for workspace
	workspaceENV := os.Getenv(defaultWorkspaceEnvName)
	switch {
	case c.flagWorkspace != "":
		return c.flagWorkspace, nil
	case workspaceENV != "":
		return workspaceENV, nil
	default:
		// attempt to load from CLI context storage
		defaultName, err := c.contextStorage.Default()
		if err != nil {
			return "", err
		}

		// If we have no context name, then we just return the default
		if defaultName != "" && defaultName != "-" {
			// Load the context and return the workspace value. If it's empty,
			// we'll fall through and return the default
			cfg, err := c.contextStorage.Load(defaultName)
			if err != nil {
				return "", err
			}
			if cfg.Workspace != "" {
				return cfg.Workspace, nil
			}
		}
		// default value
		return defaultWorkspace, nil
	}
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

	// matches either "project" or "project/app"
	reAppTarget = regexp.MustCompile(`^(?P<project>[-0-9A-Za-z_]+)/(?P<app>[-0-9A-Za-z_]+)$`)

	errDeprecatedProjectAppArg = strings.TrimSpace(`
The project/app argument has been deprecated. Instead, use -project and -app flags, or their
short notation -p and -a.
`)

	snapshotUnimplementedErr = strings.TrimSpace(`
The current Waypoint server does not support snapshots. Rerunning the command
with '-snapshot=false' is required, and there will be no automatic data backups
for the server.
`)

	warnProjectFlagMismatch = strings.TrimSpace(`
Warning: Currently in project directory for %q, but will operate 
against specified project %q
`)
)
