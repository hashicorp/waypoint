package cli

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/mitchellh/cli"
	"github.com/mitchellh/devflow/internal/pkg/status"
)

const (
	// EnvLogLevel is the env var to set with the log level.
	EnvLogLevel = "DF_LOG_LEVEL"
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
	args, log, stat, err := logger(args)
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
		Commands:                   commands(ctx, log, stat),
		Autocomplete:               true,
		AutocompleteNoDefaultFlags: true,
	}

	// Run the CLI
	exitCode, err := cli.Run()
	if err != nil {
		panic(err)
	}

	return exitCode
}

// commands returns the map of commands that can be used to initialize a CLI.
func commands(ctx context.Context, log hclog.Logger, stat status.Updater) map[string]cli.CommandFactory {
	baseCommand := &baseCommand{
		Ctx:     ctx,
		Log:     log,
		Updater: stat,
	}

	return map[string]cli.CommandFactory{
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

		// TODO(mitchellh): make hidden
		"plugin": func() (cli.Command, error) {
			return &PluginCommand{
				baseCommand: baseCommand,
			}, nil
		},
	}
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
func logger(args []string) ([]string, hclog.Logger, status.Updater, error) {
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

	var update status.Updater
	if level <= hclog.Warn {
		update = &status.SpinnerStatus{}
	} else {
		update = &status.HCLog{
			L:     logger,
			Level: level,
		}
	}

	return outArgs, logger, update, nil
}
