package tfc

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

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
	//
	// We have to do this because TFC doesn't support any sort of blocking
	// queries.
	//
	// This value is a default, the user can change the it with the refresh
	// config attribute.
	refreshPeriod = 10 * time.Minute

	// The shortest refresh period we allow the user to set.
	minimumRefreshPeriod = 1 * time.Minute
)

const DefaultBaseURL = "https://app.terraform.io"

// ConfigSourcer implements component.ConfigSourcer for Terraform Cloud
type ConfigSourcer struct {
	// BaseURL is the protocol, host, and port that are used to compose
	// the request to TFC. For example, the default is https://app.terraform.io

	client *http.Client

	config      sourceConfig
	cacheMu     sync.Mutex
	secretCache map[string]*cachedSecret
	lastRead    time.Time

	workspaceIds  map[string]string
	refreshPeriod time.Duration
}

type cachedSecret struct {
	Outputs map[string]string // The most recent outputs for a workspace
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

// wsLookup is used to read the information about the workspace. We only
// need the Id.
type wsLookup struct {
	Data struct {
		Id string `json:"id"`
	} `json:"data"`
}

// stateInclude is the structure of the included resources in the state version
// API response. We expect to only see ones with Type equal to "state-version-outputs"
type stateInclude struct {
	Id   string `json:"id"`
	Type string `json:"type"`

	Attributes struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"attributes"`
}

// stateLookup is used to read the information about the current state version.
// Presently we only read the Included data as it's the outputs we want.
type stateLookup struct {
	Data struct {
		Id string `json:"id"`
	} `json:"data"`

	Included []*stateInclude `json:"included"`
}

func (cs *ConfigSourcer) read(
	ctx context.Context,
	log hclog.Logger,
	reqs []*component.ConfigRequest,
) ([]*pb.ConfigSource_Value, error) {
	// Setup our lock
	cs.cacheMu.Lock()
	defer cs.cacheMu.Unlock()

	// Setup refresh period since there isn't a place to do it after
	// config parsing atm.
	if cs.refreshPeriod == 0 {
		if cs.config.RefreshInterval != "" {
			dur, err := time.ParseDuration(cs.config.RefreshInterval)
			if err == nil {
				if dur < minimumRefreshPeriod {
					dur = minimumRefreshPeriod
				}

				cs.refreshPeriod = dur
			} else {
				log.Warn("Unable to parse refresh_interval value: %s. Using default value.", err)
			}
		}

		if cs.refreshPeriod == 0 {
			cs.refreshPeriod = refreshPeriod
		}

		log.Debug("refresh period for TFC", "period", cs.refreshPeriod.String())
	}

	if cs.client == nil {
		cs.client = &http.Client{}

		if cs.config.SkipVerify {
			cs.client.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: cs.config.SkipVerify,
				},
			}
		}
	}

	baseURL := cs.config.BaseURL
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	// If we have a last read value and its before our refresh period, we
	// just returned cached values. So when the time is up, we just nuke
	// the cache to force ourselves to recalculate.
	if cs.lastRead.IsZero() || time.Now().Sub(cs.lastRead) > cs.refreshPeriod {
		log.Trace("purging cached secrets")
		cs.secretCache = nil
		cs.lastRead = time.Now()
	}

	if cs.secretCache == nil {
		cs.secretCache = map[string]*cachedSecret{}
	}

	// Go through each request and read it. The way this generally works:
	// If the variable is not in our cache, we read it from TFC.
	//
	// If a config change happens, the ConfigSourcer contract states that
	// Stop will be called. When Stop is called, we clear our full cache and
	// stop all renewers.
	//
	// In most cases, this returns cached values due to the long
	// default refresh period.
	var results []*pb.ConfigSource_Value
	for _, req := range reqs {
		result := &pb.ConfigSource_Value{Name: req.Name}
		results = append(results, result)

		// Decode our configuration
		var tfcReq reqConfig
		if err := mapstructure.WeakDecode(req.Config, &tfcReq); err != nil {
			result.Result = &pb.ConfigSource_Value_Error{
				Error: status.New(codes.Aborted, err.Error()).Proto(),
			}

			continue
		}

		L := log.With("workspace", tfcReq.Workspace, "organization", tfcReq.Organization)

		// We have to map the organization + workspace  to a workspace-id, so we do that first.
		// the workspaceIds map is never cleared beacuse the configuration about which
		// organization + workspace that is in use is static in the context of a config
		// sourcer.
		key := tfcReq.Organization + "/" + tfcReq.Workspace

		id, ok := cs.workspaceIds[key]
		if !ok {
			L.Trace("querying TFC workspace id")

			url := fmt.Sprintf("%s/api/v2/organizations/%s/workspaces/%s",
				baseURL, tfcReq.Organization, tfcReq.Workspace)

			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				L.Error("error constructing request", "error", err)

				result.Result = &pb.ConfigSource_Value_Error{
					Error: status.New(codes.Aborted, err.Error()).Proto(),
				}

				continue
			}

			req.Header.Set("Authorization", "Bearer "+cs.config.Token)
			req.Header.Set("Content-Type", "application/vnd.api+json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				L.Error("error sending request for workspace info", "error", err)

				result.Result = &pb.ConfigSource_Value_Error{
					Error: status.New(codes.Aborted, err.Error()).Proto(),
				}

				continue
			}

			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				L.Error("error reading workspace info", "status-code", resp.StatusCode)

				result.Result = &pb.ConfigSource_Value_Error{
					Error: status.New(codes.Aborted, "failed to read workspace info").Proto(),
				}

				continue
			}

			var data wsLookup

			err = json.NewDecoder(resp.Body).Decode(&data)
			if err != nil {
				L.Error("error in decoding workspace info", "error", err)

				result.Result = &pb.ConfigSource_Value_Error{
					Error: status.New(codes.Aborted, err.Error()).Proto(),
				}

				continue
			}

			id = data.Data.Id

			if cs.workspaceIds == nil {
				cs.workspaceIds = map[string]string{}
			}

			L.Info("mapped Terraform Cloud organization to workspace-id",
				"organization", tfcReq.Organization,
				"workspace", tfcReq.Workspace,
				"workspace-id", id)

			cs.workspaceIds[key] = id
		}

		// Get this secret or read it if we haven't already.
		cachedSecretVal, ok := cs.secretCache[key]
		if !ok {
			L.Trace("querying TFC workspace")

			url := fmt.Sprintf("%s/api/v2/workspaces/%s/current-state-version?include=outputs",
				baseURL, id)

			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				L.Error("error in creating request for state version", "error", err)

				result.Result = &pb.ConfigSource_Value_Error{
					Error: status.New(codes.Aborted, err.Error()).Proto(),
				}

				continue
			}

			req.Header.Set("Authorization", "Bearer "+cs.config.Token)
			req.Header.Set("Content-Type", "application/vnd.api+json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				L.Error("error in sending request for state version", "error", err)

				result.Result = &pb.ConfigSource_Value_Error{
					Error: status.New(codes.Aborted, err.Error()).Proto(),
				}

				continue
			}

			defer resp.Body.Close()

			var data stateLookup

			err = json.NewDecoder(resp.Body).Decode(&data)
			if err != nil {
				L.Error("error in decoding outputs", "error", err)

				result.Result = &pb.ConfigSource_Value_Error{
					Error: status.New(codes.Aborted, err.Error()).Proto(),
				}

				continue
			}

			outputs := map[string]string{}

			for _, output := range data.Included {
				if output.Type != "state-version-outputs" {
					L.Info("ignored unexpected type in included values", "type", output.Type)
					continue
				}

				outputs[output.Attributes.Name] = output.Attributes.Value
			}

			L.Info("refreshed outputs from Terraform Cloud", "workspace", id, "vars", len(outputs))

			cachedSecretVal = &cachedSecret{Outputs: outputs}
			cs.secretCache[id] = cachedSecretVal
		}

		value := cachedSecretVal.Outputs[tfcReq.Output]

		result.Result = &pb.ConfigSource_Value_Value{
			Value: value,
		}
	}

	return results, nil
}

func (cs *ConfigSourcer) stop() error {
	cs.cacheMu.Lock()
	defer cs.cacheMu.Unlock()

	if cs.client != nil {
		cs.client.CloseIdleConnections()
		cs.client = nil
	}

	// Reset our results tracking to empty. This will force the next call
	// to rebuild all our secret values.
	var zeroTime time.Time
	cs.lastRead = zeroTime
	cs.secretCache = nil

	return nil
}

func (cs *ConfigSourcer) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(
		docs.FromConfig(&sourceConfig{}),
		docs.RequestFromStruct(&reqConfig{}),
	)
	if err != nil {
		return nil, err
	}

	doc.Description("Read Terraform state outputs from Terraform Cloud.")

	doc.Example(`
config {
  env = {
    "DATABASE_USERNAME" = dynamic("terraform-cloud", {
			organization = "foocorp"
			workspace = "databases"
			output = "db_username"
    })

    "DATABASE_PASSWORD" = dynamic("vault", {
			organization = "foocorp"
			workspace = "databases"
			output = "db_password"
    })

    "DATABASE_HOST" = dynamic("vault", {
			organization = "foocorp"
			workspace = "databases"
			output = "db_hostname"
    })
  }
}
`)

	doc.SetRequestField(
		"organization",
		"The Terraform Cloud organization to query",
		docs.Summary(
			"within a single workspace, multiple dynamic values that use the same",
			"organization and workspace will only read the value once. This allows outputs",
			"to be extracted into multiple values. The example above shows",
			"this functionality by reading the username and password into separate values.",
		),
	)

	doc.SetRequestField(
		"workspace",
		"The Terraform Cloud workspace associated with the given organization to read the outputs of",
		docs.Summary(
			"The outputs associtaed with the most recent state version for the given workspace",
			"are the ones that are used. These values are refreshed according to",
			"refreshInternal, a source field.",
		),
	)

	doc.SetRequestField(
		"output",
		"The name of the output to read the value of",
	)

	doc.SetField(
		"token",
		"The Terraform Cloud API token",
		docs.Summary(
			"The token is used to authenticate access to the specific organization and",
			"workspace. Tokens are managed at https://app.terraform.io/app/settings/tokens.",
		),
	)

	doc.SetField(
		"base_url",
		"The scheme, host, and port to calculate the URL to fetch using",
		docs.Summary(
			"This is provided to allow users to query values from Terraform Enterprise",
			"installations",
		),
		docs.Default("https://api.terraform.io"),
	)

	doc.SetField(
		"skip_verify",
		"Do not validate the TLS cert presented by Terraform Cloud.",
		docs.Summary(
			"This is not recommended unless absolutely necessary.",
		),
		docs.EnvVar("TFC_SKIP_VERIFY"),
	)

	doc.SetField(
		"refresh_interval",
		"How often the outputs should be fetch.",
		docs.Default(refreshPeriod.String()),
		docs.Summary(
			"The format of this value is the Go time duration format. Specifically",
			"a whole number followed by: s for seconds, m for minutes, h for hours.",
			"The minimum value for this setting is 60 seconds, with no specified maximum.",
		),
	)

	return doc, nil
}

type reqConfig struct {
	Workspace    string `hcl:"workspace"`
	Organization string `hcl:"organization"`
	Output       string `hcl:"output"`
}

type sourceConfig struct {
	Token           string `hcl:"token"`
	BaseURL         string `hcl:"base_url,optional"`
	SkipVerify      bool   `hcl:"skip_verify,optional"`
	RefreshInterval string `hcl:"refresh_interval,optional"`
}
