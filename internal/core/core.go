package core

// Common environment variables.
const (
	// ServerAddr is the address for the Waypoint server. This should be
	// in the format of "ip:port" for TCP.
	EnvServerAddr = "WAYPOINT_SERVER_ADDR"

	// ServerInsecure should be any value that strconv.ParseBool parses as
	// true to connect to the server insecurely.
	EnvServerInsecure = "WAYPOINT_SERVER_INSECURE"
)
