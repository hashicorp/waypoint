package core

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/config"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// Pipeline represents a single pipeline and exposes all of the operations
// that can be performed on a pipeline.
type Pipeline struct {
	// UI is the UI that should be used for any output that is specific
	// to this pipeline vs the project or app UI.
	UI terminal.UI

	client pb.WaypointClient
	logger hclog.Logger

	config *config.Pipeline

	// References that ties this pipeline to a specific Waypoint app
	wsRef   *pb.Ref_Workspace
	ref     *pb.Ref_Pipeline
	project *Project
}

// newPipeline creates a Pipeline for the given application and configuration. This
// will initialize and configure all of the components of this pipeline. An error
// will be returned if this pipeline fails to initialize: configuration is invalid,
// a component could not be found, etc.
func newPipeline(
	ctx context.Context,
	project *Project,
	config *config.Pipeline,
) (*Pipeline, error) {

	pipeline := &Pipeline{
		wsRef:   project.WorkspaceRef(),
		project: project,

		client: project.client,
		logger: project.logger.Named("pipeline").Named(config.Name),
		config: config,
		UI:     project.UI,

		ref: &pb.Ref_Pipeline{
			Ref: &pb.Ref_Pipeline_Id{
				Id: &pb.Ref_PipelineId{
					Id: config.Name,
				},
			},
		},
	}

	return pipeline, nil
}

// Name() returns the name of the given pipeline
func (p *Pipeline) Name() string {
	return p.config.Name
}

// Ref returns the reference to this pipeline for us in API calls.
func (p *Pipeline) Ref() *pb.Ref_Pipeline {
	return p.ref
}
