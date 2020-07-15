package cli

import (
	"io/ioutil"
	"net"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/hashicorp/go-hclog"
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/server"
	"github.com/hashicorp/waypoint/internal/server/singleprocess"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

type ServerCommand struct {
	*baseCommand

	config config.ServerConfig
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
	if c.config.DBPath == "" {
		c.ui.Output(serverWarnDBPath, terminal.WithWarningStyle())
		c.config.DBPath = "data.db"
	}
	path := c.config.DBPath
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
	ln, err := net.Listen("tcp", c.config.Listeners.GRPC)
	if err != nil {
		c.ui.Output(
			"Error starting listener: %s", err.Error(),
			terminal.WithErrorStyle(),
		)
		return 1
	}
	defer ln.Close()

	var httpLn net.Listener
	if c.config.Listeners.HTTP != "" {
		httpLn, err = net.Listen("tcp", c.config.Listeners.HTTP)
		if err != nil {
			c.ui.Output(
				"Error starting listener: %s", err.Error(),
				terminal.WithErrorStyle(),
			)
			return 1
		}
		defer httpLn.Close()
	}

	options := []server.Option{
		server.WithContext(c.Ctx),
		server.WithLogger(log),
		server.WithGRPC(ln),
		server.WithHTTP(httpLn),
		server.WithImpl(impl),
	}

	var token string

	if c.config.RequireAuth {
		ac, ok := impl.(server.AuthChecker)
		if !ok {
			c.ui.Output(
				"Server implementation not capable of authentication",
				terminal.WithErrorStyle(),
			)

			return 1
		}

		token, err = ac.DefaultToken()
		if err != nil {
			c.ui.Output(
				"Error generating default token: %s", err,
				terminal.WithErrorStyle(),
			)

			return 1
		}

		options = append(options, server.WithAuthentication(ac))
	}

	httpAddr := "disabled"
	if httpLn != nil {
		httpAddr = httpLn.Addr().String()
	}

	// Output information to the user
	c.ui.Output("Server configuration:", terminal.WithHeaderStyle())
	c.ui.Output("")
	if token == "" {
		c.ui.NamedValues([]terminal.NamedValue{
			{Name: "DB Path", Value: path},
			{Name: "gRPC Address", Value: ln.Addr().String()},
			{Name: "HTTP Address", Value: httpAddr},
		})
	} else {
		c.ui.NamedValues([]terminal.NamedValue{
			{Name: "DB Path", Value: path},
			{Name: "gRPC Address", Value: ln.Addr().String()},
			{Name: "HTTP Address", Value: httpAddr},
			{Name: "Token", Value: token},
		})
	}
	c.ui.Output("")
	c.ui.Output("Server logs:", terminal.WithHeaderStyle())
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

	// Run the server
	log.Info("starting built-in server", "addr", ln.Addr().String())
	server.Run(options...)
	return 0
}

func (c *ServerCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.StringVar(&flag.StringVar{
			Name:    "db",
			Target:  &c.config.DBPath,
			Usage:   "Path to the database file.",
			Default: "",
		})

		f.StringVar(&flag.StringVar{
			Name:    "listen-grpc",
			Target:  &c.config.Listeners.GRPC,
			Usage:   "Address to bind to for gRPC connections.",
			Default: "127.0.0.1:1234", // TODO(mitchellh: change default
		})

		f.StringVar(&flag.StringVar{
			Name:    "listen-http",
			Target:  &c.config.Listeners.HTTP,
			Usage:   "Address to bind to for HTTP connections. Required for the UI.",
			Default: "127.0.0.1:1235", // TODO(mitchellh: change default
		})

		f.BoolVar(&flag.BoolVar{
			Name:   "require-authentication",
			Target: &c.config.RequireAuth,
			Usage:  "Require authentication to communicate with the server.",
		})
	})
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
Usage: waypoint server [options]

  Run the builtin server.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

const serverWarnDBPath = `Warning! Default DB path will be used. This is at the path shown below.
The server stores persistent data here so this file should be saved and
consistent for every server run otherwise data loss will occur. It is
recommended that you explicitly set the "-db" flag as acknowledgement of
the importance of the DB file.
`
