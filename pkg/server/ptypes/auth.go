// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ptypes

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/pkg/validationext"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// TestToken returns a valid token for tests
func TestToken(t testing.T, src *pb.Token) *pb.Token {
	t.Helper()

	if src == nil {
		src = &pb.Token{}
	}

	require.NoError(t, mergo.Merge(src, &pb.Token{
		Kind: &pb.Token_Login_{
			Login: &pb.Token_Login{
				UserId: "test",
			},
		},
	}))

	return src
}

func ValidateToken(v *pb.Token) error {
	if v == nil {
		return status.Error(codes.InvalidArgument, "token must not be nil")
	}
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Kind, validation.Required),
	))
}
