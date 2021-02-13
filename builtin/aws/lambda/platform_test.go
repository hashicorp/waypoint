package lambda

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

// Just to validate that the key armoring works properly
func TestKeys(t *testing.T) {
	hostkey, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err)

	hoststr := base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(hostkey))

	hostbytes, err := base64.StdEncoding.DecodeString(hoststr)
	require.NoError(t, err)

	hkey, err := x509.ParsePKCS1PrivateKey(hostbytes)
	require.NoError(t, err)

	assert.True(t, hostkey.Equal(hkey))

	userstr := base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(&hostkey.PublicKey))

	userbytes, err := base64.StdEncoding.DecodeString(userstr)
	require.NoError(t, err)

	userKey, err := x509.ParsePKCS1PublicKey(userbytes)
	require.NoError(t, err)

	_, err = ssh.NewPublicKey(userKey)
	require.NoError(t, err)
}
