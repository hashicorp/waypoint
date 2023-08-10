// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package client

import (
	"context"
	"os"
	"runtime"

	"github.com/docker/docker/client"
	"github.com/hashicorp/go-hclog"
)

// Fallback returns true if we should attempt to fallback to a daemonless
// mechanism. If the return value is (false, nil) then we should not fallback
// and Docker appears ready to use. If the return value is (false, non-nil)
// then we should not fallback but Docker does NOT appear healthy.
//
// Note that a return value of true means that we should attempt a fallback,
// but this method doesn't validate that any fallback mechanism is actually
// available.
func Fallback(
	ctx context.Context,
	log hclog.Logger,
	c *client.Client,
) (bool, error) {
	const DockerSocketPath = "/var/run/docker.sock"

	// We always nest ourselves because our logs are annoying (but TRACE)
	log = log.Named("docker_fallback_check")

	// Grab the server version, we do this to attempt a connection. If
	// this succeeds, we always use Docker.
	log.Trace("testing Docker client connection by calling ServerVersion")
	_, err := c.ServerVersion(ctx)
	if err == nil {
		log.Trace("ServerVersion succeeded, will use Docker daemon")
		return false, nil
	}
	if !client.IsErrConnectionFailed(err) {
		// If we got an error other than connection failure, we notify the user
		// because we shouldn't fallback if anything else happened.
		log.Trace("ServerVersion fallback check failed with non-connection error",
			"err", err)
		return false, err
	}

	// If the Docker host is set, the user wants to use Docker, so we never
	// fallback.
	log.Trace("testing DOCKER_HOST value")
	if os.Getenv("DOCKER_HOST") != "" {
		log.Trace("will not fallback because DOCKER_HOST is set")
		return false, err
	}

	// We never fallback on Windows currently because we have no mechanism
	// to build without a Docker daemon on Windows.
	log.Trace("testing GOOS")
	if runtime.GOOS == "windows" {
		log.Trace("will not fallback because GOOS is Windows")
		return false, err
	}

	// If the Docker socket does NOT exist, then fall back.
	log.Trace("testing for Docker socket existence", "path", DockerSocketPath)
	_, staterr := os.Stat(DockerSocketPath)
	if staterr == nil {
		// Docker socket does exist, so let's assume the user wants to use Docker.
		log.Trace("Docker socket exists, will not use fallback")
		return false, err
	}
	if !os.IsNotExist(staterr) {
		log.Trace("error during check for Docker socket", "err", err)
		return false, err
	}

	log.Trace("will fallback, no Docker socket found")
	return true, nil
}
