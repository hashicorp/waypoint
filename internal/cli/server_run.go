package cli

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net"
	"strings"
	"time"

	"contrib.go.opencensus.io/exporter/ocagent"
	datadog "github.com/DataDog/opencensus-go-exporter-datadog"
	"github.com/hashicorp/go-hclog"
	hznhub "github.com/hashicorp/horizon/pkg/hub"
	hzntest "github.com/hashicorp/horizon/pkg/testutils/central"

	wphzn "github.com/hashicorp/waypoint-hzn/pkg/server"
	"github.com/hashicorp/waypoint/internal/server/boltdbstate"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/singleprocess"

	"github.com/hashicorp/waypoint/internal/telemetry"
	serverpkg "github.com/hashicorp/waypoint/pkg/server"

	"github.com/mitchellh/go-testing-interface"
	"github.com/posener/complete"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"

	"github.com/hashicorp/waypoint/internal/pkg/cert"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/server"
	"github.com/hashicorp/waypoint/internal/serverconfig"
)

const tosStatement = `
The "-accept-tos" flag must be provided to use the Waypoint URL service.
The Waypoint URL service is the free service run by HashiCorp that provides
an automatic "waypoint.run" URL for each application and deployment. You can use
this URL to quickly view your deployed applications and share your applications
with others.

Usage of this service requires accepting the Terms of Service and a Privacy
Policy at the URLs below. If you do not feel comfortable accepting the terms,
you may disable the URL service or self-host the URL service.
Learn more about this at: https://waypointproject.io/docs/url

Privacy Policy: https://hashicorp.com/privacy
Terms of Service: https://waypointproject.io/terms

Please rerun this command using the "-accept-tos" flag to accept the terms above.
`

const acceptTOSHelp = `Pass to accept the Terms of Service and Privacy Policy to use the Waypoint URL Service. This is required if the URL service is enabled and you're using the HashiCorp-provided URL service rather than self-hosting. See the privacy policy at https://hashicorp.com/privacy and the ToS at https://waypointproject.io/terms`

const DefaultURLControlAddress = "https://control.hzn.network"

type ServerRunCommand struct {
	*baseCommand

	config        serverconfig.Config
	flagDisableUI bool
	flagURLInmem  bool

	flagAdvertiseAddr          string
	flagAdvertiseTLSEnabled    bool
	flagAdvertiseTLSSkipVerify bool
	flagAcceptTOS              bool
	flagTLSCertFile            string
	flagTLSKeyFile             string

	flagTelemetryOpenCensusAgentAddr     string
	flagTelemetryOpenCensusAgentInsecure bool

	flagTelemetryDatadogTraceAddr     string
	flagTelemetryOpenCensusZpagesAddr string
}

