// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package validationext provides helpers to extend the ozzo-validation.
// There are two primary goals with this package: (1) to ease validating
// deeply nested structures that are common with protobuf-based APIs and
// (2) to convert errors from ozzo-validation into proto InvalidArgument
// errors with field violation extra details.
package validationext
