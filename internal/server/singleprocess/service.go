package singleprocess

import (
	"context"
	"sync"

	"github.com/hashicorp/go-hclog"
	bolt "go.etcd.io/bbolt"

	wphznpb "github.com/hashicorp/waypoint-hzn/pkg/pb"
	"github.com/hashicorp/waypoint/internal/server/singleprocess/state"
	"github.com/hashicorp/waypoint/internal/serverconfig"
	wpoidc "github.com/hashicorp/waypoint/pkg/auth/oidc"
	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

// service implements the gRPC service for the server.
type service struct {
	// state is the state management interface that provides functions for
	// safely mutating server state.
	state serverstate.Interface

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

	// bgWg is incremented for every background goroutine that the
	// service starts up. When Close is called, we wait on this to ensure
	// that we fully shut down before returning.
	bgWg sync.WaitGroup

	// superuser is true if all API actions should act as if a superuser
	// made them. This is used for local mode only.
	superuser bool

	// oidcCache is the cache for OIDC providers.
	oidcCache *wpoidc.ProviderCache
}

// New returns a Waypoint server implementation that uses BotlDB plus
// in-memory locks to operate safely.
func New(opts ...Option) (pb.WaypointServer, error) {
	var s service
	s.oidcCache = wpoidc.NewProviderCache()

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
			nil,
			cfg.acceptUrlTerms,
			&cfgCopy,
		); err != nil {
			return nil, err
		}
	}

	// If we haven't initialized our server config before, do that once.
	conf, err := s.state.ServerConfigGet()
	if err != nil {
		return nil, err
	}
	if conf.Cookie == "" {
		if err := s.state.ServerConfigSet(conf); err != nil {
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

	// TODO: When more items are added, move this else where
	// pollableItems is a map of potential items Waypoint can queue a poll for.
	// Each item should implement the pollHandler interface
	pollableItems := map[string]pollHandler{
		"project":                  &projectPoll{state: s.state},
		"application_statusreport": &applicationPoll{state: s.state, workspace: "default"},
	}

	// Start our polling background goroutines.
	// See the func docs for more info.
	for pollName, pollItem := range pollableItems {
		s.bgWg.Add(1)
		go s.runPollQueuer(
			s.bgCtx, &s.bgWg, pollItem,
			log.Named("poll_queuer").Named(pollName),
		)
	}

	// Start out state pruning background goroutine. This calls
	// Prune on the state every 10 minutes.
	s.bgWg.Add(1)
	go s.runPrune(s.bgCtx, &s.bgWg, log.Named("prune"))

	return &s, nil
}

// State implements pkg/server/handlers/service Service
func (s *service) State(ctx context.Context) serverstate.Interface {
	return s.state
}

func (s *service) SuperUser() bool {
	return s.superuser
}

func (s *service) DecodeId(id string) (string, error) {
	return id, nil
}

func (s *service) EncodeId(ctx context.Context, id string) string {
	return id
}

// Close shuts down any background processes and resources that may
// be used by the service. This should be called after the service
// is no longer responding to requests.
func (s *service) Close() error {
	s.bgCtxCancel()
	s.bgWg.Wait()
	s.oidcCache.Close()
	return nil
}

type config struct {
	db           *bolt.DB
	serverConfig *serverconfig.Config
	log          hclog.Logger
	superuser    bool

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

// WithSuperuser forces all API actions to behave as if a superuser
// made them. This is usually turned on for local mode only. There is no
// option (at the time of writing) to enable this on a network-attached server.
func WithSuperuser() Option {
	return func(s *service, cfg *config) error {
		s.superuser = true
		return nil
	}
}

// WithAcceptURLTerms will set the config to either accept or reject the terms
// of service for using the URL service. Rejecting the TOS will disable the
// URL service. Note that the actual rejection does not occur until the
// waypoint horizon client attempts to register its guest account.
func WithAcceptURLTerms(accept bool) Option {
	return func(s *service, cfg *config) error {
		cfg.acceptUrlTerms = accept
		return nil
	}
}

var _ pb.WaypointServer = (*service)(nil)
