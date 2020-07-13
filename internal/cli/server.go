package cli

import (
	"fmt"
	"io/ioutil"
	"net"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/hashicorp/go-hclog"
	hznhub "github.com/hashicorp/horizon/pkg/hub"
	hzntest "github.com/hashicorp/horizon/pkg/testutils/central"
	wphzn "github.com/hashicorp/waypoint-hzn/pkg/server"
	"github.com/mitchellh/go-testing-interface"
	"github.com/posener/complete"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/server"
	"github.com/hashicorp/waypoint/internal/server/singleprocess"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

type ServerCommand struct {
	*baseCommand

	config       config.ServerConfig
	flagURLInmem bool
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

	// Run our in-memory URL service.
	if c.flagURLInmem {
		t := &testing.RuntimeT{}

		// Create the inmem Horizon server.
		setupCh := make(chan *hzntest.DevSetup, 1)
		closeCh := make(chan struct{})
		defer close(closeCh)
		go hzntest.Dev(t, func(setup *hzntest.DevSetup) {
			hubclient, err := hznhub.NewHub(log.Named("url-hub"), setup.ControlClient, setup.HubToken)
			require.NoError(t, err)
			go hubclient.Run(c.Ctx, setup.ClientListener)

			setupCh <- setup
			<-closeCh
		})
		setup := <-setupCh

		// Create the inmem waypoint-hzn server
		wphzndata := wphzn.TestServer(t)

		// Configure
		c.config.URL = &config.URL{
			Enabled:        true,
			APIAddress:     wphzndata.Addr,
			APIInsecure:    true,
			ControlAddress: fmt.Sprintf("dev://%s", setup.HubAddr),
			Token:          setup.AgentToken,
		}
	}

	// Create our server
	impl, err := singleprocess.New(
		singleprocess.WithDB(db),
		singleprocess.WithConfig(&c.config),
	)
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
	values := []terminal.NamedValue{
		{Name: "DB Path", Value: path},
		{Name: "gRPC Address", Value: ln.Addr().String()},
		{Name: "HTTP Address", Value: httpAddr},
	}
	if token != "" {
		values = append(values, terminal.NamedValue{Name: "Token", Value: token})
	}
	if !c.config.URL.Enabled {
		values = append(values, terminal.NamedValue{Name: "URL Service", Value: "disabled"})
	} else {
		value := c.config.URL.APIAddress
		if c.config.URL.APIToken == "" {
			value += " (account: guest)"
		} else {
			value += " (account: token)"
		}

		values = append(values, terminal.NamedValue{Name: "URL Service", Value: value})
	}
	c.ui.NamedValues(values)
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
		if c.config.URL == nil {
			c.config.URL = &config.URL{}
		}

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

		f.BoolVar(&flag.BoolVar{
			Name:    "url-enabled",
			Target:  &c.config.URL.Enabled,
			Usage:   "Enable the URL service.",
			Default: true,
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "url-inmem",
			Target:  &c.flagURLInmem,
			Usage:   "Run an in-memory URL service for dev purposes.",
			Default: false,
			Hidden:  true,
		})

		f.StringVar(&flag.StringVar{
			Name:    "url-api-addr",
			Target:  &c.config.URL.APIAddress,
			Usage:   "Address to Waypoint URL service API",
			Default: "", // TODO(mitchellh: change default
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "url-api-insecure",
			Target:  &c.config.URL.APIInsecure,
			Usage:   "True if TLS is not enabled for the Waypoint URL service API",
			Default: false,
		})

		f.StringVar(&flag.StringVar{
			Name:    "url-control-addr",
			Target:  &c.config.URL.ControlAddress,
			Usage:   "Address to Waypoint URL service control API",
			Default: "", // TODO(mitchellh: change default
		})

		f.StringVar(&flag.StringVar{
			Name:    "url-control-token",
			Target:  &c.config.URL.Token,
			Usage:   "Token for the Waypoint URL server control API.",
			Default: "", // TODO(mitchellh: change default
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
