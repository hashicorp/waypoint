// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

// Package execwriter contains helpers for writing "waypoint exec"
// streams via an io.Writer. Data written to the io.Writer will be
// automatically sent to the gRPC stream.
package execwriter

import (
	"io"

	grpc_net_conn "github.com/hashicorp/go-grpc-net-conn"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// Writer returns an io.Writer for writing the given channel of exec
// stream data (stderr or stdout). The writer doesn't have to be closed,
// you'll receive an EOF once the stream itself closes.
func Writer(
	client grpc.ClientStream,
	channel pb.EntrypointExecRequest_Output_Channel,
) io.Writer {
	return &grpc_net_conn.Conn{
		Stream:  client,
		Request: &pb.EntrypointExecRequest{},
		Encode: grpc_net_conn.SimpleEncoder(func(msg proto.Message) *[]byte {
			req := msg.(*pb.EntrypointExecRequest)
			if req.Event == nil {
				req.Event = &pb.EntrypointExecRequest_Output_{
					Output: &pb.EntrypointExecRequest_Output{
						Channel: channel,
					},
				}
			}

			return &req.Event.(*pb.EntrypointExecRequest_Output_).Output.Data
		}),
	}
}
