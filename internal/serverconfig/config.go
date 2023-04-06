// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package serverconfig

// Config is the configuration for the built-in server.
type Config struct {
	// DBPath is the path to the database file, including the filename.
	DBPath string `hcl:"db_path,attr"`

	// GRPC is the grpc service listening configuration. This is required.
	GRPC Listener `hcl:"grpc,block"`

	// HTTP is the listening configuration for the HTTP service for grpc-web.
	HTTP Listener `hcl:"http,block"`

	// HTTPInsecure sets up a listener for HTTP that never has TLS enabled.
	// This is generally not recommended but can make sense in certain
	// environments. For example, within Kubernetes where TLS termination
	// happens at a higher level you may want an additional insecure listener
	// in addition to a secure listener.
	HTTPInsecure Listener `hcl:"http_insecure,block"`

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
