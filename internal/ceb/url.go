// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ceb

import (
	"context"
	"fmt"

	"github.com/certifi/gocertifi"
	"github.com/hashicorp/horizon/pkg/agent"
	"github.com/hashicorp/horizon/pkg/discovery"
	hznpb "github.com/hashicorp/horizon/pkg/pb"
	"github.com/pkg/errors"

	"github.com/hashicorp/waypoint/internal/version"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

var ErrURLSetup = errors.New("error configuring url service")

const (
	urlLabelVersion  = "waypoint.hashicorp.com/ceb-version"
	urlLabelRevision = "waypoint.hashicorp.com/ceb-revision"
)

func (ceb *CEB) initURLService(ctx context.Context, port int, cfg *pb.EntrypointConfig_URLService) error {
	if cfg == nil || len(cfg.Labels) == 0 {
		return nil
	}

	L := ceb.logger.Named("url")

	ceb.urlAgentMu.Lock()
	defer ceb.urlAgentMu.Unlock()

	if ceb.urlAgentCancel != nil {
		L.Debug("detected old agent, requesting it close")
		ceb.urlAgentCancel()
	}

	ceb.urlAgentCtx, ceb.urlAgentCancel = context.WithCancel(ctx)

	ctx = ceb.urlAgentCtx

	L.Debug("url service enabled, configuring",
		"addr", cfg.ControlAddr,
		"service_port", port,
		"labels", cfg.Labels,
	)

	g, err := agent.NewAgent(L.Named("agent"))
	if err != nil {
		// NewAgent should never fail, this just sets some fields and returns
		// a struct. Therefore, we won't ever retry on this.
		return errors.Wrapf(err, "error configuring agent")
	}
	g.Token = cfg.Token

	// Setup the Mozilla CA cert bundle. We can ignore the error because
	// this never fails, it only returns an error for backwards compat reasons.
	g.RootCAs, _ = gocertifi.CACerts()

	// Parse our labels and add some additional labels using local data
	vsn := version.GetVersion()
	labels := hznpb.ParseLabelSet(cfg.Labels).
		Add(urlLabelVersion, vsn.FullVersionNumber(true)).
		Add(urlLabelRevision, vsn.Revision)

	// Add our service to route to.
	target := fmt.Sprintf(":%d", port)
	_, err = g.AddService(&agent.Service{
		Type:    "http",
		Labels:  labels,
		Handler: agent.HTTPHandler("http://" + target),
	})
	if err != nil {
		// This can also never fail.
		return errors.Wrapf(err, "error registering service")
	}

	L.Debug("discovering hubs")
	dc, err := discovery.NewClient(cfg.ControlAddr)
	if err != nil {
		// This shouldn't fail so we don't have to retry at the time of writing.
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
