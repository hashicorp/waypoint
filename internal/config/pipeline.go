package config

import (
	"strings"

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

	// An optional embedded pipeline stanza
	PipelineRaw *hclPipeline `hcl:"pipeline,block"`
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

		// Parse a nested pipeline step if defined
		// TODO(briancain): At the moment, we're supporting singly nestested Pipeline
		// stanzas. If we want to support fully nested embedded pipelines we'll need
		// to turn this into a recursive function.
		if stepRaw.PipelineRaw != nil {
			var embedPipeline Pipeline
			if diag := gohcl.DecodeBody(stepRaw.PipelineRaw.Body, finalizeContext(ctx), &embedPipeline); diag.HasErrors() {
				return nil, diag
			}

			s.Pipeline = &Pipeline{
				Name:   stepRaw.PipelineRaw.Name,
				ctx:    ctx,
				Config: c,
			}
			if s.Pipeline.Config != nil {
				s.Pipeline.Config.ctx = ctx
			}

			// Parse all the steps
			var embSteps []*Step
			for _, embedStepRaw := range embedPipeline.StepRaw {
				// turn stepRaw into a staged Step
				s := Step{
					ctx:       ctx,
					Name:      embedStepRaw.Name,
					DependsOn: embedStepRaw.DependsOn,
					ImageURL:  embedStepRaw.ImageURL,
					Use:       embedStepRaw.Use,
				}
				embSteps = append(embSteps, &s)
			}

			s.Pipeline.Steps = embSteps
		}

		steps = append(steps, &s)
	}

	pipeline.Steps = steps

	return &pipeline, nil
}

// PipelineProtos will take the existing HCL eval context, eval the config
// and translate the HCL result into a Pipeline Proto to be returned for
// operations such as ConfigSync.
// If a pipeline has an embedded pipeline defined, PipelineProtos will return
// each as its own separate Pipeline proto message where the step that defined
// the embedded pipeline is actually a Pipeline Step reference.
func (c *Config) PipelineProtos() ([]*pb.Pipeline, error) {
	if c == nil {
		// This is likely an internal error if this happens.
		panic("attempted to construct pipeline proto on a nil genericConfig")
	}

	// Load HCL config and convert to a Pipeline proto
	var result []*pb.Pipeline
	for _, pl := range c.hclConfig.Pipelines {
		pipes, err := c.buildPipelineProto(pl)
		if err != nil {
			return nil, err
		}

		result = append(result, pipes...)
	}

	return result, nil
}

