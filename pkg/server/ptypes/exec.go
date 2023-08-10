// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ptypes

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/pkg/validationext"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func TestExecStreamRequestStart(t testing.T, src *pb.ExecStreamRequest_Start) *pb.ExecStreamRequest_Start {
	t.Helper()

	if src == nil {
		src = &pb.ExecStreamRequest_Start{}
	}

	require.NoError(t, mergo.Merge(src, &pb.ExecStreamRequest_Start{
		Target: &pb.ExecStreamRequest_Start_DeploymentId{
			DeploymentId: "1",
		},

		Args: []string{"/bin/bash"},
	}))

	return src
}

// ValidateExecStreamRequestStart
func ValidateExecStreamRequestStart(v *pb.ExecStreamRequest_Start) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Target, validation.Required),
		validation.Field(&v.Args, validation.Required),
	))
}
