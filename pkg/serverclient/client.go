package serverclient

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"

	"github.com/hashicorp/waypoint/internal/clicontext"
	"github.com/hashicorp/waypoint/internal/env"
	"github.com/hashicorp/waypoint/pkg/inlinekeepalive"
	"github.com/hashicorp/waypoint/pkg/protocolversion"
	"github.com/hashicorp/waypoint/pkg/serverconfig"
	"github.com/hashicorp/waypoint/pkg/tokenutil"
)

// ErrNoServerConfig is the error when there is no server configuration
// found for connection.
var (
	ErrNoServerConfig = errors.New("no server connection configuration found")
	ErrInvalidToken   = errors.New("invalid token detected")
)

// ConnectOption is used to configure how Waypoint server connection
// configuration is sourced.
type ConnectOption func(*connectConfig) error

// Connect connects to the Waypoint server. This returns the raw gRPC connection.
// You'll have to wrap it in NewWaypointClient to get the Waypoint client.
// We return the raw connection so that you have control over how to close it,
// and to support potentially alternate services in the future.
//
// Authentication is done using a tokenutil.ContextToken type. You can replace the
// token at runtime by changing the context in use on a per-RPC basis.
func Connect(ctx context.Context, opts ...ConnectOption) (*grpc.ClientConn, error) {
	// Defaults
	var cfg connectConfig
	cfg.Timeout = 10 * time.Second
	cfg.Log = hclog.NewNullLogger()

	// Set config
	for _, opt := range opts {
		if err := opt(&cfg); err != nil {
			return nil, err
		}
	}

	if cfg.Addr == "" {
		if cfg.Optional {
			return nil, nil
		}

		return nil, ErrNoServerConfig
	}

	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	// Build our options
	grpcOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(protocolversion.UnaryClientInterceptor(protocolversion.Current())),
		grpc.WithChainStreamInterceptor([]grpc.StreamClientInterceptor{
			protocolversion.StreamClientInterceptor(protocolversion.Current()),

			// Send and receive keepalive messages along grpc streams.
			// Some loadbalancers (ALBs) don't respect http2 pings.
			// (https://stackoverflow.com/questions/66818645/http2-ping-frames-over-aws-alb-grpc-keepalive-ping)
			// This interceptor keeps low-traffic streams active and not timed out.
			// NOTE(izaak): long-term, we should ensure that all of our
			// streaming endpoints are robust to disconnect/resume.
			inlinekeepalive.KeepaliveClientStreamInterceptor(time.Duration(5) * time.Second),
		}...),
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				// ping after this amount of time of inactivity
				Time: 30 * time.Second,
				// send keepalive pings even if there is no active streams
				PermitWithoutStream: true,
			}),
	}

	if !cfg.Tls {
		grpcOpts = append(grpcOpts, grpc.WithInsecure())
	} else if cfg.TlsSkipVerify {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(
			credentials.NewTLS(&tls.Config{InsecureSkipVerify: true}),
		))
	} else {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(
			credentials.NewTLS(&tls.Config{}),
		))
	}

	// We always add the tokenutil.ContextToken to the gRPC options so that it can be
	// overridden. If the token is empty, it does nothing.
	var token string
	if cfg.Auth {
		token = cfg.Token
		if v := os.Getenv(EnvServerToken); v != "" {
			token = v
		}

		if token == "" {
			return nil, fmt.Errorf("No token available at the WAYPOINT_SERVER_TOKEN environment variable")
		}

		// We allow the token of "-" to explicitly be a guest.
		if token == "-" {
			token = ""
		}
	}

	var perRPC credentials.PerRPCCredentials = tokenutil.ContextToken(token)

	logArgs := []interface{}{
		"address", cfg.Addr,
		"tls", cfg.Tls,
		"tls_skip_verify", cfg.TlsSkipVerify,
		"send_auth", cfg.Auth,
		"has_token", token != "",
	}

	/*

				A wonderful flow diagram of what is happening when the Waypoint token contains
				OAuth2 creds:

		                                      (Not in the PR)
		                         ┌───────────8. OAuth dance to ─────────────────┐
		                         │            validate token.                   │
		                         │                                              ▼
		                         │                                        ┌──────────┐
		                         │                                        │  OAuth   │
		                         │                                     ┌──│  Server  │◀─┐
		                         │                                     │  │          │  │
		                         │                                     │  └──────────┘  │
		                         │                                     │                │
		                         │                                     │                │
		                         │                                6. Yeah ok, I         │
		                         │                                  know you.      5. Can I have
		                         │   ┌─────────────────────┐           │             an OAuth
		                         │   │                     │           │            token pls?
		                         ▼   ▼                     │           └────────┐    Here's my
		                      ┌────────────┐             7. Hi! Can             │      info.
		                      │  Waypoint  │             you handle             │       │
		                   ┌─▶│   Server   │──────┐        │this                │       │
		                   │  │            │      │       request,              │       │
		                   │  └────────────┘      │      here's my              ▼       │
		                   │                      │        OAuth                Λ       │
		                   │                2. Here you    token.              ╱ ╲      │
		                   │                   go, I       │                  ╱   ╲     │
		                   │                smuggled an    │                 ╱     ╲    │
		                   │                   OAuth       │                ╱       ╲   │
		               1. I want              server       └───────────────▕ client  ▏──┘
		               a Waypoint            address &                      ╲       ╱
		                 token              credentials                      ╲    4. Oh, there's
		                   │                  in it.                          ╲   OAuth info in
		                   │                      │                            ╲ ╱    this!
		                   │                      │                             V
		                   │                      │                             ▲
		                   │                      │                             │
		                   │            .─────.   │                        3. Client, do
		                   │           ; User  :  │                         your thing.
		                   └───────────:       ;◀─┘                          Here's my
		                                ╲     ╱                              Waypoint
		                                 `───'                             access token
		                                   │                                    │
		                                   │                                    │
		                                   └────────────────────────────────────┘
	*/

	if token != "" {
		externalRPC, err := tokenutil.SetupExternalCreds(ctx, cfg.Log, token)
		if err != nil {
			return nil, err
		}

		if externalRPC != nil {
			perRPC = externalRPC
		}
	}

	grpcOpts = append(grpcOpts, grpc.WithPerRPCCredentials(perRPC))

	cfg.Log.Debug("connection information", logArgs...)

	// Connect to this server
	return grpc.DialContext(ctx, cfg.Addr, grpcOpts...)
}

