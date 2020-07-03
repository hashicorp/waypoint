package serverclient

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"
)

// ConnectOption is used to configure how Waypoint server connection
// configuration is sourced.
type ConnectOption func(*connectConfig) error

// Connect connects to the Waypoint server. This returns the raw gRPC connection.
// You'll have to wrap it in NewWaypointClient to get the Waypoint client.
// We return the raw connection so that you have control over how to close it,
// and to support potentially alternate services in the future.
func Connect(ctx context.Context, opts ...ConnectOption) (*grpc.ClientConn, error) {
	var cfg connectConfig
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
		grpc.WithTimeout(5 * time.Second),
	}
	if cfg.Insecure {
		grpcOpts = append(grpcOpts, grpc.WithInsecure())
	}
	if cfg.Auth {
		token := os.Getenv(EnvServerToken)
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
	Optional bool // See Optional func
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
