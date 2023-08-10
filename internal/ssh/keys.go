// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

const (
	// Contains the port that the ssh server should listen on. The number should be
	// base 10 encoded.
	ENVSSHPort = "WAYPOINT_EXEC_PLUGIN_SSH"

	// hostKey contains an SSH RSA private key, marshaled as PKCS1 and armored
	// with base64. This will be used as the servers host key and verified
	// by the client when it connects.
	ENVHostKey = "WAYPOINT_EXEC_PLUGIN_SSH_HOST_KEY"

	// key contains an SSH RSA public key, marshaled as PKCS1 and armored
	// with base64. This will be used to authenticate the ssh client.

	ENVUserKey = "WAYPOINT_EXEC_PLUGIN_SSH_KEY"
)

// MarshalPrivateKey converts the key to a string, such that UnmarshalPrivateKey can
// return the same key.
func MarshalPrivateKey(key *rsa.PrivateKey) string {
	return base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(key))
}

// UnmarshalPrivateKey parses the string into a rsa.PrivateKey.
func UnmarshalPrivateKey(str string) (*rsa.PrivateKey, error) {
	hostbytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, errors.Wrapf(err, "decoding host key")
	}

	hkey, err := x509.ParsePKCS1PrivateKey(hostbytes)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing host key")
	}

	return hkey, nil
}

// MarshalPublicKey converts a PubilcKey into a string that can be decoded by
// UnmarshalPublicKey.
func MarshalPublicKey(key *rsa.PublicKey) string {
	return base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(key))
}

// UnmarshalPublicKey parses a string into a PubilcKey. Both keys are the same
// value, just different representations.
func UnmarshalPublicKey(str string) (*rsa.PublicKey, ssh.PublicKey, error) {
	userbytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "decoding user key")
	}

	userKey, err := x509.ParsePKCS1PublicKey(userbytes)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "parsing user key")
	}

	authorizedKey, err := ssh.NewPublicKey(userKey)

	if err != nil {
		return nil, nil, err
	}

	return userKey, authorizedKey, nil
}

// SSHKeyMaterial holds the key material required to setup an SSH connection between
// a server and client. These are commonly used by exec plugins and the waypoint entrypoint
// to create adhoc ssh servers that can run a users command.
type SSHKeyMaterial struct {
	// The rsa host key to use for the SSH server. Armored as a string for easy passage.
	HostPrivate string

	// The public half of the host key. Use this to authenticate the server when connecting.
	HostPublic ssh.PublicKey

	// The private key of the client. Use this to authenticate with the server as the client.
	UserPrivate ssh.Signer

	// The public half of the client key. The server uses this to authenticate the client.
	UserPublic string

	// The raw user key, provided in for further usage.
	UserKey *rsa.PrivateKey

	// The raw host key, provided in for further usage.
	HostKey *rsa.PrivateKey
}

// GenerateKeys generates a new SSHKeyMaterial with random keys.
func GenerateKeys() (*SSHKeyMaterial, error) {
	var mat SSHKeyMaterial
	hostkey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	mat.HostKey = hostkey
	mat.HostPrivate = MarshalPrivateKey(hostkey)

	hostSigner, err := ssh.NewSignerFromKey(hostkey)
	if err != nil {
		return nil, err
	}

	mat.HostPublic = hostSigner.PublicKey()

	userkey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	mat.UserKey = userkey
	mat.UserPrivate, err = ssh.NewSignerFromSigner(userkey)
	if err != nil {
		return nil, err
	}

	mat.UserPublic = MarshalPublicKey(&userkey.PublicKey)

	return &mat, nil
}

var ErrMissingSSHKey = errors.New("missing ssh key information in environment")

// DecodeFromEnv reads the processes environment data and decodes the host
// and user keys from it, returning ready to use representations of those keys.
func DecodeFromEnv() (ssh.Signer, ssh.PublicKey, error) {
	host := os.Getenv(ENVHostKey)
	if host == "" {
		return nil, nil, errors.Wrapf(ErrMissingSSHKey, "missing host key")
	}

	user := os.Getenv(ENVUserKey)
	if user == "" {
		return nil, nil, errors.Wrapf(ErrMissingSSHKey, "missing user key")
	}

	hostKey, err := UnmarshalPrivateKey(host)
	if err != nil {
		return nil, nil, err
	}

	signer, err := ssh.NewSignerFromKey(hostKey)
	if err != nil {
		return nil, nil, err
	}

	_, auth, err := UnmarshalPublicKey(user)
	if err != nil {
		return nil, nil, err
	}

	return signer, auth, nil
}
