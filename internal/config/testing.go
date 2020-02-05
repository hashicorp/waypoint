package config

import (
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

// TestConfig returns a Config from a string source and fails the test if
// parsing the configuration fails.
func TestConfig(t testing.T, src string) *Config {
	var result Config
	require.NoError(t, hclsimple.Decode("test.hcl", []byte(src), nil, &result))
	return &result
}
