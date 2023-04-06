// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package k8s

import (
	"context"
	base64 "encoding/base64"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
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

// ConfigSourcer implements component.ConfigSourcer for K8s
type ConfigSourcer struct {
	cacheMu     sync.Mutex
	secretCache map[string]*cachedSecret
	lastRead    time.Time
}

type cachedSecret struct {
	Data   interface{} // either a v1.ConfigMap or v1.Secret
	Cancel func()      // Non-nil to cancel the watcher
	Err    error       // Error on last renew
}

// ReadFunc implements component.ConfigSourcer
func (cs *ConfigSourcer) ReadFunc() interface{} {
	return cs.read
}

// StopFunc implements component.ConfigSourcer
func (cs *ConfigSourcer) StopFunc() interface{} {
	return cs.stop
}

func (cs *ConfigSourcer) read(
	ctx context.Context,
	log hclog.Logger,
	reqs []*component.ConfigRequest,
) ([]*pb.ConfigSource_Value, error) {
	// Setup our lock
	cs.cacheMu.Lock()
	defer cs.cacheMu.Unlock()

	// Create our cache if this is our first time
	if cs.secretCache == nil {
		cs.secretCache = map[string]*cachedSecret{}
	}

	clientset, ns, _, err := ClientsetInCluster()
	if err != nil {
		return nil, err
	}

	var results []*pb.ConfigSource_Value
	for _, req := range reqs {
		result := &pb.ConfigSource_Value{Name: req.Name}
		results = append(results, result)

		// Decode our configuration
		var k8sReq reqConfig
		if err := mapstructure.WeakDecode(req.Config, &k8sReq); err != nil {
			result.Result = &pb.ConfigSource_Value_Error{
				Error: status.New(codes.Aborted, err.Error()).Proto(),
			}

			continue
		}
		if k8sReq.Namespace == "" {
			k8sReq.Namespace = ns
		}
		L := log.With("name", k8sReq.Name, "key", k8sReq.Key, "secret", k8sReq.Secret)

		// Get this config or read it if we haven't already.
		cachedSecretVal, ok := cs.secretCache[k8sReq.CacheKey()]
		if !ok {
			L.Trace("querying K8S configuration")
			data, err := k8sReq.Get(ctx, clientset)
			if err != nil {
				result.Result = &pb.ConfigSource_Value_Error{
					Error: status.New(codes.Aborted, err.Error()).Proto(),
				}

				continue
			}

			// Store our initial value
			cachedSecretVal = &cachedSecret{Data: data}
			cs.secretCache[k8sReq.CacheKey()] = cachedSecretVal

			// Start refresher
			cs.startRefresher(clientset, &k8sReq)
		}

		// If the secret has an error, return that
		if err := cachedSecretVal.Err; err != nil {
			result.Result = &pb.ConfigSource_Value_Error{
				Error: status.New(codes.Aborted, err.Error()).Proto(),
			}

			continue
		}

		var value string
		switch d := cachedSecretVal.Data.(type) {
		case *corev1.ConfigMap:
			value, ok = d.Data[k8sReq.Key]

		case *corev1.Secret:
			var secretValue []byte
			secretValue, ok = d.Data[k8sReq.Key]

			if ok {
				// Encode secretValue byte array into string so that it can be decoded
				// and returned as string for k8s config
				sEnc := base64.StdEncoding.EncodeToString(secretValue)
				decValue, err := base64.StdEncoding.DecodeString(sEnc)
				if err != nil {
					L.Trace("failed to decode secret: ", err)
					result.Result = &pb.ConfigSource_Value_Error{
						Error: status.New(codes.Aborted, err.Error()).Proto(),
					}

					// break from outer loop early since we err decoding key, but key was found in Secret
					continue
				}

				value = string(decValue)
			}

		default:
			ok = false
		}
		if !ok {
			result.Result = &pb.ConfigSource_Value_Error{
				Error: status.New(codes.NotFound, "key not found: "+k8sReq.Key).Proto(),
			}

			continue
		}

		result.Result = &pb.ConfigSource_Value_Value{
			Value: value,
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

func (cs *ConfigSourcer) startRefresher(
	clientset *kubernetes.Clientset,
	req *reqConfig,
) {
	// The secret should be in the cache. If it isn't then just ignore.
	// The reason it should be in the cache is because we only call startRenewer
	// after querying the initial secret and inserting it into the cache.
	key := req.CacheKey()
	cachedVal, ok := cs.secretCache[key]
	if !ok {
		return
	}

	// Setup our cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cachedVal.Cancel = cancel

	// Start goroutine that actually refreshes the data. NOTE: this doesn't
	// use actual long-polling APIs that K8S provides currently, we can do
	// that in the future as an improvement.
	go func() {
		// Calculate a sleep period with a 30% jitter added to it.
		const factor = 0.5
		min := int64(math.Floor(float64(refreshPeriod) * (1 - factor)))
		max := int64(math.Ceil(float64(refreshPeriod) * (1 + factor)))

		for {
			// Calculate our sleep period. We add a jitter to it to prevent
			// applications that all started at the same time to stampede
			// dynamic sources.
			refreshDur := time.Duration(rand.Int63n(max-min) + min)

			select {
			case <-ctx.Done():
				return

			case <-time.After(refreshDur):
			}

			// Read our value
			data, err := req.Get(ctx, clientset)

			// Update our value
			cs.cacheMu.Lock()
			value, ok := cs.secretCache[key]
			if !ok {
				cs.cacheMu.Unlock()
				return
			}
			value.Data = data
			value.Err = err
			cs.cacheMu.Unlock()
		}
	}()
}

func (cs *ConfigSourcer) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.RequestFromStruct(&reqConfig{}))
	if err != nil {
		return nil, err
	}

	doc.Description("Read configuration values from Kubernetes ConfigMap or Secret resources. " +
		"Note that to read a config value from a Secret, you must set `secret = true`. Otherwise " +
		"Waypoint will load a dynamic value from a ConfigMap.")

	doc.Example(`
config {
  env = {
    PORT = dynamic("kubernetes", {
	  name = "my-config-map"
	  key = "port"
	})

    DATABASE_PASSWORD = dynamic("kubernetes", {
	  name = "database-creds"
	  key = "password"
	  secret = true
	})
  }
}
`)

	doc.SetRequestField(
		"name",
		"the name of the ConfigMap of Secret",
	)

	doc.SetRequestField(
		"namespace",
		"the namespace to load the ConfigMap or Secret from.",
		docs.Summary(
			"by default this will use the namespace of the running pod.",
			"If this config source is used outside of a pod, this will use the",
			"namespace from the kubeconfig.",
		),
	)

	doc.SetRequestField(
		"key",
		"the key in the ConfigMap or Secret to read the value from",
		docs.Summary(
			"ConfigMaps and Secrets store data in key/value format. This specifies",
			"the key to read from the resource. If you want multiple values you must",
			"specify multiple dynamic values.",
		),
	)

	doc.SetRequestField(
		"secret",
		"This must be set to true to read from a Secret. If it is false we read from a ConfigMap.",
	)

	return doc, nil
}

type reqConfig struct {
	Name      string `hcl:"name,attr"`          // config map name
	Namespace string `hcl:"namespace,optional"` // namespace for the config
	Key       string `hcl:"key,attr"`           // key in the config map to read
	Secret    bool   `hcl:"secret,optional"`    // true if this is a secret (not a configmap)
}

func (c *reqConfig) Get(ctx context.Context, clientset *kubernetes.Clientset) (interface{}, error) {
	if c.Secret {
		return clientset.CoreV1().Secrets(c.Namespace).Get(
			ctx, c.Name, metav1.GetOptions{})
	} else {
		return clientset.CoreV1().ConfigMaps(c.Namespace).Get(
			ctx, c.Name, metav1.GetOptions{})
	}
}

func (c *reqConfig) CacheKey() string {
	if c.Secret {
		return "secret/" + c.Name
	}

	return "config/" + c.Name
}
