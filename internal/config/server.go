package config

// ServerConfig is the configuration for the built-in server.
type ServerConfig struct {
	// DBPath is the path to the database file, including the filename.
	DBPath string `hcl:"db_path,attr"`

	// Listeners sets up the listeners
	Listeners Listeners `hcl:"listeners,block"`

	// Require clients to authenticate themselves
	RequireAuth bool `hcl:"require_auth,optional"`

	// URL configures a server to use a URL service.
	URL *URL `hcl:"url,block"`
}

// Listeners is the configuration for the listeners.
type Listeners struct {
	GRPC string `hcl:"grpc,attr"`
	HTTP string `hcl:"http,attr"`
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
