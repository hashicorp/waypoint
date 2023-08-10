// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package validationext

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error takes an error and turns an ozzo-validation.Errors into a
// gRPC status error with field violations populated. If the error is
// nil or not an ozzo-validation error, it is returned as-is.
//
// Note that validate.Validate doesn't return a validate.Errors. Only validation
// on structs and other containers will return the proper structure that will
// be wrapped by this call. This should be used against request structures.
func Error(err error) error {
	// Nil error returns nil directly
	if err == nil {
		return nil
	}

	// If it isn't a validation error, then return it as-is
	verr, ok := err.(validation.Errors)
	if !ok {
		return err
	}

	// Build up the status and accumulate the errors.
	st := status.New(codes.InvalidArgument, verr.Error())

	// This should NEVER fail, we verified with the code that this should
	// never fail. If it does, we panic because it should be impossible.
	st, err = st.WithDetails(&errdetails.BadRequest{
		FieldViolations: errorAppend(nil, "", verr),
	})
	if err != nil {
		panic(err)
	}

	return st.Err()
}

// errorAppend accumulates field violations by recursively nesting into the
// validation errors. We have to recurse to get nested structs/maps/etc.
// With each recursion, we prefix the errors with the field path to that
// error.
func errorAppend(
	v []*errdetails.BadRequest_FieldViolation,
	prefix string,
	verr validation.Errors) []*errdetails.BadRequest_FieldViolation {
	for k, err := range verr {
		field := k
		if prefix != "" {
			field = prefix + "." + field
		}

		// If we have another validation error, then recurse
		if verr, ok := err.(validation.Errors); ok {
			v = errorAppend(v, field, verr)
			continue
		}

		v = append(v, &errdetails.BadRequest_FieldViolation{
			Field:       field,
			Description: err.Error(),
		})
	}

	return v
}
