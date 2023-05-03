// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ssh

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/subtle"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/gliderlabs/ssh"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"
	gossh "golang.org/x/crypto/ssh"
)

func TestServer(t *testing.T) {
	hostkey, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err)

	userkey, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err)

	hostSigner, err := gossh.NewSignerFromKey(hostkey)
	require.NoError(t, err)

	userSigner, err := gossh.NewSignerFromKey(userkey)
	require.NoError(t, err)

	check := func(ctx ssh.Context, inputKey ssh.PublicKey) bool {
		return ssh.KeysEqual(inputKey, userSigner.PublicKey())
	}

	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	var server *ssh.Server

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		ssh.Serve(l,
			createHandler(ctx, hclog.L(), &server),
			ssh.Option(func(serv *ssh.Server) error {
				server = serv
				serv.PublicKeyHandler = check
				serv.AddHostKey(hostSigner)
				return nil
			}),
		)
	}()

	time.Sleep(time.Second)

	var cfg gossh.ClientConfig
	cfg.User = "waypoint"
	cfg.Auth = []gossh.AuthMethod{
		gossh.PublicKeys(userSigner),
	}

	expectedHost := hostSigner.PublicKey().Marshal()

	cfg.HostKeyCallback = func(hostname string, remote net.Addr, key gossh.PublicKey) error {
		// Weirdly this is how you make sure the host key is what you think it should be.
		// Think of this as where normal ssh client would do the "Do you want to trust this
		// host?" popup.
		if subtle.ConstantTimeCompare(expectedHost, key.Marshal()) == 1 {
			return nil
		}

		return fmt.Errorf("wrong host key detected")
	}

	cfg.Timeout = 5 * time.Second

	client, err := gossh.Dial("tcp", l.Addr().String(), &cfg)
	require.NoError(t, err)

	sess, err := client.NewSession()
	require.NoError(t, err)

	var buf bytes.Buffer

	sess.Stdout = &buf

	err = sess.Run("sh -c 'echo hello'")
	require.NoError(t, err)
	require.Eventually(
		t,
		func() bool { return buf.String() == "hello\n" },
		time.Second,
		10*time.Millisecond,
	)
}
