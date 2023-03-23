// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hcerr

import (
	"context"
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
// detail. If the error includes a grpc status, that, along with the msg and
// args, will be returned to the client.
//
// Args follow the same pattern as `hclog`. They are expected to be in order of
// "label, variable, label, variable", etc. Example:
// `hcerr.Externalize(log, err, "failed doing thing", "id", thing.id, "organization id", org.id)`
// All args will be printed as strings to transmit to clients, so rather than adding a big complex
// struct as an arg, pull out the fields of interest and add them as multiple args.
// These will be displayed as key/value pairs to the client. If there are an odd number of args,
// this assumes it's a mistake and adds "EXTRA_VALUE_AT_END" as the label for the final arg.
//
// Any UserError errors found in the chain will have their messages added to the final error.
func Externalize(log hclog.Logger, err error, msg string, args ...interface{}) error {

	if errors.Is(err, context.Canceled) {
		log.Trace(msg, append(args, "error", err)...)
	} else {
		log.Error(msg, append(args, "error", err)...)
	}

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

	// Extract all user-facing messages from error chain
	// and add them to the final error message
	currentError := err
	for currentError != nil {
		var userErr *UserError
		if errors.As(currentError, &userErr) {
			msg = fmt.Sprintf("%s\n%s", msg, userErr.UserMessage)
			currentError = userErr.Unwrap()
			// Maybe there are more user errors in the chain under this one
			continue
		}
		// No more user errors in the chain
		break
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

	stat := status.FromProto(&spb.Status{
		Code:    int32(code),
		Message: msg,
		Details: details,
	})

	ret := stat.Err()

	return ret
}
