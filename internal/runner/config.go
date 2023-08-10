// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package runner

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/appconfig"
	"github.com/hashicorp/waypoint/internal/clierrors"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// appConfigRefreshPeriod is the interval between checking for new
// config values. In a steady state, configuration NORMALLY doesn't
// change so this is set fairly high to avoid unnecessary load on
// dynamic config sources.
//
// NOTE(mitchellh): In the future, we'd like to build a way for
// config sources to edge-trigger when changes happen to prevent
// this refresh.
var appConfigRefreshPeriod = 15 * time.Second

// initConfigStream starts the RunnerConfig stream in the background. This will
// automatically retry if the server is currently unavailable or goes down
// mid-stream.
func (r *Runner) initConfigStream(ctx context.Context) error {
	log := r.logger.Named("config")

	// Start the watcher. This will do nothing until anything is sent on the
	// channel so we can start it early. We share the same channel across
	// config reconnects.
	ch := make(chan *pb.RunnerConfig)
	go r.watchConfig(ctx, log, ch)

	// Start the config receiver. This will connect to the RunnerConfig
	// endpoint and start receiving data. This will reconnect on failure.
	go r.initConfigStreamReceiver(ctx, log, ch, false)

	return nil
}

func (r *Runner) initConfigStreamReceiver(
	ctx context.Context,
	log hclog.Logger,
	ch chan<- *pb.RunnerConfig,
	isRetry bool,
) error {
	// Open our log stream
	log.Debug("registering instance, requesting config")
	client, err := r.client.RunnerConfig(ctx, grpc.WaitForReady(isRetry))
	if err != nil {
		// If the connection failed, we'll just log that and retry in the
		// background so the remainder of runner startup can continue.
		if status.Code(err) == codes.Unavailable {
			log.Error("error connecting to Waypoint server, will retry")
			go r.initConfigStreamReceiver(ctx, log, ch, true)
			return nil
		}

		return err
	}

	// Start the goroutine that receives messages from the stream.
	// This will detect any stream disconections and perform a retry.
	go r.recvConfig(ctx, client, ch, func() error {
		return r.initConfigStreamReceiver(ctx, log, ch, true)
	})

	// Send our open request.
	if err := client.Send(&pb.RunnerConfigRequest{
		Event: &pb.RunnerConfigRequest_Open_{
			Open: &pb.RunnerConfigRequest_Open{
				Runner: r.runner,
			},
		},
	}); err != nil {
		if status.Code(err) == codes.Unavailable {
			// Ignore this error, recvConfig will get it too and will
			// trigger the reconnect.
			return nil
		}

		return err
	}

	return nil
}

// watchConfig sits in a goroutine receiving the new configurations from the
// server.
func (r *Runner) watchConfig(
	ctx context.Context,
	log hclog.Logger,
	ch <-chan *pb.RunnerConfig,
) {
	log = log.Named("watcher")
	defer log.Trace("exiting goroutine")

	// Start the app config watcher. This runs in its own goroutine so that
	// stuff like dynamic config fetching doesn't block starting things like
	// exec sessions.
	appCfgCh := make(chan *appconfig.UpdatedConfig)
	w, err := appconfig.NewWatcher(
		appconfig.WithLogger(log),
		appconfig.WithNotify(appCfgCh),
		appconfig.WithRefreshInterval(appConfigRefreshPeriod),
		appconfig.WithOriginalEnv(os.Environ()),
		appconfig.WithPlugins(r.configPlugins),
	)
	if err != nil {
		log.Error("error starting app config watcher", "err", err)
		return
	}
	defer w.Close()

	for {
		select {
		case <-ctx.Done():
			log.Warn("exiting due to context ended")
			return

		case config := <-ch:
			// Update our app config watcher with the latest vars and sources.
			w.UpdateSources(ctx, config.ConfigSources)
			w.UpdateVars(ctx, config.ConfigVars)

		case appCfg := <-appCfgCh:
			log.Trace("received new app config")

			// This will keep track of our error setting config. For the CEB,
			// app config is best effort: if it isn't set, we ignore. For the
			// runner, we actually exit on config errors so that we don't accept
			// any jobs.
			var merr error

			if appCfg.UpdatedFiles {
				merr = r.writeFiles(log, appCfg)
			}

			if !appCfg.UpdatedEnv {
				log.Trace("updated env did not include new env vars, skipping restart")
				continue
			}

			// Process it for any keys that we handle differently (such as
			// WAYPOINT_LOG_LEVEL)
			r.processAppEnv(log, appCfg.EnvVars)

			// Set our env vars
			if err := r.setEnv(log, appCfg); err != nil {
				merr = multierror.Append(merr, err)
			}

			// Note that we processed our config at least once. People can
			// wait on this state to know that success or fail, one config
			// was received.
			r.incrState(&r.stateConfigOnce)

			// If we have an error, then we exit our loop. For runners,
			// not being able to set config is a fatal error. We do this so
			// that we don't run any jobs in a broken state.
			if merr != nil {
				log.Warn("error setting app config for runner", "err", merr)
				return
			}

		}
	}
}

