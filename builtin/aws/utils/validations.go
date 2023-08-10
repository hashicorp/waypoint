// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package utils

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

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

var fargateResources = map[int][]int{
	512:  {256},
	1024: {256, 512},
	2048: {256, 512, 1024},
	3072: {512, 1024},
	4096: {512, 1024},
	5120: {1024},
	6144: {1024},
	7168: {1024},
	8192: {1024},
}

func init() {
	for i := 4096; i < 16384; i += 1024 {
		fargateResources[i] = append(fargateResources[i], 2048)
	}

	for i := 8192; i <= 30720; i += 1024 {
		fargateResources[i] = append(fargateResources[i], 4096)
	}
}

func ValidateEcsMemCPUPair(mem, cpu int) error {
	cpuValues, ok := fargateResources[mem]
	if !ok {
		var (
			allValues  []int
			goodValues []string
		)

		for k := range fargateResources {
			allValues = append(allValues, k)
		}

		sort.Ints(allValues)

		for _, k := range allValues {
			goodValues = append(goodValues, strconv.Itoa(k))
		}

		return fmt.Errorf("invalid memory value: %d (valid values: %s)", mem,
			strings.Join(goodValues, ", "))
	}

	if cpu == 0 {
		// if cpu is 0 a default will likely be chosen by which ever AWS service
		// is being used, based on the memory value
		return nil
	}

	var (
		valid      bool
		goodValues []string
	)

	for _, c := range cpuValues {
		goodValues = append(goodValues, strconv.Itoa(c))
		if c == cpu {
			valid = true
			break
		}
	}

	if !valid {
		return fmt.Errorf("invalid cpu value: %d (valid values: %s)",
			mem, strings.Join(goodValues, ", "))
	}

	return nil
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
