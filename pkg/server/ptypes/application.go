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

// TestApplication returns a valid project for tests.
func TestApplication(t testing.T, src *pb.Application) *pb.Application {
	t.Helper()

	if src == nil {
		src = &pb.Application{}
	}

	require.NoError(t, mergo.Merge(src, &pb.Application{
		Project: &pb.Ref_Project{
			Project: "test",
		},

		Name: "test",
	}))

	return src
}

// ValidateUpsertApplicationRequest
func ValidateUpsertApplicationRequest(v *pb.UpsertApplicationRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Project, validation.Required),
		validation.Field(&v.Name, validation.Required, validation.By(validatePathToken)),
	))
}

// ValidateGetApplicationRequest
func ValidateGetApplicationRequest(v *pb.GetApplicationRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Application, validation.Required),
	))
}
