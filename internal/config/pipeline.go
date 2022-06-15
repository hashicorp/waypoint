package config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// Pipeline represents a single pipeline definition
type Pipeline struct {
	Name string `hcl:",label"`

	StepRaw []*hclStep `hcl:"step,block"`
	Steps   []*Step

	ctx    *hcl.EvalContext
	Config *Config
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
	pipeline.Config = c
	if pipeline.Config != nil {
		pipeline.Config.ctx = ctx
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

// PipelineProtos will take the existing HCL eval context, eval the config
// and translate the HCL result into a Pipeline Proto to be returned for
// operations such as ConfigSync.
func (c *Config) PipelineProtos() ([]*pb.Pipeline, error) {
	if c == nil {
		// This is likely an internal error if this happens.
		panic("attempted to construct pipeline proto on a nil genericConfig")
	}

	// Load HCL config and convert to a Pipeline proto
	var result []*pb.Pipeline
	for _, pl := range c.hclConfig.Pipelines {
		pipe := &pb.Pipeline{
			Name: pl.Name,
			Owner: &pb.Pipeline_Project{
				Project: &pb.Ref_Project{
					Project: c.hclConfig.Project,
				},
			},
		}

		// TODO(briancain): Use the default ODR image for image_url so that it isn't required
		// We do this already maybe in the ODR task launcher.
		steps := make(map[string]*pb.Pipeline_Step)
		for i, step := range pl.StepRaw {
			s := &pb.Pipeline_Step{
				Name:      step.Name,
				DependsOn: step.DependsOn,
				Image:     step.ImageURL, //TODO(briancain): actually use this when executing steps
			}

			// If no dependency was explictily set, we rely on the previous step
			if i != 0 && len(step.DependsOn) == 0 {
				s.DependsOn = []string{pl.StepRaw[i-1].Name}
			}

			// We currently only support one kind of Step plugin. But in the future
			// maybe this would be a switch on step.Type? Or maybe we get our
			// own Step Eval func that returns the step proto instead and does the
			// switches there.
			// NOTE(briancain): This is what you'd change to support future Step plugins
			switch step.Use.Type {
			case "build":
				var buildBody struct {
					DisablePush bool `hcl:"disable_push,optional"`
				}

				if diag := gohcl.DecodeBody(step.Use.Body, finalizeContext(c.ctx), &buildBody); diag.HasErrors() {
					return nil, diag
				}

				s.Kind = &pb.Pipeline_Step_Build_{
					Build: &pb.Pipeline_Step_Build{
						DisablePush: buildBody.DisablePush,
					},
				}

			case "exec":
				var execBody struct {
					Command string   `hcl:"command,optional"`
					Args    []string `hcl:"args,optional"`
				}

				// Evaluate the step body hcl to get options
				if diag := gohcl.DecodeBody(step.Use.Body, finalizeContext(c.ctx), &execBody); diag.HasErrors() {
					return nil, diag
				}

				s.Kind = &pb.Pipeline_Step_Exec_{
					Exec: &pb.Pipeline_Step_Exec{
						Image:   step.ImageURL,
						Command: execBody.Command,
						Args:    execBody.Args,
					},
				}
			default:
				return nil, status.Errorf(codes.Internal, "unsupported step plugin type: %q", step.Use.Type)
			}

			steps[step.Name] = s
		}

		pipe.Steps = steps

		result = append(result, pipe)
	}

	return result, nil
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
