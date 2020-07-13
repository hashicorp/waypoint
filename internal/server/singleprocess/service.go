package singleprocess

import (
	"github.com/boltdb/bolt"

	hzncontrol "github.com/hashicorp/horizon/pkg/control"
	wphznpb "github.com/hashicorp/waypoint-hzn/pkg/pb"
	"google.golang.org/grpc"

	configpkg "github.com/hashicorp/waypoint/internal/config"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/singleprocess/state"
)

//go:generate sh -c "protoc -I proto/ proto/*.proto --go_out=plugins=grpc:gen/"

// service implements the gRPC service for the server.
type service struct {
	// state is the state management interface that provides functions for
	// safely mutating server state.
	state *state.State

	// urlConfig is not nil if the URL service is enabled. This is guaranteed
	// to have the configs set.
	urlConfig *configpkg.URL
	urlClient wphznpb.WaypointHznClient
}

// New returns a Waypoint server implementation that uses BotlDB plus
// in-memory locks to operate safely.
func New(opts ...Option) (pb.WaypointServer, error) {
	var s service
	var cfg config
	for _, opt := range opts {
		if err := opt(&s, &cfg); err != nil {
			return nil, err
		}
	}

	// Initialize our state
	st, err := state.New(cfg.db)
	if err != nil {
		return nil, err
	}
	s.state = st

	// Setup our URL service config if it is enabled.
	if scfg := cfg.serverConfig; scfg != nil && scfg.URL != nil && scfg.URL.Enabled {
		opts := []grpc.DialOption{
			grpc.WithPerRPCCredentials(hzncontrol.Token(scfg.URL.APIToken)),
		}
		if scfg.URL.APIInsecure {
			opts = append(opts, grpc.WithInsecure())
		}

		conn, err := grpc.Dial(scfg.URL.APIAddress, opts...)
		if err != nil {
			return nil, err
		}

		s.urlConfig = scfg.URL
		s.urlClient = wphznpb.NewWaypointHznClient(conn)
	}

	return &s, nil
}

type config struct {
	db           *bolt.DB
	serverConfig *configpkg.ServerConfig
}

type Option func(*service, *config) error

// WithDB sets the Bolt DB for use with the server.
func WithDB(db *bolt.DB) Option {
	return func(s *service, cfg *config) error {
		cfg.db = db
		return nil
	}
}

// WithConfig sets the server config in use with this server.
func WithConfig(scfg *configpkg.ServerConfig) Option {
	return func(s *service, cfg *config) error {
		cfg.serverConfig = scfg
		return nil
	}
}

var _ pb.WaypointServer = (*service)(nil)
