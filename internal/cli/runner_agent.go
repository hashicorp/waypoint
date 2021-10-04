package cli

import (
	"context"
	"errors"
	"io/ioutil"
	"net"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/posener/complete"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/plugin"
	runnerpkg "github.com/hashicorp/waypoint/internal/runner"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/serverclient"
)

type RunnerAgentCommand struct {
	*baseCommand

	// Indicates if exec plugins run by this runner should read dynamic
	// config. This requires the runner to have credentials to the dynamic
	// config sources.
	flagDynConfig bool

	// Specifies an address to setup a noop TCP server on that can be
	// used for liveness probes.
	flagLivenessTCPAddr string

	// The specific ID that the runner should use. When not set, the runner
	// generates an random.
	flagId string

	// This indicates that the runner is an on-demand runner. This information
	// is made available to the plugins so they can alter their behavior for
	// this unique context.
	flagODR bool
}

// This is how long a runner in ODR mode will wait for its job assignment before
// timing out.
var defaultRunnerODRAcceptTimeout = 60 * time.Second

func (c *RunnerAgentCommand) Run(args []string) int {
	defer c.Close()
	ctx := c.Ctx
	log := c.Log.Named("runner").Named("agent")

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
		WithNoAutoServer(),
	); err != nil {
		return 1
	}

	plugin.InsideODR = c.flagODR

	// Connect to the server
	log.Info("sourcing credentials and connecting to the Waypoint server")
	conn, err := serverclient.Connect(ctx,
		serverclient.FromContext(c.contextStorage, ""),
		serverclient.FromEnv(),
	)
	if err != nil {
		c.ui.Output(
			"Error connecting to the Waypoint server: %s", err.Error(),
			terminal.WithErrorStyle(),
		)
		return 1
	}
	client := pb.NewWaypointClient(conn)

	// Output information to the user
	c.ui.Output("Runner configuration:", terminal.WithHeaderStyle())
	c.ui.NamedValues([]terminal.NamedValue{
		{Name: "Server address", Value: conn.Target()},
	})
	c.ui.Output("Runner logs:", terminal.WithHeaderStyle())
	if c.flagODR {
		c.ui.Output("Operating as an On-demand Runner")
	} else {
		c.ui.Output("Operating as a static Runner")
	}

	c.ui.Output("")

	// Set our log output higher if it's not already so that it begins showing.
	if !log.IsInfo() {
		log.SetLevel(hclog.Info)
	}

	// If our output is to discard, then we want to redirect the output
	// to the console. We should be able to do this as long as our logger
	// supports the OutputResettable interface.
	if c.LogOutput == ioutil.Discard {
		if lr, ok := log.(hclog.OutputResettable); ok {
			output, _, err := c.ui.OutputWriters()
			if err != nil {
				c.ui.Output(
					"Error setting up logger: %s", err.Error(),
					terminal.WithErrorStyle(),
				)
				return 1
			}

			lr.ResetOutput(&hclog.LoggerOptions{
				Output: output,
				Color:  hclog.AutoColor,
			})
		}
	}

	// Create our runner
	log.Info("initializing the runner")

	options := []runnerpkg.Option{
		runnerpkg.WithClient(client),
		runnerpkg.WithLogger(log.Named("runner")),
		runnerpkg.WithDynamicConfig(c.flagDynConfig),
	}

	if c.flagId != "" {
		options = append(options, runnerpkg.WithId(c.flagId))
	}

	if c.flagODR {
		options = append(options,
			runnerpkg.WithODR(),
			runnerpkg.ByIdOnly(),
			runnerpkg.WithAcceptTimeout(defaultRunnerODRAcceptTimeout),
		)
	}

	runner, err := runnerpkg.New(options...)
	if err != nil {
		c.ui.Output(
			"Error initializing the runner: %s", err.Error(),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	// Start the runner
	log.Info("starting runner", "id", runner.Id())
	if err := runner.Start(); err != nil {
		log.Error("error starting runner", "err", err)
		return 1
	}

	// If we have a liveness address setup, start the liveness server.
	if addr := c.flagLivenessTCPAddr; addr != "" {
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			c.ui.Output(
				"Error starting liveness server: %s", err.Error(),
				terminal.WithErrorStyle(),
			)
			return 1
		}
		defer ln.Close()

		go func() {
			for {
				conn, err := ln.Accept()
				if err != nil {
					log.Warn("error accepting liveness connection", "err", err)
					if errors.Is(err, net.ErrClosed) {
						log.Warn("liveness server exiting")
						return
					}

					continue
				}

				// Immediately close. The liveness check only ensures a
				// connection can be established, we close immediately
				// thereafter.
				conn.Close()
			}
		}()
	}

	// Accept jobs in goroutine so that we can interrupt it.
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		defer cancel()

		for {
			err := runner.Accept(ctx)
			if err == nil && c.flagODR {
				log.Debug("handled our one job in ODR mode, exiting")
				return
			}

			if err != nil {
				if err == runnerpkg.ErrClosed {
					return
				}

				if err == runnerpkg.ErrTimeout {
					log.Error("timed out waiting for a job")
					return
				}

				log.Error("error running job", "err", err)

				switch status.Code(err) {
				case codes.NotFound:
					// This error code means that the runner is deregistered.
					// There is no recover from this and we have to restart
					// the runner.
					log.Error("runner unexpectedly deregistered, exiting")
					return

				case codes.Unavailable:
					// Server became unavailable. We retry on this after
					// a short sleep to allow the server to come back online.
					log.Warn("server unavailable, sleeping before retry")
					time.Sleep(2 * time.Second)
				}
			}
		}
	}()

	// Wait for end. This ends either via an interrupt (parent context)
	// or via the runner accept loop erroring in some way.
	<-ctx.Done()

	// Gracefully close
	log.Info("quit request received, gracefully stopping runner")
	if err := runner.Close(); err != nil {
		log.Error("error stopping runner", "err", err)
		return 1
	}

	return 0
}

func (c *RunnerAgentCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:   "enable-dynamic-config",
			Target: &c.flagDynConfig,
			Usage:  "Allow dynamic config to be created when an exec plugin is used.",
		})

		f.StringVar(&flag.StringVar{
			Name:   "liveness-tcp-addr",
			Target: &c.flagLivenessTCPAddr,
			Usage: "If this is set, the runner will open a TCP listener on this " +
				"address when it is running. This can be used as a liveness probe " +
				"endpoint. The TCP server serves no other purpose.",
		})

		f.StringVar(&flag.StringVar{
			Name:   "id",
			Target: &c.flagId,
			Usage:  "If this is set, the runner will use the specified id.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:   "odr",
			Target: &c.flagODR,
			Usage:  "Indicates to the runner it's operating as an on-demand runner.",
		})
	})
}

func (c *RunnerAgentCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *RunnerAgentCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *RunnerAgentCommand) Synopsis() string {
	return "Run a runner for executing remote operations."
}

func (c *RunnerAgentCommand) Help() string {
	return formatHelp(`
Usage: waypoint runner agent [options]

  Run a runner for executing remote operations.

` + c.Flags().Help())
}
