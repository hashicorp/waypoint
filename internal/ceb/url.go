package ceb

import (
	"context"
	"fmt"

	"github.com/hashicorp/horizon/pkg/agent"
	"github.com/hashicorp/horizon/pkg/discovery"
	hznpb "github.com/hashicorp/horizon/pkg/pb"
	"github.com/pkg/errors"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

var ErrURLSetup = errors.New("error configuring url service")

func (ceb *CEB) initURLService(ctx context.Context, port int, cfg *pb.EntrypointConfig_URLService) error {
	if cfg == nil || len(cfg.Labels) == 0 {
		return nil
	}

	L := ceb.logger.Named("url")
	L.Info("url service enabled, configuring",
		"addr", cfg.ControlAddr,
		"service_port", port,
		"labels", cfg.Labels,
	)

	g, err := agent.NewAgent(L.Named("agent"))
	if err != nil {
		return errors.Wrapf(err, "error configuring agent")
	}

	g.Token = cfg.Token
	target := fmt.Sprintf(":%d", port)

	labels := hznpb.ParseLabelSet(cfg.Labels)

	_, err = g.AddService(&agent.Service{
		Type:    "http",
		Labels:  labels,
		Handler: agent.HTTPHandler("http://" + target),
	})

	if err != nil {
		return errors.Wrapf(err, "error registering service")
	}

	L.Debug("discovering hubs")

	dc, err := discovery.NewClient(cfg.ControlAddr)
	if err != nil {
		return errors.Wrapf(err, "error conecting to waypoint control service")
	}

	L.Debug("refreshing data")

	err = dc.Refresh(ctx)
	if err != nil {
		return errors.Wrapf(err, "error discovering network endpoints")
	}

	err = g.Start(ctx, dc)
	if err != nil {
		return errors.Wrapf(err, "error serving traffic")
	}

	go func() {
		err := g.Wait(ctx)
		if err != nil {
			L.Error("error in background connection to url service", "error", err)
		}
	}()

	return nil
}
