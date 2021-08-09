package cli

//go:generate go-bindata -nomemcopy -nometadata -pkg datagen -o datagen/datagen.go -prefix data/ data/...

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/mitchellh/cli"
	"github.com/mitchellh/go-glint"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/env"
	"github.com/hashicorp/waypoint/internal/pkg/signalcontext"
	"github.com/hashicorp/waypoint/internal/version"
)

const (
	// EnvLogLevel is the env var to set with the log level.
	EnvLogLevel = "WAYPOINT_LOG_LEVEL"

	// EnvPlain is the env var that can be set to force plain output mode.
	EnvPlain = "WAYPOINT_PLAIN"
)

var (
	// cliName is the name of this CLI.
	cliName = "waypoint"

	// commonCommands are the commands that are deemed "common" and shown first
	// in the CLI help output.
	commonCommands = []string{
		"login",
		"build",
		"deploy",
		"release",
		"status",
		"up",
	}

	// hiddenCommands are not shown in CLI help output.
	hiddenCommands = map[string]struct{}{
		"plugin": {},

		// Deprecated:
		"token": {}, // replaced by "user"
	}

	ExposeDocs bool
)

// Main runs the CLI with the given arguments and returns the exit code.
// The arguments SHOULD include argv[0] as the program name.
func Main(args []string) int {
	// Clean up all our plugins so we don't leave any dangling processes.
	// Note that this is a "just in case" catch. We should be properly cleaning
	// up plugin processes by calling Close on all the resources we use.
	defer plugin.CleanupClients()

	// NOTE: This is only for running `waypoint -v` and expecting it to return
	// a version. Any other subcommand will expect `-v` to be around verbose
	// logging rather than printing a version
	if len(args) == 2 && args[1] == "-v" {
		args[1] = "-version"
	}

	// Initialize our logger based on env vars
	args, log, logOutput, err := logger(args)
	if err != nil {
		panic(err)
	}

	// Log our versions
	vsn := version.GetVersion()
	log.Info("waypoint version",
		"full_string", vsn.FullVersionNumber(true),
		"version", vsn.Version,
		"prerelease", vsn.VersionPrerelease,
		"metadata", vsn.VersionMetadata,
		"revision", vsn.Revision,
	)

	// Build our cancellation context
	ctx, closer := signalcontext.WithInterrupt(context.Background(), log)
	defer closer()

	// Get our base command
	base, commands := Commands(ctx, log, logOutput)
	defer base.Close()

	// Build the CLI. We use a CLI factory function because to modify the
	// args once you call a func on CLI you need to create a new CLI instance.
	cliFactory := func() *cli.CLI {
		return &cli.CLI{
			Name:                       args[0],
			Args:                       args[1:],
			Version:                    vsn.FullVersionNumber(true),
			Commands:                   commands,
			Autocomplete:               true,
			AutocompleteNoDefaultFlags: true,
			HelpFunc:                   GroupedHelpFunc(cli.BasicHelpFunc(cliName)),
		}
	}

	// Copy the CLI to check if it is a version call. If so, we modify
	// the args to just be the version subcommand. This ensures that
	// --version behaves by calling `waypoint version` and we get consistent
	// behavior.
	cli := cliFactory()
	if cli.IsVersion() {
		// We need to reinit because you can't modify fields after calling funcs
		cli = cliFactory()
		cli.Args = []string{"version"}
	}

	// Run the CLI
	exitCode, err := cli.Run()
	if err != nil {
		panic(err)
	}

	return exitCode
}

