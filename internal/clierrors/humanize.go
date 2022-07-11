package clierrors

import (
	"fmt"

	"github.com/mitchellh/go-wordwrap"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func Humanize(err error) string {
	if err == nil {
		return ""
	}

	if IsCanceled(err) {
		return "operation canceled"
	}

	v := err.Error()
	if s, ok := status.FromError(err); ok {
		v = s.Message()

		// Include details, if any
		for _, detail := range s.Details() {
			if d, ok := detail.(*pb.ErrorDetail); ok {
				v += fmt.Sprintf("\n%q: %q", d.Key, d.Value)
			}
		}
	}

	return wordwrap.WrapString(v, 80)
}
