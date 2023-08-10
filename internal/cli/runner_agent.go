// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cli

import (
	"context"
	"errors"
	"io/ioutil"
	"net"
	"runtime"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/posener/complete"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/plugin"
	runnerpkg "github.com/hashicorp/waypoint/internal/runner"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverclient"
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

	// If this is an ODR runner, this should be the ODR profile that it was created
	// from.
	flagOdrProfileId string

	// Cookie to use for API requests. Importantly, this enables runner adoption.
	flagCookie string

	// State directory for runner.
	flagStateDir string

	// Labels for the runner.
	flagLabels map[string]string

	// The amount of concurrent jobs that can be running.
	flagConcurrency int
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
		WithNoLocalServer(),
	); err != nil {
		return 1
	}

	plugin.InsideODR = c.flagODR

	// Flag defaults
	if c.flagConcurrency == 0 {
		c.flagConcurrency = runtime.NumCPU() * 3
	}

	// Check again in case it was set to 0.
	if c.flagConcurrency < 1 {
		log.Warn("concurrency flag less than 1 has no effect, using 1")
		c.flagConcurrency = 1
	}

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

	// Build the values we'll show in the runner config table
	infoValues := []terminal.NamedValue{
		{Name: "Server address", Value: conn.Target()},
	}
	if c.flagODR {
		infoValues = append(infoValues, terminal.NamedValue{
			Name: "Type", Value: "on-demand",
		})
	} else {
		infoValues = append(infoValues, terminal.NamedValue{
			Name: "Type", Value: "remote",
		})
	}

	// Output information to the user
	c.ui.Output("Runner configuration:", terminal.WithHeaderStyle())
	c.ui.NamedValues(infoValues)
	c.ui.Output("Runner logs:", terminal.WithHeaderStyle())

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
		runnerpkg.WithStateDir(c.flagStateDir),
		runnerpkg.WithLabels(c.flagLabels),
	}

	if c.flagId != "" {
		options = append(options, runnerpkg.WithId(c.flagId))
	}

	if c.flagCookie != "" {
		options = append(options, runnerpkg.WithCookie(c.flagCookie))
	}

	if c.flagODR {
		options = append(options,
			runnerpkg.WithODR(c.flagOdrProfileId),
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

	// We always defer the close, but in happy paths this should be a noop
	// because we do a close later in this function more gracefully.
	defer runner.Close()

	// If we have a liveness address setup, start the liveness server.
	// We need to do this before starting the runner, because the runner
	// startup might block on waiting for adoption, and we need
	// the underlying platform to keep the runner alive until that completes
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

	// Start the runner
	log.Info("starting runner", "id", runner.Id())
	if err := runner.Start(ctx); err != nil {
		log.Error("error starting runner", "err", err)
		return 1
	}

	// Accept jobs in goroutine so that we can interrupt it.
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		defer cancel()

		// In non-ODR mode, we accept many jobs in parallel.
		if !c.flagODR {
			runner.AcceptParallel(ctx, c.flagConcurrency)
			return
		}

		// In ODR mode, we accept a single job.
		for {
			err := runner.Accept(ctx)
			if err == nil {
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

				case codes.PermissionDenied, codes.Unauthenticated:
					// The runner was rejected after the fact or our token
					// was revoked. We exit and expect an init process or
					// something to restart us if they want to retry.
					log.Error("no permission to request a job, exiting")
					return
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
			Hidden: true,
		})

		f.StringVar(&flag.StringVar{
			Name:   "odr-profile-id",
			Target: &c.flagOdrProfileId,
			Usage:  "The ID of the odr profile used to create the task that is running this runner.",
			Hidden: true,
		})

		f.StringVar(&flag.StringVar{
			Name:   "cookie",
			Target: &c.flagCookie,
			Usage: "The cookie value of the server to validate API requests. " +
				"This is required for runner adoption. If you do not already have a " +
				"runner token, this must be set.",
			EnvVar: serverclient.EnvServerCookie,
		})

		f.StringVar(&flag.StringVar{
			Name:   "state-dir",
			Target: &c.flagStateDir,
			Usage: "Directory to store state between restarts. This is optional. If " +
				"this is set, then a runner can restart without re-triggering the adoption " +
				"process.",
		})

		f.StringMapVar(&flag.StringMapVar{
			Name:   "label",
			Target: &c.flagLabels,
			Usage:  "Labels to set for this runner in 'k=v' format. Can be specified multiple times.",
		})

		f.IntVar(&flag.IntVar{
			Name:   "concurrency",
			Target: &c.flagConcurrency,
			Usage: "The number of concurrent jobs that can be running at one time. " +
				"This has no effect if `-odr` is set. The default value applied will be " +
				"(total number of logical cpus available * 3). A value of less than 1 will " +
				"default to 1.",

			// Most jobs that a non-ODR runner runs are IO bound, so we use
			// just a heuristic here of allowing some multiple above the CPUs.
			//Default: runtime.NumCPU() * 3,
			// NOTE(briancain): We set a default of 0 here, but when the CLI goes
			// to use this value, if set to 0, we'll attempt to set it to the default
			// runtime.NumCPU()*3.
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

  Run a remote runner for executing remote operations.

  Runners are named or identified via the ID and the label set. The ID
  can be manually specified or automatically generated. The label set is
  specified using "-label" flags.

  A runner can be registered with the server in two ways. First, a
  runner token can be created with "waypoint runner token" and used with
  this command (using the WAYPOINT_SERVER_TOKEN environment variable,
  "waypoint context", etc.). This will allow the runner to begin accepting
  jobs immediately since it is preauthorized.

  The second approach is to specify only the cookie value (acquired using
  the "waypoint server cookie" command) and the server address. This will
  trigger a process that puts the runner in a pending state until a human
  manually verifies it. This is useful for easily installing runners.

  The "-state-dir" flag is optional, but important. This flag allows runners
  to restart gracefully without regenerating a new ID or losing a rotated
  authentication token. Runners can be run without a state directory but it is
  not generally recommended.

` + c.Flags().Help())
}
