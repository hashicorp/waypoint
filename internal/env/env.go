// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// GetBool Extracts a boolean from an env var. Falls back to the default
// if the key is unset or not a valid boolean.
func GetBool(key string, defaultValue bool) (bool, error) {
	envVal := os.Getenv(key)
	if envVal == "" {
		return defaultValue, nil
	}
	value, err := strconv.ParseBool(strings.ToLower(envVal))
	if err != nil {
		return defaultValue, fmt.Errorf("failed to parse a boolean from environment variable %s=%s", key, envVal)
	}
	return value, nil
}
