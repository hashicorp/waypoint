// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package hcerr

import (
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func TestExternalize(t *testing.T) {
	require := require.New(t)

	t.Run("status is preserved", func(t *testing.T) {
		statusErr := status.Error(codes.NotFound, "original error")
		wrappedStatusErr := errors.Wrapf(statusErr, "intermediate error")

		externalError := Externalize(hclog.Default(), wrappedStatusErr, "message")

		externalizedErr, isStatusErr := status.FromError(externalError)
		require.True(isStatusErr)

		require.Equal(codes.NotFound, externalizedErr.Code())
	})

	t.Run("Non-status input errors turn into internal status errors", func(t *testing.T) {
		externalError := Externalize(hclog.Default(), errors.New("non-status error"), "message")
		externalizedErr, isStatusErr := status.FromError(externalError)
		require.True(isStatusErr)
		require.Equal(codes.Internal, externalizedErr.Code())
	})

	t.Run("Args are surfaced in status details", func(t *testing.T) {
		externalError := Externalize(hclog.Default(), errors.New("original error"), "message", "key", "value")

		statusErr, isStatusErr := status.FromError(externalError)
		require.True(isStatusErr)
		details := statusErr.Details()
		require.Len(details, 1)
		detail, ok := details[0].(*pb.ErrorDetail)
		require.True(ok)
		require.Equal("key", detail.Key)
		require.Equal("value", detail.Value)
	})

	t.Run("Can handle uneven number of args", func(t *testing.T) {
		externalError := Externalize(hclog.Default(), errors.New("original error"), "message", "key", "value", "whoops")
		statusErr, isStatusErr := status.FromError(externalError)
		require.True(isStatusErr)
		details := statusErr.Details()
		require.Len(details, 2)
	})
}

func TestUserError(t *testing.T) {
	log := hclog.New(&hclog.LoggerOptions{Level: hclog.Debug})
	require := require.New(t)

	t.Run("One level of user errors", func(t *testing.T) {
		baseError := errors.New("failed to connect to db at 10.0.0.sensitive")
		userError := UserErrorf(baseError, "missing required field 'foo'. Try resubmitting after specifying 'foo'")

		finalError := Externalize(log, userError, "failed to perform action foobarbaz")

		require.NotContainsf(finalError.Error(), "failed to connect to db at 10.0.0.sensitive", "contains base error")
		require.Contains(finalError.Error(), "missing required field 'foo'. Try resubmitting after specifying 'foo'")
		require.Contains(finalError.Error(), "failed to perform action foobarbaz")
	})

	t.Run("Multiple levels of wrapping", func(t *testing.T) {
		baseError := errors.New("failed to connect to db at 10.0.0.sensitive")
		userError1 := UserErrorf(baseError, "missing required field 'foo'. Try resubmitting after specifying 'foo'")
		middleError := errors.Wrapf(userError1, "failed to do specific internal thing")
		userError2 := UserErrorf(middleError, "invalid value 'baz' for argument 'bar'")

		finalError := Externalize(log, userError2, "failed to perform action foobarbaz")

		require.NotContainsf(finalError.Error(), "failed to connect to db at 10.0.0.sensitive", "contains base error")
		require.Contains(finalError.Error(), "missing required field 'foo'. Try resubmitting after specifying 'foo'")
		require.NotContainsf(finalError.Error(), "failed to do specific internal thing", "contains middle internal error")
		require.Contains(finalError.Error(), "invalid value 'baz' for argument 'bar'")
		require.Contains(finalError.Error(), "failed to perform action foobarbaz")
	})

	t.Run("Supports status codes", func(t *testing.T) {
		internalStatusErr := status.Errorf(codes.Internal, "unknown internal failure")
		userErr := UserErrorWithCodef(codes.InvalidArgument, internalStatusErr, "it's not an internal error actually, you made a mistake o user")

		finalErr := Externalize(log, userErr, "no new information for the user here")
		require.Equal(status.Code(finalErr), codes.InvalidArgument)
	})
}
