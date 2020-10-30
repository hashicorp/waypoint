package singleprocess

import (
	"crypto/tls"

	"github.com/boltdb/bolt"

	"github.com/hashicorp/go-hclog"
	grpctoken "github.com/hashicorp/horizon/pkg/grpc/token"
	wphznpb "github.com/hashicorp/waypoint-hzn/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/singleprocess/state"
	"github.com/hashicorp/waypoint/internal/serverconfig"
)

//go:generate sh -c "protoc -I proto/ proto/*.proto --go_out=plugins=grpc:gen/"

// service implements the gRPC service for the server.
type service struct {
	// state is the state management interface that provides functions for
	// safely mutating server state.
	state *state.State

	// id is our unique server ID.
	id string

	// urlConfig is not nil if the URL service is enabled. This is guaranteed
	// to have the configs set.
	urlConfig *serverconfig.URL
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

	log := cfg.log
	if log == nil {
		log = hclog.L()
	}

	// Initialize our state
	st, err := state.New(log, cfg.db)
	if err != nil {
		return nil, err
	}
	s.state = st

	// If we don't have a server ID, set that.
	id, err := st.ServerIdGet()
	if err != nil {
		return nil, err
	}
	if id == "" {
		id, err = server.Id()
		if err != nil {
			return nil, err
		}

		if err := st.ServerIdSet(id); err != nil {
			return nil, err
		}
	}
	s.id = id

	// Setup our URL service config if it is enabled.
	if scfg := cfg.serverConfig; scfg != nil && scfg.URL != nil && scfg.URL.Enabled {
		// Set our config
		s.urlConfig = scfg.URL

		// If we have no API token, get our guest account token.
		if scfg.URL.APIToken == "" {
			if err := s.initURLGuestAccount(cfg.acceptUrlTerms); err != nil {
				return nil, err
			}
		}

		// Now that we have a token, connect to the API service.
		opts := []grpc.DialOption{
			grpc.WithPerRPCCredentials(grpctoken.Token(scfg.URL.APIToken)),
		}
		if scfg.URL.APIInsecure {
			opts = append(opts, grpc.WithInsecure())
		} else {
			opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
		}

		conn, err := grpc.Dial(scfg.URL.APIAddress, opts...)
		if err != nil {
			return nil, err
		}

		s.urlClient = wphznpb.NewWaypointHznClient(conn)
	}

	// Set specific server config for the deployment entrypoint binaries
	if scfg := cfg.serverConfig; scfg != nil && scfg.CEBConfig != nil && scfg.CEBConfig.Addr != "" {
		// only one advertise address can be configured
		addr := &pb.ServerConfig_AdvertiseAddr{
			Addr:          scfg.CEBConfig.Addr,
			Tls:           scfg.CEBConfig.TLSEnabled,
			TlsSkipVerify: scfg.CEBConfig.TLSSkipVerify,
		}

		conf := &pb.ServerConfig{
			AdvertiseAddrs: []*pb.ServerConfig_AdvertiseAddr{addr},
		}

		if err := s.state.ServerConfigSet(conf); err != nil {
			return nil, err
		}
	}

	return &s, nil
}

type config struct {
	db           *bolt.DB
	serverConfig *serverconfig.Config
	log          hclog.Logger

	acceptUrlTerms bool
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
func WithConfig(scfg *serverconfig.Config) Option {
	return func(s *service, cfg *config) error {
		cfg.serverConfig = scfg
		return nil
	}
}

// WithLogger sets the logger for use with the server.
func WithLogger(log hclog.Logger) Option {
	return func(s *service, cfg *config) error {
		cfg.log = log
		return nil
	}
}

func WithAcceptURLTerms(accept bool) Option {
	return func(s *service, cfg *config) error {
		cfg.acceptUrlTerms = true
		return nil
	}
}

var _ pb.WaypointServer = (*service)(nil)
