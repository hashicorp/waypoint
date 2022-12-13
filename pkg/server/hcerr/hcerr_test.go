package hcerr

import (
	"fmt"
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

	t.Run("One level of user errors", func(t *testing.T) {
		baseError := errors.New("base")

		userError := UserErrorf(baseError, "user facing message")

		finalError := Externalize(log, userError, "top-level message")
		fmt.Println(finalError)
	})

	t.Run("Multiple levels of user error", func(t *testing.T) {
		baseError := errors.New("base")
		userError1 := UserErrorf(baseError, "user facing message 1")
		userError2 := UserErrorf(userError1, "user facing message 2")

		finalError := Externalize(log, userError2, "top-level message")
		fmt.Println(finalError)
	})

}
