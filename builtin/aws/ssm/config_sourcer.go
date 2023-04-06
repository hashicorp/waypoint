// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
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
		log.Trace("returning cached values", "len", len(cs.values))
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
	log.Debug("retrieving AWS session")
	sess, err := awsbase.GetSession(&awsConfig)
	if err != nil {
		log.Warn("error initializing AWS session", "err", err)
		return nil, err
	}
	ssmsvc := ssm.New(sess)
	log.Debug("AWS session initialized")

	// If this is our first read after a stop, then we need to populate
	// the requests we want to get. This doesn't actually fetch anything yet.
	cs.values = map[string]*pb.ConfigSource_Value{}
	requests := map[string]*reqConfig{}
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
		requests[req.Name] = &decodedReq
	}

	// Start the refresher. Note for the log parameter we can't use "log"
	// because that is only alive for the duration of the plugin RPC call.
	// Once the RPC completes we'll get errors. So we instantiate a new
	// logger here which will go to the plugin process stderr.
	refreshCtx, cancel := context.WithCancel(context.Background())
	cs.refreshCancel = cancel
	initCh := make(chan struct{})
	go cs.startRefresher(refreshCtx,
		hclog.L().Named("ssm-refresher"),
		initCh, requests, cs.values, ssmsvc)

	// Unlock so that the refresher can populate.
	log.Debug("waiting for first set of values")
	cs.cacheMu.Unlock()
	<-initCh
	cs.cacheMu.Lock()

	// Return our results
	log.Debug("first set of values received, returning")
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
	log hclog.Logger,
	initCh chan<- struct{},
	wpreqs map[string]*reqConfig,
	values map[string]*pb.ConfigSource_Value,
	client *ssm.SSM,
) {
	// We need to keep a map of our parameter names to the
	// value keys that require it.
	paramToReq := map[string][]string{}

	// Build our actual request parameter. This remains constant so
	// we calculate this once.
	req := &ssm.GetParametersInput{}
	req.SetWithDecryption(true)
	for k, wpreq := range wpreqs {
		name := wpreq.Path
		req.Names = append(req.Names, &name)
		paramToReq[name] = append(paramToReq[name], k)
	}

	// Calculate a sleep period with a 30% jitter added to it.
	const factor = 0.5
	min := int64(math.Floor(float64(refreshPeriod) * (1 - factor)))
	max := int64(math.Ceil(float64(refreshPeriod) * (1 + factor)))

	for {
		// Read our value
		log.Trace("querying parameters")
		resp, err := client.GetParameters(req)
		if err != nil {
			// Just skip these errors and reuse the last known values for
			// all requests that we have.
			log.Warn("error querying parameters", "err", err)
			continue
		}
		log.Trace("parameters received", "len", len(resp.Parameters))

		// Update our values
		cs.cacheMu.Lock()
		for _, param := range resp.Parameters {
			if param.Name == nil || param.Value == nil {
				// should never happen
				log.Warn("param name or value is nil", "param", param)
				continue
			}

			for _, k := range paramToReq[*param.Name] {
				v, ok := values[k]
				if !ok {
					// skip params we don't know, since we should prepopulate it all.
					log.Warn("param not found in values map, should never happen",
						"key", *param.Name)
					continue
				}

				v.Result = &pb.ConfigSource_Value_Value{Value: *param.Value}
			}
		}
		cs.cacheMu.Unlock()

		// For our first request, we close the init channel.
		if initCh != nil {
			log.Trace("closing init channel")
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

func (cs *ConfigSourcer) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(
		docs.FromConfig(&sourceConfig{}),
		docs.RequestFromStruct(&reqConfig{}),
	)
	if err != nil {
		return nil, err
	}

	doc.Description("Read configuration values from AWS SSM Parameter Store.")

	doc.Example(`
config {
  env = {
    PORT = dynamic("aws-ssm", {
	  path = "port"
	})
  }
}
`)

	doc.SetRequestField(
		"path",
		"the path for the parameter to read from the parameter store.",
	)

	doc.SetField(
		"access_key",
		"This is the AWS access key. It must be provided, but it can also be sourced from the `AWS_ACCESS_KEY_ID` environment variable, or via a shared credentials file if `profile` is specified",
	)

	doc.SetField(
		"assume_role_arn",
		"Amazon Resource Name (ARN) of the IAM Role to assume.",
	)

	doc.SetField(
		"assume_role_duration_seconds",
		"Number of seconds to restrict the assume role session duration.",
	)

	doc.SetField(
		"assume_role_external_id",
		"External identifier to use when assuming the role.",
	)

	doc.SetField(
		"assume_role_policy",
		"IAM Policy JSON describing further restricting permissions for the IAM Role being assumed.",
	)

	doc.SetField(
		"assume_role_session_name",
		"Session name to use when assuming the role.",
	)

	doc.SetField(
		"shared_credentials_file",
		"This is the path to the shared credentials file. If this is not set and a profile is specified, `~/.aws/credentials` will be used.",
	)

	doc.SetField(
		"iam_endpoint",
		"Custom endpoint address for the IAM service.",
	)

	doc.SetField(
		"insecure",
		"Explicitly allow the provider to perform \"insecure\" SSL requests.",
		docs.Default("false"),
	)

	doc.SetField(
		"max_retries",
		"This is the maximum number of times an API call is retried, in the case where requests are being throttled or experiencing transient failures. The delay between the subsequent API calls increases exponentially.",
		docs.Default("25"),
	)

	doc.SetField(
		"profile",
		"This is the AWS profile name as set in the shared credentials file.",
	)

	doc.SetField(
		"region",
		"This is the AWS region. It must be provided, but it can also be sourced from the `AWS_DEFAULT_REGION` environment variables, or via a shared credentials file if profile is specified.",
	)

	doc.SetField(
		"secret_key",
		"This is the AWS secret key. It must be provided, but it can also be sourced from the `AWS_SECRET_ACCESS_KEY` environment variable, or via a shared credentials file if `profile` is specified.",
	)

	doc.SetField(
		"skip_credentials_validation",
		"Skip the credentials validation via the STS API. Useful for AWS API implementations that do not have STS available or implemented.",
	)

	doc.SetField(
		"skip_metadata_api_check",
		"Skip the AWS Metadata API check. Useful for AWS API implementations that do not have a metadata API endpoint. Setting to true prevents Terraform from authenticating via the Metadata API. You may need to use other authentication methods like static credentials, configuration variables, or environment variables.",
	)

	doc.SetField(
		"skip_requesting_account_id",
		"Skip requesting the account ID. Useful for AWS API implementations that do not have the IAM, STS API, or metadata API.",
	)

	doc.SetField(
		"sts_endpoint",
		"Custom endpoint for the STS service.",
	)

	return doc, nil
}

type reqConfig struct {
	Path string `hcl:"path,attr"` // parameter path
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
	CredsFilename           string `hcl:"shared_credentials_file,optional"`
	IamEndpoint             string `hcl:"iam_endpoint,optional"`
	Insecure                bool   `hcl:"insecure,optional"`
	MaxRetries              int    `hcl:"max_retries,optional"`
	Profile                 string `hcl:"profile,optional"`
	Region                  string `hcl:"region,optional"`
	SecretKey               string `hcl:"secret_key,optional"`
	SkipCredsValidation     bool   `hcl:"skip_credentials_validation,optional"`
	SkipMetadataApiCheck    bool   `hcl:"skip_metadata_api_check,optional"`
	SkipRequestingAccountId bool   `hcl:"skip_requesting_account_id,optional"`
	StsEndpoint             string `hcl:"sts_endpoint,optional"`
	Token                   string `hcl:"token,optional"`
}
