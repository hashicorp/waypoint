// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package oidc

import (
	"context"
	"errors"
	"io"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/cap/oidc"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// ProviderConfig returns the OIDC provider configuration for an OIDC auth method.
// The ServerConfig argument can be nil. If it is not nil, then the server
// advertise addresses will be added as valid redirect URLs.
func ProviderConfig(am *pb.AuthMethod_OIDC, sc *pb.ServerConfig) (*oidc.Config, error) {
	var algs []oidc.Alg
	if len(am.SigningAlgs) > 0 {
		for _, alg := range am.SigningAlgs {
			algs = append(algs, oidc.Alg(alg))
		}
	} else {
		algs = []oidc.Alg{oidc.RS256}
	}

	// Modify our allowed uris to always have loopback addresses for our CLI.
	// Note that "http" instead of "https" is okay for loopback addresses since
	// the protocol is ignored anyways for loopback.
	allowedUris := make([]string, len(am.AllowedRedirectUris))
	copy(allowedUris, am.AllowedRedirectUris)
	allowedUris = append(allowedUris,
		// Loopback addresses used by the CLI
		"http://localhost/oidc/callback",
		"http://127.0.0.1/oidc/callback",
		"http://[::1]/oidc/callback",
	)

	// We also add all the addresses our UI might use with the server advertise addrs.
	if sc != nil {
		for _, addr := range sc.AdvertiseAddrs {
			var u url.URL
			u.Scheme = "https"
			u.Host = addr.Addr
			if !addr.Tls {
				u.Scheme = "http"
			}
			u.Path = "/auth/oidc-callback"

			allowedUris = append(allowedUris, u.String())
		}
	}

	return oidc.NewConfig(
		am.DiscoveryUrl,
		am.ClientId,
		oidc.ClientSecret(am.ClientSecret),
		algs,
		allowedUris,
		oidc.WithAudiences(am.Auds...),
		oidc.WithProviderCA(strings.Join(am.DiscoveryCaPem, "\n")),
	)
}

// ProviderCache is a cache for OIDC providers. OIDC providers are something
// you don't want to recreate per-request since they make HTTP requests
// themselves.
//
// The ProviderCache purges a provider under two scenarios: (1) the
// provider config is updated and it is different and (2) after a set
// amount of time (see cacheExpiry for value) in case the remote provider configuration
// changed.
type ProviderCache struct {
	providers map[string]*oidc.Provider
	mu        sync.RWMutex
	cancel    context.CancelFunc
}

// NewProviderCache should be used to initialize a provider cache. This
// will start up background resources to manage the cache.
func NewProviderCache() *ProviderCache {
	ctx, cancel := context.WithCancel(context.Background())
	result := &ProviderCache{
		providers: map[string]*oidc.Provider{},
		cancel:    cancel,
	}

	// Start the cleanup timer
	go result.runCleanupLoop(ctx)

	return result
}

// Get returns the OIDC provider for the given auth method configuration.
// This will initialize the provider if it isn't already in the cache or
// if the configuration changed.
func (c *ProviderCache) Get(
	ctx context.Context, am *pb.AuthMethod, sc *pb.ServerConfig,
) (*oidc.Provider, error) {
	amMethod, ok := am.Method.(*pb.AuthMethod_Oidc)
	if !ok {
		return nil, errors.New("auth method must be OIDC")
	}

	// No matter what we'll use the config of the arg method since we'll
	// use it to compare to existing (if exists) or initialize a new provider.
	oidcCfg, err := ProviderConfig(amMethod.Oidc, sc)
	if err != nil {
		return nil, err
	}

	// Normalize name
	name := strings.ToLower(am.Name)

	// Get our current value
	var current *oidc.Provider
	ok = false
	c.mu.RLock()
	if c.providers != nil {
		current, ok = c.providers[name]
	}
	c.mu.RUnlock()

	// If we have a current value, we want to compare hashes to detect changes.
	if ok {
		currentHash, err := current.ConfigHash()
		if err != nil {
			return nil, err
		}

		newHash, err := oidcCfg.Hash()
		if err != nil {
			return nil, err
		}

		// If the hashes match, this is cached.
		if currentHash == newHash {
			return current, nil
		}
	}

	// If we made it here, the provider isn't in the cache OR the config changed.
	// Initialize a new provider.
	newProvider, err := oidc.NewProvider(oidcCfg)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// If we have an old provider, clean up resources.
	if current != nil {
		current.Done()
	}

	c.providers[name] = newProvider

	return newProvider, nil
}

// Delete force deletes a single auth method from the cache by name.
func (c *ProviderCache) Delete(ctx context.Context, name string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	name = strings.ToLower(name)
	p, ok := c.providers[name]
	if ok {
		p.Done()
		delete(c.providers, name)
	}
}

// Clear is called to delete all the providers in the cache.
func (c *ProviderCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, p := range c.providers {
		p.Done()
	}
	c.providers = map[string]*oidc.Provider{}
}

// Close implements io.Closer. This just calls Clear, but we implement
// io.Closer so that things that look for this interface implementation for
// "cleanup" will call this.
func (c *ProviderCache) Close() error {
	c.cancel()
	c.Clear()
	return nil
}

// runCleanupLoop runs an infinite loop that clears the cache every
// "cacheExpiry" duration. This ensures that we force refresh our provider
// info periodically in case anything changes. In practice, this is very
// rare so we don't refresh very often.
func (c *ProviderCache) runCleanupLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		// note(mitchellh): we could be more clever and do a per-entry
		// expiry but in practice a Waypoint server probably has one, maybe
		// two auth methods, and its just not worth the complexity.
		case <-time.After(cacheExpiry):
			c.Clear()
		}
	}
}

// cacheExpiry is the duration after which the provider cache is reset.
const cacheExpiry = 6 * time.Hour

var _ io.Closer = (*ProviderCache)(nil)
