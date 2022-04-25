package ptypes

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/pkg/validationext"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

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
