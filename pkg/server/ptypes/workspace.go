// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ptypes

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/pkg/validationext"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// WorkspaceNameRegexp is the valid Workspace name regular expression. At this
// time the only restriction is to not allow spaces.
var WorkspaceNameRegexp = regexp.MustCompile(`^[\p{L}\p{N}]+[\p{L}\p{N}\-_]*[^\-_]?$`)

// TestWorkspace returns a valid workspace for tests.
func TestWorkspace(t testing.T, src *pb.Workspace) *pb.Workspace {
	t.Helper()

	if src == nil {
		src = &pb.Workspace{}
	}

	require.NoError(t, mergo.Merge(src, &pb.Workspace{
		Name: "test",
	}))

	return src
}

func TestGetWorkspaceRequest(t testing.T, src *pb.GetWorkspaceRequest) *pb.GetWorkspaceRequest {
	t.Helper()

	if src == nil {
		src = &pb.GetWorkspaceRequest{}
	}

	require.NoError(t, mergo.Merge(src, &pb.GetWorkspaceRequest{
		Workspace: &pb.Ref_Workspace{
			Workspace: "w_test",
		},
	}))

	return src
}

// ValidateGetWorkspaceRequest
func ValidateGetWorkspaceRequest(v *pb.GetWorkspaceRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Workspace, validation.Required),
		validationext.StructField(&v.Workspace, func() []*validation.FieldRules {
			return ValidateRefWorkspaceRules(v.Workspace)
		}),
	))
}

// ValidateUpdateUserRequest
func ValidateUpsertWorkspaceRequest(v *pb.UpsertWorkspaceRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Workspace, validation.Required),
		validationext.StructField(&v.Workspace, func() []*validation.FieldRules {
			return append(ValidateWorkspaceRules(v.Workspace),
				validation.Field(&v.Workspace.Name, validation.Required, validation.By(validatePathToken)),
			)
		}),
	))
}

// ValidateWorkspace validates the Workspace structure.
func ValidateWorkspace(v *pb.Workspace) error {
	return validationext.Error(validation.ValidateStruct(v,
		ValidateWorkspaceRules(v)...,
	))
}

// ValidateWorkspaceRules
func ValidateWorkspaceRules(v *pb.Workspace) []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&v.Name,
			validation.Required,
			validation.Match(WorkspaceNameRegexp),
			validation.Length(1, 38)),
	}
}

func ValidateWorkspaceName(str string) error {
	return validationext.Error(validation.Validate(str,
		validation.Match(WorkspaceNameRegexp),
		validation.Length(1, 38)),
	)
}
