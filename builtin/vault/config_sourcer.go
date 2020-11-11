package vault

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/mitchellh/mapstructure"
	"github.com/mitchellh/pointerstructure"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	pb "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
)

var (
	// refreshPeriod is the interval between refreshing secret values that
	// aren't renewable or have a lease associated with them. If a Read is
	// called again during this period, we will return cached values.
	//
	// We have to do this because Vault doesn't support any sort of blocking
	// queries so unless it tells us a lease period, we can't safely refresh.
	refreshPeriod = 30 * time.Second
)

// ConfigSourcer implements component.ConfigSourcer for Vault
type ConfigSourcer struct {
	// Client, if set, will be used as the client instead of initializing
	// based on the config. This is only used for tests.
	Client *vaultapi.Client

	cacheMu     sync.Mutex
	secretCache map[string]*cachedSecret
	lastRead    time.Time
}

type cachedSecret struct {
	Secret *vaultapi.Secret // The secret itself
	Cancel func()           // Non-nil to cancel the renewer
	Err    error            // Error on last renew
}

// ReadFunc implements component.ConfigSourcer
func (cs *ConfigSourcer) ReadFunc() interface{} {
	return cs.read
}

// StopFunc implements component.ConfigSourcer
func (cs *ConfigSourcer) StopFunc() interface{} {
	return nil
}

func (cs *ConfigSourcer) read(
	ctx context.Context,
	log hclog.Logger,
	reqs []*component.ConfigRequest,
) ([]*pb.ConfigSource_Value, error) {
	// Setup our lock
	cs.cacheMu.Lock()
	defer cs.cacheMu.Unlock()

	// If we have a last read value and its before our refresh period, we
	// just returned cached values. Note that cached values may still be
	// updated in the background for secrets that have leases. Additionally,
	// when Stop is called, we reset lastRead to zero.
	if cs.lastRead.IsZero() || time.Now().Sub(cs.lastRead) > refreshPeriod {
		log.Trace("purging cached secrets that aren't renewable")
		for k, s := range cs.secretCache {
			if s.Cancel == nil {
				delete(cs.secretCache, k)
			}
		}

		cs.lastRead = time.Now()
	}

	// Create our cache if this is our first time
	if cs.secretCache == nil {
		cs.secretCache = map[string]*cachedSecret{}
	}

	// Initialize our Vault client
	client := cs.Client
	if client == nil {
		log.Debug("initializing the Vault client")
		clientConfig := vaultapi.DefaultConfig()
		err := clientConfig.ReadEnvironment()
		if err != nil {
			return nil, err
		}

		client, err = vaultapi.NewClient(clientConfig)
		if err != nil {
			return nil, err
		}
	} else {
		log.Debug("using preconfigured client on struct")
	}

	// Go through each request and read it. The way this generally works:
	// If the variable is not in our cache, we re-read it from Vault. In the
	// above where we purge the cache, we keep any with Cancel set. This keeps
	// long-running dynamic secrets around so that they don't flap every refresh
	// period. Instead, those are still in the cache and we use whatever value
	// they have. A background goroutine will update those (see startRenewer).
	//
	// If a config change happens, the ConfigSourcer contract states that
	// Stop will be called. When Stop is called, we clear our full cache and
	// stop all renewers.
	//
	// Therefore, in most cases, this is re-reading static values from Vault
	// and just loading cached dynamic values.
	var results []*pb.ConfigSource_Value
	for _, req := range reqs {
		result := &pb.ConfigSource_Value{Name: req.Name}
		results = append(results, result)

		// Decode our configuration
		var vaultReq reqConfig
		if err := mapstructure.WeakDecode(req.Config, &vaultReq); err != nil {
			result.Result = &pb.ConfigSource_Value_Error{
				Error: status.New(codes.Aborted, err.Error()).Proto(),
			}

			continue
		}
		L := log.With("path", vaultReq.Path, "key", vaultReq.Key)

		// Get this secret or read it if we haven't already.
		cachedSecretVal, ok := cs.secretCache[vaultReq.Path]
		if !ok {
			L.Trace("querying Vault secret")
			secret, err := client.Logical().Read(vaultReq.Path)
			if err != nil {
				result.Result = &pb.ConfigSource_Value_Error{
					Error: status.New(codes.Aborted, err.Error()).Proto(),
				}

				continue
			}
			cachedSecretVal = &cachedSecret{Secret: secret}
			cs.secretCache[vaultReq.Path] = cachedSecretVal

			// If this secret is renewable, we will start a background
			// renewer to watch it. This more efficiently updates this secret
			// and prevents flapping values on every refresh.
			if secret.Renewable {
				L.Debug("secret is renewable, starting renewer")
				cs.startRenewer(client, vaultReq.Path, secret)
			}
		}

		// If the secret has an error, return that
		if err := cachedSecretVal.Err; err != nil {
			result.Result = &pb.ConfigSource_Value_Error{
				Error: status.New(codes.Aborted, err.Error()).Proto(),
			}

			continue
		}

		// Get the value
		if !strings.HasPrefix(vaultReq.Key, "/") {
			vaultReq.Key = "/" + vaultReq.Key
		}
		value, err := pointerstructure.Get(cachedSecretVal.Secret.Data, vaultReq.Key)
		if err != nil {
			result.Result = &pb.ConfigSource_Value_Error{
				Error: status.New(codes.Aborted, err.Error()).Proto(),
			}

			continue
		}

		// Convert it to a string
		var valueStr string
		if err := mapstructure.WeakDecode(value, &valueStr); err != nil {
			L.Warn("vault secret value couldn't be converted to string")
			result.Result = &pb.ConfigSource_Value_Error{
				Error: status.New(codes.Aborted, err.Error()).Proto(),
			}

			continue
		}

		result.Result = &pb.ConfigSource_Value_Value{
			Value: valueStr,
		}
	}

	return results, nil
}

