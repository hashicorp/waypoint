package config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/waypoint-plugin-sdk/component"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// Pipeline represents a single pipeline definition
type Pipeline struct {
	Name string `hcl:",label"`

	StepRaw []*hclStep `hcl:"step,block"`
	Steps   []*Step

	ctx    *hcl.EvalContext
	config *Config
}

// hclPipeline represents a raw HCL version of a pipeline config
type hclPipeline struct {
	Name string `hcl:",label"`

	// We need these raw values to determine the plugins need to be used.
	StepRaw []*hclStep `hcl:"step,block"`

	Body   hcl.Body `hcl:",body"`
	Remain hcl.Body `hcl:",remain"`
}

// hclStep represents a raw HCL version of a step stanza in a pipeline config
type hclStep struct {
	Name string `hcl:",label"`

	// If set, this step will depend on the defined step. The default step
	// will be the previously defined step in order that it was defined
	// in a waypoint.hcl
	DependsOn []string `hcl:"depends_on,optional"`

	// The OCI image to use for executing this step
	ImageURL string `hcl:"image_url,optional"`

	// The plugin to use for this Step
	Use *Use `hcl:"use,block"`
}

// Pipelines returns the id of all the defined pipelines
func (c *Config) Pipelines() []string {
	var result []string
	for _, p := range c.hclConfig.Pipelines {
		result = append(result, p.Name)
	}

	return result
}

// Pipeline returns the configured pipeline named n. If the pipeline doesn't
// exist, this will return (nil, nil).
func (c *Config) Pipeline(id string, ctx *hcl.EvalContext) (*Pipeline, error) {
	ctx = appendContext(c.ctx, ctx)

	// Find the pipeline by progressively decoding
	var rawPipeline *hclPipeline
	for _, p := range c.hclConfig.Pipelines {
		if p.Name == id {
			rawPipeline = p
			break
		}
	}
	if rawPipeline == nil {
		return nil, nil
	}

	// Full decode
	var pipeline Pipeline
	if diag := gohcl.DecodeBody(rawPipeline.Body, finalizeContext(ctx), &pipeline); diag.HasErrors() {
		return nil, diag
	}
	pipeline.Name = rawPipeline.Name
	pipeline.ctx = ctx
	pipeline.config = c
	if pipeline.config != nil {
		pipeline.config.ctx = ctx
	}

	// decode all of the defined raw steps for a pipeline
	var steps []*Step
	for _, stepRaw := range pipeline.StepRaw {
		// turn stepRaw into a staged Step
		s := Step{
			ctx:       ctx,
			Name:      stepRaw.Name,
			DependsOn: stepRaw.DependsOn,
			ImageURL:  stepRaw.ImageURL,
			Use:       stepRaw.Use,
		}

		steps = append(steps, &s)
	}
	pipeline.Steps = steps

	return &pipeline, nil
}

// Ref returns the ref for this pipeline.
func (c *Pipeline) Ref() *pb.Ref_Pipeline {
	return &pb.Ref_Pipeline{
		Ref: &pb.Ref_Pipeline_Id{
			Id: &pb.Ref_PipelineId{
				Id: c.Name,
			},
		},
	}
}

// Configure configures the plugin for a given Step with the use body of this operation.
func (s *Step) Configure(plugin interface{}, ctx *hcl.EvalContext) hcl.Diagnostics {
	ctx = appendContext(s.ctx, ctx)
	ctx = finalizeContext(ctx)

	return component.Configure(plugin, s.Use.Body, ctx)
}