// commands returns the map of commands that can be used to initialize a CLI.
func Commands(
	ctx context.Context,
	log hclog.Logger,
	logOutput io.Writer,
	opts ...Option,
) (*baseCommand, map[string]cli.CommandFactory) {
	baseCommand := &baseCommand{
		Ctx:           ctx,
		Log:           log,
		LogOutput:     logOutput,
		globalOptions: opts,
	}

	// Set plain mode if set
	outputModeBool, err := env.GetBool(EnvPlain, false)
	if err != nil {
		log.Warn(err.Error())
	}
	if outputModeBool {
		baseCommand.globalOptions = append(baseCommand.globalOptions,
			WithUI(terminal.NonInteractiveUI(ctx)))
	}

	// aliases is a list of command aliases we have. The key is the CLI
	// command (the alias) and the value is the existing target command.
	aliases := map[string]string{
		"build":   "artifact build",
		"deploy":  "deployment deploy",
		"install": "server install",
	}

	// start building our commands
	commands := map[string]cli.CommandFactory{
		"login": func() (cli.Command, error) {
			return &LoginCommand{
				baseCommand: baseCommand,
			}, nil
		},

		"init": func() (cli.Command, error) {
			return &InitCommand{
				baseCommand: baseCommand,
			}, nil
		},

		"up": func() (cli.Command, error) {
			return &UpCommand{
				baseCommand: baseCommand,
			}, nil
		},

		"destroy": func() (cli.Command, error) {
			return &DestroyCommand{
				baseCommand: baseCommand,
			}, nil
		},

		"exec": func() (cli.Command, error) {
			return &ExecCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"config": func() (cli.Command, error) {
			return &helpCommand{
				SynopsisText: helpText["config"][0],
				HelpText:     helpText["config"][1],
			}, nil
		},
		"config get": func() (cli.Command, error) {
			return &ConfigGetCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"config set": func() (cli.Command, error) {
			return &ConfigSetCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"config source-get": func() (cli.Command, error) {
			return &ConfigSourceGetCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"config source-set": func() (cli.Command, error) {
			return &ConfigSourceSetCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"config sync": func() (cli.Command, error) {
			return &ConfigSyncCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"logs": func() (cli.Command, error) {
			return &LogsCommand{
				baseCommand: baseCommand,
			}, nil
		},

		"artifact": func() (cli.Command, error) {
			return &helpCommand{
				SynopsisText: helpText["artifact"][0],
				HelpText:     helpText["artifact"][1],
			}, nil
		},
		"artifact build": func() (cli.Command, error) {
			return &ArtifactBuildCommand{
				baseCommand: baseCommand,
			}, nil
		},

		"artifact list": func() (cli.Command, error) {
			return &ArtifactListCommand{
				baseCommand: baseCommand,
			}, nil
		},

		"artifact list-builds": func() (cli.Command, error) {
			return &BuildListCommand{
				baseCommand: baseCommand,
			}, nil
		},

		"artifact push": func() (cli.Command, error) {
			return &ArtifactPushCommand{
				baseCommand: baseCommand,
			}, nil
		},

		"deployment": func() (cli.Command, error) {
			return &helpCommand{
				SynopsisText: helpText["deployment"][0],
				HelpText:     helpText["deployment"][1],
			}, nil
		},

		"deployment deploy": func() (cli.Command, error) {
			return &DeploymentCreateCommand{
				baseCommand: baseCommand,
			}, nil
		},

		"deployment destroy": func() (cli.Command, error) {
			return &DeploymentDestroyCommand{
				baseCommand: baseCommand,
			}, nil
		},

		"deployment list": func() (cli.Command, error) {
			return &DeploymentListCommand{
				baseCommand: baseCommand,
			}, nil
		},

		"release": func() (cli.Command, error) {
			return &ReleaseCreateCommand{
				baseCommand: baseCommand,
			}, nil
		},

		"server": func() (cli.Command, error) {
			return &helpCommand{
				SynopsisText: helpText["server"][0],
				HelpText:     helpText["server"][1],
			}, nil
		},
		"server bootstrap": func() (cli.Command, error) {
			return &ServerBootstrapCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"server install": func() (cli.Command, error) {
			return &InstallCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"server uninstall": func() (cli.Command, error) {
			return &UninstallCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"server run": func() (cli.Command, error) {
			return &ServerRunCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"server config-set": func() (cli.Command, error) {
			return &ServerConfigSetCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"server snapshot": func() (cli.Command, error) {
			return &SnapshotBackupCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"server restore": func() (cli.Command, error) {
			return &SnapshotRestoreCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"server upgrade": func() (cli.Command, error) {
			return &ServerUpgradeCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"status": func() (cli.Command, error) {
			return &StatusCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"plugin": func() (cli.Command, error) {
			return &PluginCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"version": func() (cli.Command, error) {
			return &VersionCommand{
				baseCommand: baseCommand,
				VersionInfo: version.GetVersion(),
			}, nil
		},

		"hostname": func() (cli.Command, error) {
			return &helpCommand{
				SynopsisText: helpText["hostname"][0],
				HelpText:     helpText["hostname"][1],
			}, nil
		},

		"hostname register": func() (cli.Command, error) {
			return &HostnameRegisterCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"hostname list": func() (cli.Command, error) {
			return &HostnameListCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"hostname delete": func() (cli.Command, error) {
			return &HostnameDeleteCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"token": func() (cli.Command, error) {
			return &helpCommand{
				SynopsisText: helpText["token"][0],
				HelpText:     helpText["token"][1],
			}, nil
		},
		"token new": func() (cli.Command, error) {
			return &GetTokenCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"token invite": func() (cli.Command, error) {
			return &GetInviteCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"token exchange": func() (cli.Command, error) {
			return &ExchangeInviteCommand{
				baseCommand: baseCommand,
			}, nil
		},

		"runner": func() (cli.Command, error) {
			return &helpCommand{
				SynopsisText: helpText["runner"][0],
				HelpText:     helpText["runner"][1],
			}, nil
		},
		"runner agent": func() (cli.Command, error) {
			return &RunnerAgentCommand{
				baseCommand: baseCommand,
			}, nil
		},

		"context": func() (cli.Command, error) {
			return &ContextHelpCommand{
				baseCommand:  baseCommand,
				SynopsisText: helpText["context"][0],
				HelpText:     helpText["context"][1],
			}, nil
		},
		"context inspect": func() (cli.Command, error) {
			return &ContextInspectCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"context create": func() (cli.Command, error) {
			return &ContextCreateCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"context delete": func() (cli.Command, error) {
			return &ContextDeleteCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"context use": func() (cli.Command, error) {
			return &ContextUseCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"context clear": func() (cli.Command, error) {
			return &ContextClearCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"context rename": func() (cli.Command, error) {
			return &ContextRenameCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"context list": func() (cli.Command, error) {
			return &ContextListCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"context verify": func() (cli.Command, error) {
			return &ContextVerifyCommand{
				baseCommand: baseCommand,
			}, nil
		},

		"ui": func() (cli.Command, error) {
			return &UICommand{
				baseCommand: baseCommand,
			}, nil
		},

		"docs": func() (cli.Command, error) {
			return &AppDocsCommand{
				baseCommand: baseCommand,
			}, nil
		},

		"project": func() (cli.Command, error) {
			return &helpCommand{
				SynopsisText: helpText["project"][0],
				HelpText:     helpText["project"][1],
			}, nil
		},
		"project list": func() (cli.Command, error) {
			return &ProjectListCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"project apply": func() (cli.Command, error) {
			return &ProjectApplyCommand{
				baseCommand: baseCommand,
			}, nil
		},

		"fmt": func() (cli.Command, error) {
			return &FmtCommand{
				baseCommand: baseCommand,
			}, nil
		},

		"auth-method": func() (cli.Command, error) {
			return &helpCommand{
				SynopsisText: helpText["auth-method"][0],
				HelpText:     helpText["auth-method"][1],
			}, nil
		},

		"auth-method set": func() (cli.Command, error) {
			return &helpCommand{
				SynopsisText: helpText["auth-method-set"][0],
				HelpText:     helpText["auth-method-set"][1],
			}, nil
		},

		"auth-method set oidc": func() (cli.Command, error) {
			return &AuthMethodSetOIDCCommand{
				baseCommand: baseCommand,
			}, nil
		},

		"auth-method inspect": func() (cli.Command, error) {
			return &AuthMethodInspectCommand{
				baseCommand: baseCommand,
			}, nil
		},

		"auth-method delete": func() (cli.Command, error) {
			return &AuthMethodDeleteCommand{
				baseCommand: baseCommand,
			}, nil
		},

		"auth-method list": func() (cli.Command, error) {
			return &AuthMethodListCommand{
				baseCommand: baseCommand,
			}, nil
		},

		"user": func() (cli.Command, error) {
			return &helpCommand{
				SynopsisText: helpText["user"][0],
				HelpText:     helpText["user"][1],
			}, nil
		},

		"user inspect": func() (cli.Command, error) {
			return &UserInspectCommand{
				baseCommand: baseCommand,
			}, nil
		},

		"user modify": func() (cli.Command, error) {
			return &UserModifyCommand{
				baseCommand: baseCommand,
			}, nil
		},

		"user invite": func() (cli.Command, error) {
			return &UserInviteCommand{
				baseCommand: baseCommand,
			}, nil
		},

		"user token": func() (cli.Command, error) {
			return &UserTokenCommand{
				baseCommand: baseCommand,
			}, nil
		},
	}

	// register our aliases
	for from, to := range aliases {
		commands[from] = commands[to]
	}

	if ExposeDocs {
		commands["cli-docs"] = func() (cli.Command, error) {
			return &DocsCommand{
				baseCommand: baseCommand,
				commands:    commands,
				aliases:     aliases,
			}, nil
		}
	}

	return baseCommand, commands
}

// logger returns the logger to use for the CLI. Output, level, etc. are
// determined based on environment variables if set.
func logger(args []string) ([]string, hclog.Logger, io.Writer, error) {
	app := args[0]

	// Determine our log level if we have any. First override we check if env var
	level := hclog.NoLevel
	if v := os.Getenv(EnvLogLevel); v != "" {
		level = hclog.LevelFromString(v)
		if level == hclog.NoLevel {
			return nil, nil, nil, fmt.Errorf("%s value %q is not a valid log level", EnvLogLevel, v)
		}
	}

	// Process arguments looking for `-v` flags to control the log level.
	// This overrides whatever the env var set.
	var outArgs []string
	for _, arg := range args {
		if len(arg) != 0 && arg[0] != '-' {
			outArgs = append(outArgs, arg)
			continue
		}

		switch arg {
		case "-v":
			if level == hclog.NoLevel || level > hclog.Info {
				level = hclog.Info
			}
		case "-vv":
			if level == hclog.NoLevel || level > hclog.Debug {
				level = hclog.Debug
			}
		case "-vvv":
			if level == hclog.NoLevel || level > hclog.Trace {
				level = hclog.Trace
			}
		default:
			outArgs = append(outArgs, arg)
		}
	}

	// Default output is nowhere unless we enable logging.
	var output io.Writer = ioutil.Discard
	color := hclog.ColorOff
	if level != hclog.NoLevel {
		output = os.Stderr
		color = hclog.AutoColor
	}

	logger := hclog.New(&hclog.LoggerOptions{
		Name:   app,
		Level:  level,
		Color:  color,
		Output: output,
	})

	return outArgs, logger, output, nil
}

func GroupedHelpFunc(f cli.HelpFunc) cli.HelpFunc {
	return func(commands map[string]cli.CommandFactory) string {
		var buf bytes.Buffer
		d := glint.New()
		d.SetRenderer(&glint.TerminalRenderer{
			Output: &buf,

			// We set rows/cols here manually. The important bit is the cols
			// needs to be wide enough so glint doesn't clamp any text and
			// lets the terminal just autowrap it. Rows doesn't make a big
			// difference.
			Rows: 10,
			Cols: 180,
		})

		// Header
		d.Append(glint.Style(
			glint.Text("Welcome to Waypoint"),
			glint.Bold(),
		))
		d.Append(glint.Layout(
			glint.Style(
				glint.Text("Docs:"),
				glint.Color("lightBlue"),
			),
			glint.Text(" "),
			glint.Text("https://waypointproject.io"),
		).Row())
		d.Append(glint.Layout(
			glint.Style(
				glint.Text("Version:"),
				glint.Color("green"),
			),
			glint.Text(" "),
			glint.Text(version.GetVersion().VersionNumber()),
		).Row())
		d.Append(glint.Text(""))

		// Usage
		d.Append(glint.Layout(
			glint.Style(
				glint.Text("Usage:"),
				glint.Color("lightMagenta"),
			),
			glint.Text(" "),
			glint.Text(cliName),
			glint.Text(" "),
			glint.Text("[-version] [-help] [-autocomplete-(un)install] <command> [args]"),
		).Row())
		d.Append(glint.Text(""))

		// Add common commands
		helpCommandsSection(d, "Common commands", commonCommands, commands)

		// Make our other commands
		ignoreMap := map[string]struct{}{}
		for k := range hiddenCommands {
			ignoreMap[k] = struct{}{}
		}
		for _, k := range commonCommands {
			ignoreMap[k] = struct{}{}
		}

		var otherCommands []string
		for k := range commands {
			if _, ok := ignoreMap[k]; ok {
				continue
			}

			otherCommands = append(otherCommands, k)
		}
		sort.Strings(otherCommands)

		// Add other commands
		helpCommandsSection(d, "Other commands", otherCommands, commands)

		d.RenderFrame()
		return buf.String()
	}
}

func helpCommandsSection(
	d *glint.Document,
	header string,
	commands []string,
	factories map[string]cli.CommandFactory,
) {
	// Header
	d.Append(glint.Style(
		glint.Text(header),
		glint.Bold(),
	))

	// Build our commands
	var b bytes.Buffer
	tw := tabwriter.NewWriter(&b, 0, 2, 6, ' ', 0)
	for _, k := range commands {
		fn, ok := factories[k]
		if !ok {
			continue
		}

		cmd, err := fn()
		if err != nil {
			panic(fmt.Sprintf("failed to load %q command: %s", k, err))
		}

		fmt.Fprintf(tw, "%s\t%s\n", k, cmd.Synopsis())
	}
	tw.Flush()

	d.Append(glint.Layout(
		glint.Text(b.String()),
	).PaddingLeft(2))
}

var helpText = map[string][2]string{
	"artifact": {
		"Artifact and build management",
		`
Artifact and build management.

The artifact commands can be used to create new builds, push those
builds to a registry, etc. The result of a build is known as an artifact.
Waypoint will search for artifacts to pass to the deployment phase.
`,
	},

	"auth-method": {
		"Auth Method Management",
		`
Auth Method Management

The auth-method commands can be used to manage how users can authenticate
into the Waypoint server. For day-to-day Waypoint users, you likely want
to use the "waypoint login" command or "waypoint user" commands. The
auth-method subcommand is primarily aimed at Waypoint server operators.
`,
	},

	"auth-method-set": {
		"Create or update an auth method",
		`
Create or update an auth method.

This command can be used to configure a new auth method or update
an existing auth method. Use the specific auth-method type subcommand.
Once the auth method is created it is immediately available for use by
end users.
`,
	},

	"config": {
		"Application configuration management",
		`
Manage application configuration.

The config commands can be used to manage the configuration that
Waypoint will inject into your application via environment variables.
This can be used to set values such as ports to listen on, database URLs,
etc.

For more information see: https://waypointproject.io/docs/app-config
`,
	},

	"context": {
		"Server access configurations",
		`
Manage configurations for accessing Waypoint servers.

A context contains all the configuration to connect to a single Waypoint
server. The Waypoint CLI can have multiple contexts to make it easy to switch
between different Waypoint servers.
`,
	},

	"deployment": {
		"Deployment creation and management",
		`
Create and manage application deployments.

A deployment is the process of taking an artifact and launching it,
potentially for public access. Waypoint deployment commands let you create
new deployments, delete existing ones, list previous deployments, and more.
`,
	},

	"hostname": {
		"Application URLs",
		`
Create and manage application URLs powered by the Waypoint URL service.

The Waypoint URL service registers publicly routable URLs to access your
deployments. These can be used to share previews with teammates, see
unreleased deployments, and more.

For more information see: https://waypointproject.io/docs/url
`,
	},

	"project": {
		"Project management",
		`
Project management.

Projects are comprised of one or more applications. A project maps
to a VCS repository (if one exists).
`,
	},

	"runner": {
		"Runner management",
		`
Runner management.

Runners are used to execute remote operations for Waypoint. If you're new
to Waypoint, you generally *do not need* runners and you can ignore this
entire section.
`,
	},

	"server": {
		"Server management",
		`
Server management.

The CLI, UI, and entrypoints all communicate to a Waypoint server. A
Waypoint server is required for logs, exec, config, and more to work.
The recommended way to run a server is "waypoint install".

This command contains further subcommands to work with servers.
`,
	},

	"token": {
		"Authenticate and invite collaborators",
		`
Authenticate and invite collaborators.

Tokens are the primary form of authentication to Waypoint. Everyone who
accesses a Waypoint server requires a token.
`,
	},

	"user": {
		"User information and management",
		`
View, manage, and invite users.

Everyone who uses Waypoint is represented as a Waypoint user, including
token authentication. This subcommand can be used to inspect information
about the currently logged in user, generate new access, and invite new
users directly into the Waypoint server.

If you are looking to log in to Waypoint, use "waypoint login".
`,
	},
}
