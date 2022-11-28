package ptypes

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/pkg/validationext"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// TestPaginationRequest returns a valid pagination request for tests.
func TestPaginationRequest(t testing.T, src *pb.PaginationRequest) *pb.PaginationRequest {
	t.Helper()

	if src == nil {
		src = &pb.PaginationRequest{}
	}

	require.NoError(t, mergo.Merge(src, &pb.PaginationRequest{}))

	return src
}

// ValidatePaginationRequest validates the pagination request structure.
func ValidatePaginationRequest(v *pb.PaginationRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		ValidatePaginationRequestRules(v)...,
	))
}

// ValidatePaginationRequestRules
func ValidatePaginationRequestRules(v *pb.PaginationRequest) []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&v.NextPageToken, validation.When(v.PreviousPageToken != "", validation.Empty.Error("Only one of NextPageToken or PreviousPageToken can be set."))),
	}
}
