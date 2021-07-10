package oidc

import (
	"context"
	"errors"
	"io"
	"strings"
	"sync"

	"github.com/hashicorp/cap/oidc"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// ProviderConfig returns the OIDC provider configuration for an OIDC auth method.
func ProviderConfig(am *pb.AuthMethod_OIDC) (*oidc.Config, error) {
	var algs []oidc.Alg
	if len(am.SigningAlgs) > 0 {
		for _, alg := range am.SigningAlgs {
			algs = append(algs, oidc.Alg(alg))
		}
	} else {
		algs = []oidc.Alg{oidc.RS256}
	}

	return oidc.NewConfig(
		am.DiscoveryUrl,
		am.ClientId,
		oidc.ClientSecret(am.ClientSecret),
		algs,
		am.AllowedRedirectUris,
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
// amount of time (default 12 hours) in case the remote provider configuration
// changed.
type ProviderCache struct {
	providers map[string]*oidc.Provider
	mu        sync.RWMutex
}

// Get returns the OIDC provider for the given auth method configuration.
// This will initialize the provider if it isn't already in the cache or
// if the configuration changed.
func (c *ProviderCache) Get(
	ctx context.Context, am *pb.AuthMethod,
) (*oidc.Provider, error) {
	amMethod, ok := am.Method.(*pb.AuthMethod_Oidc)
	if !ok {
		return nil, errors.New("auth method must be OIDC")
	}

	// No matter what we'll use the config of the arg method since we'll
	// use it to compare to existing (if exists) or initialize a new provider.
	oidcCfg, err := ProviderConfig(amMethod.Oidc)
	if err != nil {
		return nil, err
	}

	// Get our current value
	var current *oidc.Provider
	ok = false
	c.mu.RLock()
	if c.providers != nil {
		current, ok = c.providers[am.Name]
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

	if c.providers == nil {
		c.providers = map[string]*oidc.Provider{}
	}
	c.providers[am.Name] = newProvider

	return newProvider, nil
}

// Clear is called to delete all the providers in the cache.
func (c *ProviderCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, p := range c.providers {
		p.Done()
	}
	c.providers = nil
}

// Close implements io.Closer. This just calls Clear, but we implement
// io.Closer so that things that look for this interface implementation for
// "cleanup" will call this.
func (c *ProviderCache) Close() error {
	c.Clear()
	return nil
}

var _ io.Closer = (*ProviderCache)(nil)
