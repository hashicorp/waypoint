package ptypes

import (
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/pkg/validationext"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

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
		validation.Field(&v.Username, validation.Required),
	}
}

/*
// ValidateUpsertUserRequest
func ValidateUpsertUserRequest(v *pb.UpsertUserRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.User, validation.Required),
		validationext.StructField(&v.User, func() []*validation.FieldRules {
			return ValidateUserRules(v.User)
		}),
	))
}
*/
