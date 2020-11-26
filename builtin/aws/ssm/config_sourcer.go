package ssm

import (
	"context"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/service/ssm"
	awsbase "github.com/hashicorp/aws-sdk-go-base"
	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	pb "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
)

var (
	// refreshPeriod is the interval between refreshing secret values that
	// aren't renewable or have a lease associated with them. If a Read is
	// called again during this period, we will return cached values.
	refreshPeriod = 30 * time.Second
)

// ConfigSourcer implements component.ConfigSourcer for K8s
type ConfigSourcer struct {
	config        sourceConfig
	cacheMu       sync.Mutex
	values        map[string]*pb.ConfigSource_Value
	refreshCancel func()
}

// Config implements component.Configurable
func (cs *ConfigSourcer) Config() (interface{}, error) {
	return &cs.config, nil
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

	// If we have cached values just return those
	if cs.values != nil {
		result := make([]*pb.ConfigSource_Value, 0, len(cs.values))
		for _, v := range cs.values {
			result = append(result, v)
		}

		return result, nil
	}

	// Decode our aws config
	var awsConfig awsbase.Config
	if err := mapstructure.WeakDecode(cs.config, &awsConfig); err != nil {
		log.Warn("error decoding the config source config", "err", err)
		return nil, err
	}
	awsConfig.CallerName = "Waypoint"
	awsConfig.CallerDocumentationURL = "https://www.waypointproject.io/"

	// Get our session
	sess, err := awsbase.GetSession(&awsConfig)
	if err != nil {
		log.Warn("error initializing AWS session", "err", err)
		return nil, err
	}
	ssmsvc := ssm.New(sess)

	// If this is our first read after a stop, then we need to populate
	// the requests we want to get. This doesn't actually fetch anything yet.
	cs.values = map[string]*pb.ConfigSource_Value{}
	requests := make([]*reqConfig, 0, len(reqs))
	for _, req := range reqs {
		result := &pb.ConfigSource_Value{Name: req.Name}
		cs.values[req.Name] = result

		// Decode our configuration
		var decodedReq reqConfig
		if err := mapstructure.WeakDecode(req.Config, &decodedReq); err != nil {
			result.Result = &pb.ConfigSource_Value_Error{
				Error: status.New(codes.Aborted, err.Error()).Proto(),
			}

			continue
		}

		// Store our request
		requests = append(requests, &decodedReq)
	}

	// Start the refresher
	refreshCtx, cancel := context.WithCancel(context.Background())
	cs.refreshCancel = cancel
	initCh := make(chan struct{})
	go cs.startRefresher(refreshCtx, initCh, requests, cs.values, ssmsvc)

	// Unlock so that the refresher can populate.
	cs.cacheMu.Unlock()
	<-initCh
	cs.cacheMu.Lock()

	// Return our results
	result := make([]*pb.ConfigSource_Value, 0, len(cs.values))
	for _, v := range cs.values {
		result = append(result, v)
	}

	return result, nil
}

func (cs *ConfigSourcer) stop() error {
	cs.cacheMu.Lock()
	defer cs.cacheMu.Unlock()

	// Nullify everything which will force a refresh and restart
	cs.values = nil

	return nil
}

func (cs *ConfigSourcer) startRefresher(
	ctx context.Context,
	initCh chan<- struct{},
	wpreqs []*reqConfig,
	values map[string]*pb.ConfigSource_Value,
	client *ssm.SSM,
) {
	// Build our actual request parameter. This remains constant so
	// we calculate this once.
	req := &ssm.GetParametersInput{}
	req.SetWithDecryption(true)
	for _, wpreq := range wpreqs {
		name := wpreq.Path
		req.Names = append(req.Names, &name)
	}

	// Calculate a sleep period with a 30% jitter added to it.
	const factor = 0.5
	min := int64(math.Floor(float64(refreshPeriod) * (1 - factor)))
	max := int64(math.Ceil(float64(refreshPeriod) * (1 + factor)))

	for {
		// Read our value
		resp, err := client.GetParameters(req)
		if err != nil {
			// Just skip these errors and reuse the last known values for
			// all requests that we have.
			continue
		}

		// Update our values
		cs.cacheMu.Lock()
		for _, param := range resp.Parameters {
			if param.Name == nil || param.Value == nil {
				// should never happen
				continue
			}

			v, ok := values[*param.Name]
			if !ok {
				// skip params we don't know, since we should prepopulate it all.
				continue
			}

			v.Result = &pb.ConfigSource_Value_Value{Value: *param.Value}
		}
		cs.cacheMu.Unlock()

		// For our first request, we close the init channel.
		if initCh != nil {
			close(initCh)
			initCh = nil
		}

		// Calculate our sleep period. We add a jitter to it to prevent
		// applications that all started at the same time to stampede
		// dynamic sources.
		refreshDur := time.Duration(rand.Int63n(max-min) + min)

		select {
		case <-ctx.Done():
			return

		case <-time.After(refreshDur):
		}
	}
}

type reqConfig struct {
	Path string // parameter path
}

func (c *reqConfig) CacheKey() string {
	return c.Path
}

type sourceConfig struct {
	// The fields commented out below can't be supported yet cleanly
	// because we only allow primitive value types (non-containers) for
	// config source configs.

	AccessKey                 string `hcl:"access_key,optional"`
	AssumeRoleARN             string `hcl:"assume_role_arn,optional"`
	AssumeRoleDurationSeconds int    `hcl:"assume_role_duration_seconds,optional"`
	AssumeRoleExternalID      string `hcl:"assume_role_external_id,optional"`
	AssumeRolePolicy          string `hcl:"assume_role_policy,optional"`
	//AssumeRolePolicyARNs        []string          `hcl:"assume_role_policy_arns,optional"`
	AssumeRoleSessionName string `hcl:"assume_role_session_name,optional"`
	//AssumeRoleTags              map[string]string `hcl:"assume_role_tags,optional"`
	//AssumeRoleTransitiveTagKeys []string          `hcl:"assume_role_transitive_tag_keys,optional"`
	CredsFilename           string `hcl:"creds_filename,optional"`
	IamEndpoint             string `hcl:"iam_endpoint,optional"`
	Insecure                bool   `hcl:"insecure,optional"`
	MaxRetries              int    `hcl:"max_retries,optional"`
	Profile                 string `hcl:"profile,optional"`
	Region                  string `hcl:"region,optional"`
	SecretKey               string `hcl:"secret_key,optional"`
	SkipCredsValidation     bool   `hcl:"skip_creds_validation,optional"`
	SkipMetadataApiCheck    bool   `hcl:"skip_metadata_api_check,optional"`
	SkipRequestingAccountId bool   `hcl:"skip_requesting_account_id,optional"`
	StsEndpoint             string `hcl:"sts_endpoint,optional"`
	Token                   string `hcl:"token,optional"`
}
