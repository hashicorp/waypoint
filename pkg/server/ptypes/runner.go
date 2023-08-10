// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ptypes

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/mitchellh/hashstructure/v2"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/pkg/validationext"
	serverpkg "github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// RunnerLabelHash calculates a unique hash for the set of labels on the
// runner. This generates a consistent hash value for an empty set of labels.
// The result is never 0.
func RunnerLabelHash(v map[string]string) (uint64, error) {
	if v == nil {
		v = map[string]string{}
	}

	// We always set this special key so that even empty label sets have
	// a non-zero hash value. This MUST NEVER BE CHANGED otherwise the
	// hash values for all previously issued tokens will be invalidated.
	v["waypoint.hashicorp.com/runner-hash"] = "1"

	return hashstructure.Hash(v, hashstructure.FormatV2, nil)
}

func TestRunner(t testing.T, src *pb.Runner) *pb.Runner {
	t.Helper()

	if src == nil {
		src = &pb.Runner{}
	}

	id, err := serverpkg.Id()
	require.NoError(t, err)

	require.NoError(t, mergo.Merge(src, &pb.Runner{
		Id: id,
	}))

	return src
}

// ValidateAdoptRunnerRequest
func ValidateAdoptRunnerRequest(v *pb.AdoptRunnerRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.RunnerId, validation.Required, validation.By(validatePathToken)),
	))
}

// ValidateForgetRunnerRequest
func ValidateForgetRunnerRequest(v *pb.ForgetRunnerRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.RunnerId, validation.Required),
	))
}
