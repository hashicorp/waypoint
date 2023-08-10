// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint/builtin/vault/internal/auth"
	"github.com/hashicorp/waypoint/builtin/vault/internal/auth/approle"
	"github.com/hashicorp/waypoint/builtin/vault/internal/auth/aws"
	"github.com/hashicorp/waypoint/builtin/vault/internal/auth/gcp"
	"github.com/hashicorp/waypoint/builtin/vault/internal/auth/kubernetes"
)

// initAuthMethod initializes a goroutine that grabs a Vault token
// using the configured auth method and continues to update that Vault
// token as it changes.
//
// If no auth method is configured, this does nothing.
//
// This method expects the cacheMu lock to be held.
func (cs *ConfigSourcer) initAuthMethod(
	log hclog.Logger,
) error {
	// If we have no auth method, ensure our token is set to nothing.
	if cs.config.AuthMethod == "" {
		return nil
	}

	// You always need to set a mount path if you're using an auth method
	if cs.config.AuthMethodMountPath == "" {
		return fmt.Errorf(
			"auth_method_mount_path must be set with waypoint config source-set")
	}

	// Get the config
	config, ok := authMethodConfig(&cs.config)
	if !ok {
		return fmt.Errorf("unknown auth method: %s", cs.config.AuthMethod)
	}
	mountPath := cs.config.AuthMethodMountPath

	// Build our auth method
	var method auth.AuthMethod
	var err error
	authConfig := &auth.AuthConfig{
		Logger:    log.Named("auth"),
		MountPath: mountPath,
		Config:    config,
	}
	switch cs.config.AuthMethod {
	case "aws":
		method, err = aws.NewAWSAuthMethod(authConfig)
	case "gcp":
		method, err = gcp.NewGCPAuthMethod(authConfig)
	case "kubernetes":
		method, err = kubernetes.NewKubernetesAuthMethod(authConfig)
	case "approle":
		method, err = approle.NewApproleAuthMethod(authConfig)
	default:
		return fmt.Errorf("unknown auth method: %s", cs.config.AuthMethod)
	}
	if err != nil {
		return err
	}

	// Create our auth handler
	ah := auth.NewAuthHandler(&auth.AuthHandlerConfig{
		Logger: log.Named("auth.handler"),
		Client: cs.client,
	})

	// Create the context we'll use to cancel this
	ctx, cancel := context.WithCancel(context.Background())

	// Start the auth handler
	go func() {
		err := ah.Run(ctx, method)
		if err == nil {
			log.Debug("auth handler stopped")
		} else {
			log.Warn("auth handler stopped with error", "err", err)
		}
	}()

	// Start a goroutine that waits for token updates and stores them.
	firstTokenCh := make(chan struct{})
	go func(initCh chan<- struct{}) {
		for {
			select {
			case <-ctx.Done():
				return

			case token := <-ah.OutputCh:
				log.Trace("new Vault token received")

				// We usually lock cacheMu but on first run we have to lock
				// a blank mutex since we block waiting for the first token
				// on a goroutine that already holds the lock.
				mu := &cs.cacheMu
				if initCh != nil {
					mu = &sync.Mutex{}
				}

				mu.Lock()
				if cs.client != nil {
					cs.client.SetToken(token)
				}
				mu.Unlock()

				// We close initCh exactly once to note that we got our first token
				if initCh != nil {
					close(initCh)
					initCh = nil
				}
			}
		}
	}(firstTokenCh)

	// Wait for our first token to be set on the client before returning
	log.Debug("waiting for Vault token from auth method")
	select {
	case <-firstTokenCh:
		log.Debug("first auth token received and set")

	case <-time.After(10 * time.Second):
		// We do a 10 second timeout because the auth handler could get
		// stuck in a failure loop on the first token request if it is
		// misconfigured. This ensures that we don't block forever on
		// auth that is never going to succeed.
		cancel()
		return fmt.Errorf("timeout waiting for Vault token via auth method")
	}

	// Set our cancel
	cs.authCancel = cancel

	return nil
}

// authMethodConfig builds the configuration map that we need to send in
// to the Vault auth method library. The second return value is true if
// the given auth method was valid and found in the config.
func authMethodConfig(config *sourceConfig) (map[string]interface{}, bool) {
	// we'll accumulate our result here
	result := map[string]interface{}{}

	// tagPrefix is the prefix that the struct tag should have
	tagPrefix := config.AuthMethod + "_"

	// found is set to true when we find a field that is part of the
	// set auth method. We use this to determine if the auth method doesn't exist.
	found := false

	// We use reflection to look up all the fields that have the same
	// prefix as the auth method in use to build up our config.
	v := reflect.ValueOf(config).Elem()
	vt := v.Type()
	for i := 0; i < vt.NumField(); i++ {
		f := vt.Field(i)
		tag := f.Tag.Get("hcl")
		if !strings.HasPrefix(tag, tagPrefix) {
			// Ignore fields that don't have the prefix
			continue
		}
		found = true

		// Remove our prefix
		tag = strings.TrimPrefix(tag, tagPrefix)

		// If we have a , for an attribute, we have to trim up to that
		if idx := strings.Index(tag, ","); idx != -1 {
			tag = tag[:idx]
		}

		// We only record the config if it isn't the zero value
		if v := v.Field(i); !v.IsZero() {
			result[tag] = v.Interface()
		}
	}

	return result, found
}
