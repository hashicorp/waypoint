package env

import (
	"os"
	"strconv"
)

// GetEnvBool Extracts a boolean from an env var. Falls back to the default
// if the key is unset or not a valid boolean.
func GetEnvBool(key string, defaultValue bool) bool {
	envVal := os.Getenv(key)
	if envVal == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(envVal)
	if err != nil {
		return defaultValue
	}
	return value
}
