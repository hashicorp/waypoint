package ceb

import (
	"context"
	"time"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/appconfig"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
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

	// env stores the currently known list of environment vars we set on the
	// child. We need to store this since we want to launch all exec sessions
	// with the latest/current view on env vars too.
	var env []string

	// Start the app config watcher. This runs in its own goroutine so that
	// stuff like dynamic config fetching doesn't block starting things like
	// exec sessions.
	envCh := make(chan []string)
	w, err := appconfig.NewWatcher(
		appconfig.WithLogger(log),
		appconfig.WithPlugins(ceb.configPlugins),
		appconfig.WithNotify(envCh),
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
				ceb.startExecGroup(config.Exec, env)
			}

			// Configure our env vars for the child command. We always send
			// these even if they're nil since the app config watcher will
			// de-dup and we want to handle removing env vars.
			w.UpdateSources(ctx, config.ConfigSources)
			w.UpdateVars(ctx, config.EnvVars)

		case newEnv := <-envCh:
			// Store the new env vars. We could just do `env = <-envCh` above
			// but in my experience its super easy in the future for someone
			// to put a `:=` there and break things. This makes it more explicit.
			env = newEnv

			// Set our new env vars
			newCmd := ceb.copyCmd(ceb.childCmdBase)
			newCmd.Env = append(newCmd.Env, env...)

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
				log.Error("ceb disconnected from server, attempting reconnect")
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
