// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package core

import (
	"context"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// ConfigSync will evaluate the current hcl config for a given Pipeline and
// upsert the proto version based on the current evaluation
func (p *Pipeline) ConfigSync(ctx context.Context) error {
	// TODO(briancain): In the future, we can sync config vars for pipelines here

	// Sync the pipeline metadata
	p.logger.Debug("evaluating pipeline configs for syncing")
	pipelines, err := p.config.Config.PipelineProtos()
	if err != nil {
		return err
	}

	p.logger.Debug("syncing pipeline config")
	for _, pipeline := range pipelines {
		// TODO(briancain): We might want a multi-upsert for pipelines if this
		// request gets out of hand.
		_, err = p.client.UpsertPipeline(ctx, &pb.UpsertPipelineRequest{
			Pipeline: pipeline,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
