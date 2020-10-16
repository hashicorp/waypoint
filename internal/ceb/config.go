package ceb

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func (ceb *CEB) initConfigStream(ctx context.Context, cfg *config, isRetry bool) error {
	log := ceb.logger.Named("config")

	// Open our log stream
	log.Debug("registering instance, requesting config")
	client, err := ceb.client.EntrypointConfig(ctx, &pb.EntrypointConfigRequest{
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
			go ceb.initConfigStream(ctx, cfg, true)
			return nil
		}

		return err
	}
	ceb.cleanup(func() { client.CloseSend() })

	// Receive our first configuration which marks that we've registered,
	// plus we need the config for behavior.
	log.Trace("config stream connected, waiting for first config")
	resp, err := client.Recv()
	if err != nil {
		return err
	}
	log.Trace("first config received")

	// Modify childCmd to contain any passed variables as environment variables
	for _, cv := range resp.Config.EnvVars {
		ceb.childCmd.Env = append(ceb.childCmd.Env, cv.Name+"="+cv.Value)
	}

	// If we have URL service configuration, start it. We start this in a goroutine
	// since we don't need to block starting up our application on this.
	if url := resp.Config.UrlService; url != nil {
		go func() {
			if err := ceb.initURLService(ctx, cfg.URLServicePort, url); err != nil {
				log.Warn("error starting URL service", "err", err)
			}
		}()
	} else {
		log.Debug("no URL service configuration, will not register with URL service")
	}

	// Start the watcher
	ch := make(chan *pb.EntrypointConfig)
	go ceb.watchConfig(ch)

	// Send the first config which will trigger setup
	ch <- resp.Config

	// Start the goroutine that waits for all other configs
	go ceb.recvConfig(ctx, client, ch, func() error {
		return ceb.initConfigStream(ctx, cfg, true)
	})

	return nil
}

// watchConfig sits in a goroutine receiving the new configurations from the
// server.
func (ceb *CEB) watchConfig(ch <-chan *pb.EntrypointConfig) {
	for config := range ch {
		// Start the exec sessions if we have any
		if len(config.Exec) > 0 {
			ceb.startExecGroup(config.Exec)
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
	defer close(ch)

	for {
		// If the context is closed, exit
		if ctx.Err() != nil {
			return
		}

		// Wait for the next configuration
		resp, err := client.Recv()
		if err != nil {
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

			if err != nil {
				log.Error("error receiving configuration, exiting", "err", err)
				return
			}
		}

		log.Info("new configuration received")
		ch <- resp.Config
	}
}
