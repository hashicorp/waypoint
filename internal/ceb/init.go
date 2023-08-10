// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ceb

import (
	"context"
)

func (ceb *CEB) init(ctx context.Context, cfg *config, retry bool) error {
	log := ceb.logger.Named("init")

	// If the entrypoint is full disabled, just connect to our command.
	if cfg.disable {
		// Send our initial child command down.
		ceb.childCmdCh <- ceb.copyCmd(ceb.childCmdBase)
		ceb.markChildCmdReady()
		return nil
	}

	// Initialize our client. This will retry in the background on failure.
	if err := ceb.initClient(ctx, log, cfg, false); err != nil {
		return err
	}

	// Initialize our log stream. We do this first so we can set the proper
	// stdout/stderr on our base child command.
	// NOTE(mitchellh): at some point we want this to be configurable
	// but for now we're just going for it.
	if err := ceb.initLogStream(ctx, cfg); err != nil {
		return err
	}

	// Send our initial child command down.
	ceb.childCmdCh <- ceb.copyCmd(ceb.childCmdBase)

	// Get our configuration and start the long-running stream for it.
	// Goroutine since this requires the client and will wait for it.
	if err := ceb.initConfigStream(ctx, cfg); err != nil {
		return err
	}

	return nil
}
