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
	"strings"
	"text/tabwriter"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/mitchellh/cli"

	"github.com/hashicorp/waypoint/internal/pkg/signalcontext"
	"github.com/hashicorp/waypoint/internal/version"
)

const (
	// EnvLogLevel is the env var to set with the log level.
	EnvLogLevel = "WAYPOINT_LOG_LEVEL"
)

var (
	// cliName is the name of this CLI.
	cliName = "waypoint"

	// commonCommands are the commands that are deemed "common" and shown first
	// in the CLI help output.
	commonCommands = map[string]struct{}{
		"build":   {},
		"deploy":  {},
		"release": {},
		"up":      {},
	}

	// hiddenCommands are not shown in CLI help output.
	hiddenCommands = map[string]struct{}{
		"plugin": {},
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

	// Initialize our logger based on env vars
	args, log, logOutput, err := logger(args)
	if err != nil {
		panic(err)
	}

	// Build our cancellation context
	ctx, closer := signalcontext.WithInterrupt(context.Background(), log)
	defer closer()

	// Build the CLI
	cli := &cli.CLI{
		Name:                       args[0],
		Args:                       args[1:],
		Commands:                   Commands(ctx, log, logOutput),
		Autocomplete:               true,
		AutocompleteNoDefaultFlags: true,
		HelpFunc:                   GroupedHelpFunc(cli.BasicHelpFunc(cliName)),
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
) map[string]cli.CommandFactory {
	baseCommand := &baseCommand{
		Ctx:           ctx,
		Log:           log,
		LogOutput:     logOutput,
		globalOptions: opts,
	}

	// aliases is a list of command aliases we have. The key is the CLI
	// command (the alias) and the value is the existing target command.
	aliases := map[string]string{
		"build":  "artifact build",
		"deploy": "deployment deploy",
		"push":   "artifact push",
	}

	// start building our commands
	commands := map[string]cli.CommandFactory{
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
		"logs": func() (cli.Command, error) {
			return &LogsCommand{
				baseCommand: baseCommand,
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

		"install": func() (cli.Command, error) {
			return &InstallCommand{
				baseCommand: baseCommand,
			}, nil
		},

		"release": func() (cli.Command, error) {
			return &ReleaseCreateCommand{
				baseCommand: baseCommand,
			}, nil
		},

		"server": func() (cli.Command, error) {
			return &ServerCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"server config-set": func() (cli.Command, error) {
			return &ServerConfigSetCommand{
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

		"runner agent": func() (cli.Command, error) {
			return &RunnerAgentCommand{
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

	return commands
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
		if arg[0] != '-' {
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
		var b bytes.Buffer
		tw := tabwriter.NewWriter(&b, 0, 2, 6, ' ', 0)

		fmt.Fprintf(tw, "Usage: %s [-help] [-autocomplete-(un)install] <command> [args]\n\n", cliName)
		fmt.Fprintf(tw, "Common commands:\n")
		for k := range commonCommands {
			printCommand(tw, k, commands[k])
		}

		// Filter out common commands and aliased commands from the other
		// commands output
		otherCommands := make([]string, 0, len(commands))
		for k := range commands {
			if _, ok := commonCommands[k]; ok {
				continue
			}
			if _, ok := hiddenCommands[k]; ok {
				continue
			}

			otherCommands = append(otherCommands, k)
		}
		sort.Strings(otherCommands)

		fmt.Fprintf(tw, "\n")
		fmt.Fprintf(tw, "Other commands:\n")
		for _, v := range otherCommands {
			printCommand(tw, v, commands[v])
		}

		tw.Flush()

		return strings.TrimSpace(b.String())
	}
}

func printCommand(w io.Writer, name string, cmdFn cli.CommandFactory) {
	cmd, err := cmdFn()
	if err != nil {
		panic(fmt.Sprintf("failed to load %q command: %s", name, err))
	}
	fmt.Fprintf(w, "    %s\t%s\n", name, cmd.Synopsis())
}