func (r *Runner) writeFiles(log hclog.Logger, env *appconfig.UpdatedConfig) error {
	// If we have no files, then do nothing.
	if len(env.Files) == 0 {
		return nil
	}

	var result error
	log.Debug("writing app config files to disk", "count", len(env.Files))
	for _, fc := range env.Files {
		err := ioutil.WriteFile(fc.Path, fc.Data, 0644)
		if err != nil {
			log.Error("error writing app config file", "error", err, "path", fc.Path)
			result = multierror.Append(result, err)
		} else {
			log.Info("wrote app config file to disk", "path", fc.Path)
		}
	}

	return result
}

// setEnv sets and unsets the proper env vars for an appconfig update bundle.
func (r *Runner) setEnv(log hclog.Logger, appCfg *appconfig.UpdatedConfig) error {
	var merr error

	// Set our env vars
	for _, str := range appCfg.EnvVars {
		idx := strings.Index(str, "=")
		if idx == -1 {
			continue
		}

		log.Trace("setting env var", "key", str[:idx])
		if err := os.Setenv(str[:idx], str[idx+1:]); err != nil {
			merr = multierror.Append(merr, err)
		}
	}

	// Unset our deleted values. We can do this after setting because
	// appconfig guarantees that this does not contain anything in env vars.
	for _, str := range appCfg.DeletedEnvVars {
		log.Trace("unsetting env var", "key", str)
		if err := os.Unsetenv(str); err != nil {
			merr = multierror.Append(merr, err)
		}
	}

	return merr
}

// processAppEnv takes a list of env vars meant for the app and handles
// certain special cases (such as WAYPOINT_LOG_LEVEL) that also affect the
// runner.
func (r *Runner) processAppEnv(log hclog.Logger, env []string) {
	// Check if we changed our log level. We change this on the
	// root logger for the runner.
	for _, pair := range env {
		idx := strings.Index(pair, "=")
		if idx == -1 {
			// Shouldn't happen
			continue
		}

		key := pair[:idx]
		if key == envLogLevel {
			value := pair[idx+1:]
			level := hclog.LevelFromString(value)
			if level == hclog.NoLevel {
				// We warn this
				log.Warn("log level provided in env var is invalid", value)
			} else {
				// We set the log level on the root logger so it
				// affects all runner logs.
				r.logger.SetLevel(level)
			}
		}
	}
}

func (r *Runner) recvConfig(
	ctx context.Context,
	client pb.Waypoint_RunnerConfigClient,
	ch chan<- *pb.RunnerConfig,
	reconnect func() error,
) {
	log := r.logger.Named("config_recv")
	defer log.Trace("exiting receive goroutine")

	// Keep track of our first receive
	first := true

	// Any reason we exit, this client is done so we mark we're done sending, too.
	defer client.CloseSend()

	for {
		// If the context is closed, exit
		if ctx.Err() != nil {
			return
		}

		// Wait for the next configuration
		resp, err := client.Recv()

		if err != nil {
			// EOF means a graceful close, don't reconnect.
			if err == io.EOF || clierrors.IsCanceled(err) {
				log.Warn("EOF or cancellation received, graceful close of runner config stream")
				return
			}

			// If we get the unavailable error then the connection died.
			// We restablish the connection.
			if status.Code(err) == codes.Unavailable {
				log.Error("runner disconnected from server, attempting reconnect")
				err = reconnect()

				// If we successfully reconnected, then exit this.
				if err == nil {
					return
				}
			}

			log.Error("error receiving configuration, exiting", "err", err)
			return
		}

		// If this is our first receive, then mark that we're connected.
		if first {
			log.Debug("first config received, switching config state to true")
			first = false
			r.incrState(&r.stateConfig)
		}

		log.Info("new configuration received")
		ch <- resp.Config
	}
}
