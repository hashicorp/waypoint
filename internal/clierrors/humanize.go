package clierrors

import (
	"google.golang.org/grpc/status"
)

func Humanize(err error) string {
	s, ok := status.FromError(err)
	if !ok {
		return err.Error()
	}

	return s.Message()
}
