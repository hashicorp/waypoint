// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ptypes

import (
	"errors"
	"fmt"
	"go/token"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/hashicorp/go-bexpr"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/pkg/validationext"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
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
		validation.Field(&v.Name, validation.Required, validation.By(isNotToken), validation.By(validatePathToken)),
		validation.Field(&v.AccessSelector, validation.By(isBExpr)),

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
		validation.Field(&v.Oidc.ClaimMappings, validation.By(isClaimMapping)),
		validation.Field(&v.Oidc.ListClaimMappings, validation.By(isClaimMapping)),
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

// ValidateGetAuthMethodRequest
func ValidateGetAuthMethodRequest(v *pb.GetAuthMethodRequest) error {
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

func isBExpr(v interface{}) error {
	str := v.(string)
	if str == "" {
		return nil
	}

	_, err := bexpr.CreateEvaluator(str)
	if err != nil {
		return fmt.Errorf("invalid selector: %s", err)
	}

	return nil
}

func isClaimMapping(v interface{}) error {
	m := v.(map[string]string)
	for m, v := range m {
		if m == "" {
			return errors.New("mapping key cannot be empty")
		}

		if !token.IsIdentifier(v) {
			return errors.New(
				"mapping value must be valid identifier made of " +
					"alphanumeric characters and underscores.")
		}
	}

	return nil
}
