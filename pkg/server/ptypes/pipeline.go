package ptypes

import (
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/pkg/graph"
	"github.com/hashicorp/waypoint/internal/pkg/validationext"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// PipelineGraph returns the graph of steps for a pipeline. The graph
// vertices are the pipeline step names.
func PipelineGraph(v *pb.Pipeline) (*graph.Graph, error) {
	return pipelineGraph(v.Steps)
}

// TestPipeline returns a valid user for tests.
func TestPipeline(t testing.T, src *pb.Pipeline) *pb.Pipeline {
	t.Helper()

	if src == nil {
		src = &pb.Pipeline{}
	}

	require.NoError(t, mergo.Merge(src, &pb.Pipeline{
		Id:   "test",
		Name: "test",
		Owner: &pb.Pipeline_Project{
			Project: &pb.Ref_Project{
				Project: "project",
			},
		},
		Steps: map[string]*pb.Pipeline_Step{
			"root": {
				Name: "root",
				Kind: &pb.Pipeline_Step_Exec_{
					Exec: &pb.Pipeline_Step_Exec{
						Image: "hashicorp/waypoint",
					},
				},
			},
		},
	}))

	return src
}

// ValidatePipeline validates the user structure.
func ValidatePipeline(v *pb.Pipeline) error {
	return validationext.Error(validation.ValidateStruct(v,
		ValidatePipelineRules(v)...,
	))
}

// ValidatePipelineRules
func ValidatePipelineRules(v *pb.Pipeline) []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&v.Name, validation.Required),
		validation.Field(&v.Owner, validation.Required),

		validationext.StructOneof(&v.Owner, (*pb.Pipeline_Project)(nil),
			func() []*validation.FieldRules {
				v := v.Owner.(*pb.Pipeline_Project)
				return validatePipelineOwnerProjectRules(v)
			}),

		validation.Field(&v.Steps,
			validation.Required,
			validation.By(stepNameMatchesKey),
			validation.By(stepSingleRoot),
			validation.By(stepGraph),
			validationext.Each(validation.By(func(v interface{}) error {
				s := v.(*pb.Pipeline_Step)
				return validation.ValidateStruct(s, ValidateStepRules(s)...)
			})),
		),
	}
}

// ValidateStepRules
func ValidateStepRules(v *pb.Pipeline_Step) []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&v.Name, validation.Required),
		validation.Field(&v.Kind, validation.Required),

		validationext.StructOneof(&v.Kind, (*pb.Pipeline_Step_Exec)(nil),
			func() []*validation.FieldRules {
				v := v.Kind.(*pb.Pipeline_Step_Exec_)
				return validatePipelineStepExecRules(v)
			}),
	}
}

// validatePipelineOwnerProjectRules
func validatePipelineOwnerProjectRules(v *pb.Pipeline_Project) []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&v.Project.Project, validation.Required),
	}
}

// validatePipelineStepExecRules
func validatePipelineStepExecRules(v *pb.Pipeline_Step_Exec_) []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&v.Exec.Image, validation.Required),
	}
}

// ValidateUpsertPipelineRequest
func ValidateUpsertPipelineRequest(v *pb.UpsertPipelineRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Pipeline, validation.Required),
		validationext.StructField(&v.Pipeline, func() []*validation.FieldRules {
			return ValidatePipelineRules(v.Pipeline)
		}),
	))
}

// ValidateRunPipelineRequest
func ValidateRunPipelineRequest(v *pb.RunPipelineRequest) error {
	// Set the operation so that validation succeeds. We override it later.
	if v.JobTemplate != nil {
		v.JobTemplate.Operation = &pb.Job_Noop_{
			Noop: &pb.Job_Noop{},
		}
	}

	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Pipeline, validation.Required),
		validation.Field(&v.JobTemplate, validation.Required),
		validationext.StructField(&v.JobTemplate, func() []*validation.FieldRules {
			return ValidateJobRules(v.JobTemplate)
		}),
	))
}

// stepNameMatchesKey implements validation.RuleFunc to validate that
// the map key (string) matches the step name field.
func stepNameMatchesKey(v interface{}) error {
	for k, step := range v.(map[string]*pb.Pipeline_Step) {
		if step.Name == "" {
			// Validated elsewhere
			continue
		}

		if k != step.Name {
			return fmt.Errorf("step key %q doesn't match step name %q", k, step.Name)
		}
	}

	return nil
}

// stepSingleRoot implements validation.RuleFunc to validate that
// there is a single root step.
func stepSingleRoot(v interface{}) error {
	count := 0
	for _, step := range v.(map[string]*pb.Pipeline_Step) {
		if len(step.DependsOn) == 0 {
			count++
			if count > 1 {
				return errors.New("a pipeline requires exactly one root step")
			}
		}
	}

	return nil
}

// stepGraph implements validation.RuleFunc to validate that
// builds and validates the step graph.
func stepGraph(v interface{}) error {
	steps := v.(map[string]*pb.Pipeline_Step)
	_, err := pipelineGraph(steps)
	return err
}

func pipelineGraph(steps map[string]*pb.Pipeline_Step) (*graph.Graph, error) {
	var stepGraph graph.Graph
	for _, step := range steps {
		// Add our job
		stepGraph.Add(step.Name)

		// Add any dependencies
		for _, dep := range step.DependsOn {
			stepGraph.Add(dep)
			stepGraph.AddEdge(dep, step.Name)

			if _, ok := steps[dep]; !ok {
				return nil, fmt.Errorf(
					"step %q depends on non-existent step %q", step, dep)
			}
		}
	}
	if cycles := stepGraph.Cycles(); len(cycles) > 0 {
		return nil, fmt.Errorf(
			"step dependencies contain one or more cycles: %s", cycles)
	}

	return &stepGraph, nil
}
