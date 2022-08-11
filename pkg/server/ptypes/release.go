package ptypes

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/imdario/mergo"
	"github.com/stretchr/testify/require"

	"github.com/mitchellh/go-testing-interface"

	"github.com/hashicorp/waypoint/internal/pkg/validationext"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func TestRelease(t testing.T, src *pb.Release) *pb.Release {
	t.Helper()

	if src == nil {
		src = &pb.Release{}
	}

	require.NoError(t, mergo.Merge(src, &pb.Release{
		Id:        "test",
		Workspace: &pb.Ref_Workspace{Workspace: "default"},
	}))

	return src
}

// ValidateGetReleaseRequest
func ValidateGetReleaseRequest(v *pb.GetReleaseRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Ref, validation.Required),
		validationext.StructField(&v.Ref, func() []*validation.FieldRules {
			return ValidateRefOperationRules(v.Ref)
		}),
	))
}

// ValidateGetLatestReleaseRequest
func ValidateGetLatestReleaseRequest(v *pb.GetLatestReleaseRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Application, validation.Required),
	))
}

// ValidateUpsertArtifactRequest
func ValidateUpsertReleaseRequest(v *pb.UpsertReleaseRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Release, validation.Required),
	))
}
