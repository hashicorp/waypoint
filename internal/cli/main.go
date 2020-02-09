package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/cli"
)

const (
	// EnvLogLevel is the env var to set with the log level.
	EnvLogLevel = "DF_LOG_LEVEL"
)

// Main runs the CLI with the given arguments and returns the exit code.
// The arguments SHOULD include argv[0] as the program name.
func Main(args []string) int {
	// Initialize our logger based on env vars
	log, err := logger(args[0])
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

	return map[string]cli.CommandFactory{
		"up": func() (cli.Command, error) {
			return &UpCommand{
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
func logger(app string) (hclog.Logger, error) {
	level := hclog.Trace
	if v := os.Getenv(EnvLogLevel); v != "" {
		level = hclog.LevelFromString(v)
		if level == hclog.NoLevel {
			return nil, fmt.Errorf("%s value %q is not a valid log level", EnvLogLevel, v)
		}
	}

	return hclog.New(&hclog.LoggerOptions{
		Name:   app,
		Level:  level,
		Color:  hclog.AutoColor,
		Output: os.Stderr,
	}), nil
}
