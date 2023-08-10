// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cert

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/pkg/copy"
)

func TestCert_replace(t *testing.T) {
	require := require.New(t)

	// Our test cert
	crtPath := filepath.Join("testdata", "tls.crt")
	keyPath := filepath.Join("testdata", "tls.key")

	// Create it
	c, err := New(nil, crtPath, keyPath)
	require.NoError(err)
	defer c.Close()

	// Get the certificate
	cert, err := c.GetCertificate(nil)
	require.NoError(err)
	require.NotNil(cert)

	// Replace it
	require.NoError(c.Replace("testdata/tls2.crt", "testdata/tls2.key"))

	// Get and they should not be equal
	cert2, err := c.GetCertificate(nil)
	require.NoError(err)
	require.NotNil(cert2)
	require.NotEqual(cert.Certificate, cert2.Certificate)
}

func TestCert_replaceFail(t *testing.T) {
	require := require.New(t)

	// Our test cert
	crtPath := filepath.Join("testdata", "tls.crt")
	keyPath := filepath.Join("testdata", "tls.key")

	// Create it
	c, err := New(nil, crtPath, keyPath)
	require.NoError(err)
	defer c.Close()

	// Get the certificate
	cert, err := c.GetCertificate(nil)
	require.NoError(err)
	require.NotNil(cert)

	// Replace it
	require.Error(c.Replace("testdata/tls2.crt", "testdata/tls.key"))

	// Get and they should not be equal
	cert2, err := c.GetCertificate(nil)
	require.NoError(err)
	require.NotNil(cert2)
	require.Equal(cert.Certificate, cert2.Certificate)
}

func TestCert_watch(t *testing.T) {
	require := require.New(t)

	// Copy to a temporary directory
	td, err := ioutil.TempDir("", "go-cert")
	require.NoError(err)
	defer os.RemoveAll(td)
	path := filepath.Join(td, "testdata")

	// Copy
	require.NoError(copy.CopyDir("testdata", path))

	// Our test cert
	crtPath := filepath.Join(path, "tls.crt")
	keyPath := filepath.Join(path, "tls.key")

	// Create it
	c, err := New(nil, crtPath, keyPath)
	require.NoError(err)
	defer c.Close()

	// Get the certificate
	cert, err := c.GetCertificate(nil)
	require.NoError(err)
	require.NotNil(cert)

	// Move tls2 over.
	require.NoError(copy.CopyFile(
		filepath.Join(path, "tls2.key"),
		keyPath))
	require.NoError(copy.CopyFile(
		filepath.Join(path, "tls2.crt"),
		crtPath))

	// Should change
	require.Eventually(func() bool {
		cert2, err := c.GetCertificate(nil)
		require.NoError(err)
		require.NotNil(cert)
		return !bytes.Equal(cert.Certificate[0], cert2.Certificate[0])
	}, 5*time.Second, 100*time.Millisecond)
}

func TestCert_watchCrtFirst(t *testing.T) {
	require := require.New(t)

	// Copy to a temporary directory
	td, err := ioutil.TempDir("", "go-cert")
	require.NoError(err)
	defer os.RemoveAll(td)
	path := filepath.Join(td, "testdata")

	// Copy
	require.NoError(copy.CopyDir("testdata", path))

	// Our test cert
	crtPath := filepath.Join(path, "tls.crt")
	keyPath := filepath.Join(path, "tls.key")

	// Create it
	c, err := New(nil, crtPath, keyPath)
	require.NoError(err)
	defer c.Close()

	// Get the certificate
	cert, err := c.GetCertificate(nil)
	require.NoError(err)
	require.NotNil(cert)

	// Move tls2 over.
	require.NoError(copy.CopyFile(
		filepath.Join(path, "tls2.crt"),
		crtPath))
	time.Sleep(100 * time.Millisecond)
	require.NoError(copy.CopyFile(
		filepath.Join(path, "tls2.key"),
		keyPath))

	// Should change
	require.Eventually(func() bool {
		cert2, err := c.GetCertificate(nil)
		require.NoError(err)
		require.NotNil(cert)
		return !bytes.Equal(cert.Certificate[0], cert2.Certificate[0])
	}, 5*time.Second, 100*time.Millisecond)
}
