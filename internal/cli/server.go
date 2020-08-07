package cli

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"strings"
	"time"

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

	config          config.ServerConfig
	flagDisableAuth bool
	flagURLInmem    bool
}

func (c *ServerCommand) Run(args []string) int {
	defer c.Close()
	log := c.Log.Named("server")

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
		WithClient(false),
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
	ln, err := c.listenerForConfig(log.Named("grpc"), &c.config.GRPC)
	if err != nil {
		c.ui.Output(
			"Error starting listener: %s", err.Error(),
			terminal.WithErrorStyle(),
		)
		return 1
	}
	defer ln.Close()

	httpLn, err := c.listenerForConfig(log.Named("http"), &c.config.HTTP)
	if err != nil {
		c.ui.Output(
			"Error starting listener: %s", err.Error(),
			terminal.WithErrorStyle(),
		)
		return 1
	}
	defer httpLn.Close()

	options := []server.Option{
		server.WithContext(c.Ctx),
		server.WithLogger(log),
		server.WithGRPC(ln),
		server.WithHTTP(httpLn),
		server.WithImpl(impl),
	}
	auth := false
	if ac, ok := impl.(server.AuthChecker); ok && !c.flagDisableAuth {
		options = append(options, server.WithAuthentication(ac))
		auth = true
	}

	// Output information to the user
	c.ui.Output("Server configuration:", terminal.WithHeaderStyle())
	values := []terminal.NamedValue{
		{Name: "DB Path", Value: path},
		{Name: "gRPC Address", Value: ln.Addr().String()},
		{Name: "HTTP Address", Value: httpLn.Addr().String()},
	}
	if auth {
		values = append(values, terminal.NamedValue{Name: "Auth Required", Value: "yes"})
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
			Target:  &c.config.GRPC.Addr,
			Usage:   "Address to bind to for gRPC connections.",
			Default: "127.0.0.1:1234", // TODO(mitchellh: change default
		})

		f.StringVar(&flag.StringVar{
			Name:    "listen-http",
			Target:  &c.config.HTTP.Addr,
			Usage:   "Address to bind to for HTTP connections. Required for the UI.",
			Default: "127.0.0.1:1235", // TODO(mitchellh: change default
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "disable-auth",
			Target:  &c.flagDisableAuth,
			Usage:   "Disable auth requirements",
			Default: false,
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
			Default: "api.alpha.waypoint.run:443", // TODO(mitchellh: change default
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
			Default: "https://control.alpha.hzn.network",
		})

		f.StringVar(&flag.StringVar{
			Name:    "url-control-token",
			Target:  &c.config.URL.APIToken,
			Usage:   "Token for the Waypoint URL server control API.",
			Default: "",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "url-auto-app-hostname",
			Target:  &c.config.URL.AutomaticAppHostname,
			Usage:   "Whether apps automatically get a hostname on deploy.",
			Default: true,
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

func (c *ServerCommand) listenerForConfig(log hclog.Logger, cfg *config.Listener) (net.Listener, error) {
	// Start our bare listener
	log.Debug("starting listener", "addr", cfg.Addr)
	ln, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, err
	}

	// If we have TLS disabled then we're done.
	if cfg.TLSDisable {
		log.Warn("TLS is disabled for this listener")
		return ln, nil
	}

	// If we don't have a cert then we self-sign.
	var certPEM, keyPEM []byte
	if cfg.TLSCertFile != "" {
		certPEM, err = ioutil.ReadFile(cfg.TLSCertFile)
		if err != nil {
			return nil, err
		}
		keyPEM, err = ioutil.ReadFile(cfg.TLSKeyFile)
		if err != nil {
			return nil, err
		}

		log.Info("TLS certs loaded from specified files",
			"cert", cfg.TLSCertFile,
			"key", cfg.TLSKeyFile)
	}

	if certPEM == nil {
		log.Info("TLS cert wasn't specified, a self-signed certificate will be created")

		priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		if err != nil {
			return nil, err
		}

		template := x509.Certificate{
			SerialNumber: big.NewInt(1),
			Subject: pkix.Name{
				Organization: []string{"Waypoint"},
			},
			NotBefore:             time.Now(),
			NotAfter:              time.Now().Add(time.Hour * 24 * 365 * 10),
			KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			BasicConstraintsValid: true,
		}

		derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
		if err != nil {
			return nil, err
		}

		// Write the cert
		var out bytes.Buffer
		err = pem.Encode(&out, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
		if err != nil {
			return nil, err
		}
		certPEM = []byte(out.String())
		out.Reset()

		// Write the key
		if err := pem.Encode(&out, pemBlockForKey(priv)); err != nil {
			return nil, err
		}
		keyPEM = []byte(out.String())
	}

	// Setup the TLS listener
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		ln.Close()
		return nil, err
	}

	log.Info("listener is wrapped with TLS")
	return tls.NewListener(ln, &tls.Config{
		Certificates: []tls.Certificate{cert},
	}), nil
}

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			panic(err)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

const serverWarnDBPath = `Warning! Default DB path will be used. This is at the path shown below.
The server stores persistent data here so this file should be saved and
consistent for every server run otherwise data loss will occur. It is
recommended that you explicitly set the "-db" flag as acknowledgement of
the importance of the DB file.
`