// buildPipelineProto will recursively translate an hclPipeline into a protobuf
// Pipeline message.
func (c *Config) buildPipelineProto(pl *hclPipeline) ([]*pb.Pipeline, error) {
	var result []*pb.Pipeline
	pipe := &pb.Pipeline{
		Name: pl.Name,
		Owner: &pb.Pipeline_Project{
			Project: &pb.Ref_Project{
				Project: c.hclConfig.Project,
			},
		},
	}

	steps := make(map[string]*pb.Pipeline_Step)
	for i, step := range pl.StepRaw {
		s := &pb.Pipeline_Step{
			Name:      step.Name,
			DependsOn: step.DependsOn,
			Image:     step.ImageURL,
		}

		// If no dependency was explictily set, we rely on the previous step
		if i != 0 && len(step.DependsOn) == 0 {
			s.DependsOn = []string{pl.StepRaw[i-1].Name}
		}

		// We have an embeded pipeline for this step. This can either be an hclPipeline
		// defined directly in the step, or a pipeline reference to another pipeline
		// defined else where. If this is a ref, the raw hcl for the pipeline should
		// be empty.
		if step.PipelineRaw != nil {
			// TODO(briancain): It shouldn't be valid to have a "project/pipeline"
			// defined as well as a nested pipeline. We should validate this.

			// We need to determine if the embedded pipeline is defined directly
			// inside a step, or is simply a reference to a pipeline else where
			if len(step.PipelineRaw.StepRaw) > 0 {
				// This means this is an embedded pipeline, i.e. the HCL definition
				// is nested within the step PipelineRaw. we parse that pipeline
				// directly and store it as a separate pipeline, and make _this_ step
				// a reference to the pipeline

				// Parse nested pipeline steps
				pipelines, err := c.buildPipelineProto(step.PipelineRaw)
				if err != nil {
					return nil, err
				}

				result = append(result, pipelines...)
			}

			// We check if this step references a separate pipeline by Owner
			pipeName := step.PipelineRaw.Name
			pipeProject := c.hclConfig.Project
			if pn := strings.Split(pipeName, "/"); len(pn) > 1 {
				pipeProject = pn[0]
				pipeName = pn[1]
			}

			// Add pipeline reference as a pipeline ref step for parent pipeline
			s.Kind = &pb.Pipeline_Step_Pipeline_{
				Pipeline: &pb.Pipeline_Step_Pipeline{
					Ref: &pb.Ref_Pipeline{
						Ref: &pb.Ref_Pipeline_Owner{
							Owner: &pb.Ref_PipelineOwner{
								Project: &pb.Ref_Project{
									Project: pipeProject,
								},
								PipelineName: pipeName,
							},
						},
					},
				},
			}

			steps[step.Name] = s

			continue // continue to build the rest of the parent pipeline
		} // else do below

		// NOTE(briancain): This is what you'd change to support future Step plugins
		// or future built-in step operations.
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
		case "deploy":
			var deployBody struct {
				Release bool `hcl:"release,optional"`
			}

			if diag := gohcl.DecodeBody(step.Use.Body, finalizeContext(c.ctx), &deployBody); diag.HasErrors() {
				return nil, diag
			}

			s.Kind = &pb.Pipeline_Step_Deploy_{
				Deploy: &pb.Pipeline_Step_Deploy{
					Release: deployBody.Release,
				},
			}
		case "release":
			var releaseBody struct {
				DeploymentRef       uint64 `hcl:"deployment_ref,optional"` // 0 or "unset" means latest
				Prune               bool   `hcl:"prune,optional"`
				PruneRetain         int32  `hcl:"prune_retain,optional"`
				PruneRetainOverride bool   `hcl:"prune_retain_override,optional"`
			}

			if diag := gohcl.DecodeBody(step.Use.Body, finalizeContext(c.ctx), &releaseBody); diag.HasErrors() {
				return nil, diag
			}

			// For parsing pipeline configs, an unset `deployment_ref` translates
			// to the "latest" deployment. Otherwise if people want to release
			// a specific deployment by sequence number, they can set it explicity.
			deployRef := &pb.Ref_Deployment{
				Ref: &pb.Ref_Deployment_Latest{
					Latest: true,
				},
			}
			if releaseBody.DeploymentRef != 0 {
				deployRef = &pb.Ref_Deployment{
					Ref: &pb.Ref_Deployment_Sequence{
						Sequence: releaseBody.DeploymentRef,
					},
				}
			}

			s.Kind = &pb.Pipeline_Step_Release_{
				Release: &pb.Pipeline_Step_Release{
					Deployment:          deployRef,
					Prune:               releaseBody.Prune,
					PruneRetain:         releaseBody.PruneRetain,
					PruneRetainOverride: releaseBody.PruneRetainOverride,
				},
			}
		case "up":
			var upBody struct {
				Prune               bool  `hcl:"prune,optional"`
				PruneRetain         int32 `hcl:"prune_retain,optional"`
				PruneRetainOverride bool  `hcl:"prune_retain_override,optional"`
			}

			if diag := gohcl.DecodeBody(step.Use.Body, finalizeContext(c.ctx), &upBody); diag.HasErrors() {
				return nil, diag
			}

			s.Kind = &pb.Pipeline_Step_Up_{
				Up: &pb.Pipeline_Step_Up{
					Prune:               upBody.Prune,
					PruneRetain:         upBody.PruneRetain,
					PruneRetainOverride: upBody.PruneRetainOverride,
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
		case "":
			return nil, status.Error(codes.FailedPrecondition, "step use label cannot be empty")
		default:
			return nil, status.Errorf(codes.Internal, "unsupported step plugin type: %q", step.Use.Type)
		}

		steps[step.Name] = s
	}

	pipe.Steps = steps

	result = append(result, pipe)

	// Validate there are no cycles?

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