func (cs *ConfigSourcer) stop() error {
	cs.cacheMu.Lock()
	defer cs.cacheMu.Unlock()

	// Stop all our background renewers
	for _, s := range cs.secretCache {
		if s.Cancel != nil {
			s.Cancel()
		}
	}

	// Reset our results tracking to empty. This will force the next call
	// to rebuild all our secret values.
	var zeroTime time.Time
	cs.lastRead = zeroTime
	cs.secretCache = nil

	return nil
}

func (cs *ConfigSourcer) startRenewer(client *vaultapi.Client, path string, s *vaultapi.Secret) {
	// The secret should be in the cache. If it isn't then just ignore.
	// The reason it should be in the cache is because we only call startRenewer
	// after querying the initial secret and inserting it into the cache.
	cache, ok := cs.secretCache[path]
	if !ok {
		return
	}

	renewer, err := client.NewRenewer(&vaultapi.RenewerInput{
		Secret: cache.Secret,
	})
	if err != nil {
		cache.Err = err
		return
	}

	// Start the renewer in the background
	renewer.Renew()

	// Create our cancellation context
	ctx, cancel := context.WithCancel(context.Background())
	cache.Cancel = cancel

	// Start our goroutine that actually watches for changes. This
	// goroutine can no longer assume the "cache" variable is safe for
	// reading or writing and must acquire a lock.
	go func() {
		defer renewer.Stop()

		for {
			var newVal cachedSecret
			select {
			case <-ctx.Done():
				// If we're canceled, we assume something else is handling
				// our cleanup and values and so on so just exit.
				return

			case err := <-renewer.DoneCh():
				// Error during renew, mark the error value and exit.
				newVal.Err = err

			case renew := <-renewer.RenewCh():
				// Successful renewal, store the secret
				newVal.Secret = renew.Secret
			}

			// Grab a lock to update our value
			cs.cacheMu.Lock()

			value, ok := cs.secretCache[path]
			if !ok {
				// Shouldn't happen, exit.
				cs.cacheMu.Unlock()
				return
			}

			if newVal.Err != nil {
				value.Err = newVal.Err
			}
			if newVal.Secret != nil {
				value.Secret = newVal.Secret
			}

			cs.cacheMu.Unlock()
		}
	}()
}

type reqConfig struct {
	Path string
	Key  string
}
