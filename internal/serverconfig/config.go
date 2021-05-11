package serverconfig

import (
	"strconv"
)

// Client configures a client to connect to a server.
type Client struct {
	Address string `hcl:"address,attr" json:"address"`

	// Tls, if true, will connect to the server with TLS. If TlsSkipVerify
	// is true, the certificate presented by the server will not be validated.
	Tls           bool `hcl:"tls,optional" json:"tls,omitempty"`
	TlsSkipVerify bool `hcl:"tls_skip_verify,optional" json:"tls_skip_verify,omitempty"`

	// AddressInternal is a temporary config to work with local deployments
	// on platforms such as Docker for Mac. We need to discuss a more
	// long term approach to this.
	AddressInternal string `hcl:"address_internal,optional" json:"address_internal,omitempty"`

	// Indicates that we need to present a token to connect to this server.
	RequireAuth bool `hcl:"require_auth,optional" json:"require_path,omitempty"`

	// AuthToken is the token to use to authenticate to the server.
	// Note this will be stored plaintext on disk. You can also use the
	// WAYPOINT_SERVER_TOKEN env var.
	AuthToken string `hcl:"auth_token,optional" json:"auth_token,omitempty"`

	// The platform for where the server is running. Although this option should
	// be required, it's optional to support previously set contexts that did
	// not have a platform.
	Platform string `hcl:"platform,optional" json:"platform,omitempty"`
}

// Env returns a slice of environment variables in key=value settings
// that will authenticate to the server without a context set.
func (c *Client) Env() []string {
	result := []string{
		"WAYPOINT_SERVER_ADDR=" + c.Address,
		"WAYPOINT_SERVER_TLS=" + strconv.FormatBool(c.Tls),
		"WAYPOINT_SERVER_TLS_SKIP_VERIFY=" + strconv.FormatBool(c.TlsSkipVerify),
	}

	if c.RequireAuth {
		result = append(result, "WAYPOINT_SERVER_TOKEN="+c.AuthToken)
	}

	return result
}

// Config is the configuration for the built-in server.
type Config struct {
	// DBPath is the path to the database file, including the filename.
	DBPath string `hcl:"db_path,attr"`

	// GRPC is the grpc service listening configuration. This is required.
	GRPC Listener `hcl:"grpc,block"`

	// HTTP is the listening configuration for the HTTP service for grpc-web.
	HTTP Listener `hcl:"http,block"`

	// URL configures a server to use a URL service.
	URL *URL `hcl:"url,block"`

	// CEBConfig configures the entrypoint binary for deployments
	CEBConfig *CEBConfig `hcl:"entrypoint_config,block"`
}

// CEBConfig is specific configuration for the entrypoint binaries
// injected into the deployments
type CEBConfig struct {
	Addr          string `hcl:"addr,optional"`
	TLSEnabled    bool   `hcl:"tls_enabled,optional"`
	TLSSkipVerify bool   `hcl:"tls_skip_verify,optional"`
}

type Listener struct {
	Addr        string `hcl:"address,attr"`
	TLSDisable  bool   `hcl:"tls_disable,optional"`
	TLSCertFile string `hcl:"tls_cert_file,optional"`
	TLSKeyFile  string `hcl:"tls_key_file,optional"`
}

// URL is the configuration for the URL service.
type URL struct {
	Enabled              bool   `hcl:"enabled,optional"`
	APIAddress           string `hcl:"api_address,optional"`
	APIInsecure          bool   `hcl:"api_insecure,optional"`
	APIToken             string `hcl:"api_token,optional"`
	ControlAddress       string `hcl:"control_address,optional"`
	AutomaticAppHostname bool   `hcl:"automatic_app_hostname,optional"`
}
