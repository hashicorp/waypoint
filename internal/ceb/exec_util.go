package ceb

import (
	"io"

	"github.com/golang/protobuf/proto"
	grpc_net_conn "github.com/mitchellh/go-grpc-net-conn"
	"google.golang.org/grpc"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func execOutputWriter(
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
