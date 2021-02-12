package cli

import (
	"io/ioutil"

	"github.com/hashicorp/go-hclog"
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	runnerpkg "github.com/hashicorp/waypoint/internal/runner"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/serverclient"
)

type RunnerAgentCommand struct {
	*baseCommand
}

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
	c.ui.Output("")

	// Set our log output higher if its not already so that it begins showing.
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
	runner, err := runnerpkg.New(
		runnerpkg.WithClient(client),
		runnerpkg.WithLogger(log.Named("runner")),
	)
	if err != nil {
		c.ui.Output(
			"Error initializing the runner: %s", err.Error(),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	// Start the runner
	log.Info("starting runner")
	if err := runner.Start(); err != nil {
		log.Error("error starting runner", "err", err)
		return 1
	}

	// Accept jobs in goroutine so that we can interrupt it.
	go func() {
		for {
			if err := runner.Accept(ctx); err != nil {
				if err == runnerpkg.ErrClosed {
					return
				}

				log.Error("error running job", "err", err)
			}
		}
	}()

	// Wait for end
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
	return c.flagSet(0, nil)
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
