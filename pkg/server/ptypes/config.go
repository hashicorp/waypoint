// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ptypes

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hashicorp/waypoint/internal/pkg/validationext"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// ValidateSetConfigSourceRequest
func ValidateSetConfigSourceRequest(v *pb.SetConfigSourceRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.ConfigSource, validation.Required),
	))
}

// ValidateGetConfigSourceRequest
func ValidateGetConfigSourceRequest(v *pb.GetConfigSourceRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Scope, validation.Required),
	))
}

// ValidateGetConfigRequest
func ValidateGetConfigRequest(v *pb.ConfigGetRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Scope, validation.Required),
	))
}

// ValidateSetConfigRequest
func ValidateSetConfigRequest(v *pb.ConfigSetRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Variables, validation.Required),
	))
}

// ValidateDeleteConfigRequest
func ValidateDeleteConfigRequest(v *pb.ConfigDeleteRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Variables, validation.Required),
	))
}
