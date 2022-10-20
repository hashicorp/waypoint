package consul

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	pb "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// retryInterval is the base retry value
	retryInterval = 5 * time.Second

	// maximum back off time, this is to prevent
	// exponential runaway
	maxBackoffTime = 180 * time.Second
)

type reqConfig struct {
	Key        string `hcl:"key,attr"`             // kv path to retrieve
	Namespace  string `hcl:"namespace,optional"`   // namespace the kv data resides within
	Partition  string `hcl:"partition,optional"`   // partition the kv data resides within
	Datacenter string `hcl:"datacenter,optional"`  // datacenter the kv data resides within
	AllowStale bool   `hcl:"allow_stale,optional"` // whether to perform stale queries against non-leader servers
}

func (r *reqConfig) cacheKey() string {
	return fmt.Sprintf("%s/%s/%s/%s", r.Datacenter, r.Partition, r.Namespace, r.Key)
}

type cachedVal struct {
	kvPair *api.KVPair
	cancel context.CancelFunc
	err    error
}

// ConfigSourcerConfig is used to configure where to talk to Consul, and from
// where the KV data is to be retrieved
type ConfigSourcerConfig struct {
	// Configuration for where to talk to Consul
	Address   string         `hcl:"address,optional"`
	Scheme    string         `hcl:"scheme,optional"`
	HTTPAuth  consulHTTPAuth `hcl:"http_auth,optional"`
	Token     string         `hcl:"token,optional"`
	TokenFile string         `hcl:"token_file,optional"`
	TLSConfig tlsConfig      `hcl:"tls,optional"`

	// Default location of KV data
	Datacenter string `hcl:"datacenter,optional"`
	Namespace  string `hcl:"namespace,optional"`
	Partition  string `hcl:"partition,optional"`
}

func (conf *ConfigSourcerConfig) client() (*api.Client, error) {
	apiConfig := api.Config{
		Address:    conf.Address,
		Scheme:     conf.Scheme,
		Datacenter: conf.Datacenter,
		HttpAuth:   conf.HTTPAuth.toApiAuth(),
		Token:      conf.Token,
		TokenFile:  conf.TokenFile,
		Namespace:  conf.Namespace,
		Partition:  conf.Partition,
		TLSConfig: api.TLSConfig{
			Address:            conf.TLSConfig.ServerName,
			CAFile:             conf.TLSConfig.CAFile,
			CAPath:             conf.TLSConfig.CAPath,
			CAPem:              conf.TLSConfig.CAPem,
			CertFile:           conf.TLSConfig.CertFile,
			CertPEM:            conf.TLSConfig.CertPEM,
			KeyFile:            conf.TLSConfig.KeyFile,
			KeyPEM:             conf.TLSConfig.KeyPEM,
			InsecureSkipVerify: conf.TLSConfig.InsecureHTTPs,
		},
	}

	return api.NewClient(&apiConfig)
}

type tlsConfig struct {
	ServerName    string `hcl:"server_name,optional"`
	CAFile        string `hcl:"ca_file,optional"`
	CAPath        string `hcl:"ca_path,optional"`
	CAPem         []byte `hcl:"ca_pem,optional"`
	CertFile      string `hcl:"cert_file,optional"`
	CertPEM       []byte `hcl:"cert_pem,optional"`
	KeyFile       string `hcl:"key_file,optional"`
	KeyPEM        []byte `hcl:"key_pem,optional"`
	InsecureHTTPs bool   `hcl:"insecure_https,optional"`
}

type consulHTTPAuth struct {
	Username string `hcl:"username,optional"`
	Password string `hcl:"password,optional"`
}

func (a *consulHTTPAuth) toApiAuth() *api.HttpBasicAuth {
	if a.Username == "" && a.Password == "" {
		return nil
	}

	return &api.HttpBasicAuth{
		Username: a.Username,
		Password: a.Password,
	}
}

type ConfigSourcer struct {
	config ConfigSourcerConfig

	mu     sync.Mutex
	client *api.KV
	cache  map[string]*cachedVal
}

// Config implements the Configurable interface
func (cs *ConfigSourcer) Config() (interface{}, error) {
	return &cs.config, nil
}

