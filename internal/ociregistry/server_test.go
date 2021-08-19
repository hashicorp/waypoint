package ociregistry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	t.Run("can parse a www-authenticate header", func(t *testing.T) {
		typ, realm, qs, err := parseWWWAuthenticateHeader(`Bearer realm="https://api.digitalocean.com/v2/registry/auth",service="registry.digitalocean.com"`)
		require.NoError(t, err)

		assert.Equal(t, "Bearer", typ)
		assert.Equal(t, "https://api.digitalocean.com/v2/registry/auth", realm)
		assert.Equal(t, []string{"service=registry.digitalocean.com"}, qs)
	})
}
