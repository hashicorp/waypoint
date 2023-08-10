// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package config

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

// Step are the step settings for pipelines
type Step struct {
	Labels map[string]string `hcl:"labels,optional"`
	Use    *Use              `hcl:"use,block"`

	// Give this step a name
	Name string `hcl:",label"`

	// If set, this step will depend on the defined step. The default step
	// will be the previously defined step in order that it was defined
	// in a waypoint.hcl
	DependsOn []string `hcl:"depends_on,optional"`

	// The OCI image to use for executing this step
	ImageURL string `hcl:"image_url,optional"`

	// An optional embedded pipeline stanza
	Pipeline *Pipeline `hcl:"pipeline,block"`

	ctx *hcl.EvalContext

	// Optional workspace scoping
	Workspace string `hcl:"workspace,optional"`
}

// hclStep represents a raw HCL version of a step stanza in a pipeline config
type hclStep struct {
	Name string `hcl:",label"`

	// An optional embedded pipeline stanza
	PipelineRaw *hclPipeline `hcl:"pipeline,block"`

	Body   hcl.Body `hcl:",body"`
	Remain hcl.Body `hcl:",remain"`
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
// Note that currently this parsing function does not attempt to detect cycles
// between embedded pipelines.
func (c *Config) Pipeline(id string, ctx *hcl.EvalContext) (*Pipeline, error) {
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
		var step Step
		if diag := gohcl.DecodeBody(stepRaw.Body, finalizeContext(ctx), &step); diag.HasErrors() {
			return nil, diag
		}
		step.ctx = ctx
		step.Name = stepRaw.Name

		// Parse a nested pipeline step if defined
		// TODO(briancain): At the moment, we're supporting singly nestested Pipeline
		// stanzas. If we want to support fully nested embedded pipelines we'll need
		// to turn this into a recursive function.
		if stepRaw.PipelineRaw != nil {
			var embedPipeline Pipeline
			if diag := gohcl.DecodeBody(stepRaw.PipelineRaw.Body, finalizeContext(ctx), &embedPipeline); diag.HasErrors() {
				return nil, diag
			}

			step.Pipeline = &Pipeline{
				Name:   stepRaw.PipelineRaw.Name,
				ctx:    ctx,
				Config: c,
			}
			if step.Pipeline.Config != nil {
				step.Pipeline.Config.ctx = ctx
			}

			// Parse all the steps
			var embSteps []*Step
			for _, embedStepRaw := range embedPipeline.StepRaw {
				if embedStepRaw.PipelineRaw != nil {
					// NOTE(briancain): For now, we artificially don't allow for more than
					// one level of embedded pipeline nesting. That means this is invalid:
					/*
						  pipeline "example" {
								step "invalid" {
									pipeline "level-one" {
										step "nested" {
											pipeline "level-two" {
												// etc etc...
											}
										}
									}
								}
						  }
					*/
					return nil, status.Errorf(codes.FailedPrecondition,
						"Step %q defined 2 levels of nesting for an embedded pipeline %q. "+
							"Currently Waypoint only supports 1 level of nesting for embedded pipelines. "+
							"You can instead define a pipeline and refer to it as a step rather than "+
							"defining it directly inside a step.", embedStepRaw.Name, pipeline.Name)
				}

				// turn stepRaw into a staged Step
				var embedStep Step
				if diag := gohcl.DecodeBody(embedStepRaw.Body, finalizeContext(ctx), &embedStep); diag.HasErrors() {
					return nil, diag
				}
				embedStep.ctx = ctx
				embedStep.Name = embedStepRaw.Name
				embSteps = append(embSteps, &embedStep)
			}

			step.Pipeline.Steps = embSteps
		}

		steps = append(steps, &step)
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
		pipeline, err := c.Pipeline(pl.Name, c.ctx)
		if err != nil {
			return nil, err
		}

		pipes, err := c.buildPipelineProto(pipeline)
		if err != nil {
			return nil, err
		}

		result = append(result, pipes...)
	}

	// We should validate cycles across pipelines here

	return result, nil
}