// ConfigSet implements the ConfigurableNotify interface
func (cs *ConfigSourcer) ConfigSet(config interface{}) error {
	conf, ok := config.(*ConfigSourcerConfig)
	if !ok {
		// The Waypoint SDK should ensure this never gets hit
		return status.Errorf(codes.FailedPrecondition, "expected *ConfigSourcerConfig as parameter")
	}

	// attempt to create API client
	client, err := conf.client()
	if err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	cs.client = client.KV()

	return nil
}

// ReadFunc returns the function for reading configuration.
//
// The returned function can start a background goroutine to more efficiently
// watch for changes. The entrypoint will periodically call Read to check for
// updates.
//
// If the configuration changes for any dynamic configuration variable,
// the entrypoint will call Stop followed by Read, so plugins DO NOT need
// to implement config diffing. Plugins may safely assume if Read is called
// after a Stop that the config is new, and that subsequent calls have the
// same config.
//
// Read is called for ALL defined configuration variables for this source.
// If ANY change, Stop is called followed by Read again. Only one sourcer
// is active for a set of configs.
func (cs *ConfigSourcer) ReadFunc() interface{} {
	return cs.read
}

// StopFunc returns a function for stopping configuration sourcing.
// You can return nil if stopping is not necessary or supported for
// this sourcer.
//
// The stop function should stop any background processes started with Read.
func (cs *ConfigSourcer) StopFunc() interface{} {
	return cs.stop
}

func (cs *ConfigSourcer) read(ctx context.Context, log hclog.Logger, reqs []*component.ConfigRequest) ([]*pb.ConfigSource_Value, error) {
	log.Trace("Reading KV data from Consul")
	// Set up our lock
	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Create our cache if this is our first time
	if cs.cache == nil {
		cs.cache = make(map[string]*cachedVal)
	}

	// Create our Consul API client if this is our first time
	if cs.client == nil {
		client, err := cs.config.client()
		if err != nil {
			return nil, fmt.Errorf("invalid Consul client configuration: %w", err)
		}
		cs.client = client.KV()
	}

	var results []*pb.ConfigSource_Value
	for _, req := range reqs {
		result := &pb.ConfigSource_Value{Name: req.Name}
		results = append(results, result)

		// Decode our configuration
		kvReq := reqConfig{
			// default to allowing stale queries
			AllowStale: true,
		}
		if err := mapstructure.WeakDecode(req.Config, &kvReq); err != nil {
			result.Result = &pb.ConfigSource_Value_Error{
				Error: status.New(codes.Aborted, err.Error()).Proto(),
			}

			continue
		}

		opts := &api.QueryOptions{
			Namespace:  kvReq.Namespace,
			Partition:  kvReq.Partition,
			Datacenter: kvReq.Datacenter,
			AllowStale: kvReq.AllowStale,
		}
		reqLogger := log.With("key", kvReq.Key, "stale", kvReq.AllowStale)
		if kvReq.Namespace != "" {
			reqLogger = reqLogger.With("namespace", kvReq.Namespace)
		}
		if kvReq.Partition != "" {
			reqLogger = reqLogger.With("partition", kvReq.Partition)
		}
		if kvReq.Datacenter != "" {
			reqLogger = reqLogger.With("datacenter", kvReq.Partition)
		}

		cacheVal, ok := cs.cache[kvReq.cacheKey()]
		if !ok {
			reqLogger.Trace("querying Consul KV")
			kvpair, meta, err := cs.client.Get(kvReq.Key, opts.WithContext(ctx))
			if err != nil {
				result.Result = &pb.ConfigSource_Value_Error{
					Error: status.New(codes.Aborted, err.Error()).Proto(),
				}

				continue
			}

			refreshCtx, cancel := context.WithCancel(context.Background())
			cacheVal = &cachedVal{
				kvPair: kvpair,
				cancel: cancel,
				err:    nil,
			}
			cs.cache[kvReq.cacheKey()] = cacheVal
			go cs.startConsulBlockingQuery(refreshCtx, log, kvReq, meta.LastIndex, opts)
		}

		if cacheVal.err != nil {
			result.Result = &pb.ConfigSource_Value_Error{
				Error: status.New(codes.Aborted, cacheVal.err.Error()).Proto(),
			}
		}

		if cacheVal.kvPair == nil {
			result.Result = &pb.ConfigSource_Value_Error{
				Error: status.New(codes.NotFound, fmt.Sprintf("Configuration for key %s doesn't exist", kvReq.Key)).Proto(),
			}

			continue
		}

		result.Result = &pb.ConfigSource_Value_Value{
			Value: string(cacheVal.kvPair.Value),
		}
	}

	return results, nil
}

