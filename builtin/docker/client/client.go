// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package client

import (
	"net/http"

	"github.com/docker/cli/cli/connhelper"
	"github.com/docker/docker/client"
)

// NewClientWithOpts wraps Docker's NewClientWithOpts with withConnectionHelper
func NewClientWithOpts(ops ...client.Opt) (*client.Client, error) {
	ops = append(ops, withConnectionHelper)
	return client.NewClientWithOpts(ops...)
}

// withConnectionHelper applies a Docker-specific connection helper (concept from the
// Docker CLI) for a given daemon host. As an example, a connection helper makes it
// possible to use the client given a DOCKER_HOST with an ssh scheme.
func withConnectionHelper(c *client.Client) error {
	host := c.DaemonHost()
	helper, err := connhelper.GetConnectionHelper(host)
	if err != nil {
		return err
	}

	if helper == nil {
		return nil
	}
	httpClient := &http.Client{
		// No tls
		// No proxy
		Transport: &http.Transport{
			DialContext: helper.Dialer,
		},
	}

	opts := []client.Opt{
		client.WithHTTPClient(httpClient),
		client.WithHost(helper.Host),
		client.WithDialContext(helper.Dialer),
	}

	// Apply options
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return err
		}
	}

	return nil
}
