// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ptypes

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/pkg/validationext"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// TestStatusReport returns a valid user for tests.
func TestStatusReport(t testing.T, src *pb.StatusReport) *pb.StatusReport {
	t.Helper()

	if src == nil {
		src = &pb.StatusReport{}
	}

	require.NoError(t, mergo.Merge(src, &pb.StatusReport{
		Id: "test",
	}))

	return src
}

// ValidateStatusReport validates the user structure.
func ValidateStatusReport(v *pb.StatusReport) error {
	return validationext.Error(validation.ValidateStruct(v,
		ValidateStatusReportRules(v)...,
	))
}

// ValidateStatusReportRules
func ValidateStatusReportRules(v *pb.StatusReport) []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&v.TargetId, validation.Required),
		validation.Field(&v.Health, validation.Required),

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

// ValidateUpsertStatusReportRequest
func ValidateUpsertStatusReportRequest(v *pb.UpsertStatusReportRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.StatusReport, validation.Required),
	))
}

// ValidateListStatusReportsRequest
func ValidateListStatusReportsRequest(v *pb.ListStatusReportsRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validationext.StructField(&v.Application, func() []*validation.FieldRules {
			return []*validation.FieldRules{
				validation.Field(&v.Application.Application, validation.Required),
				validation.Field(&v.Application.Project, validation.Required),
			}
		})))
}

// ValidateGetLatestStatusReportRequest
func ValidateGetLatestStatusReportRequest(v *pb.GetLatestStatusReportRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Application, validation.Required),
	))
}

// ValidateGetStatusReportRequest
func ValidateGetStatusReportRequest(v *pb.GetStatusReportRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Ref, validation.Required),
		validationext.StructField(&v.Ref, func() []*validation.FieldRules {
			return ValidateRefOperationRules(v.Ref)
		}),
	))
}
