package singleprocess

import (
	"context"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/hashicorp/go-hclog"
	wphznpb "github.com/hashicorp/waypoint-hzn/pkg/pb"

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
	urlConfig    *serverconfig.URL
	urlClientMu  sync.Mutex
	urlClientVal wphznpb.WaypointHznClient

	// urlCEB has the configuration for the entrypoint. If this is nil,
	// then the configuration is not ready. The urlCEBWatchCh can be used
	// to watch for changes. All fields protected with urlCEBMu.
	urlCEBMu      sync.RWMutex
	urlCEB        *pb.EntrypointConfig_URLService
	urlCEBWatchCh chan struct{}

	// bgCtx is used for background tasks within the service. This is
	// cancelled when Close is called.
	bgCtx       context.Context
	bgCtxCancel context.CancelFunc
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
	if scfg := cfg.serverConfig; scfg != nil && scfg.URL != nil {
		// Set our config
		s.urlConfig = scfg.URL

		// Create a copy of the config that we use for initialization so
		// that we don't create races with s.urlConfig if this retries.
		cfgCopy := *scfg.URL

		// Initialize our CEB settings.
		s.urlCEBMu.Lock()
		s.urlCEB = &pb.EntrypointConfig_URLService{
			ControlAddr: cfgCopy.ControlAddress,
			Token:       cfgCopy.APIToken,
		}
		s.urlCEBWatchCh = make(chan struct{})
		s.urlCEBMu.Unlock()

		if err := s.initURLClient(
			log.Named("url_service"),
			false,
			cfg.acceptUrlTerms,
			&cfgCopy,
		); err != nil {
			return nil, err
		}
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

	// Setup the background context that is used for internal tasks
	s.bgCtx, s.bgCtxCancel = context.WithCancel(context.Background())

	// Start our polling background goroutine. We have a single goroutine
	// that we run in the background that handles the queue of all polling
	// operations. See the func docs for more info.
	go s.runPollQueuer(s.bgCtx, log.Named("poll_queuer"))

	return &s, nil
}

// Close shuts down any background processes and resources that may
// be used by the service. This should be called after the service
// is no longer responding to requests.
func (s *service) Close() error {
	s.bgCtxCancel()
	return nil
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
