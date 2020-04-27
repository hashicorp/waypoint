package execclient

import (
	"bytes"
	"context"
	"io"
	"os"

	"github.com/creack/pty"
	"github.com/golang/protobuf/proto"
	"github.com/mitchellh/go-grpc-net-conn"
	"golang.org/x/crypto/ssh/terminal"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	internalptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

type Client struct {
	Context      context.Context
	Client       pb.WaypointClient
	DeploymentId string
	Args         []string
	Stdin        io.Reader
	Stdout       io.Writer
	Stderr       io.Writer
}

func (c *Client) Run() (int, error) {
	// Determine if we should allocate a pty. If we should, we need to send
	// along a TERM value to the remote end that matches our own.
	var ptyReq *pb.ExecStreamRequest_PTY
	if f, ok := c.Stdout.(*os.File); ok && terminal.IsTerminal(int(f.Fd())) {
		ws, err := pty.GetsizeFull(f)
		if err != nil {
			return 0, err
		}

		ptyReq = &pb.ExecStreamRequest_PTY{
			Enable:     true,
			Term:       os.Getenv("TERM"),
			WindowSize: internalptypes.WinsizeProto(ws),
		}

		// We need to go into raw mode with stdin
		if f, ok := c.Stdin.(*os.File); ok {
			oldState, err := terminal.MakeRaw(int(f.Fd()))
			if err != nil {
				return 0, err
			}
			defer func() { _ = terminal.Restore(int(f.Fd()), oldState) }()
		}
	}

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
				Pty:          ptyReq,
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
