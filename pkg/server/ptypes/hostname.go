// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ptypes

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hashicorp/waypoint/internal/pkg/validationext"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// ValidateCreateHostnameRequest
func ValidateCreateHostnameRequest(v *pb.CreateHostnameRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Target, validation.Required),
	))
}
