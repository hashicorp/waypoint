package util

import (
	"os"
	"strconv"
)

// Extracts a boolean from an env var. Falls back to the default if
// The key is unset or not a valid boolean.
func GetEnvBool(key string, defaultValue bool) (value bool) {
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
