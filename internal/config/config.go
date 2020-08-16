package config

import (
	"github.com/hashicorp/hcl/v2"
)

// Config is the configuration structure.
type Config struct {
	Runner  *Runner           `hcl:"runner,block" default:"{}"`
	Project string            `hcl:"project,attr"`
	Apps    []*App            `hcl:"app,block"`
	Labels  map[string]string `hcl:"labels,optional"`
}

// Retrieve the app config for the named application
func (c *Config) AppConfig(name string) (*App, bool) {
	for _, appCfg := range c.Apps {
		if appCfg.Name == name {
			return appCfg, true
		}
	}

	return nil, false
}

// App represents a single application.
type App struct {
	Name   string            `hcl:",label"`
	Path   string            `hcl:"path,optional"`
	Labels map[string]string `hcl:"labels,optional"`
	URL    *AppURL           `hcl:"url,block" default:"{}"`

	Build   *Build   `hcl:"build,block"`
	Deploy  *Deploy  `hcl:"deploy,block"`
	Release *Release `hcl:"release,block"`
}

// AppURL configures the App-specific URL settings.
type AppURL struct {
	AutoHostname *bool `hcl:"auto_hostname,optional"`
}

// Server configures the remote server.
type Server struct {
	Address string `hcl:"address,attr"`

	// Tls, if true, will connect to the server with TLS. If TlsSkipVerify
	// is true, the certificate presented by the server will not be validated.
	Tls           bool `hcl:"tls,optional"`
	TlsSkipVerify bool `hcl:"tls_skip_verify,optional"`

	// AddressInternal is a temporary config to work with local deployments
	// on platforms such as Docker for Mac. We need to discuss a more
	// long term approach to this.
	AddressInternal string `hcl:"address_internal,optional"`

	// Indicates that we need to present a token to connect to this server.
	RequireAuth bool `hcl:"require_auth,optional"`

	// AuthToken is the token to use to authenticate to the server.
	// Note this will be stored plaintext on disk. You can also use the
	// WAYPOINT_SERVER_TOKEN env var.
	AuthToken string `hcl:"auth_token,optional"`
}

// Runner is the configuration for supporting runners in this project.
type Runner struct {
	// Enabled is whether or not runners are enabled. If this is false
	// then the "-remote" flag will not work.
	Enabled bool `hcl:"enabled,attr"`

	// DataSource is the default data source when a remote job is queued.
	DataSource *RunnerDataSource `hcl:"data_source,block" default:"{}"`
}

type RunnerDataSource struct {
	Type string `hcl:",label" default:"auto"`
}

// Hook is the configuration for a hook that runs at specified times.
type Hook struct {
	When      string   `hcl:",label"`
	Command   []string `hcl:"command,attr"`
	OnFailure string   `hcl:"on_failure,optional"`
}

func (h *Hook) ContinueOnFailure() bool {
	return h.OnFailure == "continue"
}

// Build are the build settings.
type Build struct {
	Labels   map[string]string `hcl:"labels,optional"`
	Hooks    []*Hook           `hcl:"hook,block"`
	Use      *Use              `hcl:"use,block"`
	Registry *Registry         `hcl:"registry,block"`
}

// Registry are the registry settings.
type Registry struct {
	Labels map[string]string `hcl:"labels,optional"`
	Hooks  []*Hook           `hcl:"hook,block"`
	Use    *Use              `hcl:"use,block"`
}

// Deploy are the deploy settings.
type Deploy struct {
	Labels map[string]string `hcl:"labels,optional"`
	Hooks  []*Hook           `hcl:"hook,block"`
	Use    *Use              `hcl:"use,block"`
}

// Release are the release settings.
type Release struct {
	Labels map[string]string `hcl:"labels,optional"`
	Hooks  []*Hook           `hcl:"hook,block"`
	Use    *Use              `hcl:"use,block"`
}

// Use is something in the Waypoint configuration that is executed
// using some underlying plugin. This is a general shared structure that is
// used by internal/core to initialize all the proper plugins.
type Use struct {
	Type string   `hcl:",label"`
	Body hcl.Body `hcl:",remain"`
}
