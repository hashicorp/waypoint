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
