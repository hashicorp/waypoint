package core

import (
	"context"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	// TODO(briancain): do we need a Multi upsert?
	p.logger.Debug("syncing pipeline config")
	for _, pipeline := range pipelines {
		pipeOwner, ok := pipeline.Owner.(*pb.Pipeline_Project)
		if !ok {
			return status.Error(codes.FailedPrecondition,
				"failed to determine pipeline owner. Only expecting pipeline owner to be a Pipeline_Project type.")
		}

		// Look up if this pipeline already exists by Pipeline Owner (project and pipeline name)
		// If so, use the existing id to upsert into an existing config. Otherwise
		// create a new pipeline entry.
		pipeResp, err := p.client.GetPipeline(ctx, &pb.GetPipelineRequest{
			Pipeline: &pb.Ref_Pipeline{
				Ref: &pb.Ref_Pipeline_Owner{
					Owner: &pb.Ref_PipelineOwner{
						Project:      pipeOwner.Project,
						PipelineName: pipeline.Name,
					},
				},
			},
		})
		if err != nil {
			// If not found, that's fine, we're going to upsert it
			if status.Code(err) != codes.NotFound {
				return err
			}
		}
		if pipeResp != nil {
			pipeline.Id = pipeResp.Pipeline.Id
			p.logger.Trace("existing pipeline already in db, using existing id", "pipeline_id", pipeline.Id)
		}

		_, err = p.client.UpsertPipeline(ctx, &pb.UpsertPipelineRequest{
			Pipeline: pipeline,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
