// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ptypes

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/hashicorp/waypoint/internal/pkg/validationext"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// ValidateUIListEventsRequest
func ValidateUIListEventsRequest(v *pb.UI_ListEventsRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Application, validation.Required),
		validationext.StructField(&v.Pagination, func() []*validation.FieldRules {
			return ValidatePaginationRequestRules(v.Pagination)
		}),
	))
}
