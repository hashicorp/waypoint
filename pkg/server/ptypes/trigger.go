package ptypes

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/pkg/validationext"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// TestTrigger returns a valid user for tests.
func TestTrigger(t testing.T, src *pb.Trigger) *pb.Trigger {
	t.Helper()

	if src == nil {
		src = &pb.Trigger{}
	}

	require.NoError(t, mergo.Merge(src, &pb.Trigger{
		Id: "test",
	}))

	return src
}

// ValidateTrigger validates the user structure.
func ValidateTrigger(v *pb.Trigger) error {
	return validationext.Error(validation.ValidateStruct(v,
		ValidateTriggerRules(v)...,
	))
}

// ValidateTriggerRules
func ValidateTriggerRules(v *pb.Trigger) []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&v.Id, validation.Required),
		validation.Field(&v.Name, validation.Required),

		validationext.StructField(&v.Project, func() []*validation.FieldRules {
			return []*validation.FieldRules{
				validation.Field(&v.Project.Project, validation.Required),
			}
		}),

		validationext.StructField(&v.Workspace, func() []*validation.FieldRules {
			return []*validation.FieldRules{
				validation.Field(&v.Workspace.Workspace, validation.Required),
			}
		}),
	}
}

// ValidateUpsertTriggerRequest
func ValidateUpsertTriggerRequest(v *pb.UpsertTriggerRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Trigger, validation.Required),
		validationext.StructField(&v.Trigger, func() []*validation.FieldRules {
			return []*validation.FieldRules{
				validation.Field(&v.Trigger.Project, validation.Required),
				// Trigger Name is also the "path" in the HTTP request, so we will
				// validate the name against our valid path token check
				validation.Field(&v.Trigger.Name, validation.Required, validation.By(validatePathToken)),
			}
		}),
	))
}

// ValidateGetTriggerRequest
func ValidateGetTriggerRequest(v *pb.GetTriggerRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Ref, validation.Required),
		validationext.StructField(&v.Ref, func() []*validation.FieldRules {
			return ValidateRefTriggerRules(v.Ref)
		}),
	))
}

// ValidateDeleteTriggerRequest
func ValidateDeleteTriggerRequest(v *pb.DeleteTriggerRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Ref, validation.Required),
		validationext.StructField(&v.Ref, func() []*validation.FieldRules {
			return ValidateRefTriggerRules(v.Ref)
		}),
	))
}

// ValidateRefTriggerRules
func ValidateRefTriggerRules(v *pb.Ref_Trigger) []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&v.Id, validation.Required),
	}
}

// ValidateRunTriggerRequest
func ValidateRunTriggerRequest(v *pb.RunTriggerRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Ref, validation.Required),
		validationext.StructField(&v.Ref, func() []*validation.FieldRules {
			return ValidateRefTriggerRules(v.Ref)
		}),
	))
}
