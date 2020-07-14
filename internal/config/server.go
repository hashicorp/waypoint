package config

// ServerConfig is the configuration for the built-in server.
type ServerConfig struct {
	// DBPath is the path to the database file, including the filename.
	DBPath string `hcl:"db_path,attr"`

	// Listeners sets up the listeners
	Listeners Listeners `hcl:"listeners,block"`

	// Require clients to authenticate themselves
	RequireAuth bool `hcl:"require_auth,optional"`
}

// Listeners is the configuration for the listeners.
type Listeners struct {
	GRPC string `hcl:"grpc,attr"`
	HTTP string `hcl:"http,attr"`
}
