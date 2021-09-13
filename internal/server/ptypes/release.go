package ptypes

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hashicorp/waypoint/internal/pkg/validationext"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

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
