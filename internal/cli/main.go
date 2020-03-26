package cli

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/mitchellh/cli"
)

const (
	// EnvLogLevel is the env var to set with the log level.
	EnvLogLevel = "DF_LOG_LEVEL"
)

var (
	// cliName is the name of this CLI.
	cliName = "devflow"

	// commonCommands are the commands that are deemed "common" and shown first
	// in the CLI help output.
	commonCommands = map[string]struct{}{
		"build": struct{}{},
		"push":  struct{}{},
		"up":    struct{}{},
	}

	// hiddenCommands are not shown in CLI help output.
	hiddenCommands = map[string]struct{}{
		"plugin": struct{}{},
	}
)

// Main runs the CLI with the given arguments and returns the exit code.
// The arguments SHOULD include argv[0] as the program name.
func Main(args []string) int {
	// Clean up all our plugins so we don't leave any dangling processes.
	// TODO(mitchellh): we should always keep this call just in case but
	// what we really want to do is implement io.Closer and have our
	// `internal/core` structures call that as necessary when they're done
	// with components.
	defer plugin.CleanupClients()

	// Initialize our logger based on env vars
	args, log, err := logger(args)
	if err != nil {
		panic(err)
	}

	// Build our cancellation context
	ctx, closer := interruptContext(context.Background(), log)
	defer closer()

	// Build the CLI
	cli := &cli.CLI{
		Name:                       args[0],
		Args:                       args[1:],
		Commands:                   commands(ctx, log),
		Autocomplete:               true,
		AutocompleteNoDefaultFlags: true,
		HelpFunc:                   groupedHelpFunc(cli.BasicHelpFunc(cliName)),
	}

	// Run the CLI
	exitCode, err := cli.Run()
	if err != nil {
		panic(err)
	}

	return exitCode
}

// commands returns the map of commands that can be used to initialize a CLI.
func commands(ctx context.Context, log hclog.Logger) map[string]cli.CommandFactory {
	baseCommand := &baseCommand{
		Ctx: ctx,
		Log: log,
	}

	// aliases is a list of command aliases we have. The key is the CLI
	// command (the alias) and the value is the existing target command.
	aliases := map[string]string{
		"build": "artifact build",
		"push":  "artifact push",
	}

	// start building our commands
	commands := map[string]cli.CommandFactory{
		"up": func() (cli.Command, error) {
			return &UpCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"exec": func() (cli.Command, error) {
			return &ExecCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"config-get": func() (cli.Command, error) {
			return &ConfigGetCommand{
				baseCommand: baseCommand,
			}, nil
		},
		"config-set": func() (cli.Command, error) {
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
		"plugin": func() (cli.Command, error) {
			return &PluginCommand{
				baseCommand: baseCommand,
			}, nil
		},
	}

	// register our aliases
	for from, to := range aliases {
		commands[from] = commands[to]
	}

	return commands
}

// interruptContext returns a Context that is done when an interrupt
// signal is received. It also returns a closer function that should be
// deferred for proper cleanup.
func interruptContext(ctx context.Context, log hclog.Logger) (context.Context, func()) {
	log.Trace("starting interrupt listener for context cancellation")

	// Create the cancellable context that we'll use when we receive an interrupt
	ctx, cancel := context.WithCancel(ctx)

	// Create the signal channel and cancel the context when we get a signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		log.Trace("interrupt listener goroutine started")

		select {
		case <-ch:
			log.Warn("interrupt received, cancelling context")
			cancel()
		case <-ctx.Done():
			log.Warn("context cancelled, stopping interrupt listener loop")
			return
		}
	}()

	// Return the context and a closer that cancels the context and also
	// stops any signals from coming to our channel.
	return ctx, func() {
		log.Trace("stopping signal listeners and cancelling the context")
		signal.Stop(ch)
		cancel()
	}
}

// logger returns the logger to use for the CLI. Output, level, etc. are
// determined based on environment variables if set.
func logger(args []string) ([]string, hclog.Logger, error) {
	app := args[0]

	// Determine our log level if we have any. First override we check if env var
	level := hclog.NoLevel
	if v := os.Getenv(EnvLogLevel); v != "" {
		level = hclog.LevelFromString(v)
		if level == hclog.NoLevel {
			return nil, nil, fmt.Errorf("%s value %q is not a valid log level", EnvLogLevel, v)
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

	return outArgs, logger, nil
}

func groupedHelpFunc(f cli.HelpFunc) cli.HelpFunc {
	return func(commands map[string]cli.CommandFactory) string {
		var b bytes.Buffer
		tw := tabwriter.NewWriter(&b, 0, 2, 6, ' ', 0)

		fmt.Fprintf(tw, "Usage: %s [-version] [-help] [-autocomplete-(un)install] <command> [args]\n\n", cliName)
		fmt.Fprintf(tw, "Common commands:\n")
		for k, _ := range commonCommands {
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
