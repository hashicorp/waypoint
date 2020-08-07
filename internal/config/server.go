package config

// ServerConfig is the configuration for the built-in server.
type ServerConfig struct {
	// DBPath is the path to the database file, including the filename.
	DBPath string `hcl:"db_path,attr"`

	// GRPC is the grpc service listening configuration. This is required.
	GRPC Listener `hcl:"grpc,block"`

	// HTTP is the listening configuration for the HTTP service for grpc-web.
	HTTP Listener `hcl:"http,block"`

	// URL configures a server to use a URL service.
	URL *URL `hcl:"url,block"`
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