func (c *ServerRunCommand) Run(args []string) int {
	defer c.Close()
	log := c.Log.Named("server")

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
		WithNoClient(),
	); err != nil {
		return 1
	}

	if c.config.URL.Enabled &&
		c.config.URL.ControlAddress == DefaultURLControlAddress &&
		!c.flagAcceptTOS {
		c.ui.Output(strings.TrimSpace(tosStatement), terminal.WithErrorStyle())
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
		c.config.URL = &serverconfig.URL{
			Enabled:        true,
			APIAddress:     wphzndata.Addr,
			APIInsecure:    true,
			ControlAddress: fmt.Sprintf("dev://%s", setup.HubAddr),
		}
	}

	// Set any server config
	c.config.CEBConfig = &serverconfig.CEBConfig{
		Addr:          c.flagAdvertiseAddr,
		TLSEnabled:    c.flagAdvertiseTLSEnabled,
		TLSSkipVerify: c.flagAdvertiseTLSSkipVerify,
	}

	// Initialize our state
	state, err := boltdbstate.New(log, db)
	if err != nil {
		c.ui.Output(
			"Error initializing boltdb state backend: %s", err.Error(),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	// Create our server
	impl, err := singleprocess.New(
		singleprocess.WithState(state),
		singleprocess.WithConfig(&c.config),
		singleprocess.WithLogger(log.Named("singleprocess")),
		singleprocess.WithAcceptURLTerms(c.flagAcceptTOS),
		singleprocess.WithFeatures(pb.ServerFeatures_FEATURE_INLINE_KEEPALIVES),
	)
	if c, ok := impl.(io.Closer); ok {
		defer c.Close()
	}

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

	var httpInsecureLn net.Listener
	if c.config.HTTPInsecure.Addr != "" {
		c.config.HTTPInsecure.TLSDisable = true
		httpInsecureLn, err = c.listenerForConfig(log.Named("http_insecure"), &c.config.HTTPInsecure)
		if err != nil {
			c.ui.Output(
				"Error starting insecure HTTP listener: %s", err.Error(),
				terminal.WithErrorStyle(),
			)
			return 1
		}
		defer httpInsecureLn.Close()
	}

	telemetryEnabled := false
	// If we have any export address configured, we should enable telemetry.
	if c.flagTelemetryOpenCensusAgentAddr != "" || c.flagTelemetryDatadogTraceAddr != "" {
		telemetryEnabled = true
	}

	options := []server.Option{
		server.WithContext(c.Ctx),
		server.WithLogger(log),
		server.WithGRPC(ln),
		server.WithHTTP(httpLn),
		server.WithImpl(impl),
		server.WithTelemetry(telemetryEnabled),
	}
	if httpInsecureLn != nil {
		options = append(options, server.WithHTTP(httpInsecureLn))
	}
	auth := false
	if ac, ok := impl.(serverpkg.AuthChecker); ok {
		options = append(options, server.WithAuthentication(ac))
		auth = true
	}

	ui := true
	if !c.flagDisableUI {
		options = append(options, server.WithBrowserUI(true))
	} else {
		ui = false
	}

	// Output information to the user
	c.ui.Output("Server configuration:", terminal.WithHeaderStyle())
	c.ui.Output("")
	values := []terminal.NamedValue{
		{Name: "DB Path", Value: path},
		{Name: "gRPC Address", Value: ln.Addr().String()},
		{Name: "HTTP Address", Value: httpLn.Addr().String()},
	}
	if httpInsecureLn != nil {
		values = append(values, terminal.NamedValue{
			Name: "HTTP Address (Insecure)", Value: httpInsecureLn.Addr().String(),
		})
	}
	if auth {
		values = append(values, terminal.NamedValue{Name: "Auth Required", Value: "yes"})
	}
	if ui {
		values = append(values, terminal.NamedValue{Name: "Browser UI Enabled", Value: "yes"})
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

	// If we aren't bootstrapped, let the user know
	if bs, ok := impl.(interface {
		Bootstrapped(ctx context.Context) bool
	}); ok && auth && !bs.Bootstrapped(c.Ctx) {
		c.ui.Output("Server requires bootstrapping!", terminal.WithHeaderStyle())
		c.ui.Output("")
		c.ui.Output(strings.TrimSpace(`
New servers must be bootstrapped to retrieve the initial auth token for
connections. To bootstrap this server, run the following command in your
terminal once the server is up and running.

  waypoint server bootstrap -server-addr=%s -server-tls-skip-verify

This command will bootstrap the server and setup a CLI context.
`), ln.Addr().String(), terminal.WithInfoStyle())
	}

	c.ui.Output("Server logs:", terminal.WithHeaderStyle())
	c.ui.Output("")

	// Close our UI. If we're using the interactive UI we need to end the
	// output right now so we don't redraw over our logs.
	if closer, ok := c.ui.(io.Closer); ok && closer != nil {
		closer.Close()
	}

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

	if telemetryEnabled {
		log.Debug("Telemetry enabled - configuring telemetry options.")
		telemetryOptions := []telemetry.Option{
			telemetry.WithLogger(log.Named("telemetry")),
		}

		// ZPages
		if c.flagTelemetryOpenCensusZpagesAddr != "" {
			telemetryOptions = append(telemetryOptions, telemetry.WithZpages(c.flagTelemetryOpenCensusZpagesAddr))
		}

		// OpenCensus Agent
		if c.flagTelemetryOpenCensusAgentAddr != "" {
			ocagentOptions := []ocagent.ExporterOption{
				ocagent.WithAddress(c.flagTelemetryOpenCensusAgentAddr),
				ocagent.WithServiceName("waypoint"),
			}
			if c.flagTelemetryOpenCensusAgentInsecure {
				ocagentOptions = append(ocagentOptions, ocagent.WithInsecure())
			}

			telemetryOptions = append(telemetryOptions, telemetry.WithOpenCensusExporter(ocagentOptions))
		}

		// Datadog
		if c.flagTelemetryDatadogTraceAddr != "" {
			telemetryOptions = append(telemetryOptions, telemetry.WithDatadogExporter(
				datadog.Options{
					TraceAddr: c.flagTelemetryDatadogTraceAddr,
					StatsAddr: c.flagTelemetryDatadogTraceAddr,
					Service:   "waypoint",
					Namespace: "waypoint",
				},
			))
		}

		t, err := telemetry.NewTelemetry(telemetryOptions...)
		if err != nil {
			c.ui.Output(
				"Error setting up telemetry: %s", err.Error(),
				terminal.WithErrorStyle(),
			)
			return 1
		}

		go func() {
			// Will gracefully exit when the context passed in closes.
			err := t.Run(c.Ctx)
			if err != nil {
				log.Error("Telemetry runner exited with error", "error", err)
			} else {
				log.Debug("Telemetry runner completed.")
			}
		}()
	}

	// Run the server
	log.Info("starting built-in server", "addr", ln.Addr().String())
	server.Run(options...)
	return 0
}

func (c *ServerRunCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		if c.config.URL == nil {
			c.config.URL = &serverconfig.URL{}
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
			Default: "127.0.0.1:9701",
		})

		f.StringVar(&flag.StringVar{
			Name:    "listen-http",
			Target:  &c.config.HTTP.Addr,
			Usage:   "Address to bind to for HTTP connections. Required for the UI.",
			Default: "127.0.0.1:9702",
		})

		f.StringVar(&flag.StringVar{
			Name:   "listen-http-insecure",
			Target: &c.config.HTTPInsecure.Addr,
			Usage: "Address to bind to for insecure HTTP connections. " +
				"This will not have TLS enabled. This will redirect users to the TLS " +
				"port UNLESS there is an X-Forwarded-Proto header. This makes this port " +
				"suitable for proxy backends.",
		})

		f.StringVar(&flag.StringVar{
			Name:   "tls-cert-file",
			Target: &c.flagTLSCertFile,
			Usage: "Path to a PEM-encoded certificate file for TLS. If this " +
				"isn't set, a self-signed certificate will be generated. This file " +
				"will be monitored for changes and will automatically reload on change.",
			Default: "",
		})

		f.StringVar(&flag.StringVar{
			Name:   "tls-key-file",
			Target: &c.flagTLSKeyFile,
			Usage: "Path to a PEM-encoded private key file for the TLS certificate " +
				"specified with -tls-cert-file. This is required if -tls-cert-file " +
				"is set. This file will be monitored for changes and will automatically " +
				"reload on change",
			Default: "",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "disable-ui",
			Target:  &c.flagDisableUI,
			Usage:   "Disable the embedded web interface",
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
			Usage:   "Address to Waypoint URL service API.",
			Default: "api.waypoint.run:443",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "url-api-insecure",
			Target:  &c.config.URL.APIInsecure,
			Usage:   "True if TLS is not enabled for the Waypoint URL service API.",
			Default: false,
		})

		f.StringVar(&flag.StringVar{
			Name:    "url-control-addr",
			Target:  &c.config.URL.ControlAddress,
			Usage:   "Address to Waypoint URL service control API.",
			Default: DefaultURLControlAddress,
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

		f.StringVar(&flag.StringVar{
			Name:   "advertise-addr",
			Target: &c.flagAdvertiseAddr,
			Usage: "Address to advertise for the server. This is used by the entrypoints\n" +
				"binaries to communicate back to the server. If this is blank, then\n" +
				"the entrypoints will not communicate to the server. Features such as\n" +
				"logs, exec, etc. will not work.",
		})
		f.BoolVar(&flag.BoolVar{
			Name:    "advertise-tls",
			Target:  &c.flagAdvertiseTLSEnabled,
			Usage:   "If true, the advertised address should be connected to with TLS.",
			Default: true,
		})
		f.BoolVar(&flag.BoolVar{
			Name:    "advertise-tls-skip-verify",
			Target:  &c.flagAdvertiseTLSSkipVerify,
			Usage:   "Do not verify the TLS certificate presented by the server.",
			Default: false,
		})
		f.BoolVar(&flag.BoolVar{
			Name:    "accept-tos",
			Target:  &c.flagAcceptTOS,
			Usage:   acceptTOSHelp,
			Default: false,
		})
		f.StringVar(&flag.StringVar{
			Name:   "telemetry-oc-agent-addr",
			Target: &c.flagTelemetryOpenCensusAgentAddr,
			Usage: "Address of an OpenCensus agent or collector available to receive OpenCensus formatted " +
				"telemetry, traces and stats (commonly port 55678). Example: localhost:55678",
		})
		f.BoolVar(&flag.BoolVar{
			Name:    "telemetry-oc-agent-insecure",
			Target:  &c.flagTelemetryOpenCensusAgentInsecure,
			Usage:   "Disables client transport security for the OpenCensus agent exporter's gRPC connection.",
			Default: false,
		})
		f.StringVar(&flag.StringVar{
			Name:   "telemetry-dd-trace-addr",
			Target: &c.flagTelemetryDatadogTraceAddr,
			Usage:  "Address of a DataDog agent available to accept traces (commonly port 8126). Example: localhost:8126",
		})
		f.StringVar(&flag.StringVar{
			Name:   "telemetry-oc-zpages-addr",
			Target: &c.flagTelemetryOpenCensusZpagesAddr,
			Usage: "If set, Waypoint will run an OpenCensus zPages server at this address. Typically this is " +
				"set to something like localhost:9999, where trace debug information could be viewed at " +
				"http://localhost:9999/debug/tracez, and rpc stats at http://localhost:55679/debug/rpcz. " +
				"More information at https://opencensus.io/zpages/",
		})
	})
}

func (c *ServerRunCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ServerRunCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ServerRunCommand) Synopsis() string {
	return "Manually run the builtin server"
}

func (c *ServerRunCommand) Help() string {
	return formatHelp(`
Usage: waypoint server run [options]

  Run the builtin server.

  The easier way to run a server is to use the "waypoint install" command.
  This command is for people who are manually running the server in any
  environment.

` + c.Flags().Help())
}

func (c *ServerRunCommand) listenerForConfig(log hclog.Logger, cfg *serverconfig.Listener) (net.Listener, error) {
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

	// Use the TLS cert/key specified in the config. If one isn't specified,
	// default to any set via the TLS flags.
	certFile := cfg.TLSCertFile
	keyFile := cfg.TLSKeyFile
	if certFile == "" {
		certFile = c.flagTLSCertFile
		keyFile = c.flagTLSKeyFile
	}

	var tlsConfig *tls.Config
	if certFile != "" {
		w, err := cert.New(log, certFile, keyFile)
		if err != nil {
			return nil, err
		}
		tlsConfig = w.TLSConfig()

		// NOTE(mitchellh): technically, the cert watcher wants to be closed
		// to clean up some goroutines, but our listener lives for the duration
		// of the process so we don't clean it up.

		log.Info("TLS certs loaded from specified files",
			"cert", certFile,
			"key", keyFile)
	}

	// If we don't have a cert then we self-sign.
	if tlsConfig == nil {
		log.Info("TLS cert wasn't specified, a self-signed certificate will be created")

		priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
		if err != nil {
			return nil, err
		}

		template := x509.Certificate{
			SerialNumber: big.NewInt(time.Now().Unix()),
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
		certPEM := out.Bytes()

		// Write the key
		out = bytes.Buffer{}
		if err := pem.Encode(&out, pemBlockForKey(priv)); err != nil {
			return nil, err
		}
		keyPEM := out.Bytes()

		// Setup the TLS listener
		cert, err := tls.X509KeyPair(certPEM, keyPEM)
		if err != nil {
			ln.Close()
			return nil, err
		}

		// Setup our TLS config
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}

	log.Info("listener is wrapped with TLS")
	return tls.NewListener(ln, tlsConfig), nil
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