// startConsulBlockingQuery starts a blocking query to the Consul API to poll
// for changes in the KV data store at the specified path. This is expected to
// be called via a goroutine. For more information on Consul blocking queries, see:
// https://developer.hashicorp.com/consul/api-docs/features/blocking
func (cs *ConfigSourcer) startConsulBlockingQuery(ctx context.Context, logger hclog.Logger, kvReq reqConfig, lastIndex uint64, opts *api.QueryOptions) {
	opts = opts.WithContext(ctx)
	failures := 0
	// Ideally we would use the github.com/hashicorp/consul/api/watch package. However that package doesn't support
	// namespaces and partitions except with the global client defaults and thus isn't suitable for this usage.

	for {
		// check if we are being stopped
		select {
		case <-ctx.Done():
			return
		default:
		}

		logger.Trace("Issuing blocking query", "wait-index", lastIndex)
		// set the wait index to use for the query
		opts.WaitIndex = lastIndex
		pair, meta, err := cs.client.Get(kvReq.Key, opts)

		// KV entry not updated - do nothing
		if meta != nil && meta.LastIndex == lastIndex {
			logger.Trace("KV value unchanged")
			time.Sleep(retryInterval)
			continue
		}

		// update the data within our cache
		cs.mu.Lock()
		val, ok := cs.cache[kvReq.cacheKey()]
		if !ok {
			logger.Error("KV data is not present within the cache - stopping blocking query refresher")
			cs.mu.Unlock()
			return
		}
		val.kvPair = pair
		val.err = err
		cs.mu.Unlock()

		// now determine the next wait index
		if err != nil {
			// reset the last index to 0 to do a non-blocking query
			lastIndex = 0
			// Set up for exponential backoff
			failures++

			retry := retryInterval * time.Duration(failures*failures)
			if retry > maxBackoffTime {
				retry = maxBackoffTime
			}
			logger.Error("KV Get errored", "error", err, "retry", retry.String())
			select {
			case <-time.After(retry):
			case <-ctx.Done():
				return
			}
		}
		lastIndex = meta.LastIndex
		failures = 0
	}
}

func (cs *ConfigSourcer) stop() error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Stop all our background renewers
	for _, v := range cs.cache {
		if v.cancel != nil {
			v.cancel()
		}
	}

	cs.cache = nil

	return nil
}

func (cs *ConfigSourcer) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.RequestFromStruct(&reqConfig{}))
	if err != nil {
		return nil, err
	}

	doc.Description("Read configuration values from the Consul KV store.")

	doc.SetRequestField(
		"key",
		"the KV path to retrieve",
	)

	doc.SetRequestField(
		"namespace",
		"the namespace to load the KV value from.",
		docs.Summary(
			"If not specified then it will default to the plugin's global namespace",
			"configuration. If that is also not specified then Consul will default",
			"the namespace like it would any other request.",
		),
	)

	doc.SetRequestField(
		"partition",
		"the partition to load the KV value from.",
		docs.Summary(
			"If not specified then it will default to the plugin's global partition",
			"configuration. If that is also not specified then Consul will default",
			"the partition like it would any other request.",
		),
	)

	doc.SetRequestField(
		"datacenter",
		"the datacenter to load the KV value from.",
		docs.Summary(
			"If not specified then it will default to the plugin's global datacenter",
			"configuration. If that is also not specified then Consul will default",
			"the datacenter like it would any other request.",
		),
	)

	doc.SetRequestField(
		"allow_stale",
		"whether to perform a stale query for retrieving the KV data",
		docs.Summary(
			"If not set this will default to true. It must explicitly be set to false",
			"in order to use consistent queries.",
		),
	)

	return doc, nil
}
