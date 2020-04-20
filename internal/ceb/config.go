package ceb

import (
	"context"

	pb "github.com/mitchellh/devflow/internal/server/gen"
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
	ceb.config = resp.Config

	return nil
}
