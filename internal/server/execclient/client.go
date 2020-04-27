package execclient

import (
	"bytes"
	"context"
	"io"

	"github.com/golang/protobuf/proto"
	"github.com/mitchellh/go-grpc-net-conn"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

type Client struct {
	Context      context.Context
	Client       pb.DevflowClient
	DeploymentId string
	Args         []string
	Stdin        io.Reader
	Stdout       io.Writer
	Stderr       io.Writer
}

func (c *Client) Run() (int, error) {
	// Start our exec stream
	client, err := c.Client.StartExecStream(c.Context)
	if err != nil {
		return 0, err
	}

	// Send the start event
	if err := client.Send(&pb.ExecStreamRequest{
		Event: &pb.ExecStreamRequest_Start_{
			Start: &pb.ExecStreamRequest_Start{
				DeploymentId: c.DeploymentId,
				Args:         c.Args,
			},
		},
	}); err != nil {
		return 0, err
	}

	// Build our connection. We only build the stdin sending side because
	// we can receive other message types from our recv.
	go io.Copy(&grpc_net_conn.Conn{
		Stream:  client,
		Request: &pb.ExecStreamRequest{},
		Encode: grpc_net_conn.SimpleEncoder(func(msg proto.Message) *[]byte {
			req := msg.(*pb.ExecStreamRequest)
			if req.Event == nil {
				req.Event = &pb.ExecStreamRequest_Input_{
					Input: &pb.ExecStreamRequest_Input{},
				}
			}

			return &req.Event.(*pb.ExecStreamRequest_Input_).Input.Data
		}),
	}, c.Stdin)

	// Loop for data
	for {
		// TODO: need a goroutine so we can handle the context
		resp, err := client.Recv()
		if err != nil {
			return 0, err
		}

		switch event := resp.Event.(type) {
		case *pb.ExecStreamResponse_Output_:
			// TODO: stderr
			out := c.Stdout
			io.Copy(out, bytes.NewReader(event.Output.Data))

		case *pb.ExecStreamResponse_Exit_:
			return int(event.Exit.Code), nil
		}
	}
}
