package ceb

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (ceb *CEB) init(ctx context.Context, cfg *config, retry bool) error {
	log := ceb.logger.Named("init")

RETRY_INIT:
	// First thing we need to do is connect to the server.
	if ceb.client == nil {
		if cfg.ServerAddr == "" {
			log.Info("no waypoint server configured, disabled entrypoint")
			return nil
		}

		err := ceb.dialServer(ctx, cfg, retry)
		if status.Code(err) == codes.Unavailable {
			// If we require a server connection, then just retry.
			if cfg.ServerRequired {
				log.Warn("server unavailable but ceb configured to require it, retrying synchronously")
				retry = true
				goto RETRY_INIT
			}

			// If we don't require a server connection, then we start a
			// goroutine to retry and eventually connect (hopefully).
			log.Warn("server unavailable, will retry in the background")
			go ceb.init(ctx, cfg, true)

			return nil
		}
	}

	// This should never happen
	if ceb.client == nil {
		log.Error("client is still nil, not expected, quitting init")
		return nil
	}

	// Get our configuration and start the long-running stream for it.
	if err := ceb.initConfigStream(ctx, cfg, retry); err != nil {
		return err
	}

	// Initialize our log stream
	// NOTE(mitchellh): at some point we want this to be configurable
	// but for now we're just going for it.
	if err := ceb.initLogStream(ctx, cfg); err != nil {
		return err
	}

	return nil
}
