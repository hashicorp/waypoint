package serverclient

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"

	"github.com/hashicorp/waypoint/internal/config"
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

	// Connect to this server
	return grpc.DialContext(ctx, cfg.Addr, grpcOpts...)
}

type connectConfig struct {
	Addr     string
	Insecure bool
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

// FromConfig sources connection information from the configuration.
func FromConfig(cfg *config.Config) ConnectOption {
	return func(c *connectConfig) error {
		if cfg.Server != nil && cfg.Server.Address != "" {
			c.Addr = cfg.Server.Address
			c.Insecure = cfg.Server.Insecure
		}

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
)
