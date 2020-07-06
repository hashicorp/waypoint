package serverclient

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"

	"github.com/hashicorp/waypoint/internal/clicontext"
)

// ConnectOption is used to configure how Waypoint server connection
// configuration is sourced.
type ConnectOption func(*connectConfig) error

// Connect connects to the Waypoint server. This returns the raw gRPC connection.
// You'll have to wrap it in NewWaypointClient to get the Waypoint client.
// We return the raw connection so that you have control over how to close it,
// and to support potentially alternate services in the future.
func Connect(ctx context.Context, opts ...ConnectOption) (*grpc.ClientConn, error) {
	// Defaults
	var cfg connectConfig
	cfg.Timeout = 5 * time.Second

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

		return nil, fmt.Errorf("no server credentials found")
	}

	// Build our options
	grpcOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTimeout(cfg.Timeout),
	}
	if cfg.Insecure {
		grpcOpts = append(grpcOpts, grpc.WithInsecure())
	}
	if cfg.Auth {
		token := cfg.Token
		if v := os.Getenv(EnvServerToken); v != "" {
			token = v
		}

		if token == "" {
			return nil, fmt.Errorf("No token available at the WAYPOINT_SERVER_TOKEN environment variable")
		}

		grpcOpts = append(grpcOpts, grpc.WithPerRPCCredentials(staticToken(token)))
	}

	// Connect to this server
	return grpc.DialContext(ctx, cfg.Addr, grpcOpts...)
}

type connectConfig struct {
	Addr     string
	Insecure bool
	Auth     bool
	Token    string
	Optional bool // See Optional func
	Timeout  time.Duration
}

// FromEnv sources the connection information from the environment
// using standard environment variables.
func FromEnv() ConnectOption {
	return func(c *connectConfig) error {
		if v := os.Getenv(EnvServerAddr); v != "" {
			c.Addr = v
			c.Insecure = os.Getenv(EnvServerInsecure) != ""
		}

		return nil
	}
}

// FromContextConfig loads a specific context config.
func FromContextConfig(cfg *clicontext.Config) ConnectOption {
	return func(c *connectConfig) error {
		if cfg != nil && cfg.Server.Address != "" {
			c.Addr = cfg.Server.Address
			c.Insecure = cfg.Server.Insecure
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

		// If we still have no name, then we do nothing.
		if n == "" {
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

// Common environment variables.
const (
	// ServerAddr is the address for the Waypoint server. This should be
	// in the format of "ip:port" for TCP.
	EnvServerAddr = "WAYPOINT_SERVER_ADDR"

	// ServerInsecure should be any value that strconv.ParseBool parses as
	// true to connect to the server insecurely.
	EnvServerInsecure = "WAYPOINT_SERVER_INSECURE"

	// EnvServerToken is the token for authenticated with the server.
	EnvServerToken = "WAYPOINT_SERVER_TOKEN"

	// EnvContext specifies a named context to load.
	EnvContext = "WAYPOINT_CONTEXT"
)

// This is a weird type that only exists to satisify the interface required by
// grpc.WithPerRPCCredentials. That api is designed to incorporate things like OAuth
// but in our case, we really just want to send this static token through, but we still
// need to the dance.
type staticToken string

func (t staticToken) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": string(t),
	}, nil
}

func (t staticToken) RequireTransportSecurity() bool {
	return false
}
