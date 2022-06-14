package hcerr

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// Externalize is intended to be called by top-level grpc handler that's
// about to return an error to the framework. Details from the err will be logged,
// but not returned directly to the client, to prevent leaking too much
// detail. If the err includes a grpc status, that, along with the msg and
// args, will be returned to the client.
func Externalize(log hclog.Logger, err error, msg string, args ...interface{}) error {
	log.Error(msg, append(args, "error", err)...)

	// Preserve the proto status
	// status.Status does not support errors.As (https://github.com/grpc/grpc-go/issues/2934)
	var grpcstatus interface{ GRPCStatus() *status.Status }
	var code codes.Code
	if errors.As(err, &grpcstatus) {
		// Otherwise use any code already in the error
		code = grpcstatus.GRPCStatus().Code()
	} else {
		// And if all else fails, default to internal error
		code = codes.Internal
	}

	var details []*anypb.Any
	if len(args) > 0 {

		// Even out the number of pairs
		if len(args)%2 != 0 {
			extra := args[len(args)-1]
			args = append(args[:len(args)-1], hclog.MissingKey, extra)
		}

		for i := 0; i < len(args); i = i + 2 {
			detailPb, err := anypb.New(&pb.ErrorDetail{
				Key:   fmt.Sprintf("%v", args[i]),
				Value: fmt.Sprintf("%v", args[i+1]),
			})
			if err != nil {
				log.Error("Unexpected error marshalling detail k/v pair",
					"key", args[i], "value", args[i+1], "error", err,
				)
			}
			details = append(details, detailPb)
		}
	}

	return status.FromProto(&spb.Status{
		Code:    int32(code),
		Message: msg,
		Details: details,
	}).Err()
}
