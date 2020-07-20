package clierrors

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func IsCanceled(err error) bool {
	if err == context.Canceled {
		return true
	}

	s, ok := status.FromError(err)
	if !ok {
		return false
	}

	return s.Code() == codes.Canceled
}
