package cli

import (
	"net"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/hashicorp/go-hclog"
	"github.com/posener/complete"

	"github.com/mitchellh/devflow/internal/pkg/flag"
	"github.com/mitchellh/devflow/internal/server"
	"github.com/mitchellh/devflow/internal/server/singleprocess"
	"github.com/mitchellh/devflow/sdk/terminal"
)

type ServerCommand struct {
	*baseCommand
}

func (c *ServerCommand) Run(args []string) int {
	defer c.Close()
	log := c.Log.Named("server")

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
	); err != nil {
		return 1
	}

	// Open our database
	path := "data.db"
	log.Info("opening DB", "path", path)
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		c.ui.Output(
			"Error opening database: %s", err.Error(),
			terminal.WithErrorStyle(),
		)
		return 1
	}
	defer func() {
		log.Info("gracefully closing db", "path", path)
		if err := db.Close(); err != nil {
			log.Error("error closing db", "path", path, "err", err)
		}
	}()

	// Create our server
	impl, err := singleprocess.New(db)
	if err != nil {
		c.ui.Output(
			"Error initializing server: %s", err.Error(),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	// We listen on a random locally bound port
	ln, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		c.ui.Output(
			"Error starting listener: %s", err.Error(),
			terminal.WithErrorStyle(),
		)
		return 1
	}
	defer ln.Close()

	// Output information to the user
	c.ui.Output("Server configuration:", terminal.WithHeaderStyle())
	c.ui.Output("")
	c.ui.Output(`
DB Path: %[1]s
gRPC Address: %[2]s`,
		path, ln.Addr().String(),
		terminal.WithKeyValueStyle(":"),
		terminal.WithStatusStyle(),
	)
	c.ui.Output("")
	c.ui.Output("Server logs:", terminal.WithHeaderStyle())
	c.ui.Output("")

	// Set our log output higher if its not already so that it begins showing.
	if !log.IsInfo() {
		log.SetLevel(hclog.Info)
	}

	// Run the server
	log.Info("starting built-in server", "addr", ln.Addr().String())
	server.Run(server.WithContext(c.Ctx),
		server.WithLogger(log),
		server.WithGRPC(ln),
		server.WithImpl(impl),
	)
	return 0
}

func (c *ServerCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *ServerCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ServerCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ServerCommand) Synopsis() string {
	return "Run the builtin server."
}

func (c *ServerCommand) Help() string {
	helpText := `
Usage: devflow server [options]

  Run the builtin server.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}
