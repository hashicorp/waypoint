package ptypes

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/pkg/validationext"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// TestArtifact returns a valid user for tests.
func TestArtifact(t testing.T, src *pb.PushedArtifact) *pb.PushedArtifact {
	t.Helper()

	if src == nil {
		src = &pb.PushedArtifact{}
	}

	require.NoError(t, mergo.Merge(src, &pb.PushedArtifact{
		Id: "test",
	}))

	return src
}

// ValidatePushedArtifact validates the user structure.
func ValidatePushedArtifact(v *pb.PushedArtifact) error {
	return validationext.Error(validation.ValidateStruct(v,
		ValidatePushedArtifactRules(v)...,
	))
}

// ValidatePushedArtifactRules
func ValidatePushedArtifactRules(v *pb.PushedArtifact) []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&v.Artifact, validation.Required),

		validationext.StructField(&v.Application, func() []*validation.FieldRules {
			return []*validation.FieldRules{
				validation.Field(&v.Application.Application, validation.Required),
				validation.Field(&v.Application.Project, validation.Required),
			}
		}),

		validationext.StructField(&v.Workspace, func() []*validation.FieldRules {
			return []*validation.FieldRules{
				validation.Field(&v.Workspace.Workspace, validation.Required),
			}
		}),
	}
}

// ValidateUpsertArtifactRequest
func ValidateUpsertPushedArtifactRequest(v *pb.UpsertPushedArtifactRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Artifact, validation.Required),
	))
}

// ValidateListPushedArtifactsRequest
func ValidateListPushedArtifactsRequest(v *pb.ListPushedArtifactsRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validationext.StructField(&v.Application, func() []*validation.FieldRules {
			return []*validation.FieldRules{
				validation.Field(&v.Application.Application, validation.Required),
				validation.Field(&v.Application.Project, validation.Required),
			}
		})))
}

// ValidateGetLatestPushedArtifactRequest
func ValidateGetLatestPushedArtifactRequest(v *pb.GetLatestPushedArtifactRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Application, validation.Required),
	))
}

// ValidateGetPushedArtifactRequest
func ValidateGetPushedArtifactRequest(v *pb.GetPushedArtifactRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Ref, validation.Required),
		validationext.StructField(&v.Ref, func() []*validation.FieldRules {
			return ValidateRefOperationRules(v.Ref)
		}),
	))
}
