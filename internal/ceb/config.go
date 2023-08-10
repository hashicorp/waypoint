// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ceb

import (
	"context"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"golang.org/x/sys/unix"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/appconfig"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

var (
	// appConfigRefreshPeriod is the interval between checking for new
	// config values. In a steady state, configuration NORMALLY doesn't
	// change so this is set fairly high to avoid unnecessary load on
	// dynamic config sources.
	//
	// NOTE(mitchellh): In the future, we'd like to build a way for
	// config sources to edge-trigger when changes happen to prevent
	// this refresh.
	appConfigRefreshPeriod = 15 * time.Second
)

func (ceb *CEB) initConfigStream(ctx context.Context, cfg *config) error {
	log := ceb.logger.Named("config")

	// Start the watcher. This will do nothing until anything is sent on the
	// channel so we can start it early. We share the same channel across
	// config reconnects.
	ch := make(chan *pb.EntrypointConfig)
	go ceb.watchConfig(ctx, log, cfg, ch)

	// Start the config receiver. This will connect ot the EntrypointConfig
	// endpoint and start receiving data. This will reconnect on failure.
	go ceb.initConfigStreamReceiver(ctx, log, cfg, ch, false)

	return nil
}

func (ceb *CEB) initConfigStreamReceiver(
	ctx context.Context,
	log hclog.Logger,
	cfg *config,
	ch chan<- *pb.EntrypointConfig,
	isRetry bool,
) error {
	// On retry we always mark the child process ready so we can begin executing
	// any staged child command. We don't do this on non-retries because we
	// still have hope that we can talk to the server and get our initial config.
	if isRetry {
		ceb.markChildCmdReady()
	}

	// wait for initial server connection
	serverClient := ceb.waitClient()
	if serverClient == nil {
		return ctx.Err()
	}

	// Open our log stream
	log.Debug("registering instance, requesting config")
	client, err := serverClient.EntrypointConfig(ctx, &pb.EntrypointConfigRequest{
		DeploymentId: ceb.deploymentId,
		InstanceId:   ceb.id,
		DisableExec:  ceb.execDisable,
	}, grpc.WaitForReady(isRetry || cfg.ServerRequired))
	if err != nil {
		// If the server is unavailable and this is our first time, then
		// we just start this up in the background in retry mode and allow
		// the startup to continue so we don't block the child process starting.
		if status.Code(err) == codes.Unavailable {
			log.Error("error connecting to Waypoint server, will retry but startup " +
				"child command without initial settings")
			go ceb.initConfigStreamReceiver(ctx, log, cfg, ch, true)
			return nil
		}

		return err
	}

	// We never send anything
	client.CloseSend()

	// Start the goroutine that waits for all other configs
	go ceb.recvConfig(ctx, client, ch, func() error {
		return ceb.initConfigStreamReceiver(ctx, log, cfg, ch, true)
	})

	return nil
}

// watchConfig sits in a goroutine receiving the new configurations from the
// server.
func (ceb *CEB) watchConfig(
	ctx context.Context,
	log hclog.Logger,
	cfg *config,
	ch <-chan *pb.EntrypointConfig,
) {
	log = log.Named("watcher")

	// We only init the URL service once. In the future, we can do diffing
	// and support automatically reinitializing if the URL service changes.
	didInitURL := false

	// appCfg stores the currently known application configuration (which includes
	// a list of environment vars and config files) we set on the
	// child. We need to store this since we want to launch all exec sessions
	// with the latest/current view on appCfg vars too.
	var appCfg *appconfig.UpdatedConfig

	// Start the app config watcher. This runs in its own goroutine so that
	// stuff like dynamic config fetching doesn't block starting things like
	// exec sessions.
	appCfgCh := make(chan *appconfig.UpdatedConfig)
	w, err := appconfig.NewWatcher(
		appconfig.WithLogger(log),
		appconfig.WithPlugins(ceb.configPlugins),
		appconfig.WithNotify(appCfgCh),
		appconfig.WithRefreshInterval(appConfigRefreshPeriod),
	)
	if err != nil {
		log.Error("error starting app config watcher", "err", err)
		return
	}
	defer w.Close()

	for {
		select {
		case <-ctx.Done():
			log.Warn("exiting, context ended")
			return

		case config := <-ch:
			// TODO(mitchellh): we need to handle changes to the URL settings
			// and stop/restart the URL service gracefully.
			if !didInitURL {
				didInitURL = true

				// If we have URL service configuration, start it. We start this in a goroutine
				// since we don't need to block starting up our application on this.
				if url := config.UrlService; url != nil {
					go func() {
						if err := ceb.initURLService(ctx, cfg.URLServicePort, url); err != nil {
							log.Warn("error starting URL service", "err", err)
						}
					}()
				} else {
					log.Debug("no URL service configuration, will not register with URL service")
				}
			}

			// Start the exec sessions if we have any
			if len(config.Exec) > 0 {
				ceb.startExecGroup(config.Exec, appCfg.EnvVars)
			}

			// Respect any value sent down right away.
			cfg.FileRewriteSignal = config.FileChangeSignal

			// Configure our env vars for the child command. We always send
			// these even if they're nil since the app config watcher will
			// de-dup and we want to handle removing env vars.
			w.UpdateSources(ctx, config.ConfigSources)
			w.UpdateVars(ctx, config.EnvVars)

		case newEnv := <-appCfgCh:
			// Store the new env vars. We could just do `env = <-envCh` above
			// but in my experience its super easy in the future for someone
			// to put a `:=` there and break things. This makes it more explicit.
			appCfg = newEnv

			log.Trace("received new config")

			if appCfg.UpdatedFiles && len(appCfg.Files) > 0 {
				ceb.writeFiles(log, cfg, appCfg)
			}

			if !appCfg.UpdatedEnv {
				log.Trace("updated env did not include new env vars, skipping restart")
				continue
			}

			// Process it for any keys that we handle differently (such as
			// WAYPOINT_LOG_LEVEL)
			ceb.processAppEnv(appCfg.EnvVars)

			// Set our new env vars
			newCmd := ceb.copyCmd(ceb.childCmdBase)
			newCmd.Env = append(newCmd.Env, appCfg.EnvVars...)

			// Note: we purposely do not process appCfg.DeletedEnvVars since
			// our env vars are only used for _new_ child command executions
			// so if they're deleted, they simply won't be present anymore
			// in the next invocations.

			// Restart
			log.Info("env vars changed, sending new child command")
			select {
			case ceb.childCmdCh <- newCmd:
			case <-ctx.Done():
			}

			// Always mark the child command ready at this point. This is
			// a noop if its already done. If its not, then we're ready now
			// because readiness is waiting for that initial set of config.
			ceb.markChildCmdReady()
		}
	}
}

var sigMap = map[string]os.Signal{
	"SIGABRT":   unix.SIGABRT,
	"SIGALRM":   unix.SIGALRM,
	"SIGBUS":    unix.SIGBUS,
	"SIGCHLD":   unix.SIGCHLD,
	"SIGCONT":   unix.SIGCONT,
	"SIGHUP":    unix.SIGHUP,
	"SIGINT":    unix.SIGINT,
	"SIGIO":     unix.SIGIO,
	"SIGKILL":   unix.SIGKILL,
	"SIGPIPE":   unix.SIGPIPE,
	"SIGPROF":   unix.SIGPROF,
	"SIGQUIT":   unix.SIGQUIT,
	"SIGSEGV":   unix.SIGSEGV,
	"SIGSTOP":   unix.SIGSTOP,
	"SIGSYS":    unix.SIGSYS,
	"SIGTERM":   unix.SIGTERM,
	"SIGTRAP":   unix.SIGTRAP,
	"SIGTSTP":   unix.SIGTSTP,
	"SIGTTIN":   unix.SIGTTIN,
	"SIGTTOU":   unix.SIGTTOU,
	"SIGUSR1":   unix.SIGUSR1,
	"SIGUSR2":   unix.SIGUSR2,
	"SIGVTALRM": unix.SIGVTALRM,
	"SIGWINCH":  unix.SIGWINCH,

	"ABRT":   unix.SIGABRT,
	"ALRM":   unix.SIGALRM,
	"BUS":    unix.SIGBUS,
	"CHLD":   unix.SIGCHLD,
	"CONT":   unix.SIGCONT,
	"HUP":    unix.SIGHUP,
	"INT":    unix.SIGINT,
	"IO":     unix.SIGIO,
	"KILL":   unix.SIGKILL,
	"PIPE":   unix.SIGPIPE,
	"PROF":   unix.SIGPROF,
	"QUIT":   unix.SIGQUIT,
	"SEGV":   unix.SIGSEGV,
	"STOP":   unix.SIGSTOP,
	"SYS":    unix.SIGSYS,
	"TERM":   unix.SIGTERM,
	"TRAP":   unix.SIGTRAP,
	"TSTP":   unix.SIGTSTP,
	"TTIN":   unix.SIGTTIN,
	"TTOU":   unix.SIGTTOU,
	"USR1":   unix.SIGUSR1,
	"USR2":   unix.SIGUSR2,
	"VTALRM": unix.SIGVTALRM,
	"WINCH":  unix.SIGWINCH,
}

func (ceb *CEB) writeFiles(log hclog.Logger, cfg *config, env *appconfig.UpdatedConfig) {
	log.Debug("writing application files to disk", "count", len(env.Files))

	var sendSignal bool

	for _, fc := range env.Files {
		err := ioutil.WriteFile(fc.Path, fc.Data, 0644)
		if err != nil {
			log.Error("error writing application file", "error", err, "path", fc.Path)
		} else {
			log.Info("wrote application file to disk", "path", fc.Path)
			sendSignal = true
		}
	}

	if sendSignal && cfg.FileRewriteSignal != "" {
		if sig, ok := sigMap[strings.ToUpper(cfg.FileRewriteSignal)]; ok {
			ceb.childSigCh <- sig
		} else {
			log.Error("unknown signal defined for file restart", "signal", cfg.FileRewriteSignal)
		}
	}
}

func (ceb *CEB) recvConfig(
	ctx context.Context,
	client pb.Waypoint_EntrypointConfigClient,
	ch chan<- *pb.EntrypointConfig,
	reconnect func() error,
) {
	log := ceb.logger.Named("config_recv")
	defer log.Trace("exiting receive goroutine")

	// Keep track of our first receive
	first := true

	for {
		// If the context is closed, exit
		if ctx.Err() != nil {
			return
		}

		// Wait for the next configuration
		resp, err := client.Recv()
		if err != nil {
			// We're disconnected
			ceb.setState(&ceb.stateConfig, false)

			// If we get the unavailable error then the connection died.
			// We restablish the connection.
			if status.Code(err) == codes.Unavailable {
				log.Warn("ceb disconnected from server, attempting reconnect")
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
			ceb.setState(&ceb.stateConfig, true)
		}

		log.Debug("new configuration received")
		ch <- resp.Config
	}
}

// processAppEnv takes a list of env vars meant for the app and handles
// certain special cases (such as WAYPOINT_LOG_LEVEL) that also affect the
// entrypoint.
func (ceb *CEB) processAppEnv(env []string) {
	// Check if we changed our log level. We change this on the
	// root logger for the CEB.
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
				ceb.logger.Warn("log level provided in env var is invalid", value)
			} else {
				// We set the log level on the root logger so it
				// affects all CEB logs.
				ceb.logger.SetLevel(level)
			}
		}
	}
}
