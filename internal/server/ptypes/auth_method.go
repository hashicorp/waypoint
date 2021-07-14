package ptypes

import (
	"errors"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
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
		Name:        "test",
		DisplayName: "test",

		Method: &pb.AuthMethod_Oidc{
			Oidc: &pb.AuthMethod_OIDC{
				ClientId:     "A",
				ClientSecret: "B",
				DiscoveryUrl: "https://example.com/discovery",
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
		validation.Field(&v.Name, validation.Required, validation.By(isNotToken)),
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
		validation.Field(&v.Oidc.DiscoveryUrl, validation.Required, is.URL),
	}
}

// ValidateUpsertAuthMethodRequest
func ValidateUpsertAuthMethodRequest(v *pb.UpsertAuthMethodRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.AuthMethod, validation.Required),
		validationext.StructField(&v.AuthMethod, func() []*validation.FieldRules {
			return ValidateAuthMethodRules(v.AuthMethod)
		}),
	))
}

// ValidateDeleteAuthMethodRequest
func ValidateDeleteAuthMethodRequest(v *pb.DeleteAuthMethodRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.AuthMethod, validation.Required),
	))
}

// ValidateGetOIDCAuthURLRequest
func ValidateGetOIDCAuthURLRequest(v *pb.GetOIDCAuthURLRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.AuthMethod, validation.Required),
		validation.Field(&v.RedirectUri, validation.Required),
	))
}

// ValidateCompleteOIDCAuthRequest
func ValidateCompleteOIDCAuthRequest(v *pb.CompleteOIDCAuthRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.AuthMethod, validation.Required),
		validation.Field(&v.RedirectUri, validation.Required),
		validation.Field(&v.State, validation.Required),
		validation.Field(&v.Code, validation.Required),
		validation.Field(&v.Nonce, validation.Required),
	))
}

func isNotToken(v interface{}) error {
	if v.(string) == "token" {
		return errors.New("name 'token' is reserved and cannot be used")
	}

	return nil
}
