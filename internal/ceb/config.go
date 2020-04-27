package ceb

import (
	"context"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func (ceb *CEB) initConfigStream(ctx context.Context, cfg *config) error {
	log := ceb.logger.Named("config")

	// Open our log stream
	log.Debug("registering instance, requesting config")
	client, err := ceb.client.EntrypointConfig(ctx, &pb.EntrypointConfigRequest{
		DeploymentId: cfg.DeploymentId,
		InstanceId:   ceb.id,
	})
	if err != nil {
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

	// Start the watcher
	ch := make(chan *pb.EntrypointConfig)
	go ceb.watchConfig(ch)
	ceb.cleanup(func() { close(ch) })

	// Send the first config which will trigger setup
	ch <- resp.Config

	// Start the goroutine that waits for all other configs
	go ceb.recvConfig(ctx, client, ch)

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
			log.Error("error receiving configuration, exiting", "err", err)
			return
		}

		log.Info("new configuration received")
		ch <- resp.Config
	}
}
