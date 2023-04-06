// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cert

import (
	"context"
	"crypto/tls"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/hashicorp/go-hclog"
)

// Cert represents a single certificate.
type Cert struct {
	log hclog.Logger

	// All the fields below must be protected with this lock.
	lock              *sync.RWMutex
	certFile, keyFile string
	cert              *tls.Certificate
	watcher           *fsnotify.Watcher
	watcherStop       context.CancelFunc
}

// New initializes a certificate from a PEM-encoded certificate and private key
// written to disk. This loads the initial certificate and sets up file watchers
// to watch for any changes to reload the certificate.
func New(log hclog.Logger, crtPath, keyPath string) (*Cert, error) {
	if log == nil {
		log = hclog.L()
	}

	var lock sync.RWMutex
	c := &Cert{
		log:      log,
		certFile: crtPath,
		keyFile:  keyPath,
		lock:     &lock,
	}

	// Do an initial reload
	if err := c.reload(); err != nil {
		return nil, err
	}

	// Initialize watcher to watch for file changes. We just watch for
	// certificate file changes because any change in the certificate
	// must have a change in the key (and vice versa) so whenever the
	// certificate changes we reload everything.
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	if err := w.Add(crtPath); err != nil {
		w.Close()
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	c.watcher = w
	c.watcherStop = cancel
	go c.watch(ctx, w)

	return c, nil
}

// Replace replaces this certificate with a new path. This is done
// atomically so active TLS connections are unaffected and new connections
// will use the new certificate.
func (c *Cert) Replace(crtPath, keyPath string) error {
	// Create a new cert, it is easier to handle errors.
	newCert, err := New(hclog.NewNullLogger(), crtPath, keyPath)
	if err != nil {
		return err
	}

	// Lock and replace some internal state
	c.lock.Lock()
	defer c.lock.Unlock()

	// End our watcher
	if v := c.watcherStop; v != nil {
		v()
	}

	// Copy over some state
	c.certFile = newCert.certFile
	c.keyFile = newCert.keyFile
	c.cert = newCert.cert
	c.watcher = newCert.watcher
	c.watcherStop = newCert.watcherStop
	return nil
}

// Close implements io.Closer. This must be called to properly clean up
// resources associated with watching for certificate changes.
func (c *Cert) Close() error {
	if c.watcherStop != nil {
		c.watcherStop()
	}

	return nil
}

// Paths returns the paths to the certificate and key that are currently in use.
func (c *Cert) Paths() (crt, key string) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.certFile, c.keyFile
}

// TLSConfig returns a TLS configuration struct that can be used that
// uses this certificate.
func (c *Cert) TLSConfig() *tls.Config {
	return &tls.Config{
		GetCertificate:           c.GetCertificate,
		PreferServerCipherSuites: true,
	}
}

// GetCertificate implements the GetCertificate callback for tls.Config
// and can be used to get the latest certificate at all times.
func (c *Cert) GetCertificate(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.cert, nil
}

func (c *Cert) reload() error {
	cert, err := tls.LoadX509KeyPair(c.certFile, c.keyFile)
	if err != nil {
		return err
	}

	// Replace our certificate
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cert = &cert
	return nil
}

func (c *Cert) watch(ctx context.Context, w *fsnotify.Watcher) {
	// When we're done, stop the watcher
	defer w.Close()

	for {
		select {
		case <-ctx.Done():
			return

		case event := <-w.Events:
			// A change happened, do a reload
			c.log.Warn("certificate change detected",
				"name", event.Name,
				"op", event.Op.String(),
			)

			// Reload and retry a few times. We retry because sometimes
			// we get a notification for the cert before the key is ready.
			for i := 0; i < 50; i++ {
				if err := c.reload(); err != nil {
					c.log.Warn("error during reload", "err", err, "i", i)
					time.Sleep(100 * time.Millisecond)
					continue
				}

				break
			}

		case err := <-w.Errors:
			c.log.Warn("error in filesystem watch", "err", err)
		}
	}
}