// buildPipelineProto will recursively translate an hclPipeline into a protobuf
// Pipeline message.
func (c *Config) buildPipelineProto(pl *Pipeline) ([]*pb.Pipeline, error) {
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
	for i, step := range pl.Steps {
		s := &pb.Pipeline_Step{
			Name:      step.Name,
			DependsOn: step.DependsOn,
			Image:     step.ImageURL,
		}

		if step.Workspace != "" {
			s.Workspace = &pb.Ref_Workspace{
				Workspace: step.Workspace,
			}
		}

		// If no dependency was explictily set, we rely on the previous step
		if i != 0 && len(step.DependsOn) == 0 {
			s.DependsOn = []string{pl.Steps[i-1].Name}
		}

		// We have an embeded pipeline for this step. This can either be an hclPipeline
		// defined directly in the step, or a pipeline reference to another pipeline
		// defined else where. If this is a ref, the raw hcl for the pipeline should
		// be a "built-in" step of type "pipeline"
		if step.Pipeline != nil {
			// Parse the embedded pipeline assuming it has steps
			if len(step.Pipeline.Steps) > 0 {
				// This means this is an embedded pipeline, i.e. the HCL definition
				// is nested within the step PipelineRaw. we parse that pipeline
				// directly and store it as a separate pipeline, and make _this_ step
				// a reference to the pipeline

				// Parse nested pipeline steps
				pipelines, err := c.buildPipelineProto(step.Pipeline)
				if err != nil {
					return nil, err
				}

				result = append(result, pipelines...)

				// We check if this step references a separate pipeline by Owner
				pipeName := step.Pipeline.Name
				pipeProject := c.hclConfig.Project

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
			}

			continue // continue to build the rest of the parent pipeline
		} // else handle any "built-in" steps

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
				Prune               *bool  `hcl:"prune,optional"`          // nil means unset, we default to "true"
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

			// unset, so default to true
			if releaseBody.Prune == nil {
				b := true
				releaseBody.Prune = &b
			}

			s.Kind = &pb.Pipeline_Step_Release_{
				Release: &pb.Pipeline_Step_Release{
					Deployment:          deployRef,
					Prune:               *releaseBody.Prune,
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
				Command string   `hcl:"command"`
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
		case "pipeline":
			var pipelineBody struct {
				Project string `hcl:"project"`
				Name    string `hcl:"name"`
			}

			// Evaluate the step body hcl to get options
			if diag := gohcl.DecodeBody(step.Use.Body, finalizeContext(c.ctx), &pipelineBody); diag.HasErrors() {
				return nil, diag
			}

			// set to *Current* project if not set
			if pipelineBody.Project == "" {
				pipelineBody.Project = c.hclConfig.Project
			}

			// Add pipeline reference as a pipeline ref step for parent pipeline
			s.Kind = &pb.Pipeline_Step_Pipeline_{
				Pipeline: &pb.Pipeline_Step_Pipeline{
					Ref: &pb.Ref_Pipeline{
						Ref: &pb.Ref_Pipeline_Owner{
							Owner: &pb.Ref_PipelineOwner{
								Project: &pb.Ref_Project{
									Project: pipelineBody.Project,
								},
								PipelineName: pipelineBody.Name,
							},
						},
					},
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

	return result, nil
}

// Ref returns the ref for this pipeline.
func (c *Pipeline) Ref() *pb.Ref_Pipeline {
	return &pb.Ref_Pipeline{
		Ref: &pb.Ref_Pipeline_Id{
			Id: c.Name,
		},
	}
}

// Configure configures the plugin for a given Step with the use body of this operation.
func (s *Step) Configure(plugin interface{}, ctx *hcl.EvalContext) hcl.Diagnostics {
	ctx = appendContext(s.ctx, ctx)
	ctx = finalizeContext(ctx)

	return component.Configure(plugin, s.Use.Body, ctx)
}