// ContextConfig will return the context configuration for the given connection
// options.
func ContextConfig(opts ...ConnectOption) (*clicontext.Config, error) {
	// Setup config
	var cfg connectConfig
	for _, opt := range opts {
		if err := opt(&cfg); err != nil {
			return nil, err
		}
	}

	// Build it
	return &clicontext.Config{
		Server: serverconfig.Client{
			Address:       cfg.Addr,
			Tls:           cfg.Tls,
			TlsSkipVerify: cfg.TlsSkipVerify,
			RequireAuth:   cfg.Token != "",
			AuthToken:     cfg.Token,
		},
	}, nil
}

type connectConfig struct {
	Addr          string
	Tls           bool
	TlsSkipVerify bool
	Auth          bool
	Token         string
	Optional      bool // See Optional func
	Timeout       time.Duration
	Log           hclog.Logger
}

// FromEnv sources the connection information from the environment
// using standard environment variables.
func FromEnv() ConnectOption {
	return func(c *connectConfig) error {
		if v := os.Getenv(EnvServerAddr); v != "" {
			c.Addr = v

			var err error
			c.Tls, err = env.GetBool(EnvServerTls, false)
			if err != nil {
				return err
			}

			c.TlsSkipVerify, err = env.GetBool(EnvServerTlsSkipVerify, false)
			if err != nil {
				return err
			}

			c.Auth = os.Getenv(EnvServerToken) != ""
		}

		return nil
	}
}

// FromContextConfig loads a specific context config.
func FromContextConfig(cfg *clicontext.Config) ConnectOption {
	return func(c *connectConfig) error {
		if cfg != nil && cfg.Server.Address != "" {
			c.Addr = cfg.Server.Address
			c.Tls = cfg.Server.Tls
			c.TlsSkipVerify = cfg.Server.TlsSkipVerify
			if cfg.Server.RequireAuth {
				c.Auth = true
				c.Token = cfg.Server.AuthToken
			}
		}

		return nil
	}
}

// FromContext loads the context. This will prefer the given name. If name
// is empty, we'll respect the WAYPOINT_CONTEXT env var followed by the
// default context.
func FromContext(st *clicontext.Storage, n string) ConnectOption {
	return func(c *connectConfig) error {
		// Figure out what context to load. We prefer to load a manually
		// specified one. If that isn't set, we prefer the env var. If that
		// isn't set, we load the default.
		if n == "" {
			if v := os.Getenv(EnvContext); v != "" {
				n = v
			} else {
				def, err := st.Default()
				if err != nil {
					return err
				}

				n = def
			}
		}

		// If we still have no name, then we do nothing. We also accept
		// "-" as a valid name that means "do nothing".
		if n == "" || n == "-" {
			return nil
		}

		// Load it and set it.
		cfg, err := st.Load(n)
		if err != nil {
			return err
		}

		opt := FromContextConfig(cfg)
		return opt(c)
	}
}

// Auth specifies that this server should require auth and therefore
// a token should be sourced from the environment and sent.
func Auth() ConnectOption {
	return func(c *connectConfig) error {
		c.Auth = true
		return nil
	}
}

// Optional specifies that getting server connection information is
// optional. If this is specified and no credentials are found, Connect
// will return (nil, nil). If this is NOT specified and no credentials are
// found, it is an error.
func Optional() ConnectOption {
	return func(c *connectConfig) error {
		c.Optional = true
		return nil
	}
}

// Timeout specifies a connection timeout. This defaults to 5 seconds.
func Timeout(t time.Duration) ConnectOption {
	return func(c *connectConfig) error {
		c.Timeout = t
		return nil
	}
}

// Logger is the logger to use.
func Logger(v hclog.Logger) ConnectOption {
	return func(c *connectConfig) error {
		c.Log = v
		return nil
	}
}

// Common environment variables.
const (
	// ServerAddr is the address for the Waypoint server. This should be
	// in the format of "ip:port" for TCP.
	EnvServerAddr = "WAYPOINT_SERVER_ADDR"

	// ServerTls should be any value that strconv.ParseBool parses as
	// true to connect to the server with TLS.
	EnvServerTls           = "WAYPOINT_SERVER_TLS"
	EnvServerTlsSkipVerify = "WAYPOINT_SERVER_TLS_SKIP_VERIFY"

	// EnvServerToken is the token for authenticated with the server.
	EnvServerToken = "WAYPOINT_SERVER_TOKEN"

	// EnvServerCookie is the optional cookie value for API requests.
	EnvServerCookie = "WAYPOINT_SERVER_COOKIE"

	// EnvContext specifies a named context to load.
	EnvContext = "WAYPOINT_CONTEXT"
)
