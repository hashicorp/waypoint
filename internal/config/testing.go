// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"
)

// TestConfig returns a Config from a string source and fails the test if
// parsing the configuration fails.
func TestConfig(t testing.T, src string) *Config {
	t.Helper()

	// Write our test config to a temp location
	td, err := ioutil.TempDir("", "waypoint")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(td) })

	path := filepath.Join(td, "waypoint.hcl")
	require.NoError(t, ioutil.WriteFile(path, []byte(src), 0644))

	result, err := Load(path, &LoadOptions{})
	require.NoError(t, err)

	return result
}

// TestSource returns valid configuration.
func TestSource(t testing.T) string {
	return testSourceVal
}

// TestSourceJSON returns valid configuration in JSON format.
func TestSourceJSON(t testing.T) string {
	return testSourceValJson
}

// TestConfigFile writes the default Waypoint configuration file with
// the given contents.
func TestConfigFile(t testing.T, src string) {
	require.NoError(t, ioutil.WriteFile(Filename, []byte(src), 0644))
}

const testSourceVal = `
project = "test"
`

const testSourceValJson = `{
  "project": "test"
}`
