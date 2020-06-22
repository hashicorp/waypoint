package config

import (
	"io/ioutil"

	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

// TestConfig returns a Config from a string source and fails the test if
// parsing the configuration fails.
func TestConfig(t testing.T, src string) *Config {
	t.Helper()

	var result Config
	require.NoError(t, hclsimple.Decode("test.hcl", []byte(src), nil, &result))
	return &result
}

// TestSource returns valid configuration.
func TestSource(t testing.T) string {
	return testSourceVal
}

// TestConfigFile writes the default Waypoint configuration file with
// the given contents.
func TestConfigFile(t testing.T, src string) {
	require.NoError(t, ioutil.WriteFile(Filename, []byte(src), 0644))
}

const testSourceVal = `
project = "test"
`
