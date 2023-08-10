// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package singleprocess

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/go-hclog"
	grpctoken "github.com/hashicorp/horizon/pkg/grpc/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"

	wphznpb "github.com/hashicorp/waypoint-hzn/pkg/pb"

	"github.com/hashicorp/waypoint/internal/pkg/grpcready"
	"github.com/hashicorp/waypoint/internal/serverconfig"
)

// urlClient returns the URL service client. This may return nil if the
// client isn't yet ready.
func (s *Service) urlClient() wphznpb.WaypointHznClient {
	s.urlClientMu.Lock()
	defer s.urlClientMu.Unlock()
	return s.urlClientVal
}

// initURLClient initializes the URL service client. This will get a guest
// account token if necessary, and will retry in the background on failure.
// If the URL service is not available, then `s.urlClient()` may be nil.
//
// Prior to retrying in the background, this attempts once to synchronously
// connect and initialize the client. Therefore, the happy path expects
// that `urlClient()` never returns nil if the URL service is enabled.
func (s *Service) initURLClient(
	log hclog.Logger,
	bo backoff.BackOff,
	acceptURLTerms bool,
	cfg *serverconfig.URL,
) error {
	if cfg == nil || !cfg.Enabled {
		log.Info("URL service is not configured or explicitly disabled")
		return nil
	}

	// We are in a retry only if we have backoff settings set.
	isRetry := bo != nil

	// We don't currently have a context here to thread through. The
	// remainder of the logic is context-aware so when we do have one, we
	// just need to replace this and everything will just work.
	ctx := context.Background()

	// If we aren't retrying, setup our backoff settings.
	if !isRetry {
		bo = backoff.NewExponentialBackOff()
		bo = backoff.WithContext(bo, ctx)
	}

	// Perform the blocking attempt to initialize.
	err := s.initURLClientBlocking(ctx, log, isRetry, acceptURLTerms, cfg)
	if status.Code(err) == codes.Unavailable {
		// Sleep the backoff duration
		boSleep := bo.NextBackOff()
		if boSleep == backoff.Stop {
			log.Warn("URL service unavailable, backoff canceled")
			return status.New(codes.DeadlineExceeded, "URL service reconnect timed out").Err()
		}
		log.Warn("URL service unavailable, will retry in the background",
			"sleep", boSleep.String())
		time.Sleep(boSleep)

		// Start a goroutine to keep retrying in the background.
		go s.initURLClient(log, bo, acceptURLTerms, cfg)
		return nil
	}

	if err != nil {
		log.Warn("failed to initialize URL service", "err", err)
	} else {
		log.Info("URL service client successfully initialized")
	}

	return err
}

func (s *Service) initURLClientBlocking(
	ctx context.Context,
	log hclog.Logger,
	isRetry bool,
	acceptURLTerms bool,
	cfg *serverconfig.URL,
) error {
	// If we have no API token, get our guest account token.
	if cfg.APIToken == "" {
		log.Debug("API token not set in config, initializing guest account")
		token, err := s.initURLGuestAccount(ctx, log, isRetry, acceptURLTerms, cfg)
		if err != nil {
			return err
		}

		// Set the API token, if logic later in this func fails and we retry
		// we will reuse the API token we already have.
		cfg.APIToken = token

		// Set our URL CEB settings. It is always initialized with the API
		// token if it is set so we only have to do this on this code path.
		s.urlCEBMu.Lock()
		s.urlCEB.Token = token
		close(s.urlCEBWatchCh) // notify any watchers we have changes
		s.urlCEBWatchCh = make(chan struct{})
		s.urlCEBMu.Unlock()
	}

	// Now that we have a token, connect to the API service with that token.
	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(grpctoken.Token(cfg.APIToken)),
	}
	if cfg.APIInsecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	}

	conn, err := grpc.Dial(cfg.APIAddress, opts...)
	if err != nil {
		return err
	}
	if err := grpcready.Conn(ctx, log, conn, isRetry); err != nil {
		return err
	}

	s.urlClientMu.Lock()
	defer s.urlClientMu.Unlock()
	s.urlClientVal = wphznpb.NewWaypointHznClient(conn)
	return nil
}

func (s *Service) initURLGuestAccount(
	ctx context.Context,
	log hclog.Logger,
	isRetry bool,
	acceptURLTerms bool,
	cfg *serverconfig.URL,
) (string, error) {
	// Check if URL Token already exists, if so, no reason to
	// re-register and generate a new hostname
	urlToken, err := s.state(ctx).ServerURLTokenGet(ctx)
	if err != nil {
		return "", err
	} else if urlToken != "" {
		log.Debug("using saved URL guest token")
		return urlToken, nil
	}

	// Connect without auth to our API client
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTimeout(10*time.Second))
	if cfg.APIInsecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		// If it isn't insecure, then we have to specify that we're using TLS
		opts = append(opts, grpc.WithTransportCredentials(
			credentials.NewTLS(&tls.Config{}),
		))
	}

	log.Debug("connecting to URL service to retrieve guest token",
		"addr", cfg.APIAddress,
		"tls", !cfg.APIInsecure,
	)
	conn, err := grpc.DialContext(ctx, cfg.APIAddress, opts...)
	if err != nil {
		// This error should NOT happen on connection failure, since
		// we connect in the background. This should happen if there is
		// a local configuration error.
		log.Warn("failed to connect to the URL service", "err", err)
		return "", err
	}

	// Verify we're connected. We do this loop so that we can
	// fail fast in the non-isRetry case.
	log.Debug("waiting on server connection state to become ready")
	if err := grpcready.Conn(ctx, log, conn, isRetry); err != nil {
		return "", err
	}

	// Init our client
	client := wphznpb.NewWaypointHznClient(conn)

	// Request a guest account
	accountResp, err := client.RegisterGuestAccount(
		context.Background(),
		&wphznpb.RegisterGuestAccountRequest{
			ServerId:  s.id,
			AcceptTos: acceptURLTerms,
		},
	)
	if err != nil {
		return "", err
	}

	if err := s.state(ctx).ServerURLTokenSet(ctx, accountResp.Token); err != nil {
		return "", err
	}

	return accountResp.Token, nil
}
