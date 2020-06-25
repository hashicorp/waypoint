package ceb

import (
	"context"
	"fmt"

	"github.com/hashicorp/horizon/pkg/agent"
	"github.com/hashicorp/horizon/pkg/discovery"
	hznpb "github.com/hashicorp/horizon/pkg/pb"
	"github.com/pkg/errors"
)

var ErrURLSetup = errors.New("error configuring url service")

func (ceb *CEB) initURLService(ctx context.Context, cfg *config) error {
	if len(cfg.URLServiceLabels) == 0 {
		return nil
	}

	L := ceb.logger.Named("url")

	g, err := agent.NewAgent(L.Named("agent"))
	if err != nil {
		return errors.Wrapf(err, "error configuring agent")
	}

	g.Token = cfg.URLToken
	target := fmt.Sprintf(":%d", cfg.URLServicePort)

	labels := hznpb.ParseLabelSet(cfg.URLServiceLabels)

	_, err = g.AddService(&agent.Service{
		Type:    "http",
		Labels:  labels,
		Handler: agent.HTTPHandler("http://" + target),
	})

	if err != nil {
		return errors.Wrapf(err, "error registering service")
	}

	L.Debug("discovering hubs")

	dc, err := discovery.NewClient(cfg.URLControlAddr)
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
