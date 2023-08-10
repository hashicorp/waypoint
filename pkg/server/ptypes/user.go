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

// UsernameRegexp is the valid username regular expression. This is
// somewhat arbitrary but exactly matches the GitHub username requirements.
// We can always loosen this later.
var UsernameRegexp = regexp.MustCompile(`(?i)^[a-z\d][a-z\d_-]*[a-z\d]+$`)

// TestUser returns a valid user for tests.
func TestUser(t testing.T, src *pb.User) *pb.User {
	t.Helper()

	if src == nil {
		src = &pb.User{}
	}

	require.NoError(t, mergo.Merge(src, &pb.User{
		Username: "test",
	}))

	return src
}

// ValidateUser validates the user structure.
func ValidateUser(v *pb.User) error {
	return validationext.Error(validation.ValidateStruct(v,
		ValidateUserRules(v)...,
	))
}

// ValidateUserRules
func ValidateUserRules(v *pb.User) []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&v.Username,
			validation.Required,
			validation.Match(UsernameRegexp),
			validation.Length(1, 38)),
	}
}

// ValidateUpdateUserRequest
func ValidateUpdateUserRequest(v *pb.UpdateUserRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.User, validation.Required),
		validationext.StructField(&v.User, func() []*validation.FieldRules {
			return append(ValidateUserRules(v.User),
				validation.Field(&v.User.Id, validation.Required),
			)
		}),
	))
}

// ValidateDeleteUserRequest
func ValidateDeleteUserRequest(v *pb.DeleteUserRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.User, validation.Required),
	))
}
