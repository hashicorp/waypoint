package ptypes

import (
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/pkg/validationext"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// TestAuthMethod returns a valid user for tests.
func TestAuthMethod(t testing.T, src *pb.AuthMethod) *pb.AuthMethod {
	t.Helper()

	if src == nil {
		src = &pb.AuthMethod{}
	}

	require.NoError(t, mergo.Merge(src, &pb.AuthMethod{
		Name: "test",

		Method: &pb.AuthMethod_Oidc{
			Oidc: &pb.AuthMethod_OIDC{
				ClientId:     "A",
				ClientSecret: "B",
			},
		},
	}))

	return src
}

// ValidateAuthMethod validates the user structure.
func ValidateAuthMethod(v *pb.AuthMethod) error {
	return validationext.Error(validation.ValidateStruct(v,
		ValidateAuthMethodRules(v)...,
	))
}

// ValidateAuthMethodRules
func ValidateAuthMethodRules(v *pb.AuthMethod) []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&v.Name, validation.Required),
		validation.Field(&v.Method, validation.Required),

		validationext.StructOneof(&v.Method, (*pb.AuthMethod_Oidc)(nil),
			func() []*validation.FieldRules {
				v := v.Method.(*pb.AuthMethod_Oidc)
				return validateAuthMethodOIDCRules(v)
			}),
	}
}

// validateAuthMethodOIDCRules
func validateAuthMethodOIDCRules(v *pb.AuthMethod_Oidc) []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&v.Oidc.ClientId, validation.Required),
		validation.Field(&v.Oidc.ClientSecret, validation.Required),
	}
}
