// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ptypes

import (
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/pkg/validationext"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// TestDeployment returns a valid project for tests.
func TestDeployment(t testing.T, src *pb.Deployment) *pb.Deployment {
	t.Helper()

	if src == nil {
		src = &pb.Deployment{}
	}

	require.NoError(t, mergo.Merge(src, &pb.Deployment{
		Id: "test",
	}))

	return src
}

// Type wrapper around the proto type so that we can add some methods.
type Deployment struct{ *pb.Deployment }

func (v *Deployment) URLFragment() string {
	// For older deployments (pre WP 0.4.0) we use the sequence. If
	// we have a generation set, we use the generation initial sequence.
	seq := v.Sequence

	// By using the generation sequence, we ensure that all deployments
	// in a generation share the same URL.
	if g := v.Generation; g != nil {
		seq = g.InitialSequence
	}

	return "v" + strconv.FormatUint(seq, 10)
}

// ValidateDeployment validates the project structure.
func ValidateDeployment(v *pb.Deployment) error {
	return validationext.Error(validation.ValidateStruct(v,
		ValidateDeploymentRules(v)...,
	))
}

// ValidateDeploymentRules
func ValidateDeploymentRules(v *pb.Deployment) []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&v.Application, validation.Required),
		validation.Field(&v.Workspace, validation.Required),
	}
}

// ValidateUpsertDeploymentRequest
func ValidateUpsertDeploymentRequest(v *pb.UpsertDeploymentRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Deployment, validation.Required),
		validationext.StructField(&v.Deployment, func() []*validation.FieldRules {
			return ValidateDeploymentRules(v.Deployment)
		}),
	))
}

// ValidateGetDeploymentRequest
func ValidateGetDeploymentRequest(v *pb.GetDeploymentRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Ref, validation.Required),
		validationext.StructField(&v.Ref, func() []*validation.FieldRules {
			return ValidateRefOperationRules(v.Ref)
		}),
	))
}
