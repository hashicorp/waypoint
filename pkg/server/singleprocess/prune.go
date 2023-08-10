// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package singleprocess

import (
	"context"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint/pkg/serverstate"
)

func (s *Service) runPrune(
	ctx context.Context,
	wg *sync.WaitGroup,
	funclog hclog.Logger,
) {
	defer wg.Done()

	funclog.Info("starting")
	defer funclog.Info("exiting")

	pruner, ok := s.state(ctx).(serverstate.Pruner)
	if !ok {
		funclog.Info("state background doesn't require pruning")
		return
	}

	tk := time.NewTicker(10 * time.Minute)
	defer tk.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tk.C:
			err := pruner.Prune()
			if err != nil {
				funclog.Error("error pruning data", "error", err)
			}
		}
	}
}
