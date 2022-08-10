package ptypes

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/pkg/validationext"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// TestBuild returns a valid user for tests.
func TestBuild(t testing.T, src *pb.Build) *pb.Build {
	t.Helper()

	if src == nil {
		src = &pb.Build{}
	}

	require.NoError(t, mergo.Merge(src, &pb.Build{
		Id: "test",
	}))

	return src
}

// ValidateBuild validates the user structure.
func ValidateBuild(v *pb.Build) error {
	return validationext.Error(validation.ValidateStruct(v,
		ValidateBuildRules(v)...,
	))
}

// ValidateBuildRules
func ValidateBuildRules(v *pb.Build) []*validation.FieldRules {
	return []*validation.FieldRules{
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

// ValidateGetBuildRequest
func ValidateGetBuildRequest(v *pb.GetBuildRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Ref, validation.Required),
		validationext.StructField(&v.Ref, func() []*validation.FieldRules {
			return ValidateRefOperationRules(v.Ref)
		}),
	))
}

// ValidateDeleteBuildRequest
func ValidateDeleteBuildRequest(v *pb.DeleteBuildRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Ref, validation.Required),
		validationext.StructField(&v.Ref, func() []*validation.FieldRules {
			return ValidateRefOperationRules(v.Ref)
		}),
	))
}

// ValidateListBuildsRequest
func ValidateListBuildsRequest(v *pb.ListBuildsRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validationext.StructField(&v.Application, func() []*validation.FieldRules {
			return []*validation.FieldRules{
				validation.Field(&v.Application.Application, validation.Required),
				validation.Field(&v.Application.Project, validation.Required),
			}
		})))
}

// ValidateGetLatestBuildRequest
func ValidateGetLatestBuildRequest(v *pb.GetLatestBuildRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Application, validation.Required),
	))
}

// ValidateUpsertBuildRequest
func ValidateUpsertBuildRequest(v *pb.UpsertBuildRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Build, validation.Required),
		validationext.StructField(&v.Build, func() []*validation.FieldRules {
			return ValidateBuildRules(v.Build)
		}),
	))
}
